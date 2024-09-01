package main

import (
	"log"

	"github.com/a-x-a/gophkeeper/internal/keeper/app"
	"github.com/a-x-a/gophkeeper/internal/keeper/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
