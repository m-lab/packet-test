package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/m-lab/go/timex"
	"github.com/m-lab/packet-test/api"
	"github.com/m-lab/packet-test/static"
	log "github.com/sirupsen/logrus"
)

// Client handles requests for packet tests.
type Client struct {
	dataDir  string
	hostname string
}

// New returns a new instance of *Client.
func New(dataDir string, hostname string) *Client {
	return &Client{
		dataDir:  dataDir,
		hostname: hostname,
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

		switch msg {
		case "pair1":
			err = c.sendPairs(conn, addr, static.PairGap)
		case "train1":
			err = c.sendTrains(conn, addr)
		}

	}
}

// HandleResult receives the measurement results from the client and writes them out to
// `datadir`.
func (c *Client) HandleResult(rw http.ResponseWriter, req *http.Request) {
	measurements := make([]api.Measurement, 0)
	err := json.NewDecoder(req.Body).Decode(&measurements)
	if err != nil {
		log.Errorf("Failed to decide measurement result: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	result := api.Result{
		Server:       c.hostname,
		Client:       req.RemoteAddr,
		Measurements: measurements,
	}

	err = c.writeMeasurements(req.URL.Query().Get("datatype"), result)
	if err != nil {
		log.Errorf("Failed to write measurement out: %v", err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (c *Client) writeMeasurements(datatype string, data interface{}) error {
	t := time.Now().UTC()
	dir := path.Join(c.dataDir, datatype, t.Format(timex.YYYYMMDDWithSlash))
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

func receiveMeasurements(listener *net.TCPListener) ([]api.Measurement, error) {
	measurements := make([]api.Measurement, 0)

	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	buf := make([]byte, math.MaxUint16)
	n, err := conn.Read(buf)

	fmt.Println(string(buf[:n]))
	fmt.Println(n)
	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf[:n], &measurements)
	if err != nil {
		return nil, err
	}

	return measurements, nil
}
