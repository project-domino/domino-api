package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/models"
)

// LoadUser loads a user into the request context
func LoadUser(db *gorm.DB, objects ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Acquire username from URL
		username := c.Param("username")

		// Set objects to be preloaded to db
		preloadedDB := db.Where("user_name = ?", username)
		for _, object := range objects {
			preloadedDB = preloadedDB.Preload(object)
		}

		// Query for user and set context
		var user models.User
		if err := preloadedDB.First(&user).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.Set("pageUser", user)

		c.Next()
	}
}
