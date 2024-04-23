package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	dbConfig "github.com/rms_api/internal/config/database"
	"github.com/rms_api/internal/config/secrets"
	L "github.com/rms_api/internal/middlewares/logger"
	mdl "github.com/rms_api/internal/models"
	"net/http"
	"time"
)

func RefreshToken(c *gin.Context) {
	db, err := dbConfig.OpenDbConnection()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer dbConfig.CloseDbConnection(db)
	// Get refresh token from request
	L.RMSLog("D", "In Refresh Token", nil)
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token is missing"})
		c.Abort()
		return
	}
	// Validate and parse refresh token
	token, err := jwt.Parse(refreshToken[len("Bearer "):], func(token *jwt.Token) (interface{}, error) {
		return []byte(secrets.RMS_DEV_API_SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		c.Abort()
		return
	}
	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return
	}
	// Check if the token is expired
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	if exp.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token has expired"})
		c.Abort()
		return
	}
	// Get user ID from claims
	userID := claims["user_id"].(string)
	// Retrieve user from database based on user ID
	var user mdl.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		c.Abort()
		return
	}
	// Generate new access token
	accessToken, err := GenerateToken(user)
	if err != nil {
		L.RMSLog("D", "Error Generating Access Token:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		c.Abort()
		return
	}
	// Update new access token in the database
	var existingToken mdl.Token
	if err := db.Where("user_id = ?", user.ID).First(&existingToken).Error; err != nil {
		// Handle error (e.g., token not found)
		L.RMSLog("D", "Error finding token in DB:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find token in DB"})
		c.Abort()
		return
	}
	// Update token value
	existingToken.Token = accessToken

	// Update token expiry
	existingToken.ExpiredAt = time.Now().Add(time.Hour * 1)

	if err := db.Save(&existingToken).Error; err != nil {
		// Handle error (e.g., update failed)
		L.RMSLog("D", "Error updating token in DB:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token in DB"})
		c.Abort()
		return
	}
	// Return new access token
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}
