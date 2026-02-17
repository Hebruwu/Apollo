package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"

	"apollo.io/clients"
	"apollo.io/objects/shared"
	"golang.org/x/crypto/argon2"
)

type UserService struct {
	pgClient clients.PostgresClient
	logger   *slog.Logger
	hashTime uint8
	memory   uint32
	threads  uint8
	keySize  uint32
}

func NewUserService(pgClient clients.PostgresClient, logger *slog.Logger) UserService {
	if pgClient == nil {
		panic("pgClient is nil")
	}
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return UserService{
		pgClient: pgClient,
		logger:   logger,
		hashTime: 3,
		memory:   64 * 1024,
		threads:  4,
		keySize:  32,
	}
}

func (us UserService) CreateUser(
	ctx context.Context,
	username string,
	email string,
	password string,
) error {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("generating salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, uint32(us.hashTime), us.memory, us.threads, us.keySize)

	newUser := shared.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Salt:         salt,
	}

	return us.pgClient.AddUser(ctx, newUser)
}
