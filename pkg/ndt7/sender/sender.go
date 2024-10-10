package sender

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/m-lab/ndt-server/ndt7/closer"
	"github.com/m-lab/ndt-server/ndt7/model"
	"github.com/m-lab/ndt-server/ndt7/ping"
	"github.com/m-lab/ndt-server/ndt7/spec"
	"github.com/m-lab/packet-test/pkg/ndt7/measurer"
)

// Params defines the parameters for the sender to end the test early.
type Params struct {
	MaxBytes       int64  // TCPInfo.BytesAcked is of type int64 (bytes).
	MaxElapsedTime int64  // TCPInfo.ElapsedTime is of type int64 (microseconds).
	MaxCwndGain    uint32 // BBRInfo.CwndGain is of type uint32.
	ImmediateExit  bool
}

func (p *Params) isEarlyExitLimit() bool {
	return p.MaxBytes > 0
}

func (p *Params) isMaxElapsedTimeLimit() bool {
	return p.MaxElapsedTime > 0
}

func (p *Params) isMaxCwndGainLimit() bool {
	return p.MaxCwndGain > 0
}

func (p *Params) isEarlyExitDone(m model.Measurement) bool {
	return m.TCPInfo.BytesAcked >= p.MaxBytes
}

func (p *Params) isMaxElapsedTimeDone(m model.Measurement) bool {
	return m.TCPInfo.ElapsedTime >= p.MaxElapsedTime
}

func (p *Params) isMaxCwndGainDone(m model.Measurement) bool {
	return m.BBRInfo.CwndGain >= p.MaxCwndGain
}

func makePreparedMessage(size int) (*websocket.PreparedMessage, error) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	return websocket.NewPreparedMessage(websocket.BinaryMessage, data)
}

// Start sends binary messages (bulk download) and measurement messages (status
// messages) to the client conn. Each measurement message will also be saved to
// data.
//
// Liveness guarantee: the sender will not be stuck sending for more than the
// MaxRuntime of the subtest. This is enforced by setting the write deadline to
// Time.Now() + MaxRuntime.
func Start(ctx context.Context, conn *websocket.Conn, data *model.ArchivalData, params *Params) error {
	// Start collecting connection measurements. Measurements will be sent to
	// src until DefaultRuntime, when the src channel is closed.
	mr := measurer.New(conn, data.UUID)
	src := mr.Start(ctx, spec.DefaultRuntime)
	defer mr.Stop(src)

	bulkMessageSize := 1 << 13
	preparedMessage, err := makePreparedMessage(bulkMessageSize)
	if err != nil {
		log.Errorf("sender: makePreparedMessage failed", err)
		return err
	}
	deadline := time.Now().Add(spec.MaxRuntime)
	err = conn.SetWriteDeadline(deadline) // Liveness!
	if err != nil {
		log.Errorf("sender: conn.SetWriteDeadline failed", err)
		return err
	}

	// Record measurement start time, and prepare recording of the endtime on return.
	data.StartTime = time.Now().UTC()
	defer func() {
		data.EndTime = time.Now().UTC()
	}()
	var totalSent int64
	for {
		select {
		case m, ok := <-src:
			if !ok { // This means that the measurer has terminated.
				closer.StartClosing(conn)
				return nil
			}

			if err := conn.WriteJSON(m); err != nil {
				log.Errorf("sender: conn.WriteJSON failed", err)
				return err
			}
			// Only save measurements sent to the client.
			data.ServerMeasurements = append(data.ServerMeasurements, m)
			if err := ping.SendTicks(conn, deadline); err != nil {
				log.Errorf("sender: ping.SendTicks failed", err)
				return err
			}

			// Check if the test should be terminated early.
			if m.TCPInfo != nil && terminateTest(params, m) {
				closer.StartClosing(conn)
				return nil
			}

		default:
			if err := conn.WritePreparedMessage(preparedMessage); err != nil {
				log.Errorf("sender: conn.WritePreparedMessage failed", err)
				return err
			}

			if terminateImmediately(params, mr, conn, data) {
				closer.StartClosing(conn)
				return nil
			}

			// The following block of code implements the scaling of message size
			// as recommended in the spec's appendix. We're not accounting for the
			// size of JSON messages because that is small compared to the bulk
			// message size. The net effect is slightly slowing down the scaling,
			// but this is currently fine. We need to gather data from large
			// scale deployments of this algorithm anyway, so there's no point
			// in engaging in fine grained calibration before knowing.
			totalSent += int64(bulkMessageSize)
			if int64(bulkMessageSize) >= spec.MaxScaledMessageSize {
				continue // No further scaling is required
			}
			if int64(bulkMessageSize) > totalSent/spec.ScalingFraction {
				continue // message size still too big compared to sent data
			}
			bulkMessageSize *= 2
			preparedMessage, err = makePreparedMessage(bulkMessageSize)
			if err != nil {
				log.Errorf("sender: makePreparedMessage failed", err)
				return err
			}
		}
	}
}

func terminateImmediately(p *Params, mr *measurer.Measurer, conn *websocket.Conn, data *model.ArchivalData) bool {
	if !p.ImmediateExit {
		return false
	}

	m, err := mr.GetMeasurement()
	if err != nil {
		return false
	}

	if !terminateTest(p, m) {
		return false
	}

	// Send and save the last snapshot.
	if err := conn.WriteJSON(m); err != nil {
		log.Errorf("sender: conn.WriteJSON failed", err)
	}
	data.ServerMeasurements = append(data.ServerMeasurements, m)

	return true
}

func terminateTest(p *Params, m model.Measurement) bool {
	switch {
	case p.isMaxCwndGainLimit() && p.isMaxElapsedTimeLimit():
		if p.isMaxCwndGainDone(m) && p.isMaxElapsedTimeDone(m) {
			log.Infof("sender: terminating test after %d CwndGain and %d ElapsedTime (µs)", m.BBRInfo.CwndGain, m.TCPInfo.ElapsedTime)
			return true
		}
	case p.isMaxCwndGainLimit() && p.isEarlyExitLimit():
		if p.isMaxCwndGainDone(m) && p.isEarlyExitDone(m) {
			log.Infof("sender: terminating test after %d CwndGain and %d BytesAcked", m.BBRInfo.CwndGain, m.TCPInfo.BytesAcked)
			return true
		}
	case p.isMaxCwndGainLimit() && p.isMaxCwndGainDone(m):
		log.Infof("sender: terminating test after %d CwndGain", m.BBRInfo.CwndGain)
		return true
	case p.isEarlyExitLimit() && p.isEarlyExitDone(m):
		log.Infof("sender: terminating test after %d BytesAcked", m.TCPInfo.BytesAcked)
		return true
	case p.isMaxElapsedTimeLimit() && p.isMaxElapsedTimeDone(m):
		log.Infof("sender: terminating test after %d ElapsedTime (µs)", m.TCPInfo.ElapsedTime)
		return true
	}

	return false
}
