package controllers

import (
	"MyChatApp/monolithic/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(us *services.UserService) *UserController {
	return &UserController{userService: us}
}

func (uc *UserController) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	users, err := uc.userService.SearchUsers(query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
