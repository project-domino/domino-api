package config

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// Database is a type for database configuration settings.
type Database struct {
	Debug bool   `toml:"debug"`
	Type  string `toml:"type"`
	URL   string `toml:"url"`
}

// DefaultDatabase is the default database configuration.
var DefaultDatabase = Database{
	Debug: false,
	Type:  "postgres",
	URL:   "dbname=domino host=domino sslmode=disable user=domino",
}

// Open opens a connection to the database.
func (config Database) Open() *gorm.DB {
	for {
		log.Printf("Connecting to database at %s...", config.URL)
		db, err := gorm.Open(config.Type, config.URL)
		if err == nil {
			log.Println("Connected to database.")
			db.LogMode(config.Debug)
			return db
		}
		log.Println(err)
		log.Println("Failed to connect to database.")
		time.Sleep(5 * time.Second)
	}
}
