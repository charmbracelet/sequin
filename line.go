package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleLine(p *ansi.Parser) (seqInfo, error) {
	var count int
	if n, ok := p.Param(0, 0); ok {
		count = n
	}

	switch p.Command() {
	case 'K':
		mnemonic := "EL"
		switch count {
		case 0:
			return seqInfo{mnemonic, "Erase line right"}, nil
		case 1:
			return seqInfo{mnemonic, "Erase line left"}, nil
		case 2:
			return seqInfo{mnemonic, "Erase entire line"}, nil
		}
	case 'L':
		return seqInfo{"IL", fmt.Sprintf("Insert %d blank lines", default1(count))}, nil
	case 'M':
		return seqInfo{"DL", fmt.Sprintf("Delete %d lines", default1(count))}, nil
	case 'S':
		return seqInfo{"SU", fmt.Sprintf("Scroll up %d lines", default1(count))}, nil
	case 'T':
		return seqInfo{"SD", fmt.Sprintf("Scroll down %d lines", default1(count))}, nil
	}
	return seqNoMnemonic(""), errUnhandled
}
