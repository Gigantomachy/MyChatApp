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

// Use a lightweight transaction
func (r *UserRepository) CreateUser(user *models.User) (bool, error) {
	usernameLower := strings.ToLower(user.Username)

	// Claim the name first, IF NOT EXISTS runs Paxos
	applied, err := r.session.Query(
		"INSERT INTO users_by_username (username_lower, user_id, username) VALUES (?, ?, ?) IF NOT EXISTS",
		usernameLower, user.UserID, user.Username,
	).MapScanCAS(map[string]interface{}{})
	if err != nil {
		return false, err
	}

	if !applied {
		return false, nil // name taken - 409 later
	}

	// We own the name. Now write the main record with a plain insert.
	if err := r.session.Query(
		"INSERT INTO users_by_id (user_id, username, password_hash, email, first_name, last_name, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.UserID, user.Username, user.PasswordHash, user.Email, user.FirstName, user.LastName, user.CreatedAt,
	).Exec(); err != nil {
		// best-effort rollback of the claim so a failed insert doesn't leave
		// the name reserved forever pointing at a user that doesn't exist
		_ = r.session.Query(
			"DELETE FROM users_by_username WHERE username_lower = ?",
			usernameLower,
		).Exec()
		return false, err
	}

	return true, nil
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

// get multiple users at once
func (r *UserRepository) FindByIDs(userIDs []gocql.UUID) ([]models.User, error) {
	if len(userIDs) == 0 {
		return []models.User{}, nil
	}

	iter := r.session.Query(
		"SELECT user_id, username, first_name, last_name FROM users_by_id WHERE user_id IN ?",
		userIDs,
	).Iter()

	var users []models.User
	var user models.User

	for iter.Scan(&user.UserID, &user.Username, &user.FirstName, &user.LastName) {
		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		return []models.User{}, err
	}

	return users, nil
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	iter := r.session.Query(
		"SELECT user_id, username, email, first_name, last_name, created_at FROM users_by_id",
	).Iter()

	var users []models.User
	var user models.User

	for iter.Scan(&user.UserID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.CreatedAt) {
		users = append(users, user)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return users, nil
}
