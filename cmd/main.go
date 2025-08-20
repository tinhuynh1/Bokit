package main

import (
	"booking-svc/internal/bootstrap"
	"log"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("failed to bootstrap app: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
