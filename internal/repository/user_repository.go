package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (m *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, status, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING user_id
	`

	err := m.DB.QueryRowContext(ctx, query, user.Username, user.Status, user.PasswordHash, user.CreatedAt).Scan(&user.UserId)
	if err != nil {
		return fmt.Errorf("error while creating a new user: %w", err)
	}

	return nil
}

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

func (m *UserRepository) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `
		SELECT user_id, username, status, password_hash
		FROM users
		WHERE username = $1 AND status = $2
	`

	var user domain.User
	err := m.DB.QueryRowContext(ctx, query, username, domain.UserStatusActive).Scan(&user.UserId, &user.Username, &user.Status, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (m *UserRepository) GetUsersByUsername(ctx context.Context, username string) ([]domain.User, error) {
	query := `
		SELECT user_id, username
		FROM users
		WHERE username LIKE '%' || $1 || '%' AND status = $2
		LIMIT 100
	`

	rows, err := m.DB.QueryContext(ctx, query, username, domain.UserStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.UserId, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return users, domain.ErrNotFound
	}

	return users, nil
}

func (m *UserRepository) CreateUserContact(ctx context.Context, contact *domain.Contact) error {
	query := `
		INSERT INTO contacts (user_id, contact_user_id, created_at)
		VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, query, contact.UserId, contact.ContactUserId, contact.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserRepository) GetUserContactsByUserID(ctx context.Context, userID int64) ([]domain.Contact, error) {
	query := `
		SELECT user_id, contact_user_id, created_at
		FROM contacts
		WHERE user_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []domain.Contact
	for rows.Next() {
		var contact domain.Contact
		err := rows.Scan(&contact.UserId, &contact.ContactUserId, &contact.CreatedAt)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
