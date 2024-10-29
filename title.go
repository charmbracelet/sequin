package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleTitle(p *ansi.Parser) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return
	}

	switch p.Cmd {
	case 0:
		fmt.Printf("Set icon name and window title to %s", parts[1])
	case 1:
		fmt.Printf("Set icon name to %s", parts[1])
	case 2:
		fmt.Printf("Set window title to %s", parts[1])
	}
}
