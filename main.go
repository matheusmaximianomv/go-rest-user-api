package main

import (
	"go-rest-user-api/api"
	"go-rest-user-api/database"
	databaseFile "go-rest-user-api/database/file"

	// databaseMemory "go-rest-user-api/database/memory"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := Init(); err != nil {
		slog.Error("Unable to start http server", "error", err)
		os.Exit(1)
	}

	slog.Info("Server start in port 8080")
}

func Init() error {

	dbFile := databaseFile.DatabaseFile{}
	// dbMemory := databaseMemory.DatabaseMemory{}
	db, err := database.InitDatabase(&dbFile)
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Handler:      api.NewHander(db),
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
