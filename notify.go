package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleNotify(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return "", errInvalid
	}

	return fmt.Sprintf("Notify %q", parts[1]), nil
}
