package api

import (
	"log"
	"os"

	"github.com/brianhumphreys/library_app/api/controllers"
	"github.com/brianhumphreys/library_app/seed"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println("Could not load you .env file: %v", err)
		server.InitializeDev()
		// seed.Load(server.DB)

		server.Run(":" + os.Getenv("PORT"))
	} else {
		log.Println("Loaded values from .env file")
		server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
		seed.Load(server.DB)

		server.Run(":8080")
	}

}
