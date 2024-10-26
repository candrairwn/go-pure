package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebsocketHandler struct {
	Log *zap.SugaredLogger
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Nama string `json:"nama"`
}

func NewWebsocketHandler(log *zap.SugaredLogger) *WebsocketHandler {
	return &WebsocketHandler{
		Log: log,
	}
}

func (wh *WebsocketHandler) Broadcast(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {

		wh.Log.Error("upgrade error", zap.String("error", err.Error()))
	}

	defer func() {
		if err := conn.Close(); err != nil {
			wh.Log.Error("close error", zap.String("error", err.Error()))
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			wh.Log.Error("read error", zap.String("error", err.Error()))
			break
		}

		wh.Log.Infow("message received", "message", string(message))

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			wh.Log.Errorw("unmarshal error", "error", err.Error())
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(
			"Hello, I'm a d4ark nich hoohoh "+msg.Nama+"!",
		)); err != nil {
			wh.Log.Errorw("write error", "error", err.Error())
			break
		}
	}
}
