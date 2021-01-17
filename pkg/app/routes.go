package app

import (
	"sync"

	"github.com/gorilla/mux"
	"github.com/sjsu-ebuddy/identity-service/pkg/middleware"
	"github.com/sjsu-ebuddy/identity-service/pkg/services"
)

var (
	router *mux.Router
	once   sync.Once
)

// GetRouter returns mux.Router with all the handlers preloaded
func GetRouter(s *services.Service) *mux.Router {

	once.Do(func() {
		router = mux.NewRouter()
	})

	router.Use(middleware.RequestLogger)
	router.Use(middleware.ContentType)
	router.HandleFunc("/", s.HealthCheckHandler)
	router.HandleFunc("/auth/register",
		s.ValidateUserHandler(s.CreateUserHandler)).Methods("POST")
	router.HandleFunc("/auth/login",
		s.ValidateLoginData(s.UserLogin)).Methods("POST")
	router.HandleFunc("/auth/verify", s.AuthHandler(s.VerifyAuthHandler)).Methods("GET")

	return router
}
