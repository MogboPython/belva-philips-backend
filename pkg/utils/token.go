package utils

import (
	"errors"
	"log"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"sessionId": id,
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token, sign and generate encoded token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config("JWT_SECRET")))
	if err != nil {
		log.Println("Error signing token:", err)
		return "", errors.New("error generating token")
	}

	return t, nil
}

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

// func validToken(t *jwt.Token, id string) bool {
// 	n, err := strconv.Atoi(id)
// 	if err != nil {
// 		return false
// 	}

// 	claims := t.Claims.(jwt.MapClaims)
// 	uid := int(claims["user_id"].(float64))

// 	return uid == n
// }

// CheckPasswordHash compare password with hash
// func CheckPasswordHash(password, hash string) bool {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	log.Println(hash, "haaaash")
// 	return err == nil
// }
