package websocket

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebsocketHandler struct {
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

func NewWebsocketHandler() *WebsocketHandler {
	return &WebsocketHandler{}
}

func (w *WebsocketHandler) Run(log *zap.SugaredLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("upgrade error", zap.String("error", err.Error()))
		}

		defer func() {
			if err := conn.Close(); err != nil {
				log.Error("close error", zap.String("error", err.Error()))
			}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("read error", zap.String("error", err.Error()))
				break
			}

			log.Infow("message received", "message", string(message))

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Errorw("unmarshal error", "error", err.Error())
				break
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(
				"Hello, I'm a d4ark nich hoohoh "+msg.Nama+"!",
			)); err != nil {
				log.Errorw("write error", "error", err.Error())
				break
			}
		}
	}
}
