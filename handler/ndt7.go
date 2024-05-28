package handler

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/m-lab/ndt-server/ndt7/model"
	"github.com/m-lab/ndt-server/ndt7/spec"
	"github.com/m-lab/ndt-server/netx"
	"github.com/m-lab/packet-test/pkg/ndt7/sender"
	"github.com/m-lab/packet-test/static"
)

// NDT7Download runs an ndt7 download test.
func (c *Client) NDT7Download(rw http.ResponseWriter, req *http.Request) {
	// Set up websocket.
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1 << 20,
		WriteBufferSize: 1 << 20,
	}
	headers := http.Header{}
	headers.Add("Sec-WebSocket-Protocol", spec.SecWebSocketProtocol)
	conn, err := upgrader.Upgrade(rw, req, headers)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get client parameters (e.g., early_exit, bbr_exit).
	params, err := getParams(req.URL.Query())

	// Get data.
	data, err := getData(conn)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Run test.
	err = sender.Start(context.Background(), conn, data, params)
}

func getData(conn *websocket.Conn) (*model.ArchivalData, error) {
	ci := netx.ToConnInfo(conn.UnderlyingConn())
	uuid, err := ci.GetUUID()
	if err != nil {
		return nil, err
	}
	data := &model.ArchivalData{
		UUID: uuid,
	}
	return data, nil
}

func getParams(values url.Values) (*sender.Params, error) {
	params := &sender.Params{}
	for name, values := range values {
		switch name {
		case static.EarlyExitParameterName:
			bytes, _ := strconv.ParseInt(values[0], 10, 64)
			params.IsEarlyExit = true
			params.MaxBytes = bytes * 1000000 // Conver MB to bytes.
		case static.BBRExitParameterName:
			params.IsBBRExit = true
		}
	}
	return params, nil
}
