package controllers

import (
	"MyChatApp/monolithic/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FriendController struct {
	friendService *services.FriendService
}

func NewFriendController(fs *services.FriendService) *FriendController {
	return &FriendController{
		friendService: fs,
	}
}

// GET /api/friends
func (f *FriendController) GetFriends(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	userIDStr := user_id.(string)
	friends, err := f.friendService.GetFriends(userIDStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friends)
}

// DELETE /api/friends/:id
func (f *FriendController) RemoveFriend(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	friend_id := c.Param("id")
	if friend_id == "" {
		// this part isn't super necessary since the router should fail to match altogether
		// if :id isn't provided
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing friend_id"})
		return
	}

	userIDStr := user_id.(string)
	err := f.friendService.RemoveFriend(userIDStr, friend_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "friend successfully removed",
		"id":      friend_id,
	})
}

// GET /friend-requests
func (f *FriendController) GetFriendRequests(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	freqs, err := f.friendService.GetFriendRequests(user_id.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, freqs)
}

// POST /friend-requests/:id
func (f *FriendController) SendFriendRequest(c *gin.Context) {
	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	recipient_id := c.Param("id")
	if recipient_id == "" {
		// this part isn't super necessary since the router should fail to match altogether
		// if :id isn't provided
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing recipient_id"})
		return
	}

	freq, err := f.friendService.SendFriendRequest(user_id.(string), recipient_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"recipient_id": freq.RecipientID,
		"status":       freq.Status,
		"created_at":   freq.CreatedAt,
	})
}

// PUT /friend-requests/:id
func (f *FriendController) AcceptFriendRequest(c *gin.Context) {

	// AKA recipient_id
	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	sender_id := c.Param("id")
	if sender_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing sender_id"})
		return
	}

	err := f.friendService.AcceptFriendRequest(user_id.(string), sender_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DELETE /api/friend-requests/outgoing/:recipient_id
func (f *FriendController) CancelFriendRequest(c *gin.Context) {
	recipient_id := c.Param("id")
	if recipient_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing recipient_id"})
		return
	}

	user_id, exists := c.Get("user_id")
	if !exists || user_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return
	}

	err := f.friendService.DeleteFriendRequest(user_id.(string), recipient_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DELETE /api/friend-requests/incoming/:sender_id
func (f *FriendController) RejectFriendRequest(c *gin.Context) {
	recipient_id, exists := c.Get("user_id")
	if !exists || recipient_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing recipient_id"})
		return
	}

	sender_id := c.Param("id")
	if sender_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing sender_id"})
		return
	}

	err := f.friendService.DeleteFriendRequest(sender_id, recipient_id.(string))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
