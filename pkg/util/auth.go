package util

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sjsu-ebuddy/identity-service/env"
)

// ClaimsValidator type
type ClaimsValidator func(claims jwt.MapClaims) bool

// JWT Constants
const (
	JWTAccessToken = "ACCESS_TOKEN"
	EbuddyIssuer   = "EBUDDY_IDENTITY_SERVICE"
)

// UserAuthClaims struct
type UserAuthClaims struct {
	jwt.StandardClaims
	Role string `json:"role,omitempty"`
}

// VerifyRole takes a comparison role and compares it with the role
// in the token
func (userClaims *UserAuthClaims) VerifyRole(role string, required bool) bool {
	if userClaims.Role == "" {
		return !required
	}
	if subtle.ConstantTimeCompare([]byte(userClaims.Role), []byte(role)) != 0 {
		return true
	}
	return false
}

// GetJWT returns a jwt string for login
func GetJWT(claims *UserAuthClaims, keyID string) string {
	hmacSecret := os.Getenv(env.HmacSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = keyID

	tokenString, err := token.SignedString([]byte(hmacSecret))

	if err != nil {
		log.Println(err.Error())
		return ""
	}

	return tokenString
}

// VerifyJWT method
func VerifyJWT(jwtToken string) (uuid.UUID, error) {
	hmacSecret := os.Getenv(env.HmacSecret)
	token, err := jwt.ParseWithClaims(jwtToken,
		&UserAuthClaims{}, func(tokenString *jwt.Token) (interface{}, error) {
			if _, ok := tokenString.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v",
					tokenString.Header["alg"])
			}
			return []byte(hmacSecret), nil
		})

	if err != nil {
		log.Println(err.Error())
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(*UserAuthClaims); ok &&
		token.Valid && claims.Valid() == nil {
		err = claims.Valid()

		userID, err := uuid.Parse(claims.Subject)

		if err != nil {
			log.Println(err.Error())
			return uuid.Nil, err
		}

		return userID, nil
	}
	return uuid.Nil, errors.New("Invalid Token")
}
