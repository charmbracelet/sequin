// https://github.com/gnachman/iterm2-website/blob/master/source/_includes/3.4/documentation-escape-codes.md#shell-integrationfinalterm
package main

import (
	"bytes"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleFinalTerm(p *ansi.Parser) (string, error) {
	parts := bytes.Split(p.Data(), []byte{';'})

	if len(parts) < 2 {
		return "", errInvalid
	}

	if len(parts[1]) != 1 {
		return "", errInvalid
	}

	var buf string
	switch parts[1][0] {
	case 'A':
		buf += "Prompt start"
		break
	case 'B':
		buf += "Command start"
		break
	case 'C':
		buf += "Command executed"
		break
	case 'D':
		buf += "Command finished"
		if len(parts) > 2 && len(parts[2]) > 1 {
			buf += fmt.Sprintf(", exit code: %s", string(parts[2]))
		}
		break
	default:
		return "", errInvalid
	}
	return buf, nil
}
