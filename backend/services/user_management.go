package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"os"

	"apollo.io/clients"
	"apollo.io/objects/servershared"
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

var hashTime uint8 = 3
var memory uint32 = 64 * 1024
var threads uint8 = 4
var keySize uint32 = 32
var saltSize uint32 = 16

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
		hashTime: hashTime,
		memory:   memory,
		threads:  threads,
		keySize:  keySize,
	}
}

func (us UserService) CreateUser(
	ctx context.Context,
	username string,
	email string,
	password string,
) error {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("generating salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, uint32(us.hashTime), us.memory, us.threads, us.keySize)

	newUser := servershared.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Salt:         salt,
	}
	err := us.pgClient.AddUser(ctx, newUser)
	if err != nil {
		return fmt.Errorf("adding user: %w", err)
	}
	return nil
}
