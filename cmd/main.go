package main

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/candrairwn/go-pure/api/config"
	"github.com/candrairwn/go-pure/api/delivery/http/controller"
)

//go:embed swagger/openapi.yaml
var openapi []byte

//go:embed swagger
var swaggerUI embed.FS

func main() {
	ctx := context.Background()

	controller.Openapi = openapi
	controller.SwaggerUI = swaggerUI

	if err := config.Run(&config.Bootstrap{Ctx: ctx}); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
