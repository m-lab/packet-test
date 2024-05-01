package handler

import (
	"encoding/json"
	"math"
	"net"
	"time"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
	log "github.com/sirupsen/logrus"
)

func (c *Client) sendTrains(conn net.PacketConn, tcpConn *net.TCPListener, addr net.Addr) (*api.Train1Result, error) {
	log.Info("Sending trains")

	result := &api.Train1Result{
		Server: c.hostname,
		Client: addr.String(),
	}
	pkt := &api.Packet{
		Sequence: 0,
		Data:     make([]byte, static.PacketBytes),
	}

	for i := 0; i < static.TrainCount; i++ {
		for j := 0; j < static.TrainLength; j++ {
			err := sendPacket(conn, addr, pkt)
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(static.TrainDelay)
		pkt.Sequence++
	}

	measurements, err := receiveMeasurements(tcpConn)
	if err != nil {
		return nil, err
	}
	result.Measurements = measurements

	return result, nil
}

func receiveMeasurements(listener *net.TCPListener) ([]api.Measurement, error) {
	measurements := make([]api.Measurement, static.TrainCount)

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	buf := make([]byte, math.MaxUint16)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf[:n], &measurements)
	if err != nil {
		return nil, err
	}

	return measurements, nil
}
