package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ed16/messenger/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (m *UserRepository) InsertUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, status, created_at, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id
	`

	// Default values for created_at and status
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	if user.Status == 0 {
		user.Status = 1 // Assuming 1 as default active status
	}

	err := m.DB.QueryRowContext(ctx, query, user.Username, user.Status, user.CreatedAt, user.PasswordHash).Scan(&user.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `
		SELECT user_id, username, status, created_at, password_hash
		FROM users
		WHERE username = $1
	`

	var user domain.User
	err := m.DB.QueryRowContext(ctx, query, username).Scan(&user.UserId, &user.Username, &user.Status, &user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}

func (m *UserRepository) GetUserByID(ctx context.Context, userID int64) (domain.User, error) {
	query := `
		SELECT user_id, username, status, created_at, password_hash
		FROM users
		WHERE user_id = $1
	`

	var user domain.User
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&user.UserId, &user.Username, &user.Status, &user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, err
	}

	return user, nil
}

// Add contact or change password
func (m *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET username = $1, status = $2, created_at = $3, password_hash = $4
		WHERE user_id = $5
	`

	_, err := m.DB.ExecContext(ctx, query, user.Username, user.Status, user.CreatedAt, user.PasswordHash, user.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserRepository) UpdateUserProfile(ctx context.Context, profile *domain.Profile) error {
	return nil
}

// GetUsersByUsername makes search with the like condition <*username*> Status = 0
func (m *UserRepository) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	query := `
		SELECT user_id, username, status, created_at, password_hash
		FROM users
		WHERE username = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.UserId, &user.Username, &user.Status, &user.CreatedAt, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
