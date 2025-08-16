package util

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//var jwtSecretKey = []byte("your-very-secret-key")

// Sinh token
func GenerateJWT(userID, jwtSecret string) (string, error) {
	clamis := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, clamis)
	return token.SignedString([]byte(jwtSecret))
}

// Parse và xác thực token
func ParseJWT(tokenStr, jwtSecret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Sử dụng đúng method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// Lấy userid từ JWT
func GetUserIDFromRequest(r *http.Request, jwtSecret string) (string, bool) {
	// Lấy token từ header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return "", false
	}

	// Lấy claim và user_id
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", false
	}
	return userID, true
}
