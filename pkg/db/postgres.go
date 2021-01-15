package db

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config takes db params
type Config struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
	SSLMode  string
}

var (
	once sync.Once
	db   *gorm.DB
)

func (cfg *Config) getDsn() string {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Host,
		cfg.Port,
		cfg.SSLMode,
	)

	return dsn
}

// GetConnection returns gorm database object
func GetConnection(cfg *Config) *gorm.DB {
	var err error
	dsn := cfg.getDsn()
	once.Do(func() {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&User{})

	return db
}
