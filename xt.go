package main

import "github.com/charmbracelet/x/ansi"

func handleXT(parser *ansi.Parser) (string, error) {
	var count int
	if parser.ParamsLen > 0 {
		count = ansi.Param(parser.Params[0]).Param()
	}

	if count != 0 {
		return "", errUnhandled
	}

	return "Request XT Version", nil
}
