package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func handleMode(p *ansi.Parser) (string, error) {
	params := p.Params()
	if len(params) == 0 {
		return "", errInvalid
	}

	cmd := ansi.Cmd(p.Command())
	private := ""
	if cmd.Prefix() == '?' {
		private = "private "
	}

	var action string
	switch cmd.Final() {
	case 'p':
		// DECRQM - Request Mode
		action = "Request"
	case 'h':
		action = "Enable"
	case 'l':
		action = "Disable"
	default:
		return "", errUnhandled
	}

	var modes []string
	for _, param := range params {
		m := param.Param(0)
		modes = append(modes, fmt.Sprintf("%q", modeDesc(m)))
	}

	return fmt.Sprintf("%s %smode %s", action, private, strings.Join(modes, ", ")), nil
}

//nolint:mnd
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
	case 1015:
		return "mouse URXVT ext"
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
		return unknown
	}
}
