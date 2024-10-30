package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleMode(parser *ansi.Parser) (string, error) {
	if parser.ParamsLen == 0 {
		// Invalid, ignore
		return "", errInvalid
	}

	mode := modeDesc(ansi.Param(parser.Params[0]).Param())
	cmd := ansi.Cmd(parser.Cmd)
	isPrivate := cmd.Marker() == '?'
	switch cmd.Command() {
	case 'p':
		// DECRQM - Request Mode
		if isPrivate {
			return fmt.Sprintf("Request private mode %q", mode), nil
		}
		return fmt.Sprintf("Request mode %q", mode), nil
	case 'h':
		if isPrivate {
			return fmt.Sprintf("Enable private mode %q", mode), nil
		}
		return fmt.Sprintf("Enable mode %q", mode), nil
	case 'l':
		if isPrivate {
			return fmt.Sprintf("Disable private mode %q", mode), nil
		}
		return fmt.Sprintf("Disable mode %q", mode), nil
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
