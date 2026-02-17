package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"apollo.io/apolloapiv1"
	"apollo.io/clients"
	"apollo.io/serverutils"
)

var readTimeout = 5 * time.Second
var writeTimeout = 10 * time.Second
var idleTimeout = 15 * time.Second

var failureToConnectToDatabase = 1
var failureToStartServer = 2

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

	hostPort := net.JoinHostPort(postgresHost, postgresPort)

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		postgresUsername,
		postgresPassword,
		hostPort,
		postgresDatabase,
	)
	logger.Info("Establishing Postgres connection")
	pgClient, err := clients.NewPostgresClient(postgresURL)

	if err != nil {
		logger.Error("Failed to connect to the database: \n" + err.Error())
		os.Exit(failureToConnectToDatabase)
	}

	// Forces defer to execute.
	err = func() error {
		defer pgClient.Close()

		serverConfig := serverutils.ServerConfig{
			ServerPort:         serverPort,
			PostgresConnection: pgClient,
			Logger:             logger,
		}

		root := http.NewServeMux()
		root.Handle("/api/v1/", http.StripPrefix("/api/v1", apolloapiv1.NewAPIV1(serverConfig)))
		srv := &http.Server{
			Addr:         serverPort,
			Handler:      root,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		}

		return srv.ListenAndServe()
	}()

	if err != nil {
		logger.Error("Failed to start server: \n" + err.Error())
		os.Exit(failureToStartServer)
	}
}
