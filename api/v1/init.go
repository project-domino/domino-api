package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/api"
)

func init() {
	api.RegisterVersion("1.0.0-alpha", func(db *gorm.DB, r gin.IRouter) {
		// Collection API
		r.POST("/collection", TODO)
		r.Group("/collection/:collection", LoadCollection(db, true)).
			GET("/", Terminal("collection")).
			POST("/", TODO)

		// Note API
		r.POST("/note", TODO)
		r.Group("/note/:note", LoadNote(db, true)).
			GET("/", Terminal("note")).
			POST("/", TODO)

		// TODO Search API
		// Or should the search server be separate?

		// Tags API
		r.GET("/tags", TODO)
		r.POST("/tags", TODO)

		// User API
		r.POST("/user/create",
			LoginUser(db, false, "admin"),
			CreateUser(db))
		r.POST("/user/auth_token/create",
			LoginUser(db, true),
			CreateAuthToken(db))
		r.Group("/user/:user", LoadUser(db)).
			GET("/", Terminal("user")).
			GET("/collections", TODO).
			GET("/notes", TODO).
			GET("/notifications", TODO)
	})
}
