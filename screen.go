package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleScreen(p *ansi.Parser) (seqInfo, error) {
	var count int
	if n, ok := p.Param(0, 0); ok {
		count = n
	}

	cmd := ansi.Cmd(p.Command())
	switch cmd.Final() {
	case 'J':
		mnemonic := "ED"
		switch count {
		case 0:
			return seqInfo{mnemonic, "Erase screen bellow"}, nil
		case 1:
			return seqInfo{mnemonic, "Erase screen above"}, nil
		case 2:
			return seqInfo{mnemonic, "Erase entire screen"}, nil
		case 3:
			return seqInfo{mnemonic, "Erase entire display"}, nil
		}
	case 'r':
		top, bot := 1, 0
		if n, ok := p.Param(0, 1); ok {
			top = n
		}
		if n, ok := p.Param(1, 0); ok {
			bot = n
		}
		return seqInfo{"DECSTBM", fmt.Sprintf(
			"Set scrolling region to top=%d bottom=%d",
			top,
			bot,
		)}, nil
	}

	return seqNoMnemonic(""), errUnhandled
}
