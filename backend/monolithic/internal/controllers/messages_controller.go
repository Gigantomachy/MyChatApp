package controllers

import (
	"MyChatApp/monolithic/internal/services"
	"MyChatApp/monolithic/internal/ws"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessagesController struct {
	messageService *services.MessageService
	hub            *ws.Hub
}

func NewMessagesController(ms *services.MessageService, hub *ws.Hub) *MessagesController {
	return &MessagesController{messageService: ms, hub: hub}
}

// GET /api/channels/:id/messages
func (mc *MessagesController) GetMessages(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	channelID := c.Param("id")

	// by default get 50 messages at a time
	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	// timeuuid is used to compute the bucket - default is today
	before := c.Query("before")

	messages, err := mc.messageService.GetMessages(channelID, uid.(string), limit, before)
	if err != nil {
		var reqErr *services.RequestError
		if errors.As(err, &reqErr) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// POST /api/channels/:id/messages
func (mc *MessagesController) SendMessage(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists || uid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	channelID := c.Param("id")

	var payload struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := mc.messageService.SendMessage(channelID, uid.(string), payload.Content)
	if err != nil {
		var reqErr *services.RequestError
		if errors.As(err, &reqErr) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memberIDs, err := mc.messageService.GetChannelMemberIDs(channelID)
	if err == nil {
		event, _ := json.Marshal(gin.H{
			"type": "message.new",
			"data": gin.H{
				"channel_id":        channelID,
				"message_id":        msg.MessageID,
				"author_id":         msg.AuthorID,
				"author_first_name": msg.AuthorFirstName,
				"author_last_name":  msg.AuthorLastName,
				"author_username":   msg.AuthorUsername,
				"content":           msg.Content,
				"created_at":        msg.CreatedAt,
			},
		})
		mc.hub.PublishToUsers(memberIDs, event)
	}

	c.JSON(http.StatusCreated, msg)
}
