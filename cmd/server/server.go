package main

import (
	"context"
	"flag"
	"net"
	"net/http"

	"github.com/m-lab/go/rtx"
	"github.com/m-lab/ndt-server/ndt7/listener"
	"github.com/m-lab/packet-test/handler"
	"github.com/m-lab/packet-test/static"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	dataDir     = flag.String("datadir", "./data", "Path to write data out to.")
	hostname    = flag.String("hostname", "localhost", "Server hostname.")
)

func main() {
	flag.Parse()

	// Set up UDP connection to run the test.
	addr, err := net.ResolveUDPAddr("udp", ":1053")
	rtx.Must(err, "ResolveUDPAddr failed")

	conn, err := net.ListenUDP("udp", addr)
	rtx.Must(err, "ListenUDP failed")
	defer conn.Close()

	h := handler.New(*dataDir, *hostname)
	go h.ProcessPacketLoop(conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/v0/result", http.HandlerFunc(h.HandleResult))
	mux.HandleFunc(static.NDT7DownloadURLPath, http.HandlerFunc(h.NDT7Download))
	srv := &http.Server{
		Addr:    ":9998",
		Handler: mux,
	}
	rtx.Must(listener.ListenAndServeAsync(srv), "Failed to start server")
	defer srv.Close()

	<-ctx.Done()
	cancel()
}
