package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleLine(parser *ansi.Parser) (string, error) {
	var count int
	if parser.ParamsLen > 0 {
		count = ansi.Param(parser.Params[0]).Param()
	}

	if count == 0 {
		// Default value is 1
		count = 1
	}

	cmd := ansi.Cmd(parser.Cmd)
	switch cmd.Command() {
	case 'K':
		switch count {
		case 1:
			return "Erase line right", nil
		case 2:
			return "Erase line left", nil
		case 3:
			return "Erase entire line", nil
		}
	case 'L':
		return fmt.Sprintf("CSI %d L: Insert %[1]d blank lines", count), nil
	case 'M':
		return fmt.Sprintf("CSI %d M: Delete %[1]d lines", count), nil
	case 'S':
		return fmt.Sprintf("CSI %d S: Scroll up %[1]d lines", count), nil
	case 'T':
		return fmt.Sprintf("CSI %d T: Scroll down %[1]d lines", count), nil
	}
	return "", errUnhandled
}
