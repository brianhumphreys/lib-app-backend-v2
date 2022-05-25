package modeltests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/brianhumphreys/library_app/api/controllers"
	"github.com/brianhumphreys/library_app/api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var server = controllers.Server{}
var userInstance = models.User{}
var bookInstance = models.Book{}
var checkoutInstance = models.Checkout{}

// Test Setup
func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	DBURL := fmt.Sprintf("%s://%s:%s@%s:%s", os.Getenv("DB_NAME"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("Additional information:", err)
	} else {
		fmt.Printf("Connection to %s database was successful", TestDbDriver)
	}
}

func refreshUserTable() error {
	server.DB.Migrator().DropTable(&models.User{})
	server.DB.AutoMigrate(&models.User{})

	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserAndBookAndCheckoutTable()

	user := models.User{
		Email:    "test@gmail.com",
		Password: "test",
		Role:     "user",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Fatalf("User table could not be seeded: %v", err)
	}
	return user, nil
}

func seedUsers() error {

	users := []models.User{
		models.User{
			Email:    "b@a.com",
			Password: "bumq123",
			Role:     "user",
		},
		models.User{
			Email:    "j@a.com",
			Password: "jwoma123",
			Role:     "admin",
		},
	}

	for i, _ := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func refreshUserAndBookAndCheckoutTable() error {

	server.DB.Migrator().DropTable(&models.Checkout{}, &models.Book{}, &models.User{})
	server.DB.AutoMigrate(&models.User{}, &models.Book{}, &models.Checkout{})

	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndTwoBookAndOneCheckout() (models.User, []models.Book, models.Checkout, error) {

	err := refreshUserAndBookAndCheckoutTable()

	// seed user
	if err != nil {
		return models.User{}, []models.Book{}, models.Checkout{}, err
	}
	user := models.User{
		Email:    "n@a.com",
		Password: "ngon9123",
		Role:     "admin",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, []models.Book{}, models.Checkout{}, err
	}

	// seed 2 books
	book1 := models.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Isbn:        "Test Isbn",
		Description: "Test Description",
	}
	err = server.DB.Model(&models.Book{}).Create(&book1).Error
	if err != nil {
		return models.User{}, []models.Book{}, models.Checkout{}, err
	}

	book2 := models.Book{
		Title:       "Test Book 2",
		Author:      "Test Author 2",
		Isbn:        "Test Isbn 2",
		Description: "Test Description 2",
	}
	err = server.DB.Model(&models.Book{}).Create(&book2).Error
	if err != nil {
		return models.User{}, []models.Book{}, models.Checkout{}, err
	}

	// seed checkout
	checkout := models.Checkout{
		UserId: 1,
		BookId: 1,
	}
	err = server.DB.Model(&models.Checkout{}).Create(&checkout).Error
	if err != nil {
		return models.User{}, []models.Book{}, models.Checkout{}, err
	}
	return user, []models.Book{book1, book2}, checkout, nil
}

func seedUsersAndBook() ([]models.User, []models.Book, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Book{}, err
	}
	var users = []models.User{
		models.User{
			Email:    "test1@gmail.com",
			Password: "test1",
			Role:     "user",
		},
		models.User{
			Email:    "test2@gmail.com",
			Password: "test2",
			Role:     "admin",
		},
	}
	var posts = []models.Book{
		models.Book{
			Title:       "Test Title 1",
			Author:      "Test Author 1",
			Isbn:        "Test Isbn 1",
			Description: "Test Description 1",
		},
		models.Book{
			Title:       "Test Title 1",
			Author:      "Test Author 1",
			Isbn:        "Test Isbn 1",
			Description: "Test Description 1",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		// posts[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Book{}).Create(&posts[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
	return users, posts, nil
}
