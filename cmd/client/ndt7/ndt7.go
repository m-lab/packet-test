package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/apex/log"

	"github.com/gorilla/websocket"
	"github.com/m-lab/go/rtx"
	"github.com/m-lab/ndt-server/ndt7/model"
)

var (
	server = flag.String("server", "localhost", "Server address")
	params = flag.String("params", "", "Client paramerters")
)

func main() {
	flag.Parse()

	dialer := websocket.Dialer{}
	ctx, cancel := context.WithCancel(context.Background())
	url := fmt.Sprintf("ws://%s:9998/v0/ndt7?%s", *server, *params)
	conn, _, err := dialer.DialContext(ctx, url, http.Header{})
	rtx.Must(err, "Dial failed", err)

	for {
		mtype, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if mtype != websocket.BinaryMessage {
			m := &model.Measurement{}
			err = json.Unmarshal(msg, m)
			if err != nil {
				log.Errorf("Failed to unmarshal measurement", err)
				continue
			}
			elapsed := float64(m.TCPInfo.ElapsedTime) / 1e06
			throughput := (8.0 * float64(m.TCPInfo.BytesAcked)) /
				elapsed / (1000.0 * 1000.0)
			log.Infof("Throughput: %f Mbit/s", throughput)
		}
	}

	<-ctx.Done()
	cancel()
}
