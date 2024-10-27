package config

import (
	"net/http"

	"github.com/candrairwn/go-pure/api/delivery/http/controller"
	"github.com/candrairwn/go-pure/api/delivery/http/routes"
	"github.com/candrairwn/go-pure/api/delivery/websocket"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BootstrapWireConfig struct {
	App   *http.ServeMux
	Log   *zap.SugaredLogger
	Viper *viper.Viper
	DB    *gorm.DB
}

func BootsrapWire(config *BootstrapWireConfig) http.Handler {

	// Set up controller
	healthController := controller.NewHealthController(config.DB, config.Viper.GetString("APP_VERSION"))
	userController := controller.NewUserController(config.Log)

	// Set up websocket
	websocketHandler := websocket.NewWebsocketHandler(config.Log)

	// Set up routes
	routeConfig := routes.RouteConfig{
		Mux:              config.App,
		Log:              config.Log,
		Version:          config.Viper.GetString("APP_VERSION"),
		UserController:   userController,
		HealthController: healthController,
		WebsocketHandler: websocketHandler,
	}

	handler := routeConfig.Setup()

	return handler
}
