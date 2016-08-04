package api

import (
	"fmt"
	"log"
	"sync"

	"github.com/blang/semver"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Func is a type for API initialization functions.
type Func func(*gorm.DB, gin.IRoutes)

type api struct {
	version semver.Version
	routes  Func
}

var apiLock sync.Mutex
var knownAPIs []api

// RegisterVersion adds an API version.
func RegisterVersion(version string, routes Func) {
	sv := semver.MustParse(version)

	apiLock.Lock()
	knownAPIs = append(knownAPIs, api{sv, routes})
	apiLock.Unlock()

	log.Println("Registered API version", sv)
}

// CurrentVersion returns the most "advanced" version of the API (the version
// with the highest major version number).
func CurrentVersion() string {
	apiLock.Lock()
	defer apiLock.Unlock()

	var max semver.Version
	for _, api := range knownAPIs {
		if api.version.GTE(max) {
			max = api.version
		}
	}
	return fmt.Sprintf("v%d", max.Major)
}

// AllVersions returns a slice of all supported API versions.
func AllVersions() []string {
	var out []string
	for v := range AllVersionRoutes() {
		out = append(out, v)
	}
	return out
}

// AllVersionRoutes returns a map of version strings to routes for all supported
// API versions.
func AllVersionRoutes() map[string]Func {
	apiLock.Lock()
	defer apiLock.Unlock()

	// First, choose the API versions.
	m := make(map[uint64]api)
	for _, api := range knownAPIs {
		o := m[api.version.Major]
		if api.version.GTE(o.version) {
			m[api.version.Major] = api
		}
	}

	// Then, build the output.
	out := make(map[string]Func, len(m))
	for _, api := range m {
		out[fmt.Sprintf("v%d", api.version.Major)] = api.routes
	}
	return out
}
