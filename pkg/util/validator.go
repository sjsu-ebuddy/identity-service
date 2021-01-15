package util

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/sjsu-ebuddy/identity-service/pkg/db"
)

var (
	v    *validator.Validate
	once sync.Once
)

// GetValidator returns validator
func GetValidator() *validator.Validate {
	once.Do(func() {
		v = validator.New()
	})
	_ = v.RegisterValidation("password", db.ValidatePassword)
	return v
}
