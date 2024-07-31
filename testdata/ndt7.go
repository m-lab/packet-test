package testdata

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-lab/go/testingx"
	"github.com/m-lab/ndt-server/netx"
	"github.com/m-lab/packet-test/handler"
	"github.com/m-lab/packet-test/static"
)

// NewNDT7Server creates a local httptest server for unit tests.
func NewNDT7Server(t *testing.T) (*handler.Client, *httptest.Server) {
	dir := t.TempDir()

	handler := handler.New(dir, "fake-hostname")
	ndt7Mux := http.NewServeMux()
	ndt7Mux.Handle(static.NDT7DownloadURLPath, http.HandlerFunc(handler.NDT7Download))

	// Create unstarted so we can setup a custom netx.Listener.
	srv := httptest.NewUnstartedServer(ndt7Mux)
	listener, err := net.Listen("tcp", ":0")
	testingx.Must(t, err, "failed to allocate a listening tcp socket")
	srv.Listener = netx.NewListener(listener.(*net.TCPListener))
	// Now that the test server has our custom listener, start it.
	srv.Start()

	return handler, srv
}
