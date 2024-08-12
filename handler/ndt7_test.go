package handler_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-lab/go/testingx"
	"github.com/m-lab/ndt-server/ndt7/spec"
	"github.com/m-lab/packet-test/static"
	"github.com/m-lab/packet-test/testdata"
)

func TestClient_NDT7Download(t *testing.T) {
	// Start the server.
	h, srv := testdata.NewNDT7Server(t)
	defer os.RemoveAll(h.DataDir)

	// Run a download test.
	URL, _ := url.Parse(srv.URL)
	URL.Scheme = "ws"
	URL.Path = static.NDT7DownloadURLPath
	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", spec.SecWebSocketProtocol)
	headers.Add("User-Agent", "fake-user-agent")
	ctx := context.Background()
	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	conn, _, err := dialer.DialContext(ctx, URL.String(), headers)
	testingx.Must(t, err, "failed to dial websocket ndt7 test")
	err = simpleDownload(ctx, t, conn)
	if err != nil && !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		testingx.Must(t, err, "failed to download")
	}
}

// WARNING: this is not a reference client.
func simpleDownload(ctx context.Context, t *testing.T, conn *websocket.Conn) error {
	defer conn.Close()
	conn.SetReadLimit(spec.MaxMessageSize)
	err := conn.SetReadDeadline(time.Now().Add(spec.MaxRuntime))
	testingx.Must(t, err, "failed to set read deadline")
	_, _, err = conn.ReadMessage()
	if err != nil {
		return err
	}
	// We only read one message, so this is an early close.
	return conn.Close()
}
