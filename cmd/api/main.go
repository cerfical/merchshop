package main

import (
	"context"
	"net/http"
	"os"

	"github.com/cerfical/merchshop/internal/config"
	"github.com/cerfical/merchshop/internal/httpserv"
	"github.com/cerfical/merchshop/internal/log"
)

func main() {
	config := config.MustLoad(os.Args)
	log := log.New(&config.Log)

	emptyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	serv := httpserv.New(&config.API, http.HandlerFunc(emptyHandler), log)
	if err := serv.Run(context.Background()); err != nil {
		log.Error("The server terminated abnormally", err)
	}
}
