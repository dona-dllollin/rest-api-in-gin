package helper

import (
	"rest-api-in-gin/internal/env"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/golang-jwt/jwt"
)

func GenerateAccessToken(userID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"expr":   time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenString, err := token.SignedString([]byte(env.GetEnvString("JWT_SECRET", "some-secret-123456")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"expr":   time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(env.GetEnvString("REFRESH_SECRET", "some-iam-fresh-2718271")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
