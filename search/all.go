package search

import (
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/models"
)

// AllResponse holds the objects returned by a search query for all
// items
type AllResponse struct {
	Notes       []models.Note
	Collections []models.Collection
	Users       []models.User
	Tags        []models.Tag
}

// All returns a struct containing a search result for all items
func All(db *gorm.DB, q string, items int) (AllResponse, error) {
	var response AllResponse
	var searchErr error

	response.Notes, searchErr = Notes(db, q, 1, items)
	response.Collections, searchErr = Collections(db, q, 1, items)
	response.Users, searchErr = Users(db, q, 1, items)
	response.Tags, searchErr = Tags(db, q, 1, items)

	return response, searchErr
}
