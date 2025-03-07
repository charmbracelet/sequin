package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleHyperlink(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 3 {
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
