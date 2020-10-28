package crconfig

import (
	"os"
)

type (
	// Switch represents one command line switch parameter
	Switch struct {
		// The command line switch. Usually starts with a '-'.
		Switch string
		// The corresponding environment variable this value overrides.
		EnvKey string
		// Description to be shown for this switch when calling -h or --help
		Description string
	}
)

// SetSwitches defines one or more cli switches.
// Use this if you plan a sophisticated command line application.
func SetSwitches(switches ...Switch) {
	args := map[string]string{}
	for _, sw := range switches {
		args[sw.Switch] = sw.EnvKey
	}
	parseSwitches(args)
}

// parses args, where args is a mapping switch -> env key
func parseSwitches(args map[string]string) {
	key, ok := "", false
	for _, arg := range os.Args {
		if ok {
			cli[key], ok = arg, false
		} else {
			key, ok = args[arg]
		}
	}
}
