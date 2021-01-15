package app

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sjsu-ebuddy/identity-service/pkg/middleware"
	"github.com/sjsu-ebuddy/identity-service/pkg/services"
	"gorm.io/gorm"
)

var (
	router *mux.Router
	once   sync.Once
)

// GetRouter returns mux.Router with all the handlers preloaded
func GetRouter(db *gorm.DB, v *validator.Validate) *mux.Router {

	once.Do(func() {
		router = mux.NewRouter()
	})

	s := &services.Service{
		DB: db,
		V:  v,
	}

	router.Use(middleware.RequestLogger)
	router.Use(middleware.ContentType)
	router.HandleFunc("/", s.HealthCheckHandler)
	router.HandleFunc("/users/register", s.ValidateUserHandler(s.CreateUserHandler)).Methods("POST")
	// router.HandleFunc("/users/login", s.ValidateLoginData(s.UserLogin)).Methods("POST")

	return router
}
