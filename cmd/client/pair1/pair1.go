package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/apex/log"
	"github.com/m-lab/go/mathx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/pkg/client"
	"github.com/m-lab/packet-test/static"
)

const (
	datatype = "pair1"
)

var (
	server = flag.String("server", "localhost", "Server address")
)

func main() {
	flag.Parse()

	// Set up UDP connection to run the test.
	udpSocket, err := net.ResolveUDPAddr("udp", *server+":1053")
	rtx.Must(err, "ResolveUDPAddr failed")

	conn, err := net.DialUDP("udp", nil, udpSocket)
	rtx.Must(err, "DialUDP failed")

	_, err = conn.Write([]byte(datatype))
	rtx.Must(err, "Kickoff failed")

	measurements, err := receivePairs(conn)
	if err != nil {
		log.Errorf("Packet pair test failed: %v", err)
	}

	err = client.SendMeasurements(*server+":8080", datatype, measurements)
	if err != nil {
		log.Errorf("Failed to send measurements to server: %v", err)
	}
}

func receivePairs(conn *net.UDPConn) ([]api.Measurement, error) {
	measurements := make([]api.Measurement, static.PairCount)
	bw := make([]int64, static.PairCount)
	var sum int64

	for i := 0; i < static.PairCount; i++ {
		pair, err := client.ReceiveTrain(conn, 2)
		if err != nil {
			return nil, fmt.Errorf("Failed to receive packet pair: %v", err)
		}

		delta := client.GetDelta(pair[0].Received, pair[1].Received)
		log.Infof("delta: %d usec", delta)
		bw[i] = pair[1].Size * 8.0 / delta
		log.Infof("bw: %d Mbps", bw[i])
		sum += bw[i]

		measurements[i] = api.Measurement{
			Packets: pair,
			Metrics: api.Metrics{
				Dispersion: delta,
				Bandwidth:  bw[i],
			},
		}
	}

	log.Infof("Bandwidth: %d Mbps", sum/static.PairCount)

	return measurements, nil
}

func receiveTrains(conn *net.UDPConn) ([]api.Measurement, error) {
	measurements := make([]api.Measurement, static.TrainCount)
	bw := make([]int64, static.TrainCount)

	for i := 0; i < static.TrainCount; i++ {
		train, err := client.ReceiveTrain(conn, static.TrainLength)
		if err != nil {
			return nil, fmt.Errorf("Failed to receive packet train: %v", err)
		}

		delta := client.GetDelta(train[1].Received, train[static.TrainLength-1].Received)
		log.Infof("delta: %d usec", delta)
		bw[i] = (train[1].Size << 3) * (static.TrainLength - 1) / delta
		log.Infof("bw: %d Mbps", bw[i])

		measurements[i] = api.Measurement{
			Packets: train,
			Metrics: api.Metrics{
				Dispersion: delta,
				Bandwidth:  bw[i],
			},
		}
	}

	mode, err := mathx.Mode(bw)
	if err != nil {
		return measurements, fmt.Errorf("Failed to calculate bandwidth: %v", err)
	}
	log.Infof("Bandwidth: %d Mbps", mode)

	return measurements, nil
}
