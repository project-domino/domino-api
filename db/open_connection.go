package db

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/config"
)

// Open opens a database connection
func Open(dbConfig config.Database) {
	opened := false
	for !opened {
		var err error
		log.Printf("Connecting to DB at %s...", dbConfig.URL)
		DB, err = gorm.Open(
			dbConfig.Type,
			dbConfig.URL,
		)
		opened = err == nil
		log.Println(err)
		time.Sleep(time.Second)
	}
	DB.LogMode(dbConfig.Debug)
}
