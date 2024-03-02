package utils

import (
	"fmt"
)

// Global version infomation

const (
	SERVICE_VERSION = "3.0.0"

	// date +%FT%T%z  // date +'%Y%m%d'
	BUILD_TIME = "2024-03-31T00:00:00+0800"

	// go version
	GO_VERSION = "1.21.0"
)

func Version(app string) string {
	return fmt.Sprintf(`{"app": "%s", "version": "%s", "build_time": "%s", "go_version": "%s"}`,
		app, SERVICE_VERSION, BUILD_TIME, GO_VERSION)
}
