package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/brianhumphreys/library_app/api/auth"
	"github.com/brianhumphreys/library_app/api/models"
	"github.com/brianhumphreys/library_app/api/responses"
	"github.com/brianhumphreys/library_app/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	fmt.Printf("Loggin in user: %d\n", user.ID)

	signedUser, token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"email": signedUser.Email,
		"role":  signedUser.Role,
		"id":    signedUser.ID,
	})
}

func (server *Server) SignIn(email, password string) (*models.User, string, error) {
	var err error

	user := models.User{}

	err = server.DB.Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return &models.User{}, "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return &models.User{}, "", err
	}
	token, err := auth.CreateToken(user)
	return &user, token, err
}
