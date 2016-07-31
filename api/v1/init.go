package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/project-domino/domino-go/api"
)

func init() {
	api.RegisterVersion("1.0.0-alpha", func(r gin.IRoutes) {
		r.GET("/", func(c *gin.Context) {
			c.JSON(200, "Hello, world!")
		})
	})
}
