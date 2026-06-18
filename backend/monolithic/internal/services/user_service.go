package services

import (
	"errors"
	"os"
	"strings"
	"time"

	"MyChatApp/monolithic/internal/models"
	"MyChatApp/monolithic/internal/repositories/db"

	"github.com/gocql/gocql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *db.UserRepository
}

func NewUserService(repo *db.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password, email, firstName, lastName string) (*models.User, string, error) {
	existing, _ := s.repo.FindByUsername(strings.ToLower(username))
	if existing != nil {
		return nil, "", errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	userID, err := gocql.RandomUUID()
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		UserID:       userID,
		Username:     username,
		PasswordHash: string(hash),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		CreatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, "", err
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(username, password string) (*models.User, string, error) {
	user, err := s.repo.FindByUsername(strings.ToLower(username))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) generateJWT(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	claims := jwt.MapClaims{
		"user_id":  user.UserID.String(),
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *UserService) SearchUsers(query string) ([]models.User, error) {
	all, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	if query == "" {
		return all, nil
	}

	lower := strings.ToLower(query)
	var filtered []models.User
	for _, u := range all {
		if strings.Contains(strings.ToLower(u.Username), lower) ||
			strings.Contains(strings.ToLower(u.FirstName), lower) ||
			strings.Contains(strings.ToLower(u.LastName), lower) {
			filtered = append(filtered, u)
		}
	}
	return filtered, nil
}