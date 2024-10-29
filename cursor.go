package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleCursor(cmd int) func(*ansi.Parser) {
	return func(parser *ansi.Parser) {
		var count int
		if parser.ParamsLen > 0 {
			count = ansi.Param(parser.Params[0]).Param()
		}

		if count == 0 {
			// Default value is 1
			count = 1
		}

		switch cmd {
		case 'A':
			// CUU - Cursor Up
			fmt.Printf("Cursor up %d", count)
		case 'B':
			// CUD - Cursor Down
			fmt.Printf("Cursor down %d", count)
		case 'C':
			// CUF - Cursor Forward
			fmt.Printf("Cursor right %d", count)
		case 'D':
			// CUB - Cursor Back
			fmt.Printf("Cursor left %d", count)
		}
	}
}
