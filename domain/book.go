package domain

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Url            string   `json:"url" gorm:"not null"` // openlibrary.org/works/{OpenLibraryKey}/Layered_Money
	OpenLibraryKey string   `json:"open_library_key" gorm:"not null;unique"`
	Title          string   `json:"title" gorm:"not null"`
	Year           uint16   `json:"year"` // first_publish_year
	NumberOfPages  uint16   `json:"number_of_pages"`
	OpenLibraryID  uint     `json:"open_library_id"` // cover_i: covers.openlibrary.org/b/id/{OpenLibraryID}
	Language       string   `json:"language" gorm:"size:20"`
	Authors        []Author `json:"authors" gorm:"many2many:authors_books"` // authors_facet: {author_key, full_name}
	Chapter        string   `json:"chapter" gorm:"-;size:256"`              // to books_ideas
}
