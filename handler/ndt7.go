package handler

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
	"github.com/m-lab/go/prometheusx"
	"github.com/m-lab/ndt-server/data"
	"github.com/m-lab/ndt-server/metadata"
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
		log.Errorf("Failed to establish connection: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get client parameters.
	params, err := GetParams(req.URL.Query())
	if err != nil {
		log.Errorf("Invalid parameters: %v", req.URL.Query())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get data.
	data, err := getData(conn)
	if err != nil {
		log.Errorf("Failed to get test data: %v", data)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	appendClientMetadata(data, req.URL.Query())

	// Set up result.
	result := setupResult(conn)
	result.StartTime = time.Now().UTC()
	result.Download = data

	defer func() {
		result.EndTime = time.Now().UTC()
		err = c.writeMeasurements("ndt7", result)
		if err != nil {
			log.Errorf("Failed to write measurement result: %v", err)
		}
	}()

	// Run test.
	err = sender.Start(context.Background(), conn, data, params)
	if err != nil {
		log.Errorf("Failed to run test: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetParams interprets and returns the client parameters.
func GetParams(urlValues url.Values) (*sender.Params, error) {
	params := &sender.Params{}

	for name, values := range urlValues {
		value := values[0]
		name = strings.TrimPrefix(name, "client_")

		switch name {
		case static.EarlyExitParameterName:
			bytes, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			params.MaxBytes = bytes * 1000000 // Convert MB to bytes.
		case static.MaxCwndGainParameterName:
			cwnd, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return nil, err
			}
			params.MaxCwndGain = uint32(cwnd)
		case static.MaxElapsedTimeParameterName:
			time, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return nil, err
			}
			params.MaxElapsedTime = time * 1e6 // Convert seconds to microseconds.
		case static.ImmediateExitParameterName:
			immExit, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			params.ImmediateExit = immExit
		}
	}
	return params, nil
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

// setupResult creates an NDT7Result from the given conn.
func setupResult(conn *websocket.Conn) *data.NDT7Result {
	// NOTE: unless we plan to run the NDT server over different protocols than TCP,
	// then we expect RemoteAddr and LocalAddr to always return net.TCPAddr types.
	clientAddr := netx.ToTCPAddr(conn.RemoteAddr())
	if clientAddr == nil {
		clientAddr = &net.TCPAddr{IP: net.ParseIP("::1"), Port: 1}
	}
	serverAddr := netx.ToTCPAddr(conn.LocalAddr())
	if serverAddr == nil {
		serverAddr = &net.TCPAddr{IP: net.ParseIP("::1"), Port: 1}
	}
	result := &data.NDT7Result{
		GitShortCommit: prometheusx.GitShortCommit,
		ClientIP:       clientAddr.IP.String(),
		ClientPort:     clientAddr.Port,
		ServerIP:       serverAddr.IP.String(),
		ServerPort:     serverAddr.Port,
	}
	return result
}

// excludeKeyRe is a regexp for excluding request parameters from client metadata.
var excludeKeyRe = regexp.MustCompile("^server_")

// appendClientMetadata adds |values| to the archival client metadata contained
// in the request parameter values. Some select key patterns will be excluded.
func appendClientMetadata(data *model.ArchivalData, values url.Values) {
	for name, values := range values {
		if matches := excludeKeyRe.MatchString(name); matches {
			continue // Skip variables that should be excluded.
		}
		data.ClientMetadata = append(
			data.ClientMetadata,
			metadata.NameValue{
				Name:  name,
				Value: values[0], // NOTE: this will ignore multi-value parameters.
			})
	}
}
