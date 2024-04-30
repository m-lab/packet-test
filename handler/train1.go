package handler

import (
	"encoding/json"
	"math"
	"net"
	"time"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
)

func (c *Client) sendTrains(conn net.PacketConn, addr net.Addr) (*api.Train1Result, error) {
	// log.Info("Sending trains")

	result := &api.Train1Result{
		Server:    conn.LocalAddr().String(),
		Client:    addr.String(),
		StartTime: time.Now(),
	}

	for i := 0; i < static.TrainCount; i++ {
		pkt := &api.Packet{
			Sequence: 0,
			Data:     make([]byte, static.PacketBytes),
		}
		for j := 0; j < static.TrainLength; j++ {
			err := sendPacket(conn, addr, pkt)
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(static.TrainDelay)
		pkt.Sequence++
	}

	measurements, err := receiveMeasurements(conn)
	if err != nil {
		return nil, err
	}
	result.Measurements = measurements
	last := measurements[len(measurements)-1].Packets
	result.EndTime = last[len(last)-1].Received

	return result, nil
}

func receiveMeasurements(conn net.PacketConn) ([]api.Measurement, error) {
	buf := make([]byte, math.MaxUint16)
	measurements := make([]api.Measurement, static.TrainCount)

	for i := 0; i < static.TrainCount; i++ {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			return nil, err
		}

		measurement := api.Measurement{}
		err = json.Unmarshal(buf[:n], &measurement)
		if err != nil {
			return nil, err
		}

		measurements[i] = measurement
	}

	return measurements, nil
}
