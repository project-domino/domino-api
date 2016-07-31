package search

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/models"
)

// Collections returns all collections that match a given query
func Collections(db *gorm.DB, q string, items int, page int) ([]models.Collection, error) {
	var collections []models.Collection

	searchQuery, err := ParseQuery(q)
	if err != nil {
		return collections, err
	}
	// qSelectors := searchQuery.Selectors
	qText := strings.Join(searchQuery.Text, " & ")

	if q != "" {
		if err := db.
			Preload("Tags").
			Where(queryFormat, qText).
			Where("published = ?", true).
			Find(&collections).
			Limit(items).
			Offset(page * items).
			Error; err != nil {
			return collections, err
		}
	}
	return collections, nil
}
