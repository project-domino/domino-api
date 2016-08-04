package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// AuthToken is a record of the tokens used for login.
type AuthToken struct {
	gorm.Model

	User    User      `json:"-"`
	UserID  uint      `json:"-"`
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
