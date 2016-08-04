package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/api"
)

func init() {
	api.RegisterVersion("1.0.0-alpha", func(db *gorm.DB, r gin.IRouter) {
		r.POST("/user/create",
			LoginUser(db, false, "admin"),
			CreateUser(db))
		r.POST("/user/auth_token/create",
			LoginUser(db, true),
			CreateAuthToken(db))
		r.Group("/user/:user", LoadUser(db)).
			GET("/", Terminal("user")).
			GET("/notifications", TODO)
	})
}
