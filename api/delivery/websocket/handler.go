package websocket

import (
	"encoding/json"
	"log/slog"
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
			log.Error("upgrade error", slog.String("error", err.Error()))
		}

		defer func() {
			if err := conn.Close(); err != nil {
				log.Error("close error", slog.String("error", err.Error()))
			}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Error("read error", slog.String("error", err.Error()))
				break
			}

			log.Info("message received", slog.String("message", string(message)))

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Error("unmarshal error", slog.String("error", err.Error()))
				break
			}

			if err := conn.WriteMessage(websocket.TextMessage, []byte(
				"Hello, I'm a d4ark nich hoohoh "+msg.Nama+"!",
			)); err != nil {
				log.Error("write error", slog.String("error", err.Error()))
				break
			}
		}
	}
}
