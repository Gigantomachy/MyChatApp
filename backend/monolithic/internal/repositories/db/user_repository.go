package db

import (
	"MyChatApp/monolithic/internal/models"
	"strings"

	"github.com/gocql/gocql"
)

type UserRepository struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) *UserRepository {
	return &UserRepository{session: session}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	batch := r.session.NewBatch(gocql.LoggedBatch)
	batch.Query(
		"INSERT INTO users_by_id (user_id, username, password_hash, email, first_name, last_name, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.UserID, user.Username, user.PasswordHash, user.Email, user.FirstName, user.LastName, user.CreatedAt,
	)
	batch.Query(
		"INSERT INTO users_by_username (username_lower, user_id, username) VALUES (?, ?, ?)",
		strings.ToLower(user.Username), user.UserID, user.Username,
	)
	return r.session.ExecuteBatch(batch)
}

func (r *UserRepository) FindByUsername(usernameLower string) (*models.User, error) {
	var userID gocql.UUID
	var storedUsername string
	if err := r.session.Query(
		"SELECT user_id, username FROM users_by_username WHERE username_lower = ?",
		usernameLower,
	).Scan(&userID, &storedUsername); err != nil {
		return nil, err
	}
	return r.FindByID(userID)
}

func (r *UserRepository) FindByID(userID gocql.UUID) (*models.User, error) {
	var user models.User
	if err := r.session.Query(
		"SELECT user_id, username, password_hash, email, first_name, last_name, created_at FROM users_by_id WHERE user_id = ?",
		userID,
	).Scan(
		&user.UserID, &user.Username, &user.PasswordHash,
		&user.Email, &user.FirstName, &user.LastName, &user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}