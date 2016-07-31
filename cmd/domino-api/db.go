package main

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/config"
	"github.com/project-domino/domino-go/models"
)

// DB is the instance of the database.
var DB *gorm.DB

// OpenDB opens a database connection.
func OpenDB(dbConfig config.Database) {
	for {
		var err error
		log.Printf("Connecting to DB at %s...", dbConfig.URL)
		DB, err = gorm.Open(
			dbConfig.Type,
			dbConfig.URL,
		)
		if err == nil {
			break
		}
		log.Println(err)
		log.Println("Failed to connect to database.")
		time.Sleep(5 * time.Second)
	}

	DB.LogMode(dbConfig.Debug)
}

// SetupDB initializes the database with empty tables of all the needed types.
func SetupDB() error {
	// ewww sql
	if !DB.HasTable(&models.User{}) {
		DB.CreateTable(&models.User{})
		DB.Exec("ALTER TABLE users ADD COLUMN searchtext TSVECTOR")
		DB.Exec("CREATE INDEX searchtext_user_gin ON users USING GIN(searchtext)")
		DB.Exec(`CREATE TRIGGER ts_searchtext_user
			BEFORE INSERT OR UPDATE ON users
			FOR EACH ROW EXECUTE PROCEDURE
			tsvector_update_trigger('searchtext', 'pg_catalog.english', 'user_name')`)
	}
	if !DB.HasTable(&models.Note{}) {
		DB.CreateTable(&models.Note{})
		DB.Exec("ALTER TABLE notes ADD COLUMN searchtext TSVECTOR")
		DB.Exec("CREATE INDEX searchtext_note_gin ON notes USING GIN(searchtext)")
	}
	if !DB.HasTable(&models.Collection{}) {
		DB.CreateTable(&models.Collection{})
		DB.Exec("ALTER TABLE collections ADD COLUMN searchtext TSVECTOR")
		DB.Exec("CREATE INDEX searchtext_collection_gin ON collections USING GIN(searchtext)")
	}
	if !DB.HasTable(&models.Tag{}) {
		DB.CreateTable(&models.Tag{})
		DB.Exec("ALTER TABLE tags ADD COLUMN searchtext TSVECTOR")
		DB.Exec("CREATE INDEX searchtext_tag_gin ON tags USING GIN(searchtext)")
		DB.Exec(`CREATE TRIGGER ts_searchtext_tag
			BEFORE INSERT OR UPDATE ON tags
			FOR EACH ROW EXECUTE PROCEDURE
			tsvector_update_trigger('searchtext', 'pg_catalog.english', 'name', 'description')`)
	}
	setupTable(&models.AuthToken{})
	setupTable(&models.Comment{})
	setupTable(&models.CollectionNote{})
	setupTable(&models.Email{})

	return DB.Error
}

// setupTable creates a table for a specified struct if one doesn't exist.
func setupTable(val interface{}) {
	if !DB.HasTable(val) {
		DB.CreateTable(val)
	}
}
