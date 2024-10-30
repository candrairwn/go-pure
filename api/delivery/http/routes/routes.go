package routes

import (
	"net/http"

	"github.com/candrairwn/go-pure/api/delivery/http/controller"
	"github.com/candrairwn/go-pure/api/delivery/http/middleware"
	"github.com/candrairwn/go-pure/api/delivery/websocket"
	"go.uber.org/zap"
)

type RouteConfig struct {
	Mux               *http.ServeMux
	Log               *zap.SugaredLogger
	AppEnv            string
	Version           string
	HealthController  *controller.HealthController
	OpenApiController *controller.OpenApiController
	UserController    *controller.UserController
	WebsocketHandler  *websocket.WebsocketHandler
}

func (c *RouteConfig) Setup() http.Handler {
	c.SetupGuestRoutes()
	c.SetupUserRoutes()

	if c.AppEnv == "development" {
		c.Log.Info("Openapi enabled")
		c.SetupOpenapi()
	}

	handler := middleware.Accesslog(c.Mux, c.Log)
	handler = middleware.Recovery(handler, c.Log)
	handler = middleware.Cors(handler)

	return handler

}

func (c *RouteConfig) SetupGuestRoutes() {
	c.Mux.HandleFunc("GET /api/health", c.HealthController.HandleGetHealth())
	c.Mux.HandleFunc("GET /ws", c.WebsocketHandler.Broadcast)
}

func (c *RouteConfig) SetupUserRoutes() {
	c.Mux.HandleFunc("POST /api/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupOpenapi() {
	c.Mux.Handle("GET /api/openapi.yaml", c.OpenApiController.HandleGetServeOpenapi(c.Version))
	c.Mux.Handle("GET /api/docs/", c.OpenApiController.HandleGetSwaggerUI())
}
