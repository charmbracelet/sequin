package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleTerminalColor(p *ansi.Parser) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return
	}

	if string(parts[1]) == "?" {
		fmt.Print("Request")
	} else {
		fmt.Print("Set")
	}
	switch p.Cmd {
	case 10:
		fmt.Printf(" foreground color to %s", parts[1])
	case 11:
		fmt.Printf(" background color to %s", parts[1])
	case 12:
		fmt.Printf(" cursor color to %s", parts[1])
	}
}

func handleResetTerminalColor(p *ansi.Parser) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 1 {
		// Invalid, ignore
		return
	}

	switch p.Cmd {
	case 110:
		fmt.Print("Reset foreground color")
	case 111:
		fmt.Print("Reset background color")
	case 112:
		fmt.Print("Reset cursor color")
	}
}
