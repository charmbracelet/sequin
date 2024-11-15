package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleScreen(parser *ansi.Parser) (string, error) {
	var count int
	if parser.ParamsLen > 0 {
		count = ansi.Param(parser.Params[0]).Param()
	}

	cmd := ansi.Cmd(parser.Cmd)
	switch cmd.Command() {
	case 'J':
		//nolint:mnd
		switch count {
		case 0:
			return "Erase screen bellow", nil
		case 1:
			return "Erase screen above", nil
		case 2:
			return "Erase entire screen", nil
		case 3:
			return "Erase entire display", nil
		}
	case 'r':
		top := count
		bottom := 0
		if parser.ParamsLen > 1 {
			bottom = ansi.Param(parser.Params[1]).Param()
		}
		return fmt.Sprintf(
			"Set scrolling region to top=%d bottom=%d",
			top,
			bottom,
		), nil
	}

	return "", errUnhandled
}
