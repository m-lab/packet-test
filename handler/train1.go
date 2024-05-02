package handler

import (
	"net"
	"time"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
	log "github.com/sirupsen/logrus"
)

func (c *Client) sendTrains(conn net.PacketConn, addr net.Addr) error {
	log.Info("Sending trains")

	pkt := &api.Packet{
		Sequence: 0,
		Data:     make([]byte, static.PacketBytes),
	}

	for i := 0; i < static.TrainCount; i++ {
		for j := 0; j < static.TrainLength; j++ {
			err := sendPacket(conn, addr, pkt)
			if err != nil {
				return err
			}
		}
		time.Sleep(static.TrainDelay)
		pkt.Sequence++
	}

	return nil
}
