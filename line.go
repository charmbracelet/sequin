package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleLine(p *ansi.Parser) (string, error) {
	var count int
	if n, ok := p.Param(0, 0); ok {
		count = n
	}

	switch p.Command() {
	case 'K':
		switch count {
		case 0:
			return "Erase line right", nil
		case 1:
			return "Erase line left", nil
		case 2:
			return "Erase entire line", nil
		}
	case 'L':
		return fmt.Sprintf("Insert %d blank lines", default1(count)), nil
	case 'M':
		return fmt.Sprintf("Delete %d lines", default1(count)), nil
	case 'S':
		return fmt.Sprintf("Scroll up %d lines", default1(count)), nil
	case 'T':
		return fmt.Sprintf("Scroll down %d lines", default1(count)), nil
	}
	return "", errUnhandled
}
