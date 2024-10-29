package main

import (
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

var csiHandlers = map[int]func(*ansi.Parser){
	'm':                    handleSgr,
	'c':                    printf("Request primary device attributes"),
	'q' | '>'<<markerShift: handleXT,

	// kitty
	'u' | '?'<<markerShift: handleKitty,
	'u' | '>'<<markerShift: handleKitty,
	'u' | '<'<<markerShift: handleKitty,
	'u' | '='<<markerShift: handleKitty,

	// cursor
	'A':                    handleCursor,
	'B':                    handleCursor,
	'C':                    handleCursor,
	'D':                    handleCursor,
	'E':                    handleCursor,
	'F':                    handleCursor,
	'H':                    handleCursor,
	'n' | '?'<<markerShift: handleCursor,
	'n':                    handleCursor,
	's':                    handleCursor,

	// screen
	'J': handleScreen,
	'r': handleScreen,
	'K': handleLine,
	'L': handleLine,
	'M': handleLine,
	'S': handleLine,
	'T': handleLine,

	// modes
	'p' | '$'<<intermedShift:                    handleReqMode,
	'p' | '?'<<markerShift | '$'<<intermedShift: handleReqMode,
	'h' | '?'<<markerShift:                      handleReqMode,
	'l' | '?'<<markerShift:                      handleReqMode,
	'h':                                         handleReqMode,
	'l':                                         handleReqMode,
}

var oscHandlers = map[int]func(*ansi.Parser){
	0:   handleTitle,
	1:   handleTitle,
	2:   handleTitle,
	8:   handleHyperlink,
	9:   handleNotify,
	10:  handleTerminalColor,
	11:  handleTerminalColor,
	12:  handleTerminalColor,
	22:  handlePointerShape,
	52:  handleClipboard,
	110: handleResetTerminalColor,
	111: handleResetTerminalColor,
	112: handleResetTerminalColor,
}

func printf(s string) func(*ansi.Parser) {
	return func(*ansi.Parser) {
		fmt.Printf(s)
	}
}
