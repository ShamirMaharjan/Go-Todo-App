package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DATABASE_URL string
	PORT         string
	JWT_SECRET   string
}

func Load() (*Config, error) {
	var err error = godotenv.Load()

	if err != nil {
		log.Println("Waring: .env file not found!")
		return nil, err
	}

	var config *Config = &Config{
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		PORT:         os.Getenv("PORT"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
	}

	return config, err
}
