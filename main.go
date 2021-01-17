package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sjsu-ebuddy/identity-service/env"
	"github.com/sjsu-ebuddy/identity-service/pkg/app"
	"github.com/sjsu-ebuddy/identity-service/pkg/db"
	"github.com/sjsu-ebuddy/identity-service/pkg/services"
	"github.com/sjsu-ebuddy/identity-service/pkg/util"
)

func main() {

	var wait time.Duration
	var serverEnv string

	flag.DurationVar(&wait, "timeout", time.Second*15,
		"The duration for which the server will wait for the connections to close")

	flag.StringVar(&serverEnv, "env", "", "-env=dev")

	flag.Parse()

	err := godotenv.Load(fmt.Sprintf("env/%s.env", serverEnv))

	if err != nil {
		log.Fatalln(err)
	}

	port := os.Getenv(env.Port)
	host := os.Getenv(env.Host)

	dbHost := os.Getenv(env.DbHost)
	dbPort := os.Getenv(env.DbPort)
	dbName := os.Getenv(env.DbName)
	dbUser := os.Getenv(env.DbUser)
	dbPass := os.Getenv(env.DbPass)

	var sslMode string

	if serverEnv != "dev" {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}

	cfg := &db.Config{
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
		User:     dbUser,
		Password: dbPass,
		SSLMode:  sslMode,
	}

	appDB := db.GetConnection(cfg)
	appValidator := util.GetValidator()

	var namespace string

	if serverEnv != "dev" {
		namespace = os.Getenv(env.Namespace)
	} else {
		namespace = "http://localhost"
	}

	srv := &app.Server{
		Addr: fmt.Sprintf("%s:%s", host, port),
		Handler: app.GetRouter(&services.Service{
			Namespace: namespace,
			V:         appValidator,
			DB:        appDB,
		}),
		Wait: wait,
	}

	srv.Start()

}
