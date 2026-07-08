package services

import (
	"MyChatApp/monolithic/internal/models"
	"MyChatApp/monolithic/internal/repositories/db"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type MessageItem struct {
	MessageID       gocql.UUID `json:"message_id"`
	AuthorID        gocql.UUID `json:"author_id"`
	AuthorFirstName string     `json:"author_first_name"`
	AuthorLastName  string     `json:"author_last_name"`
	AuthorUsername  string     `json:"author_username"`
	Content         string     `json:"content"`
	CreatedAt       time.Time  `json:"created_at"`
}

type MessageService struct {
	messagesRepo *db.MessagesRepository
	channelsRepo *db.ChannelsRepository
	usersRepo    *db.UserRepository
}

func NewMessageService(m *db.MessagesRepository, ch *db.ChannelsRepository, u *db.UserRepository) *MessageService {
	return &MessageService{
		messagesRepo: m,
		channelsRepo: ch,
		usersRepo:    u,
	}
}

const maxBucketLookback = 30

// year*10000 + month*100 + day  - so  july 8 2026 becomes 20260708
func bucketFromTime(t time.Time) int {
	t = t.UTC()
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}

func bucketToTime(bucket int) time.Time {
	year := bucket / 10000
	month := time.Month((bucket / 100) % 100)
	day := bucket % 100
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// important - don't just naively subtract one from the integer
func previousBucket(bucket int) int {
	t := bucketToTime(bucket)
	prev := t.AddDate(0, 0, -1)
	return bucketFromTime(prev)
}

func (s *MessageService) GetMessages(channelID, userID string, limit int, before string) ([]MessageItem, error) {
	cid, err := gocql.ParseUUID(channelID)
	if err != nil {
		return nil, err
	}

	uid, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, err
	}

	_, err = s.channelsRepo.GetMemberRole(cid, uid)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, &RequestError{Message: "not a member of this channel"}
		}
		return nil, err
	}

	var beforeTime *time.Time
	if before != "" {
		parsed, err := time.Parse(time.RFC3339, before)
		if err != nil {
			return nil, &RequestError{Message: "invalid before timestamp, use RFC 3339 format"}
		}
		beforeTime = &parsed
	}

	if limit <= 0 || limit > 200 {
		limit = 50
	}

	var startingBucket int
	if beforeTime != nil {
		startingBucket = bucketFromTime(*beforeTime)
	} else {
		startingBucket = bucketFromTime(time.Now())
	}

	var allMessages []models.Message
	currentBucket := startingBucket
	remaining := limit
	lookback := 0

	for remaining > 0 && lookback < maxBucketLookback {
		var cursor *time.Time
		if lookback == 0 && beforeTime != nil {
			cursor = beforeTime
		}

		msgs, err := s.messagesRepo.GetMessagesByChannel(cid, currentBucket, remaining, cursor)
		if err != nil {
			return nil, err
		}

		allMessages = append(allMessages, msgs...)
		remaining -= len(msgs)

		if remaining <= 0 {
			break
		}

		currentBucket = previousBucket(currentBucket)
		lookback++
	}

	if len(allMessages) == 0 {
		return []MessageItem{}, nil
	}

	authorIDs := make(map[gocql.UUID]struct{})
	for _, m := range allMessages {
		authorIDs[m.AuthorID] = struct{}{}
	}

	ids := make([]gocql.UUID, 0, len(authorIDs))
	for id := range authorIDs {
		ids = append(ids, id)
	}

	users, err := s.usersRepo.FindByIDs(ids)
	if err != nil {
		return nil, err
	}

	userMap := make(map[gocql.UUID]models.User)
	for _, u := range users {
		userMap[u.UserID] = u
	}

	items := make([]MessageItem, 0, len(allMessages))
	for _, m := range allMessages {
		item := MessageItem{
			MessageID: m.MessageID,
			AuthorID:  m.AuthorID,
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		}
		if u, ok := userMap[m.AuthorID]; ok {
			item.AuthorFirstName = u.FirstName
			item.AuthorLastName = u.LastName
			item.AuthorUsername = u.Username
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *MessageService) SendMessage(channelID, userID, content string) (*MessageItem, error) {
	cid, err := gocql.ParseUUID(channelID)
	if err != nil {
		return nil, err
	}

	uid, err := gocql.ParseUUID(userID)
	if err != nil {
		return nil, err
	}

	_, err = s.channelsRepo.GetMemberRole(cid, uid)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil, &RequestError{Message: "not a member of this channel"}
		}
		return nil, err
	}

	if content == "" {
		return nil, &RequestError{Message: "content cannot be empty"}
	}

	now := time.Now()
	mid, err := gocql.RandomUUID()
	if err != nil {
		return nil, err
	}

	msg := models.Message{
		ChannelID: cid,
		Bucket:    bucketFromTime(now),
		MessageID: mid,
		AuthorID:  uid,
		Content:   content,
		CreatedAt: now,
	}

	if err := s.messagesRepo.CreateMessage(&msg); err != nil {
		return nil, err
	}

	author, err := s.usersRepo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	return &MessageItem{
		MessageID:       msg.MessageID,
		AuthorID:        msg.AuthorID,
		AuthorFirstName: author.FirstName,
		AuthorLastName:  author.LastName,
		AuthorUsername:  author.Username,
		Content:         msg.Content,
		CreatedAt:       msg.CreatedAt,
	}, nil
}
