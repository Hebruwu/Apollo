package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"apollo.io/api"
	"apollo.io/clients"
	"apollo.io/utils"
)

func main() {
	// Environment config
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	logger.Info("Starting Server...")

	serverPort := ":" + os.Getenv("SERVER_PORT")

	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresUsername := os.Getenv("POSTGRES_USERNAME")
	postgresDatabase := os.Getenv("POSTGRES_DATABASE")
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		postgresUsername,
		postgresPassword,
		postgresHost,
		postgresPort,
		postgresDatabase,
	)
	logger.Info("Establishing Postgres connection")
	pgClient, err := clients.NewPostgresClient(postgresURL)

	if err != nil {
		logger.Error("Failed to connect to the database: \n" + err.Error())
		os.Exit(1)
	}

	defer pgClient.Close()

	serverConfig := utils.ServerConfig{
		ServerPort:         serverPort,
		PostgresConnection: pgClient,
		Logger:             logger,
	}

	root := http.NewServeMux()
	root.Handle("/api/v1/", http.StripPrefix("/api/v1", api.NewAPIV1(serverConfig)))

	err = http.ListenAndServe(serverPort, root)

	if err != nil {
		logger.Error("Failed to start server: \n" + err.Error())
		os.Exit(2)
	}
}
