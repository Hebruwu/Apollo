package utils

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

type ServerConfig struct {
	ServerPort         string
	PostgresConnection *pgxpool.Pool
	Logger             *slog.Logger
}
