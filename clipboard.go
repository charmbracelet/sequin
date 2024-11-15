package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

var clipboardName = map[string]string{
	"c": "system",
	"p": "primary",
}

//nolint:mnd
func handleClipboard(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 3 {
		// Invalid, ignore
		return "", errInvalid
	}

	if string(parts[2]) == "?" {
		return fmt.Sprintf("Request %q clipboard", clipboardName[string(parts[1])]), nil
	}

	b64, err := base64.StdEncoding.DecodeString(string(parts[2]))
	if err != nil {
		// Invalid, ignore
		//nolint:wrapcheck
		return "", err
	}

	return fmt.Sprintf("Set clipboard %q to %q", clipboardName[string(parts[1])], b64), nil
}
