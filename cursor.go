package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleCursor(p *ansi.Parser) (string, error) {
	var count int
	if p.ParamsLen > 0 {
		count = ansi.Param(p.Params[0]).Param()
	}

	if count == 0 {
		// Default value is 1
		count = 1
	}

	cmd := ansi.Cmd(p.Cmd)
	isPrivate := cmd.Marker() == '?'
	switch cmd.Command() {
	case 'A':
		// CUU - Cursor Up
		return fmt.Sprintf("Cursor up %d", count), nil
	case 'B':
		// CUD - Cursor Down
		return fmt.Sprintf("Cursor down %d", count), nil
	case 'C':
		// CUF - Cursor Forward
		return fmt.Sprintf("Cursor right %d", count), nil
	case 'D':
		// CUB - Cursor Back
		return fmt.Sprintf("Cursor left %d", count), nil
	case 'E':
		return fmt.Sprintf("Cursor next line %d", count), nil
	case 'F':
		return fmt.Sprintf("Cursor previous line %d", count), nil
	case 'H':
		row := count
		col := 1
		if p.ParamsLen > 1 {
			col = p.Params[1]
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
	}
	return "", errUnhandled
}
