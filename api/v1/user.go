package v1

import (
	"encoding/base64"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/models"
)

var authorizationRegexp = regexp.MustCompile("([A-Za-z]+) (.+)")

// User returns a middleware that attempts to load a user by authentication
// token or username and password, as passed in the Authorization header, while
// also applying given constraints.
//
// Example Headers:
//     Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=
//     Authorization: Token TODO
func User(db *gorm.DB, required bool, requiredTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get and parse the authentication header.
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			if required {
				errors.AuthRequired.Apply(c)
			} else {
				c.Next()
			}
			return
		}
		parts := authorizationRegexp.FindStringSubmatch(authorization)
		if len(parts) != 3 {
			errors.InvalidAuthHeader.Apply(c)
			return
		}

		// Get the user from the Authorization header.
		var user *models.User
		var err error
		switch parts[1] {
		case "Basic":
			var decode []byte
			decode, err = base64.StdEncoding.DecodeString(parts[2])
			if err != nil {
				break
			}
			credentials := strings.SplitN(string(decode), ":", 2)
			if len(credentials) != 2 {
				err = errors.InvalidAuthHeader
			} else {
				user, err = basicAuth(db, credentials[0], credentials[1])
			}
		case "Token":
			user, err = tokenAuth(db, parts[2])
		default:
			err = errors.UnknownAuthMethod
		}
		if err != nil {
			if err2, ok := err.(*errors.Error); ok {
				err2.Apply(c)
			} else {
				c.AbortWithError(403, err)
			}
		}

		// Check that the user has one of the required types, if applicable.
		ok := (requiredTypes == nil)
		for _, reqType := range requiredTypes {
			if user.Type == reqType {
				ok = true
				break
			}
		}
		if !ok {
			errors.NoPermission.Apply(c)
			return
		}

		// Add the user to the context.
		c.Set("user", user)
	}
}

func basicAuth(db *gorm.DB, username, password string) (*models.User, error) {
	var users []models.User
	err := db.Limit(1).
		Where("user_name = ?", username).
		Find(&users).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.InvalidCredentials
	}
	return &users[0], nil
}

func tokenAuth(db *gorm.DB, token string) (*models.User, error) {
	var authEntries []models.AuthToken
	err := db.Limit(1).
		Preload("User").
		Where("token = ?", token).
		Where("expires > ?", time.Now()).
		Find(&authEntries).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if len(authEntries) == 0 {
		return nil, errors.InvalidCredentials
	}
	return &authEntries[0].User, nil
}
