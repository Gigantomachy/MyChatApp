package models

import (
	"time"

	"github.com/gocql/gocql"
)

/*
CREATE TABLE IF NOT EXISTS my_chat_app.users_by_id (
    user_id uuid PRIMARY KEY,
    created_at timestamp,
    email text,
    first_name text,
    last_name text,
    password_hash text,
    username text
);

CREATE TABLE IF NOT EXISTS my_chat_app.users_by_username (
    username_lower text PRIMARY KEY,
    user_id uuid,
    username text
);
*/

type User struct {
	UserID       gocql.UUID `json:"user_id"`
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`
	Email        string     `json:"email"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	CreatedAt    time.Time  `json:"created_at"`
}
