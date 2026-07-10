package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

// our authentication middleware, will set "user_id" that can be
// used by downstream handlers and routes
func JWTAuthMiddleware() gin.HandlerFunc {
	// JWTAuthMiddleware itself is only run once when we wire up the routes
	// so env read will only happen once during startup here
	// inner function captures the variable
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-in-production"
	}

	return func(c *gin.Context) {
		var tokenStr string

		if cookieToken, err := c.Cookie("auth_token"); err == nil && cookieToken != "" {
			tokenStr = cookieToken
		}

		// fall back to authorization header with bearer token if cookie is not found
		if tokenStr == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no auth token"})
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
				return
			}
			tokenStr = parts[1]
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {

			// verify that the signing method is what we expect
			// jwt.SigningMethodHS256 is not a type - the underlying type is jwt.SigningMethodHMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, exists := claims["user_id"]
			if !exists || userID == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
				return
			}

			c.Set("user_id", userID)
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			return
		}

		c.Next()
	}
}
