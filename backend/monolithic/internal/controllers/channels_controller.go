package controllers

import (
	"MyChatApp/monolithic/internal/services"
	"MyChatApp/monolithic/internal/ws"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

// TODO: replace with DTO later
type ChannelPayload struct {
	ChannelName string   `json:"channel_name"`
	ChannelType string   `json:"channel_type"`
	Members     []string `json:"members"`
}

type ChannelsController struct {
	channelService *services.ChannelService
	hub            *ws.Hub
}

func NewChannelsController(cs *services.ChannelService, hub *ws.Hub) *ChannelsController {
	return &ChannelsController{
		channelService: cs,
		hub:            hub,
	}
}

// GET /api/channels/discover - for searching - replace later
func (cc *ChannelsController) GetAllChannels(c *gin.Context) {
	channels, err := cc.channelService.GetChannels()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, channels)
}

// GET /api/channels
func (cc *ChannelsController) GetChannelsByUser(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	memberships, err := cc.channelService.GetChannelMembershipsByUser(uid.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, memberships)
}

// /api/channels/:id
func (cc *ChannelsController) GetChannelUsersAndInformation(c *gin.Context) {
	cid := c.Param("id")

	chn_info, err := cc.channelService.GetChannelByID(cid)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: replace with user DTO later ?
	users, err := cc.channelService.GetUsersByChannel(cid)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"channel_id":   chn_info.ChannelID,
		"channel_name": chn_info.Name,
		"channel_type": chn_info.Type,
		"created_by":   chn_info.CreatedBy,
		"created_at":   chn_info.CreatedAt,
		"users":        users,
	})
}

// POST /api/channels
func (cc *ChannelsController) CreateChannel(c *gin.Context) {
	var payload ChannelPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	chn, err := cc.channelService.CreateChannel(payload.Members, uid.(string), payload.ChannelName, payload.ChannelType)
	if err != nil {
		var reqErr *services.RequestError
		if errors.As(err, &reqErr) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if payload.ChannelType == "dm" && len(payload.Members) > 0 {
		event, _ := json.Marshal(gin.H{"type": "channel.new"})
		cc.hub.PublishToUsers([]string{payload.Members[0]}, event)
	}

	c.JSON(http.StatusCreated, chn)
}

// POST /api/channels/:id
func (cc *ChannelsController) CreateChannelMembership(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	cid := c.Param("id")

	chnMem, err := cc.channelService.CreateChannelMembership(uid.(string), cid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, chnMem)
}

// PUT /api/channels/:id
func (cc *ChannelsController) ModifyChannelMembership(c *gin.Context) {
	var payload struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	cid := c.Param("id")

	err := cc.channelService.ModifyChannelMembership(uid.(string), cid, payload.Role)
	if err != nil {
		var reqErr *services.RequestError
		if errors.As(err, &reqErr) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// delete channel
func (cc *ChannelsController) DeleteChannel(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	cid := c.Param("id")

	err := cc.channelService.DeleteChannel(uid.(string), cid)
	if err != nil {
		var reqErr *services.RequestError
		if errors.As(err, &reqErr) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// delete channel membership
func (cc *ChannelsController) DeleteChannelMembership(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	cid := c.Param("id")

	err := cc.channelService.DeleteChannelMembership(uid.(string), cid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
