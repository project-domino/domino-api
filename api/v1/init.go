package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/api"
)

func init() {
	api.RegisterVersion("1.0.0-alpha", func(db *gorm.DB, r gin.IRoutes) {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, "Hello, world!")
		})
		r.POST("/auth/create_token",
			User(db, true),
			CreateAuthToken(db))
	})
}
