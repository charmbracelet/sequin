package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func handleTermcap(p *ansi.Parser) (string, error) {
	if p.ParamsLen != 0 || p.DataLen == 0 {
		return "", errInvalid
	}

	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})

	var caps []string
	for _, part := range parts {
		capName, err := hex.DecodeString(string(part))
		if err != nil {
			return "", err
		}
		caps = append(caps, string(capName))
	}

	return fmt.Sprintf("Request termcap entry for %s", strings.Join(caps, ", ")), nil
}
