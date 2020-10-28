package meta

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Cmd executes a shell command.
// It returns stdout and either a execution error or if none, stderr.
func Cmd(cmd string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command := exec.Command("/bin/sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if err != nil {
		return "", err
	}
	if stderr.Len() > 0 {
		err = fmt.Errorf(stderr.String())
	}

	return strings.TrimSpace(stdout.String()), err
}
