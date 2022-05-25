package modeltests

import (
	"log"
	"testing"

	"github.com/brianhumphreys/library_app/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestFindAllUsers(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	users, err := userInstance.FindAllUsers(server.DB)
	if err != nil {
		t.Errorf("error running TestFindAllUsers: %v\n", err)
		return
	}
	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}
	newUser := models.User{
		Email:    "test1@gmail.com",
		Password: "password",
		Role:     "user",
	}
	savedUser, err := newUser.SaveUser(server.DB)
	if err != nil {
		t.Errorf("Error running TestSaveUser: %v\n", err)
		return
	}
	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Password, savedUser.Password)
	assert.Equal(t, newUser.Role, savedUser.Role)
}

func TestGetUserByID(t *testing.T) {

	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("User table could not be seeded: %v", err)
	}
	userInstance := models.User{}
	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)
	if err != nil {
		t.Errorf("Error running TestGetUserByID: %v\n", err)
		return
	}
	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Password, user.Password)
	assert.Equal(t, foundUser.Role, user.Role)
}

func TestDeleteAUser(t *testing.T) {
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("User table could not be seeded: %v", err)
	}

	_, err = user.DeleteAUser(server.DB, user.ID)
	if err != nil {
		t.Errorf("There was an error while deleting the user: %v\n", err)
		return
	}

	foundUser2, err := userInstance.FindUserByID(server.DB, user.ID)

	assert.Equal(t, foundUser2.ID, uint(0))

}
