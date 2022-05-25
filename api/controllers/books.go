package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/brianhumphreys/library_app/api/auth"
	"github.com/brianhumphreys/library_app/api/models"
	"github.com/brianhumphreys/library_app/api/responses"
	"github.com/brianhumphreys/library_app/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) CreateBook(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	book := models.Book{}
	err = json.Unmarshal(body, &book)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, role, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if role != "admin" {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	book.Prepare()
	err = book.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	book.Available = true
	bookCreated, err := book.SaveBook(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	fmt.Printf("Found book with ID: %d\n", bookCreated.ID)

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, bookCreated.ID))
	responses.JSON(w, http.StatusCreated, bookCreated)
}

func (server *Server) GetBooks(w http.ResponseWriter, r *http.Request) {
	book := models.Book{}

	books, err := book.FindAllBooks(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}

func (server *Server) GetBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	book, err := models.FindBookByID(server.DB, bid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if len(*book) == 0 {
		responses.ERROR(w, http.StatusNotFound, errors.New("This book was not found in the library"))
	}
	responses.JSON(w, http.StatusOK, (*book)[0])
}

func (server *Server) UpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, role, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if role != "admin" {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	book := models.Book{}
	err = server.DB.Model(models.Book{}).Where("id = ?", bid).Take(&book).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Book not found"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	bookUpdate := models.Book{}
	err = json.Unmarshal(body, &bookUpdate)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	bookUpdate.Prepare()
	err = bookUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	bookUpdate.ID = book.ID

	bookUpdated, err := bookUpdate.UpdateABook(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}

	fmt.Printf("Updated book with ID: %d\n", bookUpdated.ID)

	responses.JSON(w, http.StatusOK, bookUpdated)
}

func (server *Server) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	_, role, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if role != "admin" {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	book := models.Book{}
	err = server.DB.Model(models.Book{}).Where("id = ?", bid).Take(&book).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	fmt.Printf("Deleting book with ID: %d\n", book.ID)
	_, err = book.DeleteABook(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", bid))
	responses.JSON(w, http.StatusNoContent, "")
}
