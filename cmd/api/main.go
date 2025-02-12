package main

import (
	"context"
	"os"

	"github.com/cerfical/merchshop/internal/api"
	"github.com/cerfical/merchshop/internal/config"
	"github.com/cerfical/merchshop/internal/domain/services"
	"github.com/cerfical/merchshop/internal/httpserv"
	"github.com/cerfical/merchshop/internal/infrastructure/bcrypt"
	"github.com/cerfical/merchshop/internal/infrastructure/jwt"
	"github.com/cerfical/merchshop/internal/infrastructure/postgres"
	"github.com/cerfical/merchshop/internal/log"
)

func main() {
	config := config.MustLoad(os.Args)
	log := log.New(&config.Log)

	db, err := postgres.NewStorage(&config.DB)
	if err != nil {
		log.Fatal("Failed to open the database", err)
	}

	if err := db.UpMigrations(); err != nil {
		log.Fatal("Failed to apply migrations to the database", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Failed to close the database", err)
		}
	}()

	tokens := jwt.NewTokenAuth(&config.API.Auth.Token)
	auth := services.NewAuthService(db, bcrypt.NewHasher(), tokens)

	serv := httpserv.New(&config.API.Server, api.NewHandler(auth, log), log)
	if err := serv.Run(context.Background()); err != nil {
		log.Error("The server terminated abnormally", err)
	}
}
