package authentication

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/rms_api/internal/config/secrets"
	L "github.com/rms_api/internal/middlewares/logger"
	mdl "github.com/rms_api/internal/models"
	"strconv"
	"time"
)

func GenerateToken(user mdl.User) (string, error) {
	// Create the claims
	claims := mdl.Claims{
		UserID: strconv.Itoa(int(user.ID)),
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	jwtSecret := []byte(secrets.RMS_DEV_API_SECRET_KEY)
	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	L.RMSLog("D", "Token:"+L.PrintStruct(tokenString), err)
	return tokenString, nil
}
