package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ElvinEga/gofiber_starter/blacklist"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(userId string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}
func GenerateJWTRole(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func VerifyJWT(c *fiber.Ctx) (userID string, role string, err error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", "", errors.New("missing token")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == authHeader {
		return "", "", errors.New("invalid token format")
	}

	// Check blacklist first
	if blacklist.IsBlacklisted(tokenStr) {
		return "", "", errors.New("token revoked")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return "", "", err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, claims.Role, nil
	}

	return "", "", errors.New("invalid token claims")
}
func VerifyJWTRole(c *fiber.Ctx) (userID string, role string, err error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", "", errors.New("Missing token")
	}

	tokenStr := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("Invalid token claims")
	}

	id, okID := claims["user_id"].(string)
	roleStr, okRole := claims["role"].(string)
	if !okID || !okRole {
		return "", "", errors.New("Invalid token data")
	}

	return id, roleStr, nil
}
