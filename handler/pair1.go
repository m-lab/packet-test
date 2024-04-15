package handler

import (
	"encoding/json"
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
)

// Client handles requests for packet tests.
type Client struct{}

// ProcessPacketLoop listens for a kickoff UDP packet and then runs a packet pair test.
func (h *Client) ProcessPacketLoop(conn net.PacketConn) {
	log.Info("Listening for UDP packets")

	buf := make([]byte, static.BufferBytes)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Errorf("Failed to read UDP packet: %v", err)
			continue
		}

		log.Errorf("Received UDP packet addr: %s, n: %d, data: %s ", addr.String(), n, string(buf[:n]))
		err = sendPairs(conn, addr)
		if err != nil {
			log.Errorf("Failed packet pair: %v", err)
		}
	}
}

func sendPairs(conn net.PacketConn, addr net.Addr) error {
	log.Info("Sending pairs")
	pkt := &api.Packet{
		Sequence: 0,
		Data:     make([]byte, static.PacketBytes),
	}

	for i := 0; i < static.PairCount; i++ {
		err := sendPacket(conn, addr, pkt)
		if err != nil {
			return err
		}

		err = sendPacket(conn, addr, pkt)
		if err != nil {
			return err
		}

		time.Sleep(static.PairGap)
		pkt.Sequence++
	}

	return nil
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
