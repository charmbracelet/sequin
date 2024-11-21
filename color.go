package main

import (
	"bytes"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleTerminalColor(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return "", errInvalid
	}

	arg := string(parts[1])

	var buf string
	if arg == "?" {
		buf += "Request"
	} else {
		buf += "Set"
	}
	switch p.Cmd() {
	case 10:
		buf += " foreground color"
	case 11:
		buf += " background color"
	case 12:
		buf += " cursor color"
	}
	if arg == "?" {
		buf += " to " + arg
	}
	return buf, nil
}

//nolint:mnd
func handleResetTerminalColor(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 1 {
		// Invalid, ignore
		return "", errInvalid
	}
	var buf string
	switch p.Cmd() {
	case 110:
		buf += "Reset foreground color"
	case 111:
		buf += "Reset background color"
	case 112:
		buf += "Reset cursor color"
	}
	return buf, nil
}
