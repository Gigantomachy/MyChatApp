package db

import (
	"MyChatApp/monolithic/internal/models"
	"time"

	"github.com/gocql/gocql"
)

type MessagesRepository struct {
	session *gocql.Session
}

func NewMessagesRepository(session *gocql.Session) *MessagesRepository {
	return &MessagesRepository{session: session}
}

func (r *MessagesRepository) GetMessagesByChannel(
	channelID gocql.UUID,
	bucket int,
	limit int,
	before *time.Time,
) ([]models.Message, error) {
	var query string
	var args []interface{}

	if before != nil {
		query = "SELECT message_id, author_id, content, created_at FROM messages_by_channel WHERE channel_id = ? AND bucket = ? AND created_at < ? ORDER BY created_at DESC, message_id ASC LIMIT ?"
		args = []interface{}{channelID, bucket, *before, limit}
	} else {
		query = "SELECT message_id, author_id, content, created_at FROM messages_by_channel WHERE channel_id = ? AND bucket = ? ORDER BY created_at DESC, message_id ASC LIMIT ?"
		args = []interface{}{channelID, bucket, limit}
	}

	iter := r.session.Query(query, args...).Iter()

	var messages []models.Message
	var msg models.Message

	for iter.Scan(&msg.MessageID, &msg.AuthorID, &msg.Content, &msg.CreatedAt) {
		msg.ChannelID = channelID
		msg.Bucket = bucket
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessagesRepository) CreateMessage(msg *models.Message) error {
	return r.session.Query(
		"INSERT INTO messages_by_channel (channel_id, bucket, created_at, message_id, author_id, content) VALUES (?, ?, ?, ?, ?, ?)",
		msg.ChannelID, msg.Bucket, msg.CreatedAt, msg.MessageID, msg.AuthorID, msg.Content,
	).Exec()
}
