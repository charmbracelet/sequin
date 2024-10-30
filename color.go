package main

import (
	"bytes"

	"github.com/charmbracelet/x/ansi"
)

func handleTerminalColor(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return "", errUnknown
	}

	var buf string
	if string(parts[1]) == "?" {
		buf += "Request"
	} else {
		buf += "Set"
	}
	switch p.Cmd {
	case 10:
		buf += " foreground color to " + string(parts[1])
	case 11:
		buf += " background color to " + string(parts[1])
	case 12:
		buf += " cursor color to " + string(parts[1])
	}
	return buf, nil
}

func handleResetTerminalColor(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 1 {
		// Invalid, ignore
		return "", errUnknown
	}
	var buf string
	switch p.Cmd {
	case 110:
		buf += "Reset foreground color"
	case 111:
		buf += "Reset background color"
	case 112:
		buf += "Reset cursor color"
	}
	return buf, nil
}
