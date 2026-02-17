package clients

import (
	"context"
	"errors"
	"fmt"

	"apollo.io/objects/servershared"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCloser interface {
	Close()
}

type PostgresWriter interface {
	AddUser(ctx context.Context, user servershared.User) error
}

type PostgresClient interface {
	PostgresWriter
	PostgresCloser
}

// pool is deliberately left non-anonymouse to prevent idiots (me)
// from accidentally using it directly and breaking the dependency
// injection pattern.
type postgresClient struct {
	pool *pgxpool.Pool
}

func NewPostgresClient(dsn string) (PostgresClient, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w due to %w", ErrFailedToConnect, err)
	}
	return postgresClient{pool}, nil
}

func (p postgresClient) Close() {
	p.pool.Close()
}

func (p postgresClient) AddUser(ctx context.Context, user servershared.User) error {
	query := "INSERT INTO users (username, email, password_hash, salt) VALUES ($1, $2, $3, $4)"
	_, err := p.pool.Exec(ctx, query, user.Username, user.Email, user.PasswordHash, user.Salt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_pkey" {
			return servershared.ErrUsernameAlreadyExists
		}
		return fmt.Errorf("inserting user: %w", err)
	}
	return nil
}
