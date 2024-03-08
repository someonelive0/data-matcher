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

	IDSS_BANNER = `
	╔══╗╔═══╗╔═══╗╔═══╗
	╚╣╠╝╚╗╔╗║║╔═╗║║╔═╗║
	 ║║  ║║║║║╚══╗║╚══╗
	 ║║  ║║║║╚══╗║╚══╗║
	╔╣╠╗╔╝╚╝║║╚═╝║║╚═╝║
	╚══╝╚═══╝╚═══╝╚═══╝	
`
)

func Version(app string) string {
	return fmt.Sprintf(`{"app": "%s", "version": "%s", "build_time": "%s", "go_version": "%s"}`,
		app, SERVICE_VERSION, BUILD_TIME, GO_VERSION)
}

func ShowBanner() {
	fmt.Printf("%s\n", IDSS_BANNER)
	fmt.Printf("规则匹配 data-matcher %s  Copyright (C) 2024 IDSS\n", SERVICE_VERSION)
}

func ShowBannerForApp(app, version, build_time string) {
	fmt.Printf("%s\n", IDSS_BANNER)
	fmt.Printf("规则匹配 data-matcher 3.0  Copyright (C) 2024 IDSS\n")
	fmt.Printf("%s version %s, build on %s\n\n", app, version, build_time)
}
