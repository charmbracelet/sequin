package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handlePointerShape(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return "", errInvalid
	}

	return fmt.Sprintf("Set pointer shape to %q", parts[1]), nil
}
