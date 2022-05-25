package seed

import (
	"fmt"
	"log"

	"github.com/brianhumphreys/library_app/api/models"
	"gorm.io/gorm"
)

var users = []models.User{
	models.User{
		Email:    "batya@pt.com",
		Password: "batya123",
		Role:     "admin",
	},
	models.User{
		Email:    "rob@pt.com",
		Password: "rob123",
		Role:     "admin",
	},
	models.User{
		Email:    "brian@pt.com",
		Password: "brian123",
		Role:     "user",
	},
}

var books = []models.Book{
	models.Book{
		Title:       "A Little Life",
		Author:      "Hanya Yanagihara",
		Isbn:        "isbn",
		Description: "description",
		Available:   false,
	},
	models.Book{
		Title:       "The Anthropocene Reviewed",
		Author:      "John Green",
		Isbn:        "isbn",
		Description: "description",
		Available:   false,
	},
	models.Book{
		Title:       "The Handmaid's Tale",
		Author:      "Margaret Atwood",
		Isbn:        "isbn",
		Description: "description",
		Available:   false,
	},
	models.Book{
		Title:       "The Perks of Being a Wallflower",
		Author:      "Stephen Chbosky",
		Isbn:        "isbn",
		Description: "description",
		Available:   false,
	},
	models.Book{
		Title:       "Memoirs of a Geisha",
		Author:      "Arthur Golden",
		Isbn:        "isbn",
		Description: "description",
		Available:   true,
	},
	models.Book{
		Title:       "The Souls of Black Folk",
		Author:      "W. E. B. Du Bois",
		Isbn:        "isbn",
		Description: "description",
		Available:   true,
	},
}

var checkouts = []models.Checkout{
	models.Checkout{
		BookId: 1,
		UserId: 2,
	},
	models.Checkout{
		BookId: 2,
		UserId: 2,
	},
	models.Checkout{
		BookId: 3,
		UserId: 2,
	},
	models.Checkout{
		BookId: 4,
		UserId: 2,
	},
}

func Load(db *gorm.DB) {
	db.Migrator().DropTable(&models.Checkout{}, &models.Book{}, &models.User{})
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Checkout{})

	var createErr error
	for i, _ := range users {
		createErr = db.Create(&users[i]).Error
		if createErr != nil {
			log.Fatalf("User table could not be seeded: %v", createErr)
		}
	}
	fmt.Println("User table seeded.")

	for i, _ := range books {
		createErr = db.Create(&books[i]).Error
		if createErr != nil {
			log.Fatalf("Book table could not be seeded: %v", createErr)
		}
	}
	fmt.Println("Book table seeded.")

	for i, _ := range checkouts {
		createErr = db.Create(&checkouts[i]).Error
		if createErr != nil {
			log.Fatalf("Checkout table could not be seeded: %v", createErr)
		}
	}

	booksResult := []models.Book{}
	err := db.Table("users").Select("books.title, books.author, books.isbn, books.description").Joins("JOIN checkouts on checkouts.user_id = users.id").Joins("JOIN books on books.id = checkouts.book_id").Where("checkouts.user_id = ?", 5).Limit(100).Find(&booksResult).Error
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Checkout table seeded.")
}
