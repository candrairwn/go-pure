package main

import (
	"context"
	"fmt"
	"os"

	"github.com/candrairwn/go-pure/api/config"
)

func main() {
	ctx := context.Background()

	if err := config.Run(&config.Bootstrap{Ctx: ctx}); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
