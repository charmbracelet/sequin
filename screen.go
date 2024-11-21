package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleScreen(p *ansi.Parser) (string, error) {
	var count int
	if n, ok := p.Param(0, 0); ok {
		count = n
	}

	cmd := p.Cmd()
	switch cmd.Command() {
	case 'J':
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
		top, bot := 1, 0
		if n, ok := p.Param(0, 1); ok {
			top = n
		}
		if n, ok := p.Param(1, 0); ok {
			bot = n
		}
		return fmt.Sprintf(
			"Set scrolling region to top=%d bottom=%d",
			top,
			bot,
		), nil
	}

	return "", errUnhandled
}
