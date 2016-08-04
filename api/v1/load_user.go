package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/models"
)

// LoadUser returns a middleware that attempts to load a user based on the user
// path parameter into the request context.
func LoadUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		err := db.Limit(1).
			Where("user_name = ?", c.Param("user")).
			Find(&users).
			Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errors.Apply(c, err)
		} else if len(users) == 0 {
			errors.Apply(c, errors.UserNotFound)
		} else {
			c.Set("user", &users[0])
		}
	}
}
