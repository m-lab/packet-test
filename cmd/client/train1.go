package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/apex/log"
	"github.com/m-lab/go/mathx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
)

var (
	server = flag.String("server", ":1053", "Server address")
)

func main() {
	flag.Parse()

	udpSocket, err := net.ResolveUDPAddr("udp", *server)

	rtx.Must(err, "ResolveUDPAddr failed")

	conn, err := net.DialUDP("udp", nil, udpSocket)
	rtx.Must(err, "DialUDP failed")

	_, err = conn.Write([]byte("train1"))
	rtx.Must(err, "Kickoff failed")

	measurements, err := receiveTrains(conn)
	if err != nil {
		log.Errorf("Packed train test failed: %v", err)
	}

	err = sendMeasurements(conn, measurements)
	if err != nil {
		log.Errorf("Failed to send measurements to server: %v", err)
	}
}

func receiveTrains(conn *net.UDPConn) ([]api.Measurement, error) {
	measurements := make([]api.Measurement, static.TrainCount)
	bw := make([]int64, static.TrainCount)

	for i := 0; i < static.TrainCount; i++ {
		train, err := receiveTrain(conn)
		if err != nil {
			return nil, fmt.Errorf("Failed to receive packet train: %v", err)
		}

		delta := getDelta(train[1].Received, train[static.TrainLength-1].Received)
		log.Infof("delta: %d usec", delta)
		bw[i] = (train[1].Size << 3) * (static.TrainLength - 1) / delta
		log.Infof("bw: %d Mbps", bw[i])
		measurements[i] = api.Measurement{
			Packets:    train,
			Dispersion: delta,
			Bandwidth:  bw[i],
		}
	}

	mode, err := mathx.Mode(bw)
	if err != nil {
		return measurements, fmt.Errorf("Failed to calculate bandwidth: %v", err)
	}
	log.Infof("Bandwidth: %d Mbps", mode)

	return measurements, nil
}

func sendMeasurements(conn *net.UDPConn, measurements []api.Measurement) error {
	for _, m := range measurements {
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}

		_, err = conn.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func receiveTrain(conn *net.UDPConn) ([]*api.Received, error) {
	buf := make([]byte, 1024)
	pkts := make([]*api.Received, 0)

	for j := 0; j < static.TrainLength; j++ {
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
			Size:     int64(n),
		}

		pkts = append(pkts, rcvd)
	}
	return pkts, nil
}

// Compute the difference between two timestamps in microseconds.
func getDelta(first time.Time, last time.Time) int64 {
	return (last.Unix()-first.Unix())*1000000 +
		(last.UnixMicro() - first.UnixMicro())
}
