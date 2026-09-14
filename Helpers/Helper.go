package Helpers

import (
	"os"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var SecretKey = os.Getenv("SECRET_KEY")

func GenerateJWT(issuer string) (string, error) {

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	return claims.SignedString([]byte(SecretKey))
}

// ParseJWT validates the JWT and returns the issuer.
func ParseJWT(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {

			// Make sure the token uses HMAC/HS256.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return claims.Issuer, nil
}

// ValidateToken checks whether a JWT is valid and has not expired.
func ValidateToken(tokenString string) error {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {

			// Verify the signing algorithm.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// Explicit expiration check.
	if claims.ExpiresAt < time.Now().Unix() {
		return errors.New("token has expired")
	}

	return nil
}

// UpdateToken creates a new JWT using the issuer
// from the existing token.
func UpdateToken(tokenString string) (string, error) {

	// First validate the existing token.
	err := ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Extract the issuer from the existing token.
	issuer, err := ParseJWT(tokenString)
	if err != nil {
		return "", err
	}

	// Generate a new token with a fresh expiration time.
	return GenerateJWT(issuer)
}