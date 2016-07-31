package api

import "github.com/gin-gonic/gin"

// Version is a handler that returns the current (and all known) API versions as
// a JSON object.
func Version(c *gin.Context) {
	c.JSON(200, gin.H{
		"current":  CurrentVersion(),
		"versions": AllVersions(),
	})
}
