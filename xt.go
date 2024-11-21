package main

import "github.com/charmbracelet/x/ansi"

func handleXT(p *ansi.Parser) (string, error) {
	var count int
	if n, ok := p.Param(0, 0); ok {
		count = n
	}

	if count != 0 {
		return "", errInvalid
	}

	return "Request XT Version", nil
}
