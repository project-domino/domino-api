package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/project-domino/domino-go/errors"
)

// Terminal is a terminal handler that outputs the given variable from the
// request context as JSON.
func Terminal(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get(name)
		if ok {
			c.JSON(200, val)
		} else {
			errors.Apply(c, errors.UnknownValue)
		}
	}
}
