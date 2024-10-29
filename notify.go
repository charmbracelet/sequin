package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleNotify(p *ansi.Parser) {
	parts := bytes.Split(p.Data, []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return
	}

	fmt.Printf("Notify %q", parts[1])
}
