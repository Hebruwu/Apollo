package utils

import (
	"log/slog"

	"apollo.io/clients"
)

type ServerConfig struct {
	ServerPort         string
	PostgresConnection clients.PostgresClient
	Logger             *slog.Logger
}
