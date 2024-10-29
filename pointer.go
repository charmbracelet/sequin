package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handlePointerShape(p *ansi.Parser) {
	parts := bytes.Split(p.Data, []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return
	}

	fmt.Printf("Set pointer shape to %q", parts[1])
}
