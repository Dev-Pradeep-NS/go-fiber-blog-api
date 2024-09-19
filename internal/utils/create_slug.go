package utils

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func CreateSlug(input string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	processedString := reg.ReplaceAllString(input, " ")

	processedString = strings.TrimSpace(processedString)

	slug := strings.ReplaceAll(processedString, " ", "-")

	slug = strings.ToLower(slug)

	return slug
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	if token.Method.Alg() != jwt.SigningMethodHS256.Name {
		return nil, errors.New("invalid token signing method")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
