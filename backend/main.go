package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"apollo.io/api"
	"apollo.io/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Environment config
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

	pgPool, err := pgxpool.New(context.Background(), postgresURL)

	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	defer pgPool.Close()

	serverConfig := utils.ServerConfig{
		ServerPort:         serverPort,
		PostgresConnection: pgPool,
	}

	root := http.NewServeMux()
	root.Handle("/api/v1/", http.StripPrefix("/api/v1", api.NewAPIV1(serverConfig)))

	log.Fatal(http.ListenAndServe(serverPort, root))

}
