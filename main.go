package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sjsu-ebuddy/identity-service/env"
	"github.com/sjsu-ebuddy/identity-service/pkg/app"
	"github.com/sjsu-ebuddy/identity-service/pkg/db"
	"github.com/sjsu-ebuddy/identity-service/pkg/util"
)

func main() {

	serverEnv := os.Getenv(env.Env)

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

	srv := &app.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: app.GetRouter(appDB, appValidator),
	}

	srv.Start()

}
