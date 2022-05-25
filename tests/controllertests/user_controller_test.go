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

func TestCreateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		role         string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "brianhumphreys@gmail.com", "password": "password", "role": "user"}`,
			statusCode:   201,
			email:        "brianhumphreys@gmail.com",
			role:         "user",
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "brianhumphreys2@gmail.com", "password": "password2", "role": "admin"}`,
			statusCode:   201,
			email:        "brianhumphreys2@gmail.com",
			role:         "admin",
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "newemail@gmail.com", "password": "password2", "role": "not valid"}`,
			statusCode:   422,
			errorMessage: "Role must be 'user' or 'admin'",
		},
		{
			inputJSON:    `{"email": "brianhumphreys@gmail.com", "password": "password", "role": "admin"}`,
			statusCode:   500,
			errorMessage: "Email Already Taken",
		},
		{
			inputJSON:    `{"email": "brianhumphreys.com", "password": "password", "role": "admin"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
		{
			inputJSON:    `{"email": "", "password": "password", "role": "admin"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"email": "brye@gmail.com", "password": "", "role": "admin"}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
		{
			inputJSON:    `{"email": "brye@gmail.com", "password": "password", "role": ""}`,
			statusCode:   422,
			errorMessage: "Required Role",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["email"], v.email)
		}
		if v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = SeedUsers()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Errorf("Additional Information: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetUsers)
	handler.ServeHTTP(rr, req)

	var users []models.User
	err = json.Unmarshal([]byte(rr.Body.String()), &users)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(users), 2)
}

func TestGetUserByID(t *testing.T) {

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	userSample := []struct {
		id           string
		statusCode   int
		email        string
		errorMessage string
	}{
		{
			id:         strconv.Itoa(int(user.ID)),
			statusCode: 200,
			email:      user.Email,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("Additional Information: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, user.Email, responseMap["email"])
		}
	}
}

func TestUpdateUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := SeedUsers()
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Test update of first user
	currentID := users[0].ID
	currentEmail := users[0].Email
	currentPassword := "bumq123"

	// First login
	_, token, err := server.SignIn(currentEmail, currentPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		id           string
		updateJSON   string
		statusCode   int
		updateEmail  string
		role         string
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int32 to int first before converting to string
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "newbhumq@gmail.com", "password": "newpassword", "role": "admin"}`,
			statusCode:   200,
			updateEmail:  "newbhumq@gmail.com",
			role:         "admin",
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When password field is empty
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "bhumq@gmail.com", "password": "", "role": "admin"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Password",
		},
		{
			// When no token was passed
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "test3@gmail.com", "password": "test3", "role": "admin"}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "b@g.com", "password": "password", "role": "admin"}`,
			statusCode:   401,
			tokenGiven:   "This is incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			// Remember "kenny@gmail.com" belongs to user 2
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "j@a.com", "password": "jwoma123", "role": "admin"}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Email Already Taken",
		},
		{
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "brianhumphreys.com", "password": "password", "role": "admin"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Invalid Email",
		},
		{
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "", "password": "password", "role": "admin"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Email",
		},
		{
			id:           strconv.Itoa(int(currentID)),
			updateJSON:   `{"email": "j@a.com", "password": "password", "role": ""}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Role",
		},
		{
			id:         "bad request",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// When user 2 is using user 1 token
			id:           strconv.Itoa(int(2)),
			updateJSON:   `{"email": "j@a.com", "password": "jwoma123"}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {
		req, err := http.NewRequest("POST", "/users", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["email"], v.updateEmail)
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	users, err := SeedUsers()
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}
	// Test deletion of first user
	currentID := users[0].ID
	currentEmail := users[0].Email
	currentPassword := "bumq123"

	// First login
	_, token, err := server.SignIn(currentEmail, currentPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	userSample := []struct {
		id           string
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			id:           strconv.Itoa(int(currentID)),
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// check 401 when no token given
			id:           strconv.Itoa(int(currentID)),
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// check 401 when bad token given
			id:           strconv.Itoa(int(currentID)),
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         "badrequest",
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			// Check 401 when user uses someone elses token
			id:           strconv.Itoa(int(2)),
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range userSample {

		req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteUser)

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
