package search

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/models"
)

// Tags returns all tags that match a given query
func Tags(db *gorm.DB, q string, items int, page int) ([]models.Tag, error) {
	var tags []models.Tag

	searchQuery, err := ParseQuery(q)
	if err != nil {
		return tags, err
	}
	// qSelectors := searchQuery.Selectors
	qText := strings.Join(searchQuery.Text, " & ")

	if err := db.Where(queryFormat, qText).
		Find(&tags).
		Limit(items).
		Offset(page * items).
		Error; err != nil {
		return tags, err
	}

	return tags, nil
}
