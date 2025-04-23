package utils

import (
	"errors"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const noOfHours = 72

func GenerateToken(id, role string) (string, error) {
	claims := jwt.MapClaims{
		"sessionId": id,
		"role":      role,
		"exp":       time.Now().Add(time.Hour * noOfHours).Unix(),
	}

	// Create token, sign and generate encoded token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(config.Config("JWT_SECRET")))
	if err != nil {
		log.Error("Error signing token:", err)
		return "", errors.New("error generating token")
	}

	return t, nil
}

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
