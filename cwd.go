package main

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleWorkingDirectoryURL(p *ansi.Parser) (seqInfo, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return seqNoMnemonic(""), errInvalid
	}

	u, err := url.ParseRequestURI(string(parts[1]))

	if err != nil || u.Scheme != "file" {
		// Should be a file URL.
		return seqNoMnemonic(""), errInvalid
	}

	return seqNoMnemonic(fmt.Sprintf("Set working directory to %s (on %s)", u.Path, u.Host)), nil
}
