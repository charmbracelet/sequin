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

var screen = map[string]string{
	"enable alt buffer":  ansi.DisableAltScreenBuffer,
	"disable alt buffer": ansi.DisableAltScreenBuffer,
	"request alt buffer": ansi.DisableAltScreenBuffer,
	"passthrough":        ansi.ScreenPassthrough(ansi.SaveCursor, 0), // TODO: impl
	"erase above":        ansi.EraseScreenAbove,
	"erase below":        ansi.EraseScreenBelow,
	"erase full":         ansi.EraseEntireScreen,
	"erase display":      ansi.EraseEntireDisplay,
	"scrolling region":   ansi.SetScrollingRegion(10, 20),
}

var line = map[string]string{
	"right":       ansi.EraseLineRight,
	"left":        ansi.EraseLineLeft,
	"entire":      ansi.EraseEntireLine,
	"insert":      ansi.InsertLine(3),
	"delete":      ansi.DeleteLine(5),
	"scroll up":   ansi.ScrollUp(12),
	"scroll down": ansi.ScrollDown(12),
}

var mode = map[string]string{
	"enable cursor keys":          ansi.EnableCursorKeys,
	"disable cursor keys":         ansi.DisableCursorKeys,
	"request cursor keys":         ansi.RequestCursorKeys,
	"enable cursor visibility":    ansi.ShowCursor,
	"disable cursor visibility":   ansi.HideCursor,
	"request cursor visibility":   ansi.RequestCursorVisibility,
	"enable mouse hilite":         ansi.EnableMouseHilite,
	"disable mouse hilite":        ansi.DisableMouseHilite,
	"request mouse hilite":        ansi.RequestMouseHilite,
	"enable mouse cellmotion":     ansi.EnableMouseCellMotion,
	"disable mouse cellmotion":    ansi.DisableMouseCellMotion,
	"request mouse cellmotion":    ansi.RequestMouseCellMotion,
	"enable mouse allmotion":      ansi.EnableMouseAllMotion,
	"disable mouse allmotion":     ansi.DisableMouseAllMotion,
	"request mouse allmotion":     ansi.RequestMouseAllMotion,
	"enable report focus":         ansi.EnableReportFocus,
	"disable report focus":        ansi.DisableReportFocus,
	"request report focus":        ansi.RequestReportFocus,
	"enable mouse sgr":            ansi.EnableMouseSgrExt,
	"disable mouse sgr":           ansi.DisableMouseSgrExt,
	"request mouse sgr":           ansi.RequestMouseSgrExt,
	"enable altscreen":            ansi.EnableAltScreenBuffer,
	"disable altscreen":           ansi.DisableAltScreenBuffer,
	"request altscreen":           ansi.RequestAltScreenBuffer,
	"enable bracketed paste":      ansi.EnableBracketedPaste,
	"disable bracketed paste":     ansi.DisableBracketedPaste,
	"request bracketed paste":     ansi.RequestBracketedPaste,
	"enable synchronized output":  ansi.EnableSyncdOutput,
	"disable synchronized output": ansi.DisableSyncdOutput,
	"request synchronized output": ansi.RequestSyncdOutput,
	"enable grapheme clustering":  ansi.EnableGraphemeClustering,
	"disable grapheme clustering": ansi.DisableGraphemeClustering,
	"request grapheme clustering": ansi.RequestGraphemeClustering,
	"enable win32 input":          ansi.EnableWin32Input,
	"disable win32 input":         ansi.DisableWin32Input,
	"request win32 input":         ansi.RequestWin32Input,
}

var kitty = map[string]string{
	"set all mode 1": ansi.KittyKeyboard(ansi.KittyAllFlags, 1),
	"set all mode 2": ansi.KittyKeyboard(ansi.KittyAllFlags, 2),
	"set all mode 3": ansi.KittyKeyboard(ansi.KittyAllFlags, 3),
	"request":        ansi.RequestKittyKeyboard,
	"disable":        ansi.DisableKittyKeyboard,
	"pop":            ansi.PopKittyKeyboard(2),
	"push 1":         ansi.PushKittyKeyboard(1),
	"push 2":         ansi.PushKittyKeyboard(2),
	"push 4":         ansi.PushKittyKeyboard(4),
	"push 8":         ansi.PushKittyKeyboard(8),
	"push 16":        ansi.PushKittyKeyboard(16),
}

var others = map[string]string{
	"request primary device attrs": ansi.RequestPrimaryDeviceAttributes,
	"request xt version":           ansi.RequestXTVersion,
	"termcap":                      ansi.RequestTermcap("bw", "ccc"),
}

var sgr = map[string]string{
	"reset":   ansi.ResetStyle,
	"style 1": new(ansi.Style).Bold().Faint().Italic().CurlyUnderline().String(),
	"style 2": new(ansi.Style).SlowBlink().Reverse().Strikethrough().String(),
	"style 3": new(ansi.Style).RapidBlink().BackgroundColor(ansi.Green).ForegroundColor(ansi.BrightGreen).UnderlineColor(ansi.Blue).String(),
	"style 4": new(ansi.Style).BackgroundColor(ansi.BrightYellow).ForegroundColor(ansi.Black).UnderlineColor(ansi.BrightCyan).String(),
	"style 5": new(ansi.Style).BackgroundColor(ansi.TrueColor(0xffeeaa)).ForegroundColor(ansi.TrueColor(0xffeeaa)).UnderlineColor(ansi.TrueColor(0xffeeaa)).String(),
	"style 6": new(ansi.Style).BackgroundColor(ansi.ExtendedColor(255)).ForegroundColor(ansi.ExtendedColor(255)).UnderlineColor(ansi.ExtendedColor(255)).String(),
}

var title = map[string]string{
	"set":      ansi.SetWindowTitle("hello"),
	"set icon": ansi.SetIconName("terminal"),
	"set both": ansi.SetIconNameWindowTitle("terminal"),
}

var hyperlink = map[string]string{
	"uri only": ansi.SetHyperlink("https://charm.sh"),
	"full":     ansi.SetHyperlink("https://charm.sh", "my title"),
	"reset":    ansi.ResetHyperlink("my title"),
}

var notify = map[string]string{
	"notify": "\x1b]9;notification text\x07", // TODO: add to ansi?
}

var termcolor = map[string]string{
	"set bg":         ansi.SetBackgroundColor(ansi.Black),
	"set fg":         ansi.SetForegroundColor(ansi.Red),
	"set cursor":     ansi.SetCursorColor(ansi.Blue),
	"request bg":     ansi.RequestBackgroundColor,
	"request fg":     ansi.RequestForegroundColor,
	"request cursor": ansi.RequestCursorColor,
	"reset bg":       ansi.ResetBackgroundColor,
	"reset fg":       ansi.ResetForegroundColor,
	"reset cursor":   ansi.ResetCursorColor,
}

var clipboard = map[string]string{
	"request system":  ansi.RequestSystemClipboard,
	"request primary": ansi.RequestPrimaryClipboard,
	"set system":      ansi.SetSystemClipboard("hello"),
	"set primary":     ansi.SetPrimaryClipboard("hello"),
}

func TestSequences(t *testing.T) {
	for name, table := range map[string]map[string]string{
		"cursor":    cursor,
		"screen":    screen,
		"line":      line,
		"mode":      mode,
		"kitty":     kitty,
		"sgr":       sgr,
		"title":     title,
		"hyperlink": hyperlink,
		"notify":    notify,
		"termcolor": termcolor,
		"clipboard": clipboard,
		"others":    others,
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
