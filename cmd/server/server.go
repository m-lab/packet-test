package main

import (
	"context"
	"flag"
	"net"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/packet-test/handler"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	dataDir     = flag.String("datadir", "./data", "Path to write data out to.")
	hostname    = flag.String("hostname", "localhost", "Server hostname.")
)

func main() {
	flag.Parse()

	// Set up TCP connection for results.
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	rtx.Must(err, "ResolveTCPAddr failed")

	tcpConn, err := net.ListenTCP("tcp", tcpAddr)
	rtx.Must(err, "ListenTCP")
	defer tcpConn.Close()

	// Set up UDP connection to run the test.
	addr, err := net.ResolveUDPAddr("udp", ":1053")
	rtx.Must(err, "ResolveUDPAddr failed")

	conn, err := net.ListenUDP("udp", addr)
	rtx.Must(err, "ListenUDP failed")
	defer conn.Close()

	h := handler.New(*dataDir, *hostname)
	go h.ProcessPacketLoop(conn, tcpConn)

	<-ctx.Done()
	cancel()
}
