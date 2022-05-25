package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/brianhumphreys/library_app/api/auth"
	"github.com/brianhumphreys/library_app/api/models"
	"github.com/brianhumphreys/library_app/api/responses"
	"github.com/gorilla/mux"
)

func (s *Server) CheckoutABook(w http.ResponseWriter, r *http.Request) {

	var checkout models.Checkout

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&checkout)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check that the user has the correct ID
	tokenID, _, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(checkout.UserId) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	// check that book exists
	book, err := models.FindBookByID(s.DB, checkout.BookId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if len(*book) == 0 {
		responses.ERROR(w, http.StatusNotFound, errors.New("This book was not found in the library"))
		return
	}

	// check user has maxed out the number of books they are allowed to checkout
	currentlyCheckedOutBooks, err := models.GetCurrentlyCheckedOutBooksOfUserWithID(s.DB, checkout.UserId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if len(*currentlyCheckedOutBooks) >= 5 {
		responses.ERROR(w, http.StatusTooManyRequests, errors.New("You have checked out too many books."))
		return
	}

	// check if the book is already checked out
	user_ids, err := models.GetCurrentOwnerOfBookWithID(s.DB, checkout.BookId)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	if len(user_ids) > 0 {
		responses.ERROR(w, http.StatusTeapot, errors.New("Someone has checked this book out"))
		return
	}

	fmt.Printf("Checking out book with ID: %d, by user: %d\n", checkout.BookId, checkout.UserId)
	// check out the book
	checkout.MakeACheckout(s.DB)
	responses.JSON(w, http.StatusCreated, checkout)
}

func (server *Server) CheckinABook(w http.ResponseWriter, r *http.Request) {
	var checkin models.Checkout
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&checkin)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	// check that the user has the correct ID
	tokenID, _, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(checkin.UserId) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	fmt.Printf("Returning book withh ID: %d, by user: %d\n", checkin.BookId, checkin.UserId)
	// Make sure this user has checked out the book that they are attempting to check in
	err = checkin.HasUserCheckedBook(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusForbidden, errors.New("You do not currently have this book checked out"))
		return
	}

	// check in the book
	checkin.CheckinABook(server.DB)
	responses.JSON(w, http.StatusAccepted, checkin)
}

func (server *Server) GetBookCheckoutHistoryOfUserWithID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check that the user has the correct ID
	tokenID, _, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	books, err := models.GetBookCheckoutHistoryOfUserWithID(server.DB, uint(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}

func (server *Server) GetUserCheckoutHistoryOfBookWithID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	bid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	users, err := models.GetUserCheckoutHistoryOfBookWithID(server.DB, bid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) GetCurrentlyCheckedOutBooksOfUserWithID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// check that the user has the correct ID
	tokenID, _, err := auth.ExtractTokenIDAndRole(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	books, err := models.GetCurrentlyCheckedOutBooksOfUserWithID(server.DB, uint(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, books)
}
