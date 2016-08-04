package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// These constants are the valid values for User.Type.
const (
	Admin   string = "admin"
	Writer         = "writer"
	General        = "general"
)

// A User is a user of the website. They can either be a admin, writer, or general.
type User struct {
	gorm.Model

	Type     string `json:"type"`
	Name     string `json:"name"`
	UserName string `json:"username"`
	Passhash string `json:"-"`

	// Only for writer
	UniversityID uint         `json:"-"`
	University   University   `json:"-"`
	Notes        []Note       `gorm:"ForeignKey:AuthorID" json:"-"`
	Collections  []Collection `gorm:"ForeignKey:AuthorID" json:"-"`

	Email          string `json:"-"` // `json:"email"`
	EmailVerified  bool   `json:"-"` // `json:"email_verified"`
	SendNewsletter bool   `json:"-"` // `json:"email_newsletter"`

	// Ranking Info
	UpvoteCollections   []Collection `gorm:"many2many:upvotecollection_user;" json:"-"`
	DownvoteCollections []Collection `gorm:"many2many:downvotecollection_user;" json:"-"`

	UpvoteComments   []Comment `gorm:"many2many:upvotecomment_user;" json:"-"`
	DownvoteComments []Comment `gorm:"many2many:downvotecomment_user;" json:"-"`

	UpvoteNotes   []Note `gorm:"many2many:upvotenote_user;" json:"-"`
	DownvoteNotes []Note `gorm:"many2many:downvotenote_user;" json:"-"`
}

// GetUser extracts a pointer to a User struct (or nil) from a Gin context.
func GetUser(c *gin.Context) *User {
	u, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := u.(*User)
	if !ok {
		return nil
	}
	return user
}

// CheckPassword checks if the provided password is correct. Note that it will
// return false whether the password was incorrect or an error is encountered,
// with no means to disambiguate the two.
func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Passhash), []byte(password)) == nil
}

// SetPassword hashes the provided password with bcrypt with the default cost,
// currently 10.
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		return err
	}
	u.Passhash = string(hash)
	return nil
}
