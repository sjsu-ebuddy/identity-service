package util

import (
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/sjsu-ebuddy/identity-service/env"
)

// ClaimsValidator type
type ClaimsValidator func(claims jwt.MapClaims) bool

// GetJWT returns a jwt string for login
func GetJWT(claims jwt.MapClaims) string {
	hmacSecret := os.Getenv(env.HmacSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(hmacSecret)

	if err != nil {
		log.Println(err.Error())
		return ""
	}

	return tokenString
}

// VerifyJWT method
func VerifyJWT(jwtToken string, validator ClaimsValidator) bool {
	return false
}
