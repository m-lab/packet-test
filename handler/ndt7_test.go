package handler_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/m-lab/go/testingx"
	"github.com/m-lab/ndt-server/ndt7/spec"
	"github.com/m-lab/packet-test/handler"
	"github.com/m-lab/packet-test/pkg/ndt7/sender"
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

func Test_getParams(t *testing.T) {
	type args struct {
		urlValues url.Values
	}
	tests := []struct {
		name    string
		vals    url.Values
		want    *sender.Params
		wantErr bool
	}{
		{
			name: static.EarlyExitParameterName,
			vals: url.Values{
				static.EarlyExitParameterName: []string{"10"}, // 10 MB.
			},
			want: &sender.Params{
				MaxBytes: 10000000, // 10000000 Bytes.
			},
			wantErr: false,
		},
		{
			name: static.MaxCwndGainParameterName,
			vals: url.Values{
				static.MaxCwndGainParameterName: []string{"512"},
			},
			want: &sender.Params{
				MaxCwndGain: 512,
			},
			wantErr: false,
		},
		{
			name: static.MaxElapsedTimeParameterName,
			vals: url.Values{
				static.MaxElapsedTimeParameterName: []string{"5"}, // 5 seconds.
			},
			want: &sender.Params{
				MaxElapsedTime: 5000000, // 5000000 microseconds.
			},
			wantErr: false,
		},
		{
			name: static.ImmediateExitParameterName,
			vals: url.Values{
				static.ImmediateExitParameterName: []string{"true"},
			},
			want: &sender.Params{
				ImmediateExit: true,
			},
			wantErr: false,
		},
		{
			name: "error",
			vals: url.Values{
				static.EarlyExitParameterName: []string{"foo"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.GetParams(tt.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handler.GetParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
