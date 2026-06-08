package models

import (
	"github.com/gocql/gocql"
	"time"
)

/*
-- Incoming requests (partition by recipient so "my pending" is a single query)
CREATE TABLE friend_requests (
    recipient_id uuid,
    sender_id uuid,
    status text,        -- 'pending' | 'accepted' | 'rejected'
    created_at timestamp,
    PRIMARY KEY (recipient_id, sender_id)
);

-- Mutual friendships, stored bidirectionally
CREATE TABLE friendships (
    user_id uuid,
    friend_id uuid,
    created_at timestamp,
    PRIMARY KEY (user_id, friend_id)
);
*/

type FriendShip struct {
	UserID    gocql.UUID `json:"user_id"`
	FriendID  gocql.UUID `json:"friend_id"`
	CreatedAt time.Time  `json:"created_at"`
}

type FriendRequest struct {
	RecipientID gocql.UUID `json:"recipient_id"`
	SenderID    gocql.UUID `json:"sender_id"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
}
