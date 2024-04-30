package handler

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"path"
	"time"

	"github.com/m-lab/go/timex"
	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
	log "github.com/sirupsen/logrus"
)

var errNoData = errors.New("failed to receive measurement result")

// Client handles requests for packet tests.
type Client struct {
	dataDir string
}

// New returns a new instance of *Client.
func New(dataDir string) *Client {
	return &Client{
		dataDir: dataDir,
	}
}

// ProcessPacketLoop listens for a kickoff UDP packet and then runs a packet test.
func (c *Client) ProcessPacketLoop(conn net.PacketConn) {
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

		var result interface{}

		switch msg {
		case "pair1":
			err = c.sendPairs(conn, addr)
		case "train1":
			result, err = c.sendTrains(conn, addr)
		}

		err = c.handleResult(conn, msg, err, result)
		if err != nil {
			log.Errorf("Failed %s test: %v", msg, err)
		}
	}
}

func (c *Client) handleResult(conn net.PacketConn, datatype string, err error, data interface{}) error {
	log.Info(data)
	if err != nil {
		return err
	}

	if data == nil {
		return errNoData
	}

	return c.writeMeasurements(conn, datatype, data)
}

func (c *Client) writeMeasurements(conn net.PacketConn, datatype string, data interface{}) error {
	t := time.Now().UTC()
	dir := path.Join(c.dataDir, datatype, t.Format(timex.YYYYMMDDWithSlash))
	log.Info(dir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	filename := path.Join(dir, datatype+"-"+t.Format("20060102T150405.000000000Z")+".json")
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonResult, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	_, err = file.Write(jsonResult)
	return err
}

func sendPacket(conn net.PacketConn, addr net.Addr, pkt *api.Packet) error {
	pkt.Sent = time.Now().UTC().UnixMicro()

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
