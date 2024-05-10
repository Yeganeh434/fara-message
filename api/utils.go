package api

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
)

func hash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hashedBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString
}

func generateID() string {
	const charset = "0123456789"
	rand.NewSource(10)
	id := make([]byte, 5)
	for idx := range id {
		id[idx] = charset[rand.Intn(len(charset))]
	}

	return string(id)
}

func GetUserIDFromToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims format")
	}

	userID := claims[TokenUserID].(string)
	expirationTimeUnix := claims[TokenExpireTime].(float64)
	expirationTime := time.Unix(int64(expirationTimeUnix), 0)
	if expirationTime.Before(time.Now()) {
		return "", errors.New("token has expired")
	}
	return userID, nil
}

func GetUserID(authorizationHeader string) (string, error) {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errors.New("Unauthorized")
	}
	accessToken := parts[1]
	userID, err := GetUserIDFromToken(accessToken)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func IsStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasSpecialChar := false
	specialChars := "!@#$%^&*()_-|{}[]?<>.,"
	for _, char := range password {
		for _, specialChar := range specialChars {
			if char == specialChar {
				hasSpecialChar = true
				break
			}
		}
	}
	if !hasSpecialChar {
		return hasSpecialChar
	}

	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return hasDigit
	}

	hasUppercase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}
	if !hasUppercase {
		return hasUppercase
	}

	hasLowercase := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowercase = true
			break
		}
	}
	if !hasLowercase {
		return hasLowercase
	}

	return true
}
