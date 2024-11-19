package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
)

func executeCommand(ctx context.Context, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint: gosec
	pty, err := runInPty(cmd)
	if err != nil {
		return nil, err
	}
	defer pty.Close() //nolint: errcheck
	var out bytes.Buffer
	var errorOut bytes.Buffer
	go func() {
		_, _ = io.Copy(&out, pty)
		errorOut.Write(out.Bytes())
	}()

	err = cmd.Wait()
	if err != nil {
		return errorOut.Bytes(), err //nolint: wrapcheck
	}
	return out.Bytes(), nil
}

// runInPty opens a new pty and runs the given command in it.
// The returned file is the pty's file descriptor and must be closed by the
// caller.
func runInPty(c *exec.Cmd) (*os.File, error) {
	//nolint: wrapcheck
	return pty.StartWithAttrs(c, &pty.Winsize{
		Cols: 80,
		Rows: 10,
	}, &syscall.SysProcAttr{})
}
