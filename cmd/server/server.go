package main

import (
	"context"
	"flag"
	"net/http"

	"github.com/m-lab/access/controller"
	"github.com/m-lab/access/token"
	"github.com/m-lab/go/flagx"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/ndt-server/ndt7/listener"
	"github.com/m-lab/packet-test/handler"
	"github.com/m-lab/packet-test/static"
)

var (
	ctx, cancel    = context.WithCancel(context.Background())
	dataDir        = flag.String("datadir", "./data", "Path to write data out to.")
	hostname       = flag.String("hostname", "localhost", "Server hostname.")
	address        = flag.String("address", ":80", "Listen address/port for connections.")
	addressSecure  = flag.String("address-secure", ":443", "Listen address/port for secure connections.")
	certFile       = flag.String("cert", "", "The file with server certificates in PEM format.")
	keyFile        = flag.String("key", "", "The file with server key in PEM format.")
	tokenVerifyKey = flagx.FileBytesArray{}
	tokenVerify    bool
	tokenMachine   string
)

func init() {
	flag.Var(&tokenVerifyKey, "token.verify-key", "Public key for verifying access tokens")
	flag.BoolVar(&tokenVerify, "token.verify", false, "Verify access tokens")
	flag.StringVar(&tokenMachine, "token.machine", "", "Use given machine name to verify token claims")
}

func main() {
	flag.Parse()

	v, err := token.NewVerifier(tokenVerifyKey.Get()...)
	if tokenVerify {
		rtx.Must(err, "Failed to load verifier")
	}
	paths := controller.Paths{
		static.NDT7DownloadURLPath: true,
	}
	acm, _ := controller.Setup(ctx, v, tokenVerify, tokenMachine, paths, paths)

	h := handler.New(*dataDir, *hostname)

	mux := http.NewServeMux()
	mux.HandleFunc(static.NDT7DownloadURLPath, http.HandlerFunc(h.NDT7Download))
	srv := &http.Server{
		Addr:    *address,
		Handler: acm.Then(mux),
	}
	rtx.Must(listener.ListenAndServeAsync(srv), "Failed to start server")
	defer srv.Close()

	if *certFile != "" && *keyFile != "" {
		srvSecure := &http.Server{
			Addr:    *addressSecure,
			Handler: acm.Then(mux),
		}
		rtx.Must(listener.ListenAndServeTLSAsync(srvSecure, *certFile, *keyFile), "Failed to start secure server")
		defer srvSecure.Close()
	}

	<-ctx.Done()
	cancel()
}
