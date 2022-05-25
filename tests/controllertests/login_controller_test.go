package controllertests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSignIn(t *testing.T) {

	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("There was an error seeding the user table %v\n", err)
	}

	samples := []struct {
		email        string
		password     string
		errorMessage string
	}{
		{
			email:        user.Email,
			password:     "test", //Note the password has to be this, not the hashed one from the database
			errorMessage: "",
		},
		{
			email:        user.Email,
			password:     "Wrong password",
			errorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			email:        "Wrong email",
			password:     "password",
			errorMessage: "record not found",
		},
	}

	for _, v := range samples {
		signedUser, token, err := server.SignIn(v.email, v.password)
		if err != nil {
			assert.Equal(t, err, errors.New(v.errorMessage))
		} else {
			assert.NotEqual(t, token, "")
			assert.Equal(t, signedUser.Email, "test@gmail.com")
			assert.Equal(t, signedUser.Role, "admin")
		}
	}
}

func TestLogin(t *testing.T) {

	_, err := seedOneUser()
	if err != nil {
		fmt.Printf("There was an error seeding the user table %v\n", err)
	}
	samples := []struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "test@gmail.com", "password": "test"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "test@gmail.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Password",
		},
		{
			inputJSON:    `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Required Email",
		},
		{
			inputJSON:    `{"email": "brian@gmail.com", "password": ""}`,
			statusCode:   422,
			errorMessage: "Required Password",
		},
		{
			inputJSON:    `{"email": "wrongemail@gmail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Incorrect Details",
		},
		{
			inputJSON:    `{"email": "invalidemail.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: "Invalid Email",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("Additional information: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}

		if v.statusCode == 422 && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
