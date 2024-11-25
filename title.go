package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleTitle(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return "", errInvalid
	}
	switch p.Cmd() {
	case 0:
		return fmt.Sprintf("Set icon name and window title to %q", parts[1]), nil
	case 1:
		return fmt.Sprintf("Set icon name to %q", parts[1]), nil
	case 2:
		return fmt.Sprintf("Set window title to %q", parts[1]), nil
	}
	return "", errUnhandled
}
