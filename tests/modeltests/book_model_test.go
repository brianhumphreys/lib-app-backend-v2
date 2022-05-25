package modeltests

import (
	"log"
	"testing"

	"github.com/brianhumphreys/library_app/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllBooks(t *testing.T) {
	_, books, _, err := seedOneUserAndTwoBookAndOneCheckout()
	if err != nil {
		log.Fatalf("Tables could not be seeded %v\n", err)
	}
	foundBooks, err := books[0].FindAllBooks(server.DB)
	if err != nil {
		t.Errorf("There was an error getting the books: %v\n", err)
		return
	}
	assert.Equal(t, len(*foundBooks), 2)
}

func TestSaveBook(t *testing.T) {

	refreshUserAndBookAndCheckoutTable()

	newBook := models.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Isbn:        "Test Isbn",
		Description: "Test Description",
	}
	savedBook, err := newBook.SaveBook(server.DB)
	if err != nil {
		t.Errorf("There was an error while saving the book: %v\n", err)
		return
	}
	assert.Equal(t, savedBook.ID, newBook.ID)
	assert.Equal(t, savedBook.Title, newBook.Title)
	assert.Equal(t, savedBook.Author, newBook.Author)
	assert.Equal(t, savedBook.Isbn, newBook.Isbn)
	assert.Equal(t, savedBook.Description, newBook.Description)

}

func TestGetPostByID(t *testing.T) {

	_, books, _, err := seedOneUserAndTwoBookAndOneCheckout()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	foundBook, err := models.FindBookByID(server.DB, uint64(books[0].ID))
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, len(*foundBook), 1)
	assert.Equal(t, (*foundBook)[0].ID, books[0].ID)
	assert.Equal(t, (*foundBook)[0].Title, books[0].Title)
	assert.Equal(t, (*foundBook)[0].Author, books[0].Author)
	assert.Equal(t, (*foundBook)[0].Isbn, books[0].Isbn)
	assert.Equal(t, (*foundBook)[0].Description, books[0].Description)
}

func TestUpdateAPost(t *testing.T) {

	_, _, _, err := seedOneUserAndTwoBookAndOneCheckout()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}
	book1 := models.Book{
		Title:       "Test Book",
		Author:      "CORRECTED Test Author",
		Isbn:        "NEW Test Isbn",
		Description: "NEW Test Description",
	}
	book1.ID = 1

	updatedBook, err := book1.UpdateABook(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedBook.ID, book1.ID)
	assert.Equal(t, updatedBook.Title, book1.Title)
	assert.Equal(t, updatedBook.Author, book1.Author)
	assert.Equal(t, updatedBook.Isbn, book1.Isbn)
	assert.Equal(t, updatedBook.Description, book1.Description)
}

func TestDeleteAPost(t *testing.T) {

	_, books, _, err := seedOneUserAndTwoBookAndOneCheckout()
	if err != nil {
		log.Fatalf("Error Seeding tables")
	}
	_, err = books[0].DeleteABook(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}

	foundUser2, err := models.FindBookByID(server.DB, uint64(books[0].ID))

	assert.Equal(t, len(*foundUser2), 0)
}
