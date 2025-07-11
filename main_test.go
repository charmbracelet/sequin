package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/stretchr/testify/require"
)

var c0c1 = map[string]string{
	// c0
	"c0": func() string {
		var c0codes string
		for controlCode := 0x00; controlCode <= 0x1f; controlCode++ {
			c0codes = fmt.Sprintf("%s%c", c0codes, controlCode)
		}
		return c0codes
	}(),
	// c1
	"c1": func() string {
		var controlCode byte
		c1codes := []byte("")
		for controlCode = 0x80; controlCode <= 0x9f; controlCode++ {
			// Skip DCS, SOS, CSI, OSC, PM, and APC
			switch controlCode {
			case ansi.DCS, ansi.SOS, ansi.CSI, ansi.OSC, ansi.PM, ansi.APC:
				continue
			}
			c1codes = append(c1codes, controlCode)
		}
		return string(c1codes[:])
	}(),
}

var ascii = map[string]string{
	// space and del
	"ascii": fmt.Sprintf(" %c", ansi.DEL),
}

var cursor = map[string]string{
	// cursor
	"save":                         ansi.SaveCursor,
	"restore":                      ansi.RestoreCursor,
	"request pos":                  ansi.RequestCursorPosition,
	"request extended pos":         ansi.RequestExtendedCursorPosition,
	"invalid request extended pos": strings.Replace(ansi.RequestExtendedCursorPosition, "6", "7", 1),
	"up 1":                         ansi.CursorUp1,
	"up":                           ansi.CursorUp(5),
	"down 1":                       ansi.CursorDown1,
	"down":                         ansi.CursorDown(3),
	"right 1":                      ansi.CursorRight1,
	"right":                        ansi.CursorRight(3),
	"left 1":                       ansi.CursorLeft1,
	"left":                         ansi.CursorLeft(3),
	"next line":                    ansi.CursorNextLine(3),
	"previous line":                ansi.CursorPreviousLine(3),
	"set pos":                      ansi.SetCursorPosition(10, 20),
	"home pos":                     ansi.CursorHomePosition,
	"save pos":                     ansi.SaveCursorPosition,
	"restore pos":                  ansi.RestoreCursorPosition,
	"style 0":                      ansi.SetCursorStyle(0),
	"style 1":                      ansi.SetCursorStyle(1),
	"style 2":                      ansi.SetCursorStyle(2),
	"style 3":                      ansi.SetCursorStyle(3),
	"style 4":                      ansi.SetCursorStyle(4),
	"style 5":                      ansi.SetCursorStyle(5),
	"style 6":                      ansi.SetCursorStyle(6),
	"style 7":                      ansi.SetCursorStyle(7),
	"pointer shape":                ansi.SetPointerShape("crosshair"),
	"invalid pointer shape":        strings.Replace(ansi.SetPointerShape(""), ";", "", 1),
}

var screen = map[string]string{
	"enable alt buffer":  ansi.SetAltScreenBufferMode,
	"disable alt buffer": ansi.ResetAltScreenBufferMode,
	"request alt buffer": ansi.RequestAltScreenBufferMode,
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
	"enable cursor keys":          ansi.SetCursorKeysMode,
	"disable cursor keys":         ansi.ResetCursorKeysMode,
	"request cursor keys":         ansi.RequestCursorKeysMode,
	"enable cursor visibility":    ansi.ShowCursor,
	"disable cursor visibility":   ansi.HideCursor,
	"request cursor visibility":   ansi.RequestCursorVisibility,
	"enable mouse":                ansi.EnableMouse,
	"disable mouse":               ansi.DisableMouse,
	"request mouse":               ansi.RequestMouse,
	"enable mouse hilite":         ansi.EnableMouseHilite,
	"disable mouse hilite":        ansi.DisableMouseHilite,
	"request mouse hilite":        ansi.RequestMouseHilite,
	"enable mouse cellmotion":     ansi.EnableMouseCellMotion,
	"disable mouse cellmotion":    ansi.DisableMouseCellMotion,
	"request mouse cellmotion":    ansi.RequestMouseCellMotion,
	"enable mouse allmotion":      ansi.EnableMouseAllMotion,
	"disable mouse allmotion":     ansi.DisableMouseAllMotion,
	"request mouse allmotion":     ansi.RequestMouseAllMotion,
	"enable report focus":         ansi.SetFocusEventMode,
	"disable report focus":        ansi.ResetFocusEventMode,
	"request report focus":        ansi.RequestFocusEventMode,
	"enable mouse sgr":            ansi.SetSgrExtMouseMode,
	"disable mouse sgr":           ansi.ResetSgrExtMouseMode,
	"request mouse sgr":           ansi.RequestSgrExtMouseMode,
	"enable altscreen":            ansi.SetAltScreenBufferMode,
	"disable altscreen":           ansi.ResetAltScreenBufferMode,
	"request altscreen":           ansi.RequestAltScreenBufferMode,
	"enable bracketed paste":      ansi.SetBracketedPasteMode,
	"disable bracketed paste":     ansi.ResetBracketedPasteMode,
	"request bracketed paste":     ansi.RequestBracketedPasteMode,
	"enable synchronized output":  ansi.SetSynchronizedOutputMode,
	"disable synchronized output": ansi.ResetSynchronizedOutputMode,
	"request synchronized output": ansi.RequestSynchronizedOutputMode,
	"enable grapheme clustering":  ansi.SetGraphemeClusteringMode,
	"disable grapheme clustering": ansi.ResetGraphemeClusteringMode,
	"request grapheme clustering": ansi.RequestGraphemeClusteringMode,
	"enable win32 input":          ansi.SetWin32InputMode,
	"disable win32 input":         ansi.ResetWin32InputMode,
	"request win32 input":         ansi.RequestWin32InputMode,
	"invalid":                     strings.Replace(ansi.ShowCursor, "25", "27", 1),
	"non private":                 strings.Replace(ansi.ShowCursor, "?", "", 1),
}

var kitty = map[string]string{
	"set all mode 1":   ansi.KittyKeyboard(ansi.KittyAllFlags, 1),
	"set all mode 2":   ansi.KittyKeyboard(ansi.KittyAllFlags, 2),
	"set all mode 3":   ansi.KittyKeyboard(ansi.KittyAllFlags, 3),
	"set invalid mode": ansi.KittyKeyboard(ansi.KittyAllFlags, 4),
	"request":          ansi.RequestKittyKeyboard,
	"disable":          "\x1b[>0u",
	"pop":              ansi.PopKittyKeyboard(2),
	"push 1":           ansi.PushKittyKeyboard(1),
	"push 2":           ansi.PushKittyKeyboard(2),
	"push 4":           ansi.PushKittyKeyboard(4),
	"push 8":           ansi.PushKittyKeyboard(8),
	"push 16":          ansi.PushKittyKeyboard(16),
}

var others = map[string]string{
	"request primary device attrs": ansi.RequestPrimaryDeviceAttributes,
	"request xt version":           ansi.RequestXTVersion,
	"termcap":                      ansi.RequestTermcap("bw", "ccc"),
	"invalid termcap":              strings.Replace(ansi.RequestTermcap("a"), hex.EncodeToString([]byte("a")), "", 1),
	"invalid termcap hex":          strings.Replace(ansi.RequestTermcap("a"), hex.EncodeToString([]byte("a")), "a", 1),
	"invalid xt":                   "\x1b[>1q",
	"text":                         "some text",
	"bold text":                    new(ansi.Style).Bold().String() + "some text" + ansi.ResetStyle,
	"esc":                          fmt.Sprintf("%c", ansi.ESC),
	"file sep":                     fmt.Sprintf("%c", ansi.FS),
	"apc":                          "\x1b_Hello World\x1b\\",
	"pm":                           "\x1b^Hello World\x1b\\",
	"sos":                          "\x1bXHello World\x1b\\",
}

var sgr = map[string]string{
	"reset":                        ansi.ResetStyle + strings.Replace(ansi.ResetStyle, "m", "0m", 1),
	"style 1":                      new(ansi.Style).Bold().Faint().Italic().CurlyUnderline().String(),
	"style 2":                      new(ansi.Style).SlowBlink().Reverse().Strikethrough().String(),
	"style 3":                      new(ansi.Style).RapidBlink().BackgroundColor(ansi.Green).ForegroundColor(ansi.BrightGreen).UnderlineColor(ansi.Blue).String(),
	"style 4":                      new(ansi.Style).BackgroundColor(ansi.BrightYellow).ForegroundColor(ansi.Black).UnderlineColor(ansi.BrightCyan).String(),
	"style 5":                      new(ansi.Style).BackgroundColor(ansi.TrueColor(0xffeeaa)).ForegroundColor(ansi.TrueColor(0xffeeaa)).UnderlineColor(ansi.TrueColor(0xffeeaa)).String(),
	"style 6":                      new(ansi.Style).BackgroundColor(ansi.ExtendedColor(255)).ForegroundColor(ansi.ExtendedColor(255)).UnderlineColor(ansi.ExtendedColor(255)).String(),
	"style 7":                      new(ansi.Style).NoUnderline().NoItalic().NormalIntensity().NoBlink().NoConceal().NoReverse().NoStrikethrough().String(),
	"style 8":                      new(ansi.Style).UnderlineStyle(ansi.NoUnderlineStyle).DefaultBackgroundColor().String(),
	"style 9":                      strings.Replace(new(ansi.Style).UnderlineStyle(ansi.SingleUnderlineStyle).DefaultForegroundColor().String(), "[4", "[4:1", 1),
	"style 10":                     new(ansi.Style).UnderlineStyle(ansi.DoubleUnderlineStyle).String(),
	"style 11":                     new(ansi.Style).UnderlineStyle(ansi.CurlyUnderlineStyle).String(),
	"style 12":                     new(ansi.Style).UnderlineStyle(ansi.DottedUnderlineStyle).String(),
	"style 13":                     new(ansi.Style).UnderlineStyle(ansi.DashedUnderlineStyle).Conceal().String(),
	"empty values":                 strings.Replace(new(ansi.Style).Bold().String(), "[", "[;;;", 1),
	"underlined text, but no bold": new(ansi.Style).UnderlineStyle(ansi.CurlyUnderlineStyle).Bold().String(),
	"mittchels tweet":              "\033[;4:3;38;2;175;175;215;58:2::190:80:70m",
}

var title = map[string]string{
	"set":         ansi.SetWindowTitle("hello"),
	"set icon":    ansi.SetIconName("terminal"),
	"set both":    ansi.SetIconNameWindowTitle("terminal"),
	"invalid":     strings.Replace(ansi.SetWindowTitle("hello"), ";hello", "", 1),
	"invalid cmd": strings.Replace(ansi.SetWindowTitle("hello"), "2", "5", 1),
}

var cwd = map[string]string{
	"single part":    ansi.NotifyWorkingDirectory("localhost", "foo"),
	"multiple parts": ansi.NotifyWorkingDirectory("localhost", "foo", "bar"),
	"invalid":        strings.Replace(ansi.NotifyWorkingDirectory("localhost", "foo"), ";", "", 1),
	"invalid url":    strings.Replace(ansi.NotifyWorkingDirectory("localhost", "foo"), "file://localhost/foo", "foooooo:/bar", 1),
}

var hyperlink = map[string]string{
	"uri only":        ansi.SetHyperlink("https://charm.sh"),
	"full":            ansi.SetHyperlink("https://charm.sh", "my title"),
	"reset":           ansi.ResetHyperlink("my title"),
	"multiple params": ansi.SetHyperlink("https://charm.sh", "my title", "some description"),
	"invalid":         strings.Replace(ansi.ResetHyperlink(), ";", "", 1),
}

var notify = map[string]string{
	"notify":  ansi.Notify("notification body"),
	"invalid": strings.Replace(ansi.Notify(""), ";", "", 1),
}

var termcolor = map[string]string{
	"set bg":         ansi.SetBackgroundColor("#000000"),
	"set fg":         ansi.SetForegroundColor("#800000"),
	"set cursor":     ansi.SetCursorColor("#000080"),
	"request bg":     ansi.RequestBackgroundColor,
	"request fg":     ansi.RequestForegroundColor,
	"request cursor": ansi.RequestCursorColor,
	"reset bg":       ansi.ResetBackgroundColor,
	"reset fg":       ansi.ResetForegroundColor,
	"reset cursor":   ansi.ResetCursorColor,
	"invalid set":    strings.Replace(ansi.SetBackgroundColor("#000000"), ";", "", 1),
	"invalid reset":  strings.Replace(ansi.ResetBackgroundColor, "111", "111;1", 1),
}

var clipboard = map[string]string{
	"request system":  ansi.RequestSystemClipboard,
	"request primary": ansi.RequestPrimaryClipboard,
	"set system":      ansi.SetSystemClipboard("hello"),
	"set primary":     ansi.SetPrimaryClipboard("hello"),
	"incomplete":      strings.Replace(ansi.RequestPrimaryClipboard, ";?", "", 1),
	"invalid":         strings.Replace(ansi.SetPrimaryClipboard("hello"), "=", "", 1),
}

var finalterm = map[string]string{
	"prompt start":               ansi.FinalTermPrompt(),
	"prompt start invalid":       ansi.FinalTerm("AB"),
	"command start":              ansi.FinalTermCmdStart(),
	"command executed":           ansi.FinalTermCmdExecuted(),
	"command finished":           ansi.FinalTermCmdFinished(),
	"command finished exit code": ansi.FinalTermCmdFinished("127"),
	"invalid":                    ansi.FinalTerm("Q"),
}

var keypad = map[string]string{
	"normal keypad":      ansi.KeypadNumericMode,
	"application keypad": ansi.KeypadApplicationMode,
}

func TestSequences(t *testing.T) {
	for name, table := range map[string]map[string]string{
		"c0c1":      c0c1,
		"ascii":     ascii,
		"cursor":    cursor,
		"screen":    screen,
		"line":      line,
		"mode":      mode,
		"kitty":     kitty,
		"sgr":       sgr,
		"title":     title,
		"cwd":       cwd,
		"hyperlink": hyperlink,
		"notify":    notify,
		"termcolor": termcolor,
		"clipboard": clipboard,
		"others":    others,
		"finalterm": finalterm,
		"keypad":    keypad,
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
