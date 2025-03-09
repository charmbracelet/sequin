package main

import "github.com/charmbracelet/x/ansi"

func handleXT(p *ansi.Parser) (seqInfo, error) {
	var count int
	if n, ok := p.Param(0, 0); ok {
		count = n
	}

	if count != 0 {
		return seqNoMnemonic(""), errInvalid
	}

	return seqInfo{"XTVERSION", "Request XT Version"}, nil
}
