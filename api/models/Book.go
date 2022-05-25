package models

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string `gorm:"size:255;not null" json:"title"`
	Author      string `gorm:"size:255;not null" json:"author"`
	Isbn        string `gorm:"size:255;not null" json:"isbn"`
	Description string `gorm:"size:4096;not null" json:"description"`
	Available   bool   `gorm:"not null" json:"available"`
}

func (b *Book) Prepare() {
	b.Title = html.EscapeString(strings.TrimSpace(b.Title))
	b.Author = html.EscapeString(strings.TrimSpace(b.Author))
	b.Isbn = html.EscapeString(strings.TrimSpace(b.Isbn))
	b.Description = html.EscapeString(strings.TrimSpace(b.Description))
	// b.Description = html.EscapeString(strings.TrimSpace(b.Description))
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.New("Required Title")
	}
	if b.Author == "" {
		return errors.New("Required Author")
	}
	if b.Isbn == "" {
		return errors.New("Required Isbn")
	}
	if b.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}

func (b *Book) SaveBook(db *gorm.DB) (*Book, error) {
	var err error
	err = db.Create(&b).Error
	if err != nil {
		return &Book{}, err
	}
	return b, nil
}

func (p *Book) FindAllBooks(db *gorm.DB) (*[]Book, error) {
	var err error
	books := []Book{}
	err = db.Order("updated_at desc").Limit(100).Find(&books).Error
	// err = db.Table("books").Select("books.id as id, books.title as title, books.author as author, books.isbn as isbn, books.description as description, books.created_at as created_at, books.updated_at as updated_at, checkouts.user_id as user_id, checkouts.checked_in as checked_in").Order("books.updated_at desc").Limit(100).Joins("LEFT JOIN checkouts on checkouts.book_id = books.id").Find(&booksAvailable).Error
	if err != nil {
		return &[]Book{}, err
	}
	return &books, nil
}

func FindBookByID(db *gorm.DB, bid uint64) (*[]Book, error) {
	var err error

	books := []Book{}
	err = db.Where("id = ?", bid).Find(&books).Error
	if err != nil {
		return &[]Book{}, err
	}

	return &books, nil
}

func (b *Book) UpdateABook(db *gorm.DB) (*Book, error) {
	var err error

	books, err := FindBookByID(db, uint64(b.ID))
	if err != nil {
		return &Book{}, err
	}
	if len(*books) == 0 {
		return &Book{}, errors.New("book does not exist")
	}

	b.Available = (*books)[0].Available
	err = db.Model(&Book{}).Where("id = ?", b.ID).Updates(Book{
		Title:       b.Title,
		Author:      b.Author,
		Isbn:        b.Isbn,
		Description: b.Description,
		Available:   (*books)[0].Available,
	}).Error
	fmt.Println(err)
	if err != nil {
		return &Book{}, err
	}
	return b, nil
}

func (b *Book) DeleteABook(db *gorm.DB) (int64, error) {
	db.Delete(&b)

	if db.Error != nil {
		return 0, db.Error
	}

	return db.RowsAffected, nil
}
