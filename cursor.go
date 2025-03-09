package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleCursor(p *ansi.Parser) (seqInfo, error) {
	count := 1
	if n, ok := p.Param(0, 1); ok && n > 0 {
		count = n
	}

	cmd := ansi.Cmd(p.Command())
	isPrivate := cmd.Prefix() == '?'
	switch cmd.Final() {
	case 'A':
		// CUU - Cursor Up
		return seqInfo{
			"CUU",
			fmt.Sprintf("Cursor up %d", default1(count)),
		}, nil
	case 'B':
		// CUD - Cursor Down
		return seqInfo{
			"CUD",
			fmt.Sprintf("Cursor down %d", default1(count)),
		}, nil
	case 'C':
		// CUF - Cursor Forward
		return seqInfo{
			"CUF",
			fmt.Sprintf("Cursor right %d", default1(count)),
		}, nil
	case 'D':
		// CUB - Cursor Back
		return seqInfo{
			"CUB",
			fmt.Sprintf("Cursor left %d", default1(count)),
		}, nil
	case 'E':
		return seqInfo{
			"CNL",
			fmt.Sprintf("Cursor next line %d", default1(count)),
		}, nil
	case 'F':
		return seqInfo{
			"CPL",
			fmt.Sprintf("Cursor previous line %d", default1(count)),
		}, nil
	case 'H':
		row, col := 1, 1
		if n, ok := p.Param(0, 1); ok && n > 0 {
			row = n
		}
		if n, ok := p.Param(1, 1); ok && n > 0 {
			col = n
		}
		return seqInfo{
			"CUP",
			fmt.Sprintf("Set cursor position row=%[1]d col=%[2]d", row, col),
		}, nil
	case 'n':
		if count != 6 {
			return seqNoMnemonic(""), errInvalid
		}
		if isPrivate {
			return seqInfo{
				"DECXCPR",
				"Request extended cursor position",
			}, nil
		}
		return seqInfo{"CPR", "Request cursor position"}, nil
	case 's':
		return seqInfo{"SCOSC", "Save cursor position"}, nil
	case 'u':
		return seqInfo{"SCORC", "Restore cursor position"}, nil
	case 'q':
		return seqInfo{
			"DECSCUSR",
			fmt.Sprintf("Set cursor style %s", descCursorStyle(count)),
		}, nil
	}
	return seqNoMnemonic(""), errUnhandled
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
