package utils

import (
	"errors"
	"log"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"sessionId": id,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token, sign and generate encoded token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config("JWTSecretKey")))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", errors.New("error generating token")
	}

	return t, nil
}
