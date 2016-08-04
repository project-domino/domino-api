package v1

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/project-domino/domino-go/errors"
	"github.com/project-domino/domino-go/models"
)

// CreateAuthToken creates an authentication token for the authenticated user.
//
// It optionally takes a single query parameter, "expires", that contains the
// time at which the token should expire, in seconds since the Unix epoch. If
// this parameter is not provided, it defaults to 168 hours in the future.
//
// It then responds with the token's JSON representation.
func CreateAuthToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		expires := time.Now()
		if expiresUnixString, ok := c.GetQuery("expires"); ok {
			expiresUnix, err := strconv.ParseInt(expiresUnixString, 64, 10)
			if err != nil {
				errors.Apply(c, errors.BadParameters)
				return
			}
			expires = time.Unix(expiresUnix, 0)
		}

		token := make([]byte, 16)
		if _, err := rand.Read(token); err != nil {
			c.AbortWithError(500, err)
			return
		}

		authToken := models.AuthToken{
			User:    *c.MustGet("user").(*models.User),
			Token:   hex.EncodeToString(token),
			Expires: expires,
		}
		if err := db.Create(&authToken).Error; err != nil {
			c.AbortWithError(500, err)
			return
		}
		c.JSON(200, authToken)
	}
}
