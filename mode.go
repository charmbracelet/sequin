package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleReqMode(isPrivate bool) func(*ansi.Parser) {
	// DECRQM - Request Mode
	return func(parser *ansi.Parser) {
		if parser.ParamsLen == 0 {
			// Invalid, ignore
			return
		}

		mode := ansi.Param(parser.Params[0]).Param()
		if isPrivate {
			fmt.Printf("Request private mode %d", mode)
		} else {
			fmt.Printf("Request mode %d", mode)
		}
	}
}
