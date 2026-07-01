package db

import (
	"MyChatApp/monolithic/internal/models"
	"time"

	"github.com/gocql/gocql"
)

type ChannelsRepository struct {
	session *gocql.Session
}

func NewChannelsRepository(session *gocql.Session) *ChannelsRepository {
	return &ChannelsRepository{
		session,
	}
}

func (c *ChannelsRepository) GetChannelByID(channel_id gocql.UUID) (*models.Channel, error) {
	var chn models.Channel

	err := c.session.Query(
		"SELECT channel_id, name, type, created_by, created_at FROM channels WHERE channel_id = ?",
		channel_id,
	).Scan(&chn.ChannelID, &chn.Name, &chn.Type, &chn.CreatedBy, &chn.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &chn, nil
}

func (c *ChannelsRepository) GetChannels() ([]models.Channel, error) {
	iter := c.session.Query("SELECT channel_id, name, type, created_by, created_at FROM channels").Iter()

	var channels []models.Channel
	var channel_id gocql.UUID
	var name string
	var _type string
	var created_by gocql.UUID
	var created_at time.Time

	for iter.Scan(&channel_id, &name, &_type, &created_by, &created_at) {
		channels = append(channels, models.Channel{
			ChannelID: channel_id,
			Name:      name,
			Type:      _type,
			CreatedBy: created_by,
			CreatedAt: created_at,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return channels, nil
}

func (c *ChannelsRepository) GetChannelsByUser(user_id gocql.UUID) ([]models.ChannelMembership, error) {
	iter := c.session.Query(
		"SELECT channel_id, channel_name, channel_type, joined_at FROM channels_by_user WHERE user_id = ?", user_id,
	).Iter()

	var chn_memberships []models.ChannelMembership
	var chn_id gocql.UUID
	var chn_name string
	var chn_type string
	var joined_at time.Time

	for iter.Scan(&chn_id, &chn_name, &chn_type, &joined_at) {
		chn_memberships = append(chn_memberships, models.ChannelMembership{
			ChannelID:   chn_id,
			ChannelName: chn_name,
			ChannelType: chn_type,
			JoinedAt:    joined_at,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return chn_memberships, nil
}

func (c *ChannelsRepository) GetMembersByChannel(chn_id gocql.UUID) ([]gocql.UUID, error) {
	iter := c.session.Query(
		"SELECT user_id from members_by_channel WHERE channel_id = ?",
		chn_id,
	).Iter()

	var users []gocql.UUID
	var id gocql.UUID

	for iter.Scan(&id) {
		users = append(users, id)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *ChannelsRepository) GetMemberRole(channel_id, user_id string) (string, error) {
	iter := c.session.Query(
		"SELECT role FROM members_by_channel WHERE channel_id = ? AND user_id = ?",
		channel_id, user_id,
	).Iter()

	var role string
	iter.Scan(&role)

	if err := iter.Close(); err != nil {
		return "", err
	}

	return role, nil
}

func (c *ChannelsRepository) CreateChannel(chn *models.Channel) error {
	//return c.session.Query(
	//	"INSERT INTO channels (channel_id, name, type, created_by, created_at) VALUES (?, ?, ?, ?, ?)",
	//	chn.ChannelID, chn.Name, chn.Type, chn.CreatedBy, chn.CreatedAt,
	//).Exec()

	batch := c.session.NewBatch(gocql.LoggedBatch)

	batch.Query(
		"INSERT INTO channels (channel_id, name, type, created_by, created_at) VALUES (?, ?, ?, ?, ?)",
		chn.ChannelID, chn.Name, chn.Type, chn.CreatedBy, chn.CreatedAt,
	)

	batch.Query(
		"INSERT into channels_by_user (user_id, channel_id, channel_name, channel_type, joined_at) VALUES (?, ?, ?, ?, ?)",
		chn.CreatedBy, chn.ChannelID, chn.Name, chn.Type, chn.CreatedAt,
	)

	batch.Query(
		"INSERT into members_by_channel (channel_id, user_id, role, joined_at) VALUES (?, ?, ?, ?)",
		chn.ChannelID, chn.CreatedBy, "owner", chn.CreatedAt,
	)

	return c.session.ExecuteBatch(batch)
}

func (c *ChannelsRepository) CreateChannelMembership(user_id gocql.UUID, chn *models.ChannelMembership, role string) error {
	batch := c.session.NewBatch(gocql.LoggedBatch)

	batch.Query(
		"INSERT into channels_by_user (user_id, channel_id, channel_name, channel_type, joined_at) VALUES (?, ?, ?, ?, ?)",
		user_id, chn.ChannelID, chn.ChannelName, chn.ChannelType, chn.JoinedAt,
	)

	batch.Query(
		"INSERT into members_by_channel (channel_id, user_id, role, joined_at) VALUES (?, ?, ?, ?)",
		chn.ChannelID, user_id, role, chn.JoinedAt,
	)

	return c.session.ExecuteBatch(batch)
}

// modify channel membership - only useful for changing owners / admins ?
func (c *ChannelsRepository) ModifyChannelMembership(user_id gocql.UUID, chn *models.ChannelMembership, role string) error {
	return c.session.Query(
		"INSERT into members_by_channel (channel_id, user_id, role, joined_at) VALUES (?, ?, ?, ?)",
		chn.ChannelID, user_id, role, chn.JoinedAt,
	).Exec()
}

// deletions
func (c *ChannelsRepository) DeleteChannelMembership(user_id gocql.UUID, channel_id gocql.UUID) error {
	batch := c.session.NewBatch(gocql.LoggedBatch)

	batch.Query(
		"DELETE from channels_by_user WHERE user_id = ? AND channel_id = ?",
		user_id, channel_id,
	)

	batch.Query(
		"DELETE from members_by_channel WHERE channel_id = ? AND user_id = ?",
		channel_id, user_id,
	)

	return c.session.ExecuteBatch(batch)
}

func (c *ChannelsRepository) DeleteChannel(chn_id gocql.UUID) error {
	return c.session.Query(
		"DELETE from channels where channel_id = ?",
		chn_id,
	).Exec()
}
