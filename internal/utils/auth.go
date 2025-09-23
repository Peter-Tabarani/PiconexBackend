package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte("your-secret-key")

// Claims used for JWT tokens
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// CreateJWT generates a new JWT for a user
func CreateJWT(userID int, role string) (string, error) {
	expiration := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtSecret)
}
