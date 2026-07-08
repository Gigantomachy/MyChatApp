package models

import (
	"github.com/gocql/gocql"
	"time"
)

type Channel struct {
	ChannelID gocql.UUID `json:"channel_id"`
	Name      string     `json:"name"`
	Type      string     `json:"type"` // public, dm, group
	CreatedBy gocql.UUID `json:"created_by"`
	CreatedAt time.Time  `json:"created_at"`
}

type ChannelMembership struct {
	ChannelID   gocql.UUID `json:"channel_id"`
	ChannelName string     `json:"channel_name"`
	ChannelType string     `json:"channel_type"`
	JoinedAt    time.Time  `json:"joined_at"`
}

type Member struct {
	UserID   gocql.UUID `json:"user_id"`
	Role     string     `json:"role"`
	JoinedAt time.Time  `json:"joined_at"`
}

type Message struct {
	ChannelID gocql.UUID `json:"channel_id"`
	Bucket    int        `json:"-"`
	MessageID gocql.UUID `json:"message_id"`
	AuthorID  gocql.UUID `json:"author_id"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
}
