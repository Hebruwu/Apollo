package utils

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServerConfig struct {
	ServerPort         string
	PostgresConnection *pgxpool.Pool
}
