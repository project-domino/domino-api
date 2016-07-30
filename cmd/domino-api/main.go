package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/project-domino/domino-go/db"
	"github.com/project-domino/domino-go/middleware"
)

func main() {
	db.Open(Config.Database)
	if err := db.Setup(); err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

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

	// Add routes.

	// Start serving.
	Must(r.Run(fmt.Sprintf(":%d", Config.HTTP.Port)))
}
