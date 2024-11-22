package main

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleWorkingDirectoryURL(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data(), []byte{';'})
	if len(parts) != 2 {
		// Invalid, ignore
		return "", errInvalid
	}

	u, err := url.ParseRequestURI(string(parts[1]))

	if err != nil || u.Scheme != "file" {
		// Should be a file URL.
		return "", errInvalid
	}

	return fmt.Sprintf("Set working directory to %s (on %s)", u.Path, u.Host), nil
}
