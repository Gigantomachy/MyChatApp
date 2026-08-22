package controllers

import (
	"MyChatApp/monolithic/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{userService: userService}
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.userService.Register(
		req.Username, req.Password, req.Email, req.FirstName, req.LastName,
	)
	if err != nil {
		if errors.Is(err, services.ErrUsernameTaken) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("auth_token", token, 86400, "/", "", false, true) // secure flag false for now

	ctx.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.userService.Login(req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("auth_token", token, 86400, "/", "", false, true) // secure flag false for now

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("auth_token", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (c *AuthController) Me(ctx *gin.Context) {
	uid, exists := ctx.Get("user_id")
	if !exists || uid == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing user id",
		})
		return
	}

	usr, err := c.userService.FindByID(uid.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: replace with actual user DTO later - with no password field to be extra safe
	// password should still be filtered out as of right now due to json annotations though
	ctx.JSON(http.StatusOK, gin.H{"user": usr})
}
