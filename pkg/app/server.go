package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

// Server struct for our identity server
type Server struct {
	Addr       string
	Handler    *mux.Router
	Wait       time.Duration
	HTTPServer *http.Server
}

// Start method for the server
func (srv *Server) Start() {

	srv.HTTPServer = &http.Server{
		Addr:         srv.Addr,
		Handler:      srv.Handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		log.Printf("Starting Server on http://%s \n", srv.Addr)
		if err := srv.HTTPServer.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), srv.Wait)

	defer cancel()

	srv.HTTPServer.Shutdown(ctx)

	log.Println("shutting down")

	os.Exit(0)
}
