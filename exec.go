package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/x/term"
	"github.com/charmbracelet/x/xpty"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

func executeCommand(ctx context.Context, args []string) ([]byte, error) {
	width, height, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		width = defaultWidth
		height = defaultHeight
	}

	pty, err := xpty.NewPty(width, height)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = pty.Close()
	}()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...) //nolint: gosec
	if err := pty.Start(cmd); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	var errorOut bytes.Buffer
	go func() {
		_, _ = io.Copy(&out, pty)
		errorOut.Write(out.Bytes())
	}()

	if err := xpty.WaitProcess(ctx, cmd); err != nil {
		return errorOut.Bytes(), err //nolint: wrapcheck
	}
	return out.Bytes(), nil
}
