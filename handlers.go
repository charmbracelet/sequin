package main

import "github.com/charmbracelet/x/ansi"

var csiHandlers = map[int]func(*ansi.Parser){
	'm':                      handleSgr,
	'A':                      handleCursor('A'),
	'B':                      handleCursor('B'),
	'C':                      handleCursor('C'),
	'D':                      handleCursor('D'),
	'p' | '$'<<intermedShift: handleReqMode(false),
	'p' | '?'<<markerShift | '$'<<intermedShift: handleReqMode(true),
}

var oscHandlers = map[int]func(*ansi.Parser){
	0:  handleTitle,
	1:  handleTitle,
	2:  handleTitle,
	8:  handleHyperlink,
	9:  handleNotify,
	22: handlePointerShape,
	52: handleClipboard,
}
