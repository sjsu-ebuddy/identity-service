package app

import (
	"context"
	"flag"
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
	httpServer *http.Server
	wait       time.Duration
}

// Start method for the server
func (srv *Server) Start() {

	srv.httpServer = &http.Server{
		Addr:         srv.Addr,
		Handler:      srv.Handler,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	flag.DurationVar(&srv.wait, "graceful-timeout", time.Second*15,
		"The duration for which the server will wait for the connections to close")
	flag.Parse()

	go func() {
		log.Printf("Starting Server on http://%s \n", srv.Addr)
		if err := srv.httpServer.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), srv.wait)

	defer cancel()

	srv.httpServer.Shutdown(ctx)

	log.Println("shutting down")

	os.Exit(0)
}
