package measurer

import (
	"context"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/m-lab/ndt-server/ndt7/measurer"
	"github.com/m-lab/ndt-server/ndt7/model"
	"github.com/m-lab/ndt-server/netx"
)

// Monitor measures network metrics.
type Monitor struct {
	*measurer.Measurer
	conn     *websocket.Conn
	connInfo *model.ConnectionInfo
	start    time.Time
}

// New creates a new Monitor instance.
func New(conn *websocket.Conn, UUID string) *Monitor {
	return &Monitor{
		Measurer: measurer.New(conn, UUID),
		conn:     conn,
		connInfo: &model.ConnectionInfo{
			Client: conn.RemoteAddr().String(),
			Server: conn.LocalAddr().String(),
			UUID:   UUID,
		},
	}
}

// Start runs the measurement loop and records the start time.
func (m *Monitor) Start(ctx context.Context, timeout time.Duration) <-chan model.Measurement {
	m.start = time.Now()
	return m.Measurer.Start(ctx, timeout)
}

// GetMeasurement retrieves the latest measurement snapshot.
func (m *Monitor) GetMeasurement() (model.Measurement, error) {
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
