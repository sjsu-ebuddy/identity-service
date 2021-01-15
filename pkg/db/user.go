package db

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

// User object
type User struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:uuid;primary_key;" json:"ID"`
	FirstName  string    `json:"firstName" validate:"required"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName" validate:"required"`
	Email      string    `json:"email" gorm:"uniqueIndex:idx_email" validate:"email"`
	Password   string    `json:"password,omitempty" validate:"password"`
}

// BeforeCreate method
func (user *User) BeforeCreate(tx *gorm.DB) error {
	salt := randomHex(16)
	derivedKey, err := scrypt.Key([]byte(user.Password), []byte(salt), 1<<15, 8, 1, 32)

	if err != nil {
		return err
	}

	hash := hex.EncodeToString(derivedKey)
	user.Password = fmt.Sprintf("%s:%s", salt, hash)
	user.ID = genUUID()
	return nil
}

// VerifyPassword takes plain string password and compares it with hash
func (user *User) VerifyPassword(password string) bool {
	hashedPwd := strings.Split(user.Password, ":")
	derivedKey, err := scrypt.Key([]byte(user.Password), []byte(hashedPwd[0]), 1<<15, 8, 1, 32)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return hex.EncodeToString(derivedKey) == hashedPwd[1]
}

// ValidateUser validates user details.
func (user *User) ValidateUser(usrValidator *validator.Validate) error {
	err := usrValidator.Struct(user)

	log.Println(err)

	return err
}

// FormatResponse returns user as a json value
func (user *User) FormatResponse() *User {
	user.Password = ""
	return user
}

func randomHex(n int) string {
	bytes := make([]byte, n)

	if _, err := rand.Read(bytes); err != nil {
		log.Println(err.Error())
		return ""
	}

	return hex.EncodeToString(bytes)
}

func genUUID() uuid.UUID {
	uuid, err := uuid.NewRandom()

	if err != nil {
		log.Fatalln(err.Error())
	}

	return uuid
}

// ValidatePassword is a custom validator
func ValidatePassword(fl validator.FieldLevel) bool {

	pwd := fl.Field().String()

	if len(pwd) < 8 || len(pwd) > 32 {
		return false
	}

	var upper bool
	var lower bool
	var num bool
	var splChar bool
	var unExpected bool

	for _, ch := range pwd {
		switch {
		case unicode.IsNumber(ch):
			num = true
		case unicode.IsUpper(ch):
			upper = true
		case unicode.IsLower(ch):
			lower = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			splChar = true
		default:
			unExpected = true
		}
	}

	return !unExpected && upper && lower && num && splChar
}
