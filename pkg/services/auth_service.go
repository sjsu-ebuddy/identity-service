package services

import (
	"encoding/json"
	"log"
	"net/http"
)

type loginFormData struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

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

		}

		next(w, r)
	}
}

// UserLogin method for authentication
func (s *Service) UserLogin(w http.ResponseWriter, r *http.Request) {

}
