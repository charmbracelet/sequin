package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleCursor(p *ansi.Parser) (string, error) {
	count := 1
	if n, ok := p.Param(0, 1); ok && n > 0 {
		count = n
	}

	cmd := p.Cmd()
	isPrivate := cmd.Marker() == '?'
	switch cmd.Command() {
	case 'A':
		// CUU - Cursor Up
		return fmt.Sprintf("Cursor up %d", default1(count)), nil
	case 'B':
		// CUD - Cursor Down
		return fmt.Sprintf("Cursor down %d", default1(count)), nil
	case 'C':
		// CUF - Cursor Forward
		return fmt.Sprintf("Cursor right %d", default1(count)), nil
	case 'D':
		// CUB - Cursor Back
		return fmt.Sprintf("Cursor left %d", default1(count)), nil
	case 'E':
		return fmt.Sprintf("Cursor next line %d", default1(count)), nil
	case 'F':
		return fmt.Sprintf("Cursor previous line %d", default1(count)), nil
	case 'H':
		row, col := 1, 1
		if n, ok := p.Param(0, 1); ok && n > 0 {
			row = n
		}
		if n, ok := p.Param(1, 1); ok && n > 0 {
			col = n
		}
		return fmt.Sprintf("Set cursor position row=%[1]d col=%[2]d", row, col), nil
	case 'n':
		if count != 6 {
			return "", errInvalid
		}
		if isPrivate {
			return "Request extended cursor position", nil
		}
		return "Request cursor position", nil
	case 's':
		return "Save cursor position", nil
	case 'u':
		return "Restore cursor position", nil
	case 'q':
		return fmt.Sprintf("Set cursor style %s", descCursorStyle(count)), nil
	}
	return "", errUnhandled
}

//nolint:mnd
func descCursorStyle(i int) string {
	switch i {
	case 0, 1:
		return "Blinking block"
	case 2:
		return "Steady block"
	case 3:
		return "Blinking underline"
	case 4:
		return "Steady underline"
	case 5:
		return "Blinking bar"
	case 6:
		return "Steady bar"
	default:
		return "Unknown"
	}
}
