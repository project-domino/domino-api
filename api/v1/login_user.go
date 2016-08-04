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

// LoginUser returns a middleware that attempts to load a user by authentication
// token or username and password, as passed in the Authorization header, while
// also applying given constraints.
//
// If required is true, the request will fail if any error is encountered. If
// requiredTypes is present, an error will be encountered if the user's type is
// not one of the requiredTypes. If required is false, errors will be ignored.
//
// TODO(dev policy): Should errors be logged?
//
// Example Headers:
//     Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=
//     Authorization: Token 1ae79b24250d6eea3c10033f013af79a
func LoginUser(db *gorm.DB, required bool, requiredTypes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		user := models.GetUser(c)
		if user == nil {
			user, err = authUser(db, c.Request.Header.Get("Authorization"))
		}
		if err == nil {
			err = checkUserPermission(user, requiredTypes)
		}

		if err == nil {
			c.Set("user", user)
		} else if required {
			errors.Apply(c, err)
		}
	}
}

func authUser(db *gorm.DB, authorization string) (*models.User, error) {
	// Get and parse the Authorization header.
	if authorization == "" {
		return nil, errors.AuthRequired
	}
	parts := authorizationRegexp.FindStringSubmatch(authorization)
	if len(parts) != 3 {
		return nil, errors.InvalidAuthHeader
	}

	// Get the user from the Authorization header.
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
			return basicAuth(db, credentials[0], credentials[1])
		}
	case "Token":
		return tokenAuth(db, parts[2])
	default:
		err = errors.UnknownAuthMethod
	}
	return nil, err
}

func checkUserPermission(user *models.User, requiredTypes []string) error {
	ok := (requiredTypes == nil)
	for _, reqType := range requiredTypes {
		if user.Type == reqType {
			ok = true
			break
		}
	}
	if !ok {
		return errors.NoPermission
	}
	return nil
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
	user := &users[0]
	if !user.CheckPassword(password) {
		return user, errors.InvalidCredentials
	}
	return user, err
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
