package cmd_test

import (
	"fmt"
	"os"
	"testing"

	"cleverreach.com/crtools/cmd"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	m.Run()
}

type FirstCMD struct {
	str []string
}

func (c *FirstCMD) Init() {
	c.str = append(c.str, "init")
}
func (c *FirstCMD) Start() error {
	c.str = append(c.str, "start")
	return fmt.Errorf("we're cool")
}
func (c *FirstCMD) Clean() {
	c.str = append(c.str, "clean")
}

type SecondCMD struct {
	str []string
}

func (c *SecondCMD) Init() {
	c.str = append(c.str, "2ndinit")
}
func (c *SecondCMD) Start() error {
	c.str = append(c.str, "2ndstart")
	return fmt.Errorf("we're fine")
}
func (c *SecondCMD) Clean() {
	c.str = append(c.str, "2ndclear")
}

func TestStartCommands(t *testing.T) {
	test := assert.New(t)

	{
		os.Args = []string{"cmd_test", "first"}

		cmd1 := FirstCMD{}

		cmd.RegisterCmd("first", &cmd1)
		cmd.RegisterCmd("second", &SecondCMD{})

		err := cmd.Start()
		ok := test.NotNil(err)
		ok = ok && test.Equal("we're cool", err.Error())

		ok = ok && test.Len(cmd1.str, 3)
		ok = ok && test.Equal("init", cmd1.str[0])
		ok = ok && test.Equal("start", cmd1.str[1])
		ok = ok && test.Equal("clean", cmd1.str[2])
	}

	{
		os.Args = []string{"cmd_test", "second"}

		cmd2 := SecondCMD{}

		cmd.RegisterCmd("first", &FirstCMD{})
		cmd.RegisterCmd("second", &cmd2)

		err := cmd.Start()
		ok := test.NotNil(err)
		ok = ok && test.Equal("we're fine", err.Error())

		ok = ok && test.Len(cmd2.str, 3)
		ok = ok && test.Equal("2ndinit", cmd2.str[0])
		ok = ok && test.Equal("2ndstart", cmd2.str[1])
		ok = ok && test.Equal("2ndclear", cmd2.str[2])
	}

}

func TestStartWrongCommand(t *testing.T) {
	test := assert.New(t)

	cmd.Reset()
	{
		os.Args = []string{"cmd_test", "third"}

		cmd1 := FirstCMD{}
		cmd2 := SecondCMD{}

		cmd.RegisterCmd("first", &cmd1)
		cmd.RegisterCmd("second", &cmd2)

		err := cmd.Start()
		ok := test.NotNil(err)
		ok = ok && test.Equal("No matching command found", err.Error())
	}

	{
		os.Args = []string{"cmd_test", "third"}

		cmd1 := FirstCMD{}
		cmd2 := SecondCMD{}

		cmd.PanicEmptyCommand = true
		cmd.RegisterCmd("first", &cmd1)
		cmd.RegisterCmd("second", &cmd2)

		test.Panics(func() { cmd.Start() })
	}
}

type NoPtrCMD struct{}

func (c NoPtrCMD) Init() {}
func (c NoPtrCMD) Start() error {
	return nil
}
func (c NoPtrCMD) Clean() {}

func TestNoPointer(t *testing.T) {
	test := assert.New(t)

	os.Args = []string{"cmd_test", "first"}

	test.Panics(func() { cmd.RegisterCmd("first", NoPtrCMD{}) })
}
