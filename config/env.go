package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

var DbHost string
var DbPort string
var DbDatabase string
var DbUsername string
var DbPassword string
var ListeningPort string
var JWTSecretKey string
var LifeTime string
var UrlImage string

func Initialize() {
	rootEnvPath := filepath.Join("..", ".env") // Adjust if needed
	localEnvPath := ".env"                     // Default location

	if err := godotenv.Load(localEnvPath); err != nil {
		log.Println("⚠️  Could not load local .env, trying parent directory...")
		if err := godotenv.Load(rootEnvPath); err != nil {
			log.Fatal("❌ Failed to load .env file from any location:", err)
		}
	}

	DbHost = os.Getenv("DB_HOST")
	DbPort = os.Getenv("DB_PORT")
	DbDatabase = os.Getenv("DB_DATABASE")
	DbUsername = os.Getenv("DB_USERNAME")
	DbPassword = os.Getenv("DB_PASSWORD")
	ListeningPort = os.Getenv("PORT")
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")
	LifeTime = os.Getenv("LIFETIME")
	UrlImage = os.Getenv("URL_IMAGE")
}
