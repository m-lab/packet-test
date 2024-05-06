package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/m-lab/packet-test/api"
)

func ReceiveTrain(conn *net.UDPConn, length int) ([]*api.Received, error) {
	buf := make([]byte, 1024)
	pkts := make([]*api.Received, 0)

	for j := 0; j < length; j++ {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}

		var pkt = &api.Packet{}
		err = json.Unmarshal(buf[:n], pkt)
		if err != nil {
			return nil, err
		}

		t := time.Now().UTC()
		rcvd := &api.Received{
			Sequence: pkt.Sequence,
			Sent:     time.UnixMicro(pkt.Sent),
			Received: t,
			Gap:      pkt.Gap.Microseconds(),
			Size:     int64(n),
		}

		pkts = append(pkts, rcvd)
	}
	return pkts, nil
}

// GetDelta computes the difference between two timestamps in microseconds.
func GetDelta(first time.Time, last time.Time) int64 {
	return (last.Unix()-first.Unix())*1000000 +
		(last.UnixMicro() - first.UnixMicro())
}

func SendMeasurements(addr string, datatype string, measurements []api.Measurement) error {
	b, err := json.Marshal(measurements)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s/v0/result?datatype=%s", addr, datatype)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(b))
	return err
}
