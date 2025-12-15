package configuration

import (
	"flag"
)

type flags struct {
	driver string
	script string
	key    string
}

var options = new(flags)

func init() {
	flag.StringVar(&options.driver, "driver", "", "config driver")
	flag.StringVar(&options.script, "script", "", "config driver initialize script")
	flag.StringVar(&options.key, "key", "config.yaml", "config name")
}
