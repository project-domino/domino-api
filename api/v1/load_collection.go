package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/models"
)

// LoadCollection returns a middleware that attempts to load a collection based
// on the collection path parameter into the request context.
func LoadCollection(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var collections []models.Collection
		err := db.Limit(1).
			Where("id = ?", c.Param("collection")).
			Find(&collections).
			Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errors.Apply(c, err)
		} else if len(collections) == 0 {
			errors.Apply(c, errors.CollectionNotFound)
		} else {
			c.Set("collection", &collections[0])
		}
	}
}
