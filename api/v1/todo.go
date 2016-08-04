package v1

import "github.com/gin-gonic/gin"

// TODO is a terminal handler that outputs "TODO".
func TODO(c *gin.Context) {
	c.Data(200, "text/plain", []byte("TODO"))
}
