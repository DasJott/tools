package meta_test

import (
	"testing"

	"cleverreach.com/crtools/meta"
	"github.com/stretchr/testify/assert"
)

func TestCmd(t *testing.T) {
	test := assert.New(t)

	res, err := meta.Cmd("echo 'foo bar'")
	test.Nil(err)
	test.Equal("foo bar", res)

	msg := ""
	res, err = meta.Cmd("ls cmd*")
	if err != nil {
		msg = err.Error()
	}
	test.Nil(err, msg)
	test.Equal("cmd.go\ncmd_test.go", res)
}
