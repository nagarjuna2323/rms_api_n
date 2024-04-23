package authorization

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rms_api/internal/config/database"
	"github.com/rms_api/internal/config/secrets"
	L "github.com/rms_api/internal/middlewares/logger"
	mdl "github.com/rms_api/internal/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}
		// Extract JWT token from Authorization header
		tokenString = tokenString[len("Bearer "):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secrets.RMS_DEV_API_SECRET_KEY), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Check if the token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Check token expiry
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := time.Unix(int64(claims["exp"].(float64)), 0)
			if exp.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}
		}
		// Check if the token is blacklisted
		if isRevoked := IsTokenRevoked(tokenString); isRevoked {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}
		// Set user information in the context for further processing
		c.Set("user", token.Claims)
		c.Next()
	}
}

// IsTokenRevoked checks if the given token is blacklisted in the database
func IsTokenRevoked(token string) bool {
	// Get database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		// If unable to connect to the database, assume token is not revoked
		return false
	}
	defer func(db *gorm.DB) {
		err := database.CloseDbConnection(db)
		if err != nil {
			L.RMSLog("E", "Unable to Close database Connection", err)
			return
		}
	}(db)
	// Check if token exists in blacklist
	var count int64
	if err := db.Model(&mdl.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error; err != nil {
		// If an error occurred during the database query, assume token is not revoked
		return false
	}
	// If the count is greater than 0, token is revoked
	return count > 0
}
