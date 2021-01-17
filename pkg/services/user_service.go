package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sjsu-ebuddy/identity-service/pkg/db"
	"gorm.io/gorm"
)

// ValidateUserHandler method
func (s *Service) ValidateUserHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user db.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			log.Println(err.Error())
			s.handleResponse(&Response{
				StatusCode: http.StatusBadRequest,
				Error:      BadRequest,
			}, w)

			return
		}

		err = user.ValidateUser(s.V)

		if err != nil {

			var validationErrors []string

			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, fmt.Sprintf("%s:%s", err.Field(), err.Tag()))
			}

			s.handleResponse(&Response{
				StatusCode: http.StatusBadRequest,
				Error:      ValidatationError,
				Messages:   validationErrors,
			}, w)

			return
		}

		ctxWithUser := context.WithValue(r.Context(), userContextKey, &user)
		rWithUser := r.WithContext(ctxWithUser)
		next(w, rWithUser)
	}
}

// CreateUserHandler method
func (s *Service) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(*db.User)

	var userExists db.User

	err := s.DB.First(&userExists, db.User{
		Email: user.Email,
	}).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		s.DB.Create(&user)

		s.handleResponse(&Response{
			StatusCode: http.StatusCreated,
			Data: &UserData{
				User: user.FormatResponse(),
			},
			Message: Success,
		}, w)

		return
	}

	s.handleResponse(&Response{
		StatusCode: http.StatusConflict,
		Error:      Conflict,
		Message:    UserAlreadyExists,
	}, w)

}
