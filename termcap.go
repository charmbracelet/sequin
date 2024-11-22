package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func handleTermcap(p *ansi.Parser) (string, error) {
	data := p.Data()
	if len(data) == 0 {
		return "", errInvalid
	}

	parts := bytes.Split(data, []byte{';'})
	if len(parts) == 0 {
		return "", errInvalid
	}

	caps := make([]string, 0, len(parts))
	for _, part := range parts {
		capName, err := hex.DecodeString(string(part))
		if err != nil {
			//nolint:wrapcheck
			return "", err
		}
		caps = append(caps, string(capName))
	}

	return fmt.Sprintf("Request termcap entry for %s", strings.Join(caps, ", ")), nil
}
