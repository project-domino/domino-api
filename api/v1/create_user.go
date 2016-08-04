package v1

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/models"
)

var usernameRegexp = regexp.MustCompile("[A-Za-z0-9_.-]+")

// CreateUser creates a new user.
//
// The handler takes the body parameters name, username and password, which are
// each strings corresponding to the appropriate values. The username must be
// composed solely of the ASCII alphanumerics, underscores, periods, and dashes.
// Optionally, a body parameter type may be provided. It must be one of the
// values "admin", "writer", or "general". If its value is not "general", the
// request must be performed by an authenticated user whose type is "admin".
//
// The handler then responds with the user's JSON representation.
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the mandatory query parameters.
		name, ok := c.GetPostForm("name")
		if !ok {
			errors.Apply(c, errors.MissingParameters)
			return
		}
		username, ok := c.GetPostForm("username")
		if !ok {
			errors.Apply(c, errors.MissingParameters)
			return
		}
		if !usernameRegexp.MatchString(username) {
			errors.Apply(c, errors.BadParameters)
			return
		}
		password, ok := c.GetPostForm("password")
		if !ok {
			errors.Apply(c, errors.MissingParameters)
			return
		}

		// Try getting type.
		userType, ok := c.GetPostForm("type")
		if !ok {
			userType = models.General
		}
		if userType != models.Admin && userType != models.Writer && userType != models.General {
			errors.Apply(c, errors.BadParameters)
			return
		}
		if _, ok := c.Get("user"); userType != models.General && !ok {
			errors.Apply(c, errors.NoPermission)
			return
		}

		// Check if any users have the same username.
		var checkUsers []models.User
		err := db.Where("user_name = ?", username).
			Find(&checkUsers).
			Error
		if err != nil && err != gorm.ErrRecordNotFound {
			errors.Apply(c, err)
			return
		}
		if len(checkUsers) != 0 {
			errors.Apply(c, errors.UserExists)
			return
		}

		// Create the user.
		user := &models.User{
			Type:     userType,
			Name:     name,
			UserName: username,
		}
		if err := user.SetPassword(password); err != nil {
			errors.Apply(c, err)
			return
		}
		if err := db.Create(user).Error; err != nil {
			errors.Apply(c, err)
			return
		}

		// Respond with the user's JSON.
		c.JSON(200, user)
	}
}
