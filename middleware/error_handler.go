package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/project-domino/domino-go/errors"
)

// ErrorHandler is a middleware that handles errors.
//
//     err := functionCall()
//     if err != nil {
//         errors.Apply(c, err)
//         return
//     }
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			var status = 500
			if err, ok := c.Errors.Last().Err.(*errors.Error); ok {
				status = err.Status
			}
			if c.Writer.Written() {
				status = c.Writer.Status()
			}

			c.JSON(status, c.Errors.JSON())
		}
	}
}
