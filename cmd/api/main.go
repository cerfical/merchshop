package main

import (
	"context"
	"os"

	"github.com/cerfical/merchshop/internal/api"
	"github.com/cerfical/merchshop/internal/config"
	"github.com/cerfical/merchshop/internal/httpserv"
	"github.com/cerfical/merchshop/internal/log"
	"github.com/cerfical/merchshop/internal/postgres"
)

func main() {
	config := config.MustLoad(os.Args)
	log := log.New(&config.Log)

	db, err := postgres.Open(&config.DB)
	if err != nil {
		log.Fatal("Failed to open the database", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Failed to close the database", err)
		}
	}()

	serv := httpserv.New(
		&config.Server,
		api.New(db, log),
		log,
	)

	if err := serv.Run(context.Background()); err != nil {
		log.Error("The server terminated abnormally", err)
	}
}
