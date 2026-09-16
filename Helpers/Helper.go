package Helpers

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SecretKey = os.Getenv("SECRET_KEY")

type CustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(userID, email, role string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "restaurant-api",
			Subject:   userID,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return claims.SignedString([]byte(SecretKey))
}

// ParseJWT validates the JWT and returns the issuer.
func ParseJWT(tokenString string) (string, error) {

	claims, err := parseClaims(tokenString)

	if err != nil {
		return "", err
	}

	return claims.Issuer, nil
}

// ValidateToken checks whether a JWT is valid and has not expired.
func ValidateToken(tokenString string) error {

	claims, err := parseClaims(tokenString)

	if err != nil {
		return err
	}

	// Explicit expiration check.
	if claims.ExpiresAt < time.Now().Unix() {
		return errors.New("token has expired")
	}

	return nil
}

func parseClaims(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// UpdateToken creates a new JWT using the user details
// from the existing token.
func UpdateToken(tokenString string) (string, error) {

	// First validate the existing token.
	err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Extract the user details from the existing token.
	claims, err := parseClaims(tokenString)
	if err != nil {
		return "", err
	}

	// Generate a new token with a fresh expiration time.
	return GenerateJWT(claims.UserID, claims.Email, claims.Role)
}
