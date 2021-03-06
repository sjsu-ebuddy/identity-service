package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sjsu-ebuddy/identity-service/pkg/db"
	"gorm.io/gorm"
)

// Response object for identity service
type Response struct {
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Message    string      `json:"message,omitempty"`
	Messages   []string    `json:"messages,omitempty"`
}

// Service Object
type Service struct {
	Namespace string
	DB        *gorm.DB
	V         *validator.Validate
}

// UserData for response object
type UserData struct {
	User  *db.User `json:"user,omitempty"`
	Token string   `json:"accessToken,omitempty"`
}

type loginFormData struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

// HealthCheckHandler returns OK if service is running
func (s *Service) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received call")

	response := map[string]string{
		"message": "OK",
	}

	j, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		log.Println(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// HandleResponse returns json error based on the error passed
func (s *Service) handleResponse(res *Response, w http.ResponseWriter) {

	j, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
	}

	w.WriteHeader(res.StatusCode)
	w.Write(j)
}
