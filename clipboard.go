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

func handleClipboard(p *ansi.Parser) {
	parts := bytes.Split(p.Data[:p.DataLen], []byte{';'})
	if len(parts) != 3 {
		// Invalid, ignore
		return
	}

	b64, err := base64.StdEncoding.DecodeString(string(parts[2]))
	if err != nil {
		// Invalid, ignore
		return
	}

	fmt.Printf("Set clipboard %q to %q", clipboardName[string(parts[1])], b64)
}
