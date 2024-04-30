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
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":1053")
	rtx.Must(err, "ResolveUDPAddr failed")

	conn, err := net.ListenUDP("udp", addr)
	rtx.Must(err, "ListenUDP failed")
	defer conn.Close()

	h := handler.New(*dataDir)
	go h.ProcessPacketLoop(conn)

	<-ctx.Done()
	cancel()
}
