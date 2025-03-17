package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleTitle(p *ansi.Parser) (seqInfo, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return seqNoMnemonic(""), errInvalid
	}
	switch p.Command() {
	case 0:
		return seqNoMnemonic(fmt.Sprintf("Set icon name and window title to %q", parts[1])), nil
	case 1:
		return seqNoMnemonic(fmt.Sprintf("Set icon name to %q", parts[1])), nil
	case 2:
		return seqNoMnemonic(fmt.Sprintf("Set window title to %q", parts[1])), nil
	}
	return seqNoMnemonic(""), errUnhandled
}
