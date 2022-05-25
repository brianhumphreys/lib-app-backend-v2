package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/brianhumphreys/library_app/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	DBURL := fmt.Sprintf("%s://%s:%s@%s:%s", DbName, DbUser, DbPassword, DbHost, DbPort)

	val, present := os.LookupEnv("DATABASE_URL")
	if present {
		server.DB, err = gorm.Open(postgres.Open(val), &gorm.Config{})
	} else {
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	}

	if err != nil {
		fmt.Printf("Connection with %s could not be obtained: ", Dbdriver)
		log.Fatal("Additional output:", err)
	} else {
		fmt.Printf("Connection with %s was successful", Dbdriver)
	}

	server.DB.AutoMigrate(&models.User{}, &models.Book{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) InitializeDev() {

	var err error

	val, present := os.LookupEnv("DATABASE_URL")
	if present {
		server.DB, err = gorm.Open(postgres.Open(val), &gorm.Config{})
	} else {
		log.Fatalf("Could not load you .env file: %v", err)
	}

	if err != nil {
		fmt.Printf("Connection with database could not be obtained: ")
		log.Fatal("Additional output:", err)
	} else {
		fmt.Printf("Connection with database was successful")
	}

	server.DB.AutoMigrate(&models.User{}, &models.Book{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	// fmt.Println("Server is running on port 8080")
	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"http://localhost:3000"},
	// 	AllowedMethods: []string{"GET", "POST", "PATCH"},
	// 	AllowedHeaders: []string{"Bearer", "Content_Type", "Authorization"}})
	// handler := c.Handler(server.Router)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
