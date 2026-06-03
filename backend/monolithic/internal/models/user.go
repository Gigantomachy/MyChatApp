package models

import (
	"time"

	"github.com/gocql/gocql"
)

type User struct {
	UserID       gocql.UUID `json:"user_id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`
	Email        string     `json:"email"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	CreatedAt    time.Time  `json:"created_at"`
}