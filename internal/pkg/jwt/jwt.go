package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var secretKey = []byte(os.Getenv("SECRET_KEY_JWT_ACCESS_TOKEN"))

func CreateToken(username string, isRefreshToken bool) (string, error) {
	if isRefreshToken {
		secretKey = []byte(os.Getenv("SECRET_KEY_JWT_REFRESH_TOKEN"))
	}

	now := time.Now().Local()
	var exp int64
	if isRefreshToken {
		exp = now.Add(7 * 24 * time.Hour).Unix()
	} else {
		exp = now.Add(1 * time.Hour).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":        "agria-backend",
		"sub":        username,
		"exp":        exp,
		"iat":        now,
		"role":       "user",
		"is_refresh": isRefreshToken,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, isRefreshToken bool) error {
	if isRefreshToken {
		secretKey = []byte(os.Getenv("SECRET_KEY_JWT_REFRESH_TOKEN"))
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
