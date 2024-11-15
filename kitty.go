package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

func handleKitty(parser *ansi.Parser) (string, error) {
	flagDesc := func(flag int) string {
		var r []string
		if flag&1 != 0 {
			r = append(r, "Disambiguate escape codes")
		}
		if flag&2 != 0 {
			r = append(r, "Report event types")
		}
		if flag&4 != 0 {
			r = append(r, "Report alternate keys")
		}
		if flag&8 != 0 {
			r = append(r, "Report all keys as escape codes")
		}
		if flag&16 != 0 {
			r = append(r, "Report associated text")
		}
		return strings.Join(r, ", ")
	}
	modeDesc := func(mode int) string {
		//nolint:mnd
		switch mode {
		case 1:
			return "Set given flags and unset all others"
		case 2:
			return "Set given flags and keep existing flags unchanged"
		case 3:
			return "Unset given flags and keep existing flags unchanged"
		default:
			return "Unknown mode"
		}
	}

	var first int
	if parser.ParamsLen > 0 {
		first = ansi.Param(parser.Params[0]).Param()
	}

	cmd := ansi.Cmd(parser.Cmd)
	switch cmd.Marker() {
	case '?':
		return "Request Kitty keyboard", nil
	case '>':
		if first == 0 {
			return "Disable Kitty keyboard", nil
		}
		return fmt.Sprintf("Push %q Kitty keyboard flag", flagDesc(first)), nil
	case '<':
		return fmt.Sprintf("Pop %d Kitty keyboard flags", first), nil
	case '=':
		if parser.ParamsLen > 1 {
			second := ansi.Param(parser.Params[1]).Param()
			return fmt.Sprintf("Set %q Kitty keyboard flags to %q", flagDesc(first), modeDesc(second)), nil
		}
	}
	return "", errUnhandled
}
