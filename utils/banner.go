package utils

import (
	"fmt"
)

func ShowBanner() {
	fmt.Printf("%s\n", IDSS_BANNER)
	fmt.Printf("规则匹配 data-matcher 3.0  Copyright (C) 2024 IDSS\n")
}

func ShowBannerForApp(app, version, build_time string) {
	fmt.Printf("%s\n", IDSS_BANNER)
	fmt.Printf("规则匹配 data-matcher 3.0  Copyright (C) 2024 IDSS\n")
	fmt.Printf("%s version %s, build on %s\n\n", app, version, build_time)
}

const (
	IDSS_BANNER = `
	╔══╗╔═══╗╔═══╗╔═══╗
	╚╣╠╝╚╗╔╗║║╔═╗║║╔═╗║
	 ║║  ║║║║║╚══╗║╚══╗
	 ║║  ║║║║╚══╗║╚══╗║
	╔╣╠╗╔╝╚╝║║╚═╝║║╚═╝║
	╚══╝╚═══╝╚═══╝╚═══╝	
`
)
