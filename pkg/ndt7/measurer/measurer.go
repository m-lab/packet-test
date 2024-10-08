// Package measurer collects metrics from a socket connection
// and returns them for consumption.
package measurer

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/m-lab/go/memoryless"
	"github.com/m-lab/ndt-server/logging"
	"github.com/m-lab/ndt-server/ndt7/model"
	"github.com/m-lab/ndt-server/ndt7/spec"
	"github.com/m-lab/ndt-server/netx"
)

var (
	BBREnabled = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ndt7_measurer_bbr_enabled_total",
			Help: "A counter of every attempt to enable bbr.",
		},
		[]string{"status", "error"},
	)
)

// Measurer performs measurements
type Measurer struct {
	conn     *websocket.Conn
	uuid     string
	ticker   *memoryless.Ticker
	start    time.Time
	connInfo *model.ConnectionInfo
}

// New creates a new measurer instance
func New(conn *websocket.Conn, UUID string) *Measurer {
	return &Measurer{
		conn: conn,
		uuid: UUID,
	}
}

func (m *Measurer) getSocketAndPossiblyEnableBBR() (netx.ConnInfo, error) {
	ci := netx.ToConnInfo(m.conn.UnderlyingConn())
	err := ci.EnableBBR()
	success := "true"
	errstr := ""
	if err != nil {
		success = "false"
		errstr = err.Error()
		uuid, _ := ci.GetUUID() // to log error with uuid.
		logging.Logger.WithError(err).Warn("Cannot enable BBR: " + uuid)
		// FALLTHROUGH
	}
	BBREnabled.WithLabelValues(success, errstr).Inc()
	return ci, nil
}

func measure(measurement *model.Measurement, ci netx.ConnInfo, elapsed time.Duration) {
	// Implementation note: we always want to sample BBR before TCPInfo so we
	// will know from TCPInfo if the connection has been closed.
	t := int64(elapsed / time.Microsecond)
	bbrinfo, tcpInfo, err := ci.ReadInfo()
	if err == nil {
		measurement.BBRInfo = &model.BBRInfo{
			BBRInfo:     bbrinfo,
			ElapsedTime: t,
		}
		measurement.TCPInfo = &model.TCPInfo{
			LinuxTCPInfo: tcpInfo,
			ElapsedTime:  t,
		}
	}
}

func (m *Measurer) loop(ctx context.Context, timeout time.Duration, dst chan<- model.Measurement) {
	logging.Logger.Debug("measurer: start")
	defer logging.Logger.Debug("measurer: stop")
	defer close(dst)
	measurerctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	_, err := m.getSocketAndPossiblyEnableBBR()
	if err != nil {
		logging.Logger.WithError(err).Warn("getSocketAndPossiblyEnableBBR failed")
		return
	}
	m.start = time.Now()
	m.connInfo = &model.ConnectionInfo{
		Client: m.conn.RemoteAddr().String(),
		Server: m.conn.LocalAddr().String(),
		UUID:   m.uuid,
	}
	// Implementation note: the ticker will close its output channel
	// after the controlling context is expired.
	ticker, err := memoryless.NewTicker(measurerctx, memoryless.Config{
		Min:      spec.MinPoissonSamplingInterval,
		Expected: spec.AveragePoissonSamplingInterval,
		Max:      spec.MaxPoissonSamplingInterval,
	})
	if err != nil {
		logging.Logger.WithError(err).Warn("memoryless.NewTicker failed")
		return
	}
	m.ticker = ticker
	for range ticker.C {
		measurement, _ := m.GetMeasurement()
		dst <- measurement // Liveness: this is blocking
	}
}

// GetMeasurement retrieves the latest measurement snapshot.
func (m *Measurer) GetMeasurement() (model.Measurement, error) {
	t := int64(time.Now().Sub(m.start) / time.Microsecond)
	ci := netx.ToConnInfo(m.conn.UnderlyingConn())
	bbr, tcp, err := ci.ReadInfo()
	if err != nil {
		log.Errorf("measurer: ReadInfo failed", err)
		return model.Measurement{}, err
	}

	return model.Measurement{
		ConnectionInfo: m.connInfo,
		BBRInfo: &model.BBRInfo{
			BBRInfo:     bbr,
			ElapsedTime: t,
		},
		TCPInfo: &model.TCPInfo{
			LinuxTCPInfo: tcp,
			ElapsedTime:  t,
		},
	}, nil
}

// Start runs the measurement loop in a background goroutine and emits
// the measurements on the returned channel.
//
// Liveness guarantee: the measurer will always terminate after
// the given timeout, provided that the consumer continues reading from the
// returned channel. Measurer may be stopped early by canceling ctx, or by
// calling Stop.
func (m *Measurer) Start(ctx context.Context, timeout time.Duration) <-chan model.Measurement {
	dst := make(chan model.Measurement)
	go m.loop(ctx, timeout, dst)
	return dst
}

// Stop ends the measurements and drains the measurement channel. Stop
// guarantees that the measurement goroutine completes by draining the
// measurement channel. Users that call Start should also call Stop.
func (m *Measurer) Stop(src <-chan model.Measurement) {
	if m.ticker != nil {
		m.ticker.Stop()
	}
	for range src {
		// make sure we drain the channel, so the measurement loop can exit.
	}
}
