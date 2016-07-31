package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/project-domino/domino-go/handlers"
	"github.com/project-domino/domino-go/handlers/api"
	"github.com/project-domino/domino-go/handlers/redirect"
	"github.com/project-domino/domino-go/middleware"
	"github.com/project-domino/domino-go/models"
	"github.com/project-domino/domino-go/templates"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// Enable/disable gin's debug mode.
	if Config.HTTP.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create and set up router.
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.Login())

	// Load and set up templates.
	template, err := templates.Load(Config.Templates)
	if err != nil {
		log.Fatal(err)
	}
	r.SetHTMLTemplate(template)

	// Authentication Routes
	r.GET("/login", handlers.Simple("login.html"))
	r.GET("/register", handlers.Simple("register.html"))
	r.POST("/login", api.Login)
	r.POST("/register", api.Register)
	r.POST("/logout", api.Logout)

	// View Routes
	r.GET("/", handlers.Simple("home.html"))

	r.Group("/account",
		middleware.RequireAuth()).
		GET("/", redirect.Account).
		GET("/profile",
			middleware.AddPageName("profile"),
			handlers.Simple("account-profile.html")).
		GET("/security",
			middleware.AddPageName("security"),
			handlers.Simple("account-security.html")).
		GET("/notifications",
			middleware.AddPageName("notifications"),
			handlers.Simple("account-notifications.html"))

	r.GET("/search/:searchType",
		middleware.LoadSearchItems(),
		middleware.LoadSearchVars(),
		handlers.Simple("search.html"))

	r.Group("/u/:username",
		middleware.LoadUser("Notes", "Collections", "Notes.Tags", "Collections.Tags")).
		GET("/", redirect.User).
		GET("/notes", handlers.Simple("user-notes.html")).
		GET("/collections", handlers.Simple("user-collections.html"))

	r.Group("/note",
		middleware.LoadNote("Author", "Tags"),
		middleware.VerifyNotePublic()).
		GET("/:noteID", handlers.Simple("individual-note.html")).
		GET("/:noteID/:note-name", handlers.Simple("individual-note.html"))

	r.Group("/collection",
		middleware.LoadCollection("Author", "Tags"),
		middleware.VerifyCollectionPublic()).
		GET("/:collectionID",
			handlers.Simple("collection.html")).
		GET("/:collectionID/note/:noteID",
			middleware.LoadNote("Author", "Tags"),
			handlers.Simple("collection-note.html")).
		GET("/:collectionID/note/:noteID/:noteName",
			middleware.LoadNote("Author", "Tags"),
			handlers.Simple("collection-note.html"))

	r.Group("/writer-panel",
		middleware.RequireAuth(),
		middleware.RequireUserType(models.Writer, models.Admin),
		middleware.LoadRequestUser("Notes", "Collections")).
		GET("/", redirect.WriterPanel).
		GET("/note",
			middleware.AddPageName("new-note"),
			handlers.Simple("new-note.html")).
		GET("/note/:noteID/edit",
			middleware.LoadNote("Author", "Tags"),
			middleware.VerifyNoteOwner(),
			handlers.Simple("edit-note.html")).
		GET("/collection",
			middleware.AddPageName("new-collection"),
			handlers.Simple("new-collection.html")).
		GET("/collection/:collectionID/edit",
			middleware.LoadCollection("Author", "Tags"),
			middleware.VerifyCollectionOwner(),
			handlers.Simple("edit-collection.html")).
		GET("/tag",
			middleware.AddPageName("new-tag"),
			handlers.Simple("new-tag.html"))

	// Start serving.
	if err := r.Run(Config.HTTP.ServeOn()); err != nil {
		log.Fatal(err)
	}
}
