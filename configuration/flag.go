package configuration

import (
	"flag"
	"os"
)

type flags struct {
	driver string
	script string
	key    string
}

func (f *flags) parseFlags() error {
	set := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if !flag.Parsed() {
		flag.Parse()
	}
	set.StringVar(&f.driver, "driver", "", "config driver")
	set.StringVar(&f.script, "script", "", "config driver initialize script")
	set.StringVar(&f.key, "key", "config.yaml", "config name")
	return set.Parse(flag.Args())
}
