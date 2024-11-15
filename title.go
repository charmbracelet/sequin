package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleTitle(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 2 { //nolint:mnd
		// Invalid, ignore
		return "", errInvalid
	}
	switch p.Cmd {
	case 0:
		return fmt.Sprintf("Set icon name and window title to %s", parts[1]), nil
	case 1:
		return fmt.Sprintf("Set icon name to %s", parts[1]), nil
	case 2: //nolint:mnd
		return fmt.Sprintf("Set window title to %s", parts[1]), nil
	}
	return "", errUnhandled
}
