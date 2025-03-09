package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleNotify(p *ansi.Parser) (seqInfo, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return seqNoMnemonic(""), errInvalid
	}

	return seqNoMnemonic(fmt.Sprintf("Notify %q", parts[1])), nil
}
