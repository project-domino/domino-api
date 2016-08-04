package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Apply applies an error to the context.
func Apply(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	if h, ok := err.(interface {
		HTTPStatus() int
	}); ok {
		status = h.HTTPStatus()
	}
	c.AbortWithError(status, err)
}
