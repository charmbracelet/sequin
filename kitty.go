package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleKitty(p *ansi.Parser) (seqInfo, error) {
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
	if n, ok := p.Param(0, 0); ok {
		first = n
	}

	cmd := ansi.Cmd(p.Command())
	switch cmd.Prefix() {
	case '?':
		return seqNoMnemonic("Request Kitty keyboard"), nil
	case '>':
		if first == 0 {
			return seqNoMnemonic("Disable Kitty keyboard"), nil
		}
		return seqNoMnemonic(fmt.Sprintf("Push %q Kitty keyboard flag", flagDesc(first))), nil
	case '<':
		return seqNoMnemonic(fmt.Sprintf("Pop %d Kitty keyboard flags", first)), nil
	case '=':
		if n, ok := p.Param(1, 0); ok {
			return seqNoMnemonic(fmt.Sprintf("Set %q Kitty keyboard flags to %q", flagDesc(first), modeDesc(n))), nil
		}
	}
	return seqNoMnemonic(""), errUnhandled
}
