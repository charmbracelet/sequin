package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleHyperlink(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 3 { //nolint:mnd
		// Invalid, ignore
		return "", errInvalid
	}

	opts := bytes.Split(parts[1], []byte{':'})
	buf := "Set hyperlink, "
	for i, opt := range opts {
		if i > 0 {
			buf += ", "
		}
		buf += string(opt)
	}

	buf += fmt.Sprintf(" to %q", parts[2])
	return buf, nil
}
