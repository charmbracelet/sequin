package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

func handleMode(p *ansi.Parser) (seqInfo, error) {
	var m int
	if n, ok := p.Param(0, 0); ok {
		m = n
	}
	mode := modeDesc(m)
	cmd := ansi.Cmd(p.Command())
	private := ""
	if cmd.Prefix() == '?' {
		private = "private "
	}
	switch cmd.Final() {
	case 'p':
		// DECRQM - Request Mode
		return seqInfo{
			"DECRQM",
			fmt.Sprintf("Request %smode %q", private, mode),
		}, nil
	case 'h':
		var mnemonic string
		if cmd.Prefix() == '?' {
			mnemonic = "DECSET"
		} else {
			mnemonic = "SM"
		}
		return seqInfo{
			mnemonic,
			fmt.Sprintf("Enable %smode %q", private, mode),
		}, nil
	case 'l':
		var mnemonic string
		if cmd.Prefix() == '?' {
			mnemonic = "DECRST"
		} else {
			mnemonic = "RM"
		}
		return seqInfo{
			mnemonic,
			fmt.Sprintf("Disable %smode %q", private, mode),
		}, nil
	}
	return seqNoMnemonic(""), errUnhandled
}

//nolint:mnd
func modeDesc(mode int) string {
	switch mode {
	case 1:
		return "cursor keys"
	case 25:
		return "cursor visibility"
	case 1000:
		return "show mouse"
	case 1001:
		return "mouse hilite"
	case 1002:
		return "mouse cell motion"
	case 1003:
		return "mouse all motion"
	case 1004:
		return "report focus"
	case 1006:
		return "mouse SGR ext"
	case 1049:
		return "altscreen"
	case 2004:
		return "bracketed paste"
	case 2026:
		return "synchronized output"
	case 2027:
		return "grapheme clustering"
	case 9001:
		return "win32 input"
	default:
		return "unknown"
	}
}
