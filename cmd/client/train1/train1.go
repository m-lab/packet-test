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
	datatype = "train1"
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

	measurements, err := receiveTrains(conn)
	if err != nil {
		log.Errorf("Packet train test failed: %v", err)
	}

	err = client.SendMeasurements(*server+":9998", datatype, measurements)
	if err != nil {
		log.Errorf("Failed to send measurements to server: %v", err)
	}
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
