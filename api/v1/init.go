package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/api"
)

func init() {
	api.RegisterVersion("1.0.0-alpha", func(db *gorm.DB, r gin.IRoutes) {
		r.POST("/user/create",
			User(db, false, "admin"),
			CreateUser(db))
		r.POST("/user/auth_token/create",
			User(db, true),
			CreateAuthToken(db))
	})
}
