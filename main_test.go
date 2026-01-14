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
	"request pos":                  ansi.RequestCursorPositionReport,
	"request extended pos":         ansi.RequestExtendedCursorPositionReport,
	"invalid request extended pos": strings.Replace(ansi.RequestExtendedCursorPositionReport, "6", "7", 1),
	"up 1":                         ansi.CUU1,
	"up":                           ansi.CursorUp(5),
	"down 1":                       ansi.CUD1,
	"down":                         ansi.CursorDown(3),
	"right 1":                      ansi.CUF1,
	"right":                        ansi.CursorForward(3),
	"left 1":                       ansi.CUB1,
	"left":                         ansi.CursorBackward(3),
	"next line":                    ansi.CursorNextLine(3),
	"previous line":                ansi.CursorPreviousLine(3),
	"set pos":                      ansi.CursorPosition(10, 20),
	"home pos":                     ansi.CursorHomePosition,
	"save pos":                     ansi.SaveCurrentCursorPosition,
	"restore pos":                  ansi.RestoreCurrentCursorPosition,
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
	"enable alt buffer":  ansi.SetModeAltScreenSaveCursor,
	"disable alt buffer": ansi.ResetModeAltScreenSaveCursor,
	"request alt buffer": ansi.RequestModeAltScreenSaveCursor,
	"passthrough":        ansi.ScreenPassthrough(ansi.SaveCursor, 0), // TODO: impl
	"erase above":        ansi.EraseScreenAbove,
	"erase below":        ansi.EraseScreenBelow,
	"erase full":         ansi.EraseEntireScreen,
	"erase display":      ansi.EraseEntireDisplay,
	"scrolling region":   ansi.SetTopBottomMargins(10, 20),
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
	"enable cursor keys":          ansi.SetModeCursorKeys,
	"disable cursor keys":         ansi.ResetModeCursorKeys,
	"request cursor keys":         ansi.RequestModeCursorKeys,
	"enable cursor visibility":    ansi.SetModeTextCursorEnable,
	"disable cursor visibility":   ansi.ResetModeTextCursorEnable,
	"request cursor visibility":   ansi.RequestModeTextCursorEnable,
	"enable mouse":                ansi.SetModeMouseNormal,
	"disable mouse":               ansi.ResetModeMouseNormal,
	"request mouse":               ansi.RequestModeMouseNormal,
	"enable mouse hilite":         ansi.SetModeMouseHighlight,
	"disable mouse hilite":        ansi.ResetModeMouseHighlight,
	"request mouse hilite":        ansi.RequestModeMouseHighlight,
	"enable mouse cellmotion":     ansi.SetModeMouseButtonEvent,
	"disable mouse cellmotion":    ansi.ResetModeMouseButtonEvent,
	"request mouse cellmotion":    ansi.RequestModeMouseButtonEvent,
	"enable mouse allmotion":      ansi.SetModeMouseAnyEvent,
	"disable mouse allmotion":     ansi.ResetModeMouseAnyEvent,
	"request mouse allmotion":     ansi.RequestModeMouseAnyEvent,
	"enable report focus":         ansi.SetModeFocusEvent,
	"disable report focus":        ansi.ResetModeFocusEvent,
	"request report focus":        ansi.RequestModeFocusEvent,
	"enable mouse sgr":            ansi.SetModeMouseExtSgr,
	"disable mouse sgr":           ansi.ResetModeMouseExtSgr,
	"request mouse sgr":           ansi.RequestModeMouseExtSgr,
	"enable altscreen":            ansi.SetModeAltScreenSaveCursor,
	"disable altscreen":           ansi.ResetModeAltScreenSaveCursor,
	"request altscreen":           ansi.RequestModeAltScreenSaveCursor,
	"enable bracketed paste":      ansi.SetModeBracketedPaste,
	"disable bracketed paste":     ansi.ResetModeBracketedPaste,
	"request bracketed paste":     ansi.RequestModeBracketedPaste,
	"enable synchronized output":  ansi.SetModeSynchronizedOutput,
	"disable synchronized output": ansi.ResetModeSynchronizedOutput,
	"request synchronized output": ansi.RequestModeSynchronizedOutput,
	"enable grapheme clustering":  ansi.SetModeUnicodeCore,
	"disable grapheme clustering": ansi.ResetModeUnicodeCore,
	"request grapheme clustering": ansi.RequestModeUnicodeCore,
	"enable win32 input":          ansi.SetModeWin32Input,
	"disable win32 input":         ansi.ResetModeWin32Input,
	"request win32 input":         ansi.RequestModeWin32Input,
	"invalid":                     strings.Replace(ansi.SetModeTextCursorEnable, "25", "27", 1),
	"non private":                 strings.Replace(ansi.SetModeTextCursorEnable, "?", "", 1),
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
	"request xt version":           ansi.RequestNameVersion,
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
	"style 1":                      new(ansi.Style).Bold().Faint().Italic(true).UnderlineStyle(ansi.UnderlineStyleCurly).String(),
	"style 2":                      new(ansi.Style).Blink(true).Reverse(true).Strikethrough(true).String(),
	"style 3":                      new(ansi.Style).RapidBlink(true).BackgroundColor(ansi.Green).ForegroundColor(ansi.BrightGreen).UnderlineColor(ansi.Blue).String(),
	"style 4":                      new(ansi.Style).BackgroundColor(ansi.BrightYellow).ForegroundColor(ansi.Black).UnderlineColor(ansi.BrightCyan).String(),
	"style 5":                      new(ansi.Style).BackgroundColor(ansi.RGBColor{R: 0xff, G: 0xee, B: 0xaa}).ForegroundColor(ansi.RGBColor{R: 0xff, G: 0xee, B: 0xaa}).UnderlineColor(ansi.RGBColor{R: 0xff, G: 0xee, B: 0xaa}).String(),
	"style 6":                      new(ansi.Style).BackgroundColor(ansi.IndexedColor(255)).ForegroundColor(ansi.IndexedColor(255)).UnderlineColor(ansi.IndexedColor(255)).String(),
	"style 7":                      new(ansi.Style).Underline(false).Italic(false).Normal().Blink(false).Conceal(false).Reverse(false).Strikethrough(false).String(),
	"style 8":                      new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleNone).BackgroundColor(nil).String(),
	"style 9":                      strings.Replace(new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleSingle).ForegroundColor(nil).String(), "[4", "[4:1", 1),
	"style 10":                     new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleDouble).String(),
	"style 11":                     new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleCurly).String(),
	"style 12":                     new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleDotted).String(),
	"style 13":                     new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleDashed).Conceal(true).String(),
	"empty values":                 strings.Replace(new(ansi.Style).Bold().String(), "[", "[;;;", 1),
	"underlined text, but no bold": new(ansi.Style).UnderlineStyle(ansi.UnderlineStyleCurly).Bold().String(),
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

var xtmodkeys = map[string]string{
	"enable modifyOtherKeys":    "\x1b[>4;1m",
	"disable modifyOtherKeys":   "\x1b[>4;0m",
	"reset modifyOtherKeys":     "\x1b[>4m",
	"enable modifyCursorKeys":   "\x1b[>1;1m",
	"enable modifyFunctionKeys": "\x1b[>2;1m",
	"enable modifyKeypadKeys":   "\x1b[>3;1m",
	"enable modifyModifierKeys": "\x1b[>6;1m",
	"enable modifySpecialKeys":  "\x1b[>7;1m",
	"disable modifyCursorKeys":  "\x1b[>1;0m",
	"reset modifyFunctionKeys":  "\x1b[>2m",
	"unknown resource enable":   "\x1b[>5;1m",
	"unknown resource disable":  "\x1b[>9;0m",
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
		"xtmodkeys": xtmodkeys,
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
