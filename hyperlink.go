package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleHyperlink(p *ansi.Parser) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 3 {
		// Invalid, ignore
		return
	}

	opts := bytes.Split(parts[1], []byte{':'})
	fmt.Printf("Set hyperlink, ")
	for i, opt := range opts {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", opt)
	}

	fmt.Printf(" to %q", parts[2])
}
