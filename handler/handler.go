package handler

import (
	"encoding/json"
	"net"
	"time"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
	log "github.com/sirupsen/logrus"
)

// Client handles requests for packet tests.
type Client struct{}

// ProcessPacketLoop listens for a kickoff UDP packet and then runs a packet test.
func (h *Client) ProcessPacketLoop(conn net.PacketConn) {
	log.Info("Listening for UDP packets")

	buf := make([]byte, static.BufferBytes)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Errorf("Failed to read UDP packet: %v", err)
			continue
		}

		msg := string(buf[:n])
		log.Infof("Received UDP packet addr: %s, n: %d, type: %s ", addr.String(), n, msg)

		switch msg {
		case "pair1":
			err = sendPairs(conn, addr)
		case "train1":
			err = sendTrains(conn, addr)
		}

		if err != nil {
			log.Errorf("Failed %s test: %v", msg, err)
		}
	}
}

func sendPacket(conn net.PacketConn, addr net.Addr, pkt *api.Packet) error {
	pkt.Timestamp = time.Now().UTC().UnixMicro()

	m, err := json.Marshal(pkt)
	if err != nil {
		return err
	}

	_, err = conn.WriteTo(m, addr)
	if err != nil {
		return err
	}

	return nil
}
