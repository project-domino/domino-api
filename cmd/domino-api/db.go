package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/config"
	"github.com/project-domino/domino-go/models"
)

// OpenDB opens a database connection.
func OpenDB(dbConfig config.Database) *gorm.DB {
	for {
		log.Printf("Connecting to DB at %s...", dbConfig.URL)
		db, err := gorm.Open(
			dbConfig.Type,
			dbConfig.URL,
		)
		if err == nil {
			db.LogMode(dbConfig.Debug)
			return db
		}
		log.Println(err)
		log.Println("Failed to connect to database.")
		time.Sleep(5 * time.Second)
	}
}

// SetupDB initializes the database with empty tables of all the needed types.
func SetupDB(db *gorm.DB) error {
	// ewww sql
	if !db.HasTable(&models.User{}) {
		db.CreateTable(&models.User{})
		db.Exec("ALTER TABLE users ADD COLUMN searchtext TSVECTOR")
		db.Exec("CREATE INDEX searchtext_user_gin ON users USING GIN(searchtext)")
		db.Exec(`CREATE TRIGGER ts_searchtext_user
			BEFORE INSERT OR UPDATE ON users
			FOR EACH ROW EXECUTE PROCEDURE
			tsvector_update_trigger('searchtext', 'pg_catalog.english', 'user_name')`)
	}
	if !db.HasTable(&models.Note{}) {
		db.CreateTable(&models.Note{})
		db.Exec("ALTER TABLE notes ADD COLUMN searchtext TSVECTOR")
		db.Exec("CREATE INDEX searchtext_note_gin ON notes USING GIN(searchtext)")
	}
	if !db.HasTable(&models.Collection{}) {
		db.CreateTable(&models.Collection{})
		db.Exec("ALTER TABLE collections ADD COLUMN searchtext TSVECTOR")
		db.Exec("CREATE INDEX searchtext_collection_gin ON collections USING GIN(searchtext)")
	}
	if !db.HasTable(&models.Tag{}) {
		db.CreateTable(&models.Tag{})
		db.Exec("ALTER TABLE tags ADD COLUMN searchtext TSVECTOR")
		db.Exec("CREATE INDEX searchtext_tag_gin ON tags USING GIN(searchtext)")
		db.Exec(`CREATE TRIGGER ts_searchtext_tag
			BEFORE INSERT OR UPDATE ON tags
			FOR EACH ROW EXECUTE PROCEDURE
			tsvector_update_trigger('searchtext', 'pg_catalog.english', 'name', 'description')`)
	}
	setupTable(db, &models.AuthToken{})
	setupTable(db, &models.Comment{})
	setupTable(db, &models.CollectionNote{})
	setupTable(db, &models.Email{})

	return db.Error
}

// setupTable creates a table for a specified struct if one doesn't exist.
func setupTable(db *gorm.DB, val interface{}) {
	if !db.HasTable(val) {
		db.CreateTable(val)
	}
}
