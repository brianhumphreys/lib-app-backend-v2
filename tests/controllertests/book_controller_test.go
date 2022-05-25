package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/brianhumphreys/library_app/api/models"
	"github.com/gorilla/mux"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateBook(t *testing.T) {

	err := refreshUserAndBookAndCheckoutTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, err := seedUsersAndBook()
	if err != nil {
		log.Fatalf("User table could not be seeded %v\n", err)
	}
	userEmail := users[1].Email
	userPassword := "test2"

	adminEmail := users[0].Email
	adminPassword := "test1"

	_, userToken, err := server.SignIn(userEmail, userPassword)
	if err != nil {
		log.Fatalf("Could not Login: %v\n", err)
	}
	userTokenString := fmt.Sprintf("Bearer %v", userToken)

	_, adminToken, err := server.SignIn(adminEmail, adminPassword)
	if err != nil {
		log.Fatalf("Could not Login: %v\n", err)
	}
	adminTokenString := fmt.Sprintf("Bearer %v", adminToken)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		author       string
		isbn         string
		description  string
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"Memoirs of a Geisha", "author": "Arthur Golden", "isbn": "isbn", "description": "description"}`,
			statusCode:   201,
			tokenGiven:   adminTokenString,
			title:        "Memoirs of a Geisha",
			author:       "Arthur Golden",
			isbn:         "isbn",
			description:  "description",
			errorMessage: "",
		},
		{
			// duplicate creations are allowed
			inputJSON:    `{"title":"Memoirs of a Geisha", "author": "Arthur Golden", "isbn": "isbn", "description": "description"}`,
			statusCode:   201,
			tokenGiven:   adminTokenString,
			title:        "Memoirs of a Geisha",
			author:       "Arthur Golden",
			isbn:         "isbn",
			description:  "description",
			errorMessage: "",
		},
		{
			// non admin users cannot create books
			inputJSON:    `{"title":"Memoirs of a Geisha", "author": "Arthur Golden", "isbn": "isbn", "description": "description"}`,
			tokenGiven:   userTokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content"}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "author": "author", "isbn": "isbn", "description": "description"}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"title": "Another Title", "author": "", "isbn": "isbn", "description": "description"}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Author",
		},
		{
			inputJSON:    `{"title": "Another Title", "author": "author", "isbn": "", "description": "description"}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Isbn",
		},
		{
			inputJSON:    `{"title": "Another Title", "author": "author", "isbn": "isbn", "description": ""}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Description",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/books", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateBook)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["author"], v.author)
			assert.Equal(t, responseMap["isbn"], v.isbn)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetBooks(t *testing.T) {

	refreshUserAndBookAndCheckoutTable()

	_, _, err := seedUsersAndBook()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Errorf("Additional Information: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetBooks)
	handler.ServeHTTP(rr, req)

	var books []models.Book
	err = json.Unmarshal([]byte(rr.Body.String()), &books)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(books), 2)
}

func TestGetBookByID(t *testing.T) {

	err := refreshUserAndBookAndCheckoutTable()
	if err != nil {
		log.Fatal(err)
	}
	_, book, err := seedUsersAndBook()
	if err != nil {
		log.Fatal(err)
	}
	bookTestPayloads := []struct {
		id           string
		statusCode   int
		title        string
		author       string
		isbn         string
		description  string
		author_id    uint32
		errorMessage string
	}{
		{
			id:          strconv.Itoa(int(book[0].ID)),
			statusCode:  200,
			title:       book[0].Title,
			author:      book[0].Author,
			isbn:        book[0].Isbn,
			description: book[0].Description,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range bookTestPayloads {

		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetBook)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, book[0].Title, responseMap["title"])
			assert.Equal(t, book[0].Author, responseMap["author"])
			assert.Equal(t, book[0].Isbn, responseMap["isbn"])
			assert.Equal(t, book[0].Description, responseMap["description"])
		}
	}
}

func TestUpdateBook(t *testing.T) {

	refreshUserAndBookAndCheckoutTable()

	users, books, err := seedUsersAndBook()
	if err != nil {
		log.Fatal(err)
	}

	userEmail := users[1].Email
	userPassword := "test2"

	adminEmail := users[0].Email
	adminPassword := "test1"

	_, userToken, err := server.SignIn(userEmail, userPassword)
	if err != nil {
		log.Fatalf("Could not Login: %v\n", err)
	}
	userTokenString := fmt.Sprintf("Bearer %v", userToken)

	_, adminToken, err := server.SignIn(adminEmail, adminPassword)
	if err != nil {
		log.Fatalf("Could not Login: %v\n", err)
	}
	adminTokenString := fmt.Sprintf("Bearer %v", adminToken)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		title        string
		author       string
		isbn         string
		description  string
		tokenGiven   string
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "New Isbn", "description": "New Description"}`,
			statusCode:   200,
			title:        "New Title",
			author:       "New Author",
			isbn:         "New Isbn",
			description:  "New Description",
			tokenGiven:   adminTokenString,
			errorMessage: "",
		},
		{
			// Duplicate Titles are allowed
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "New Isbn", "description": "New Description"}`,
			statusCode:   200,
			title:        "New Title",
			author:       "New Author",
			isbn:         "New Isbn",
			description:  "New Description",
			tokenGiven:   adminTokenString,
			errorMessage: "",
		},
		{
			// non admin users cannot create a book
			id:           strconv.Itoa(int(books[1].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "New Isbn", "description": "New Description"}`,
			tokenGiven:   userTokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When no token is provided
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "New Isbn", "description": "New Description"}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "New Isbn", "description": "New Description"}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"", "author": "New Author", "isbn": "New Isbn", "description": "New Description"}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Title",
		},
		{
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "", "isbn": "New Isbn", "description": "New Description"}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Author",
		},
		{
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "", "description": "New Description"}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Isbn",
		},
		{
			id:           strconv.Itoa(int(books[0].ID)),
			updateJSON:   `{"title":"New Title", "author": "New Author", "isbn": "New Isbn", "description": ""}`,
			statusCode:   422,
			tokenGiven:   adminTokenString,
			errorMessage: "Required Description",
		},
		{
			id:         "bad request",
			statusCode: 400,
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/books", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateBook)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["author"], v.author)
			assert.Equal(t, responseMap["isbn"], v.isbn)
			assert.Equal(t, responseMap["description"], v.description)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteBook(t *testing.T) {

	err := refreshUserAndBookAndCheckoutTable()
	if err != nil {
		log.Fatal(err)
	}
	users, _, err := seedUsersAndBook()
	if err != nil {
		log.Fatal(err)
	}

	userEmail := users[1].Email
	userPassword := "test2"

	adminEmail := users[0].Email
	adminPassword := "test1"

	_, userToken, err := server.SignIn(userEmail, userPassword)
	if err != nil {
		log.Fatalf("Could not Login: %v\n", err)
	}
	userTokenString := fmt.Sprintf("Bearer %v", userToken)

	_, adminToken, err := server.SignIn(adminEmail, adminPassword)
	if err != nil {
		log.Fatalf("Could not Login: %v\n", err)
	}
	adminTokenString := fmt.Sprintf("Bearer %v", adminToken)

	// Get only the second post
	postSample := []struct {
		id           string
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           strconv.Itoa(1),
			tokenGiven:   adminTokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// Book does not exist
			id:           strconv.Itoa(1),
			tokenGiven:   adminTokenString,
			statusCode:   404,
			errorMessage: "",
		},
		{
			// Book does not exist
			id:           strconv.Itoa(7),
			tokenGiven:   adminTokenString,
			statusCode:   404,
			errorMessage: "",
		},
		{
			// non admin user cannot delete a book
			id:           strconv.Itoa(2),
			tokenGiven:   userTokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When empty token is passed
			id:           strconv.Itoa(2),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           strconv.Itoa(2),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "bad request",
			tokenGiven: adminTokenString,
			statusCode: 400,
		},
	}
	for _, v := range postSample {
		req, _ := http.NewRequest("GET", "/books", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteBook)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
