package main

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/x/ansi"
)

var csiHandlers = map[int]handlerFn{
	'm': handleSgr,
	'c': printf("Request primary device attributes"),

	// kitty
	'u' | '?'<<markerShift: handleKitty,
	'u' | '>'<<markerShift: handleKitty,
	'u' | '<'<<markerShift: handleKitty,
	'u' | '='<<markerShift: handleKitty,

	// cursor
	'A':                      handleCursor,
	'B':                      handleCursor,
	'C':                      handleCursor,
	'D':                      handleCursor,
	'E':                      handleCursor,
	'F':                      handleCursor,
	'H':                      handleCursor,
	'n' | '?'<<markerShift:   handleCursor,
	'n':                      handleCursor,
	's':                      handleCursor,
	'u':                      handleCursor,
	'q' | ' '<<intermedShift: handleCursor,

	// screen
	'r': handleScreen,
	'J': handleScreen,
	'K': handleLine,
	'L': handleLine,
	'M': handleLine,
	'S': handleLine,
	'T': handleLine,

	// modes
	'p' | '$'<<intermedShift:                    handleMode,
	'p' | '?'<<markerShift | '$'<<intermedShift: handleMode,
	'h' | '?'<<markerShift:                      handleMode,
	'l' | '?'<<markerShift:                      handleMode,
	'h':                                         handleMode,
	'l':                                         handleMode,

	'q' | '>'<<markerShift: handleXT,
}

var oscHandlers = map[int]handlerFn{
	0:   handleTitle,
	1:   handleTitle,
	2:   handleTitle,
	7:   handleWorkingDirectoryURL,
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

var dcsHandlers = map[int]handlerFn{
	'q' | '+'<<intermedShift: handleTermcap,
}

var escHandler = map[int]handlerFn{
	'7': printf("Save cursor"),
	'8': printf("Restore cursor"),

	// C0/7-bit ASCII variant of ST.
	// C1/8-bit extended ASCII variant handled as Ctrl.
	'\\': printf("String terminator"),
}

var (
	errUnhandled = errors.New("TODO: unhandled sequence")
	errInvalid   = errors.New("invalid sequence")
)

type handlerFn = func(*ansi.Parser) (string, error)

func printf(format string, v ...any) handlerFn { //nolint:unparam
	return func(*ansi.Parser) (string, error) {
		return fmt.Sprintf(format, v...), nil
	}
}

func default1(i int) int {
	if i == 0 {
		return 1
	}
	return i
}
