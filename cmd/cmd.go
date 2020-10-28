package cmd

import (
	"fmt"
	"os"
	"reflect"
)

var (
	// PanicEmptyCommand set cmd to panic if no matching command was found on Start
	PanicEmptyCommand bool

	cmd Command
)

// Command is to be implemented for new commands
type Command interface {
	Init()
	Start() error
	Clean()
}

// Reset empties stored values - good for testing
func Reset() {
	cmd = nil
}

// RegisterCmd routes a start command to a certain argument name.
// Provide a pointer to a struct, implementing the Command interface.
func RegisterCmd(name string, c Command) {
	arg := ""
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}

	if arg == name {
		if reflect.ValueOf(c).Kind() != reflect.Ptr {
			panic("command must be pointer")
		}
		cmd = c
	}
}

// Start starts the matching command
func Start() (err error) {
	if cmd != nil {
		cmd.Init()
		err = cmd.Start()
		cmd.Clean()
		return err
	}
	if PanicEmptyCommand {
		panic("no matching command")
	}
	return fmt.Errorf("No matching command found")
}
