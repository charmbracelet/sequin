package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

// https://invisible-island.net/xterm/ctlseqs/ctlseqs.html#h4-Functions-using-CSI-_-ordered-by-the-final-character-lparen-s-rparen:CSI-gt-Pp;Pv-m.1EB3
func handleXTModKeys(p *ansi.Parser) (string, error) {
	var resource int
	var value int
	var hasValue bool

	if len(p.Params()) > 2 || len(p.Params()) == 0 {
		return "", errInvalid
	}

	// First parameter is the resource (Pp)
	if n, ok := p.Param(0, 0); ok {
		resource = n
	}

	// Second parameter is the value to set for the resource (Pv)
	if n, ok := p.Param(1, 0); ok {
		value = n
		hasValue = true
	}

	resourceName := xtModKeysResourceName(resource)

	if !hasValue {
		return fmt.Sprintf("Reset %s to initial value", resourceName), nil
	}

	if value == 0 {
		return fmt.Sprintf("Disable %s", resourceName), nil
	}

	return fmt.Sprintf("Enable %s", resourceName), nil
}

//nolint:mnd
func xtModKeysResourceName(resource int) string {
	switch resource {
	case 0:
		return "modifyKeyboard"
	case 1:
		return "modifyCursorKeys"
	case 2:
		return "modifyFunctionKeys"
	case 3:
		return "modifyKeypadKeys"
	case 4:
		return "modifyOtherKeys"
	case 6:
		return "modifyModifierKeys"
	case 7:
		return "modifySpecialKeys"
	default:
		return fmt.Sprintf("unknown modifier resource %d", resource)
	}
}
