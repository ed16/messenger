package repository_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ed16/messenger/domain"
	"github.com/ed16/messenger/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	user := &domain.User{
		Username:     "testuser",
		Status:       1,
		CreatedAt:    time.Now(),
		PasswordHash: "hashedpassword",
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs(user.Username, strconv.Itoa(int(user.Status)), user.CreatedAt, user.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

	err = userRepo.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.UserId)
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	user := &domain.User{
		UserId:       1,
		Username:     "testuser",
		Status:       1,
		CreatedAt:    time.Now(),
		PasswordHash: "hashedpassword",
	}

	mock.ExpectExec(`UPDATE users`).
		WithArgs(user.Username, strconv.Itoa(int(user.Status)), user.CreatedAt, user.PasswordHash, user.UserId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = userRepo.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	expectedUser := domain.User{
		UserId:       1,
		Username:     "testuser",
		Status:       1,
		CreatedAt:    time.Now(),
		PasswordHash: "hashedpassword",
	}

	mock.ExpectQuery(`SELECT user_id, username, status, created_at, password_hash FROM users WHERE user_id = \$1`).
		WithArgs(expectedUser.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "status", "created_at", "password_hash"}).
			AddRow(expectedUser.UserId, expectedUser.Username, expectedUser.Status, expectedUser.CreatedAt, expectedUser.PasswordHash))

	user, err := userRepo.GetUserByID(context.Background(), expectedUser.UserId)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	expectedUser := domain.User{
		UserId:       1,
		Username:     "testuser",
		PasswordHash: "hashedpassword",
	}

	mock.ExpectQuery(`SELECT user_id, username, password_hash FROM users WHERE username = \$1`).
		WithArgs(expectedUser.Username).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username", "password_hash"}).
			AddRow(expectedUser.UserId, expectedUser.Username, expectedUser.PasswordHash))

	user, err := userRepo.GetUserByUsername(context.Background(), expectedUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}

func TestGetUsersByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	expectedUser := domain.User{
		UserId:   1,
		Username: "testuser",
	}

	mock.ExpectQuery(`SELECT user_id, username FROM users WHERE username LIKE '%' \|\| \$1 \|\| '%' and status = '1' LIMIT 100`).
		WithArgs(expectedUser.Username).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "username"}).
			AddRow(expectedUser.UserId, expectedUser.Username))

	users, err := userRepo.GetUsersByUsername(context.Background(), expectedUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, expectedUser, users[0])
}

func TestCreateUserContact(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	contact := &domain.Contact{
		UserId:        1,
		ContactUserId: 2,
		CreatedAt:     time.Now(),
	}

	mock.ExpectExec(`INSERT INTO contacts`).
		WithArgs(contact.UserId, contact.ContactUserId, contact.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = userRepo.CreateUserContact(context.Background(), contact)
	assert.NoError(t, err)
}

func TestGetUserContactsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := repository.NewUserRepo(db)

	expectedContact := domain.Contact{
		UserId:        1,
		ContactUserId: 2,
		CreatedAt:     time.Now(),
	}

	mock.ExpectQuery(`SELECT user_id, contact_user_id, created_at FROM contacts WHERE user_id = \$1`).
		WithArgs(expectedContact.UserId).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "contact_user_id", "created_at"}).
			AddRow(expectedContact.UserId, expectedContact.ContactUserId, expectedContact.CreatedAt))

	contacts, err := userRepo.GetUserContactsByUserID(context.Background(), expectedContact.UserId)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(contacts))
	assert.Equal(t, expectedContact, contacts[0])
}
