package services

import (
	"MyChatApp/monolithic/internal/models"
	"MyChatApp/monolithic/internal/repositories/db"
	"errors"
	"time"

	"github.com/gocql/gocql"
)

type ChannelService struct {
	channelsRepo *db.ChannelsRepository
	usersRepo    *db.UserRepository
}

func NewChannelService(r *db.ChannelsRepository, u *db.UserRepository) *ChannelService {
	return &ChannelService{
		channelsRepo: r,
		usersRepo:    u,
	}
}

func (c *ChannelService) GetChannelByID(chn_id string) (*models.Channel, error) {
	id, err := gocql.ParseUUID(chn_id)
	if err != nil {
		return nil, err
	}

	return c.channelsRepo.GetChannelByID(id)
}

func (c *ChannelService) GetChannels() ([]models.Channel, error) {
	return c.channelsRepo.GetChannels()
}

func (c *ChannelService) GetChannelMembershipsByUser(user_id string) ([]models.ChannelMembership, error) {
	uid, err := gocql.ParseUUID(user_id)
	if err != nil {
		return nil, err
	}

	return c.channelsRepo.GetChannelsByUser(uid)
}

func (c *ChannelService) GetUsersByChannel(chn_id string) ([]models.User, error) {
	cid, err := gocql.ParseUUID(chn_id)
	if err != nil {
		return nil, err
	}

	userIDs, err := c.channelsRepo.GetMembersByChannel(cid)
	if err != nil {
		return nil, err
	}

	return c.usersRepo.FindByIDs(userIDs)
}

func (c *ChannelService) CreateChannel(creator_id, name, _type string) (*models.Channel, error) {
	cid, err := gocql.RandomUUID()
	if err != nil {
		return nil, err
	}

	creatorID, err := gocql.ParseUUID(creator_id)
	if err != nil {
		return nil, err
	}

	var chn = models.Channel{
		ChannelID: cid,
		Name:      name,
		Type:      _type,
		CreatedBy: creatorID,
		CreatedAt: time.Now(),
	}

	err = c.channelsRepo.CreateChannel(&chn)
	if err != nil {
		return nil, err
	}

	return &chn, nil
}

// TODO: maybe make role a parameter?
func (c *ChannelService) CreateChannelMembership(user_id, channel_id string) (*models.ChannelMembership, error) {
	uid, err := gocql.ParseUUID(user_id)
	if err != nil {
		return nil, err
	}

	cid, err := gocql.ParseUUID(channel_id)
	if err != nil {
		return nil, err
	}

	chn, err := c.channelsRepo.GetChannelByID(cid)
	if err != nil {
		return nil, err
	}

	chnMem := models.ChannelMembership{
		ChannelID:   cid,
		ChannelName: chn.Name,
		ChannelType: chn.Type,
		JoinedAt:    time.Now(),
	}

	err = c.channelsRepo.CreateChannelMembership(uid, &chnMem, "member")
	if err != nil {
		return nil, err
	}

	return &chnMem, nil
}

func (c *ChannelService) ModifyChannelMembership(user_id, channel_id, role string) error {
	uid, err := gocql.ParseUUID(user_id)
	if err != nil {
		return err
	}

	cid, err := gocql.ParseUUID(channel_id)
	if err != nil {
		return err
	}

	chn, err := c.channelsRepo.GetChannelByID(cid)
	if err != nil {
		return err
	}

	chnMem := models.ChannelMembership{
		ChannelID:   cid,
		ChannelName: chn.Name,
		ChannelType: chn.Type,
		JoinedAt:    time.Now(),
	}

	return c.channelsRepo.ModifyChannelMembership(uid, &chnMem, role)
}

func (c *ChannelService) DeleteChannel(user_id, channel_id string) error {
	cid, err := gocql.ParseUUID(channel_id)
	if err != nil {
		return err
	}

	role, err := c.channelsRepo.GetMemberRole(channel_id, user_id)
	if err != nil {
		return err
	}

	if role != "owner" {
		return errors.New("Only owners can delete channels")
	}

	return c.channelsRepo.DeleteChannel(cid)
}

// TODO: If 1 person left in DM, delete DM channel too? if 0 people left in channel, delete channel?
func (c *ChannelService) DeleteChannelMembership(user_id, channel_id string) error {
	uid, err := gocql.ParseUUID(user_id)
	if err != nil {
		return err
	}

	cid, err := gocql.ParseUUID(channel_id)
	if err != nil {
		return err
	}

	return c.channelsRepo.DeleteChannelMembership(uid, cid)
}
