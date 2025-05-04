package main

import (
	"go-rest-user-api/api"
	"go-rest-user-api/database"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := Init(); err != nil {
		slog.Error("Não foi possível iniciar o servidor")
		os.Exit(1)
	}

	slog.Info("Servidor iniciado na port 8080")
}

func Init() error {
	database, err := database.InitDatabase()
	if err != nil || database == nil {
		return err
	}

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Handler:      api.NewHander(*database),
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
