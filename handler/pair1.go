package handler

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
)

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

		time.Sleep(static.PairDelay)
		pkt.Sequence++
	}

	return nil
}
