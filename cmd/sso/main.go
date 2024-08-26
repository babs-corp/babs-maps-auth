package main

import (
	"fmt"

	"gituhb.com/babs-corp/babs-maps-auth/internal/config"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger (slog)

	// TODO: init app

	// TODO: run grpc server
}