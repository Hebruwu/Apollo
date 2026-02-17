package services_test

import (
	"context"
	"testing"

	"apollo.io/objects/servershared"
	"apollo.io/services"
	"github.com/stretchr/testify/assert"
)

type mockPostgresClient struct {
	users  []servershared.User
	closed bool
	err    error
}

func (m *mockPostgresClient) AddUser(_ context.Context, user servershared.User) error {
	if m.closed {
		panic("client is closed")
	}
	if m.err != nil {
		return m.err
	}
	m.users = append(m.users, user)
	return nil
}

func (m *mockPostgresClient) Close() {
	m.closed = true
}

func TestCreateUser_SavesUsernameToDatabase(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	err := service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", mock.users[0].Username)
}

func TestCreateUser_SavesEmailToDatabase(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	err := service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "test@example.com", mock.users[0].Email)
}

func TestCreateUser_GeneratesSixteenByteSalt(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	err := service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.Len(t, mock.users[0].Salt, 16)
}

func TestCreateUser_GeneratesThirtyTwoBytePasswordHash(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	err := service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.Len(t, mock.users[0].PasswordHash, 32)
}

func TestCreateUser_DoesNotStoreRawPassword(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	err := service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEqual(t, []byte("password123"), mock.users[0].PasswordHash)
}

func TestCreateUser_GeneratesUniqueSaltPerUser(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	_ = service.CreateUser(context.Background(), "user1", "user1@example.com", "password")
	_ = service.CreateUser(context.Background(), "user2", "user2@example.com", "password")

	assert.NotEqual(t, mock.users[0].Salt, mock.users[1].Salt)
}

func TestCreateUser_SamePasswordProducesDifferentHashes(t *testing.T) {
	mock := &mockPostgresClient{}
	service := services.NewUserService(mock, nil)

	_ = service.CreateUser(context.Background(), "user1", "user1@example.com", "samepassword")
	_ = service.CreateUser(context.Background(), "user2", "user2@example.com", "samepassword")

	assert.NotEqual(t, mock.users[0].PasswordHash, mock.users[1].PasswordHash)
}

func TestCreateUser_PropagatesClientError(t *testing.T) {
	mock := &mockPostgresClient{err: servershared.ErrUsernameAlreadyExists}
	service := services.NewUserService(mock, nil)

	err := service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.ErrorIs(t, err, servershared.ErrUsernameAlreadyExists)
}

func TestCreateUser_DoesNotSaveUserWhenClientFails(t *testing.T) {
	mock := &mockPostgresClient{err: servershared.ErrUsernameAlreadyExists}
	service := services.NewUserService(mock, nil)

	_ = service.CreateUser(context.Background(), "testuser", "test@example.com", "password123")

	assert.Empty(t, mock.users)
}
