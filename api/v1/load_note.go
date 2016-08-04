package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/models"
)

// LoadNote returns a middleware that attempts to load a note based on the note
// path parameter into the request context.
func LoadNote(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var notes []models.Note
		err := db.Limit(1).
			Where("id = ?", c.Param("note")).
			Find(&notes).
			Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errors.Apply(c, err)
		} else if len(notes) == 0 {
			errors.Apply(c, errors.NoteNotFound)
		} else {
			c.Set("note", &notes[0])
		}
	}
}
