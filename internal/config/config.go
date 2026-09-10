package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)


type Config struct {
	Port 	string
	Env 	string
}

func MustLoad() Config {
	 godotenv.Load()

	// if err != nil {
	// 	log.Fatal("Error loading .env file.")
	// }
	
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		log.Fatal("ENV is required")
	}

	return Config{
		Port: port,
		Env: env,
	}	
}