package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/stretchr/testify/require"
)

var cursor = map[string]string{
	// cursor
	"save":               ansi.SaveCursor,
	"restore":            ansi.RestoreCursor,
	"request pos":        ansi.RequestCursorPosition,
	"request cursor pos": ansi.RequestExtendedCursorPosition,
	"up 1":               ansi.CursorUp1,
	"up":                 ansi.CursorUp(5),
	"down 1":             ansi.CursorDown1,
	"down":               ansi.CursorDown(3),
	"right 1":            ansi.CursorRight1,
	"right":              ansi.CursorRight(3),
	"left 1":             ansi.CursorLeft1,
	"left":               ansi.CursorLeft(3),
	"next line":          ansi.CursorNextLine(3),
	"previous line":      ansi.CursorPreviousLine(3),
	"set pos":            ansi.SetCursorPosition(10, 20),
	"origin":             ansi.CursorOrigin,
	"save pos":           ansi.SaveCursorPosition,
	"restore pos":        ansi.RestoreCursorPosition,
	"style":              ansi.SetCursorStyle(4), // TODO: bug in ansi
	"pointer shape":      ansi.SetPointerShape("crosshair"),
}

// "reset": ansi.ResetStyle,

func TestSequences(t *testing.T) {
	for name, table := range map[string]map[string]string{
		"cursor": cursor,
	} {
		t.Run(name, func(t *testing.T) {
			for name, input := range table {
				t.Run(name, func(t *testing.T) {
					var b bytes.Buffer
					cmd := cmd()
					cmd.SetOut(&b)
					cmd.SetErr(&b)
					cmd.SetIn(strings.NewReader(input))
					cmd.SetArgs([]string{})
					require.NoError(t, cmd.Execute())
					golden.RequireEqual(t, b.Bytes())
				})
			}
		})
	}
}
