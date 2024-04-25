package main

import (
	"encoding/json"
	"flag"
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

	receiveTrains(conn)
}

func receiveTrains(conn *net.UDPConn) {
	measurements := make([]int64, static.TrainCount)

	for i := 0; i < static.TrainCount; i++ {
		train, err := receiveTrain(conn)
		if err != nil {
			log.Errorf("Failed to receive packet train: %v", err)
			return
		}

		delta := getDelta(train[1].Received, train[static.TrainLength-1].Received)
		bw := (train[1].Size << 3) * (static.TrainLength - 1) / delta
		measurements[i] = bw
	}

	mode, err := mathx.Mode(measurements)
	if err != nil {
		log.Errorf("Failed to calculate bandwidth: %v", err)
		return
	}
	log.Infof("Bandwidth: %d Mbps", mode)
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
			Packet:   pkt,
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
