package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sjsu-ebuddy/identity-service/pkg/db"
	"github.com/sjsu-ebuddy/identity-service/pkg/util"
)

// ValidateLoginData method for validation
func (s *Service) ValidateLoginData(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var formData loginFormData

		err := json.NewDecoder(r.Body).Decode(&formData)

		if err != nil {
			log.Println(err.Error())
			s.handleResponse(&Response{
				StatusCode: http.StatusBadRequest,
				Error:      BadRequest,
			}, w)

			return
		}

		err = s.V.Struct(formData)

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
		}

		ctxWithData := context.WithValue(r.Context(), loginContextKey, formData)
		rWithData := r.WithContext(ctxWithData)

		next(w, rWithData)
	}
}

// UserLogin method for authentication
func (s *Service) UserLogin(w http.ResponseWriter, r *http.Request) {
	loginDetails := r.Context().Value(loginContextKey).(loginFormData)

	var user db.User

	err := s.DB.First(&user, db.User{
		Email: loginDetails.Email,
	}).Error

	if err != nil {
		log.Println(err.Error())
		s.handleResponse(&Response{
			StatusCode: http.StatusNotFound,
			Error:      NotFound,
			Message:    NoSuchUser,
		}, w)
		return
	}

	if !user.VerifyPassword(loginDetails.Password) {
		s.handleResponse(&Response{
			StatusCode: http.StatusUnauthorized,
			Error:      UnauthorizedAccess,
			Message:    InvalidPassword,
		}, w)
		return
	}

	now := time.Now()

	claims := &util.UserAuthClaims{}
	claims.Audience = s.Namespace
	claims.Issuer = util.EbuddyIssuer
	claims.IssuedAt = now.Unix()
	// 5 days
	claims.ExpiresAt = now.Add(time.Second * 60 * 60 * 24 * 5).Unix()
	claims.Subject = user.ID.String()

	token := util.GetJWT(claims, util.JWTAccessToken)

	s.handleResponse(&Response{
		StatusCode: http.StatusOK,
		Data: &UserData{
			User:  user.FormatResponse(),
			Token: token,
		},
		Message: Success,
	}, w)

}

// AuthHandler for verifying if token is valid
func (s *Service) AuthHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get(AuthorizationHeader)
		reqToken = strings.Split(reqToken, "Bearer ")[1]

		userID, err := util.VerifyJWT(reqToken)

		if err != nil {
			s.handleResponse(&Response{
				StatusCode: http.StatusUnauthorized,
				Error:      UnauthorizedAccess,
				Message:    InvalidToken,
			}, w)
			return
		}

		var user db.User

		err = s.DB.First(&user, db.User{
			ID: userID,
		}).Error
		// No such user
		if err != nil {
			s.handleResponse(&Response{
				StatusCode: http.StatusUnauthorized,
				Error:      UnauthorizedAccess,
				Message:    InvalidToken,
			}, w)
			return
		}

		ctxWithData := context.WithValue(r.Context(), userContextKey, &user)
		rWithUser := r.WithContext(ctxWithData)

		next(w, rWithUser)
	}
}

// VerifyAuthHandler deals with the /auth/verify endpoint.
func (s *Service) VerifyAuthHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContextKey).(*db.User)

	s.handleResponse(&Response{
		StatusCode: http.StatusOK,
		Data: &UserData{
			User: user.FormatResponse(),
		},
		Message: Success,
	}, w)
}
