package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleMode(parser *ansi.Parser) (string, error) {
	mode := modeDesc(ansi.Param(parser.Params[0]).Param())
	cmd := ansi.Cmd(parser.Cmd)
	private := ""
	if cmd.Marker() == '?' {
		private = "private "
	}
	switch cmd.Command() {
	case 'p':
		// DECRQM - Request Mode
		return fmt.Sprintf("Request %smode %q", private, mode), nil
	case 'h':
		return fmt.Sprintf("Enable %smode %q", private, mode), nil
	case 'l':
		return fmt.Sprintf("Disable %smode %q", private, mode), nil
	}
	return "", errUnhandled
}

func modeDesc(mode int) string {
	switch mode {
	case 1:
		return "cursor keys"
	case 25:
		return "cursor visibility"
	case 1000:
		return "show mouse"
	case 1001:
		return "mouse hilite"
	case 1002:
		return "mouse cell motion"
	case 1003:
		return "mouse all motion"
	case 1004:
		return "report focus"
	case 1006:
		return "mouse SGR ext"
	case 1049:
		return "altscreen"
	case 2004:
		return "bracketed paste"
	case 2026:
		return "synchronized output"
	case 2027:
		return "grapheme clustering"
	case 9001:
		return "win32 input"
	default:
		return "unknown"
	}
}
