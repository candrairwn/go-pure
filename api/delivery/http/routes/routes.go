package routes

import (
	"net/http"

	"github.com/candrairwn/go-pure/api/delivery/http/controller"
	"github.com/candrairwn/go-pure/api/delivery/http/middleware"
	"github.com/candrairwn/go-pure/api/delivery/websocket"
	"go.uber.org/zap"
)

type RouteConfig struct {
	Mux              *http.ServeMux
	Log              *zap.SugaredLogger
	Version          string
	HealthController *controller.HealthController
	UserController   *controller.UserController
	WebsocketHandler *websocket.WebsocketHandler
}

func (c *RouteConfig) Setup() http.Handler {
	c.SetupGuestRoutes()
	c.SetupUserRoutes()

	handler := middleware.Accesslog(c.Mux, c.Log)
	handler = middleware.Recovery(handler, c.Log)

	return handler

}

func (c *RouteConfig) SetupGuestRoutes() {
	c.Mux.HandleFunc("GET /health", c.HealthController.HandleGetHealth())
	c.Mux.HandleFunc("GET /ws", c.WebsocketHandler.Broadcast)
}

func (c *RouteConfig) SetupUserRoutes() {
	// create user with v1 version with strip prefix
	c.Mux.HandleFunc("POST /login", c.UserController.Login)
}
