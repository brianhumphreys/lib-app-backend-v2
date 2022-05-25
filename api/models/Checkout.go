package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Checkout struct {
	gorm.Model
	UserId    uint   `gorm:"size:100;not null;" json:"user_id"`
	BookId    uint64 `gorm:"size:100;not null;" json:"book_id"`
	CheckedIn bool   `json:"checked_in"`
}

type BookRecord struct {
	Title      string `gorm:"size:512;" json:"title"`
	Author     string `gorm:"size:100;" json:"author"`
	CheckedIn  string `gorm:"size:100;" json:"checked_in"`
	CheckedOut string `gorm:"size:100;" json:"checked_out"`
}

type UserRecord struct {
	Email      string `gorm:"size:512;" json:"email"`
	CheckedIn  string `gorm:"size:100;" json:"checked_in"`
	CheckedOut string `gorm:"size:100;" json:"checked_out"`
}

func GetBookCheckoutHistoryOfUserWithID(db *gorm.DB, uid uint) (*[]BookRecord, error) {
	var err error
	books := []BookRecord{}
	err = db.Table("checkouts").Order("checkouts.created_at desc").Select("books.title as title, books.author as author, checkouts.updated_at as checked_in, checkouts.created_at as checked_out").Joins("RIGHT JOIN books on books.id = checkouts.book_id").Where("checkouts.user_id = ?", uid).Limit(100).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return &books, nil
}

func GetUserCheckoutHistoryOfBookWithID(db *gorm.DB, bid uint64) (*[]UserRecord, error) {
	var err error
	users := []UserRecord{}
	err = db.Table("checkouts").Order("checkouts.created_at desc").Select("users.email as email, checkouts.updated_at as checked_in, checkouts.created_at as checked_out").Joins("RIGHT JOIN users on checkouts.user_id = users.id").Where("checkouts.book_id = ?", bid).Limit(100).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func GetCurrentOwnerOfBookWithID(db *gorm.DB, bid uint64) ([]uint64, error) {
	var err error
	var uid []uint64
	err = db.Table("checkouts").Select("user_id").Where("book_id = ? AND checked_in = false", bid).Find(&uid).Error
	if err != nil {
		return nil, err
	}
	return uid, nil
}

func GetCurrentlyCheckedOutBooksOfUserWithID(db *gorm.DB, uid uint) (*[]Book, error) {
	var err error
	books := []Book{}
	err = db.Table("users").Select("books.id, books.title, books.author, books.isbn, books.description").Joins("JOIN checkouts on checkouts.user_id = users.id").Joins("JOIN books on books.id = checkouts.book_id").Where("checkouts.user_id = ? AND checkouts.checked_in = false AND books.deleted_at is NULL", uid).Limit(100).Find(&books).Error
	if err != nil {
		return nil, err
	}
	for i := range books {
		fmt.Println(books[i].Title)
	}
	return &books, nil
}

func (c *Checkout) MakeACheckout(db *gorm.DB) error {
	var err error

	books, err := FindBookByID(db, c.BookId)
	if err != nil {
		return err
	}
	if len(*books) == 0 {
		return errors.New("book does not exist")
	}

	// todo: create transactions
	err = db.Model(&Book{}).Where("id = ?", (*books)[0].ID).Update("available", false).Error
	if err != nil {
		return err
	}

	err = db.Create(&c).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Checkout) HasUserCheckedBook(db *gorm.DB) error {
	var err error
	err = db.Table("checkouts").Where("user_id = ? AND book_id = ? AND checked_in = false", c.UserId, c.BookId).First(&c).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *Checkout) CheckinABook(db *gorm.DB) error {
	var err error

	books, err := FindBookByID(db, c.BookId)
	if err != nil {
		return err
	}
	if len(*books) == 0 {
		return errors.New("book does not exist")
	}

	// todo: create transactions
	err = db.Model(&Book{}).Where("id = ?", (*books)[0].ID).Update("available", true).Error
	if err != nil {
		return err
	}

	err = db.Model(&c).Where("book_id = ? AND checked_in = false", c.BookId).Update("checked_in", true).Error
	if err != nil {
		return err
	}
	return nil
}
