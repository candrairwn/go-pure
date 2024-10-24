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
	WebsocketHandler *websocket.WebsocketHandler
}

func NewRouteConfig(log *zap.SugaredLogger, version string) *RouteConfig {
	return &RouteConfig{
		Mux:              http.NewServeMux(),
		Log:              log,
		Version:          version,
		HealthController: &controller.HealthController{},
		WebsocketHandler: websocket.NewWebsocketHandler(),
	}
}

func Route(log *zap.SugaredLogger, version string) http.Handler {
	// Create a new RouteConfig instance
	config := NewRouteConfig(log, version)
	config.SetupGuestRoutes()

	// Setup middleware
	handler := middleware.Accesslog(config.Mux, log)
	handler = middleware.Recovery(handler, log)

	return handler
}

func (c *RouteConfig) SetupGuestRoutes() {
	c.Mux.HandleFunc("/health", c.HealthController.HandleGetHealth(c.Version))
	c.Mux.HandleFunc("/ws", c.WebsocketHandler.Run(c.Log))
}
