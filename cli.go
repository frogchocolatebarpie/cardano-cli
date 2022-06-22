package main

import (
	"bytes"
	"os"
	"os/exec"
)

type Cli struct {
	Path string
}

func (c *Cli) Exec(args ...string) ([]byte, error) {
	out := &bytes.Buffer{}

	cmd := exec.Command(c.Path, args...)
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
