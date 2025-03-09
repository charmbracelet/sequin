package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handlePointerShape(p *ansi.Parser) (seqInfo, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return seqNoMnemonic(""), errInvalid
	}

	return seqNoMnemonic(fmt.Sprintf("Set pointer shape to %q", parts[1])), nil
}
