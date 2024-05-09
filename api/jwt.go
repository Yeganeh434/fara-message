// package api

// import (
// 	"errors"
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// )

// var secretKey = []byte("farawin")

// const (
// 	TokenExpireTime = "expiration_time"
// 	TokenUserID     = "user_id"
// )

// func CreateJWTToken(userID string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		TokenUserID:     userID,
// 		TokenExpireTime: time.Now().Add(time.Hour * 24).Unix(),
// 	})
// 	tokenString, err := token.SignedString(secretKey)
// 	if err != nil {
// 		return "", fmt.Errorf("error signing token: %w", err)
// 	}
// 	return tokenString, nil

// }

// func ValidateToken(tokenString string) (string, error) {
// 	parts := strings.Split(tokenString, " ")
// 	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 		return "", errors.New("Unauthorized")
// 	}
// 	accessToken := parts[1]

// 	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("failed to parse token: %w", err)
// 	}
// 	if !token.Valid {
// 		return "", errors.New("invalid token")
// 	}
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return "", errors.New("invalid claims format")
// 	}
// 	// userID := claims[TokenUserID].(string)
// 	// expirationTime := claims[TokenExpireTime].Unix()
// 	// if expirationTime.Before(time.Now()) {
// 	// 	return "", errors.New("token has expired")
// 	// }

// 	// return userID, nil




// 	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 	// 	return secretKey, nil
// 	// })
// 	// if err != nil {
// 	// 	return "", fmt.Errorf("failed to parse token: %w", err)
// 	// }
// 	// if !token.Valid {
// 	// 	return "", errors.New("invalid token")
// 	// }
// 	// claims, ok := token.Claims.(jwt.MapClaims)
// 	// if !ok {
// 	// 	return "", errors.New("invalid claims format")
// 	// }
// 	userID, ok := claims[TokenUserID].(string)
// 	if !ok {
// 		return "", errors.New("invalid user ID format")
// 	}
// 	expireTimeClaim, ok := claims[TokenExpireTime].(float64)
// 	if !ok {
// 		return "", errors.New("invalid expiration time format")
// 	}
// 	expirationTime := time.Unix(int64(expireTimeClaim), 0)
// 	if expirationTime.Before(time.Now()) {
// 		return "", errors.New("token has expired")
// 	}

// 	return userID, nil
// }

package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey = []byte("farawin")

const (
	TokenExpireTime = "expiration_time"
	TokenUserID     = "user_id"
)

func CreateJWTToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		TokenUserID:     userID,
		TokenExpireTime: time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil

}

func ValidateToken(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.Trim(tokenString, `"`)

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
