package main

import (
	"log"
	"os"

	"ecommerce/database"
	"ecommerce/router"

	"github.com/joho/godotenv"
)

func init() {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Loading the config from .env file")
		err = godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env config file")
		}
		log.Println("Successfully loaded the config file")
	}
	database.ConnectDb()
}

func main() {
	router.ClientRoutes()
}
