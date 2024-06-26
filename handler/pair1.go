package handler

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
)

func (c *Client) sendPairs(conn net.PacketConn, addr net.Addr, gapIncr time.Duration) error {
	log.Info("Sending pairs")

	pkt := &api.Packet{
		Sequence: 0,
		Data:     make([]byte, static.PacketBytes),
	}

	var gap = 0 * time.Microsecond
	for i := 0; i < static.PairCount; i++ {
		err := sendPacket(conn, addr, pkt)
		if err != nil {
			return err
		}

		time.Sleep(gap)
		err = sendPacket(conn, addr, pkt)
		if err != nil {
			return err
		}

		pkt.Sequence++
		gap += gapIncr
		pkt.Gap = gap
		time.Sleep(static.PairDelay)
	}

	return nil
}
