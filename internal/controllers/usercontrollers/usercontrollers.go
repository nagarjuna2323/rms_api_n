package usercontrollers

import (
	"github.com/gin-gonic/gin"
	database "github.com/rms_api/internal/config/database"
	"github.com/rms_api/internal/middlewares/authentication"
	hash "github.com/rms_api/internal/middlewares/hashpassword"
	L "github.com/rms_api/internal/middlewares/logger"
	vld "github.com/rms_api/internal/middlewares/validation"
	mdl "github.com/rms_api/internal/models"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var db *gorm.DB

func SignUpService(c *gin.Context) {
	var newUser mdl.User
	L.RMSLog("D", "Opening database Connection", nil)
	db, err := database.OpenDbConnection()
	defer database.CloseDbConnection(db)
	//validate request
	if err := c.ShouldBindJSON(&newUser); err != nil {
		L.RMSLog("E", "Error Processing the Request", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate email format
	if !vld.IsValidEmail(newUser.Email) {
		L.RMSLog("E", "Here User Sent Invalid Email format :"+L.PrintStruct(newUser.Email), nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}
	// Check if email already exists
	var existingUser mdl.User
	result := db.Where("email = ?", newUser.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		L.RMSLog("E", "Email is already registered"+L.PrintStruct(newUser.Email), nil)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}
	// Hash the password before saving it
	hashedPassword, err := hash.HashPassword(newUser.Password)
	if err != nil {
		L.RMSLog("E", "Error generating hash password:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	// Create a new UUID for the user
	newUser.ID = uint(rand.Uint64()) //random number generator function required
	newUser.Password = hashedPassword
	// Create the user in the database
	newUser.UserType = strings.ToLower(newUser.UserType) // Convert user type to lowercase
	if newUser.UserType != "applicant" && newUser.UserType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user type"})
		return
	}
	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	// Return success response
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func LogInService(c *gin.Context) {
	var signInReq mdl.SignInRequest
	// Bind JSON request body into signInReq struct
	if err := c.ShouldBindJSON(&signInReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Opening database connection
	db, err := database.OpenDbConnection()
	if err != nil {
		L.RMSLog("E", "Error opening database connection", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer database.CloseDbConnection(db)
	//L.RMSLog("D", "Opening database Connection", nil)
	//db, err := database.OpenDbConnection()

	// Retrieve user from the database
	var user mdl.User
	if err := db.Where("email = ?", signInReq.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare password
	err = hash.ComparePasswords(user.Password, signInReq.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Generate JWT token
	token, err := authentication.GenerateToken(user)
	if err != nil {
		L.RMSLog("D", "Error Generating Token:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Store the token with issued time and expiry time
	newToken := mdl.Token{
		UserID:    user.ID,
		Email:     user.Email,
		Token:     token,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 1),
	}

	if err := db.Create(&newToken).Error; err != nil {
		L.RMSLog("D", "Error Storing Token:\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}
	// Return JWT token as response
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Controller function for handling the upload of resume files

func UploadResume(c *gin.Context) {
	// Check file format
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error retrieving file"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".docx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file format. Only PDF or DOCX files are allowed."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Resume file uploaded successfully"})
}
