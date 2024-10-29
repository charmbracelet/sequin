package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
)

const (
	markerShift   = 8
	intermedShift = 16
)

func main() {
	in := strings.Join(os.Args[1:], " ")
	if in == "-" || in == "" {
		bts, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		in = string(bts)
	}

	var state byte
	p := ansi.GetParser()
	defer ansi.PutParser(p)

	for len(in) > 0 {
		seq, width, n, newState := ansi.DecodeSequence(in, state, p)

		switch {
		case ansi.HasCsiPrefix(seq):
			switch p.Cmd {
			case 'm':
				// SGR
				handleSgr(p)
			case 'A':
				// CUU
			case 'B':
				// CUD
			case 'C':
				// CUF
			case 'D':
				// CUB
			case 'p' | '$'<<intermedShift:
				// ANSI mode DECRQM
			case 'p' | '?'<<markerShift | '$'<<intermedShift:
				// Private DEC mode DECRQM
			}
		case ansi.HasDcsPrefix(seq):
		case ansi.HasOscPrefix(seq):
			switch p.Cmd {
			case 0:
				// set window and icon title
			case 1:
				// set icon title
			case 2:
				// set window title
			case 8:
				// hyperlinks
			case 9:
				// notify
			}
		case ansi.HasApcPrefix(seq):
			switch {
			case ansi.HasPrefix(p.Data, []byte("G")):
				// Kitty graphics
			}
		case ansi.HasEscPrefix(seq):
			if len(seq) == 1 {
				// just an ESC
				fmt.Println("Control code: ESC")
				break
			}

			switch p.Cmd {
			case 7:
				// save cursor
			case 8:
				// restore cursor
			}
		default:
			if width == 0 && len(seq) == 1 {
				// control code
				fmt.Printf("Control code: %q\n", seq)
				break
			} else if width == 0 {
				fmt.Printf("Unknown: %q\n", seq)
				break
			}

			fmt.Printf("Print: %q\n", seq)
		}

		in = in[n:]
		state = newState
	}
}

func parse(in []byte) []string {
	p := ansi.GetParser()
	p.Reset()

	var sequences []ansi.Sequence
	p.Parse(func(s ansi.Sequence) {
		sequences = append(sequences, s)
	}, in)

	var r []string
	for _, s := range sequences {
		switch seq := s.(type) {
		case ansi.ControlCode:
			r = append(r, "Control code")
		case ansi.SosSequence:
			r = append(r, "SOS: TODO")
		case ansi.ApcSequence:
			r = append(r, "APC: TODO")
		case ansi.EscSequence:
			switch seq.Command() {
			case '7':
				r = append(r, "ESC 7: Save cursor")
			case '8':
				r = append(r, "ESC 8: Restore cursor")
			default:
				r = append(r, fmt.Sprintf("ESC %s: Restore cursor", seq.String()))
			}
		case ansi.OscSequence:
			switch seq.Command() {
			case 8:
				var uri, params string
				if len(seq.Params()) > 2 {
					params = seq.Params()[1]
					uri = seq.Params()[2]
				}
				r = append(r, fmt.Sprintf("OSC 8 ; %s ; %s ST: Set hyperlink '%[2]s' with params '%[1]s'", params, uri))
			case 10:
				var color string
				if len(seq.Params()) > 1 {
					color = seq.Params()[1]
				}
				r = append(r, fmt.Sprintf("OSC 10 ; %s ST: Set foreground color '%[1]s'", color))
			case 11:
				var color string
				if len(seq.Params()) > 1 {
					color = seq.Params()[1]
				}
				r = append(r, fmt.Sprintf("OSC 11 ; %s ST: Set background color '%[1]s'", color))
			case 12:
				var color string
				if len(seq.Params()) > 1 {
					color = seq.Params()[1]
				}
				r = append(r, fmt.Sprintf("OSC 12 ; %s ST: Set cursor color '%[1]s'", color))
			case 22:
				var shape string
				if len(seq.Params()) > 1 {
					shape = seq.Params()[1]
				}
				r = append(r, fmt.Sprintf("OSC 22 ; %s ST: Set mouse shape '%[1]s'", shape))
			case 52:
				var clip string
				if len(seq.Params()) > 1 {
					clip = seq.Params()[1]
				}
				if len(seq.Params()) > 3 {
					r = append(r, fmt.Sprintf("OSC 52 ; %s; <hidden> ; ST: Set %s clipboard", clip, clipboardDesc(clip)))
				} else {
					r = append(r, fmt.Sprintf("OSC 52 ; %s; ST: Request %s clipboard", clip, clipboardDesc(clip)))
				}
			default:
				r = append(r, fmt.Sprintf("OSC %s: TODO", seq.String()))
			}
		case ansi.CsiSequence:
			switch seq.Command() {
			case 'A':
				lines := 1
				if seq.Len() > 0 {
					lines = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d A: Cursor up %[1]d lines", lines))
			case 'B':
				lines := 1
				if seq.Len() > 0 {
					lines = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d B: Cursor down %[1]d lines", lines))
			case 'C':
				lines := 1
				if seq.Len() > 0 {
					lines = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d C: Cursor right %[1]d lines", lines))
			case 'D':
				lines := 1
				if seq.Len() > 0 {
					lines = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d D: Cursor left %[1]d lines", lines))
			case 'E':
				times := 1
				if seq.Len() > 0 {
					times = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d E: Cursor next line %[1]d times", times))
			case 'F':
				times := 1
				if seq.Len() > 0 {
					times = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d F: Cursor previous line %[1]d times", times))
			case 'H':
				row := 1
				col := 1
				if seq.Len() > 1 {
					row = seq.Param(0)
					col = seq.Param(1)
				}
				r = append(r, fmt.Sprintf("CSI %d;%d H: Set cursor position row=%[1]d col=%[2]d", row, col))
			case 'J':
				switch seq.Param(0) {
				case 0:
					r = append(r, "CSI 0 J: Erase screen below")
				case 1:
					r = append(r, "CSI 1 J: Erase screen above")
				case 2:
					r = append(r, "CSI 2 J: Erase entire screen")
				case 3:
					r = append(r, "CSI 3 J: Erase entire display")
				}
			case 'K':
				switch seq.Param(0) {
				case 0:
					r = append(r, "CSI 0 K: Erase line right")
				case 1:
					r = append(r, "CSI 1 K: Erase line left")
				case 2:
					r = append(r, "CSI 2 K: Erase entire line")
				}
			case 'L':
				r = append(r, fmt.Sprintf("CSI %d L: Insert %[1]d blank lines", seq.Param(0)))
			case 'M':
				r = append(r, fmt.Sprintf("CSI %d M: Delete %[1]d lines", seq.Param(0)))
			case 'S':
				r = append(r, fmt.Sprintf("CSI %d S: Scroll up %[1]d lines", seq.Param(0)))
			case 'T':
				r = append(r, fmt.Sprintf("CSI %d T: Scroll down %[1]d lines", seq.Param(0)))
			case 'h':
				switch seq.Marker() {
				case '?':
					switch seq.Param(0) {
					case 1:
						r = append(r, "CSI ? 1 h: Enable cursor keys")
					case 25:
						r = append(r, "CSI ? 25 h: Show cursor")
					case 1000:
						r = append(r, "CSI ? 1000 h: Enable mouse")
					case 1001:
						r = append(r, "CSI ? 1001 h: Enable mouse hilite")
					case 1002:
						r = append(r, "CSI ? 1002 h: Enable mouse cell motion")
					case 1003:
						r = append(r, "CSI ? 1003 $ h: Enable mouse all motion")
					case 1004:
						r = append(r, "CSI ? 1004 $ h: Enable report focus")
					case 1006:
						r = append(r, "CSI ? 1006 $ h: Enable mouse SGR ext")
					case 1049:
						r = append(r, "CSI ? 1049 $ h: Enable altscreen mode")
					case 2004:
						r = append(r, "CSI ? 2004 $ h: Enable bracketed paste mode")
					case 2026:
						r = append(r, "CSI ? 2026 $ h: Enable synchronized output mode")
					case 2027:
						r = append(r, "CSI ? 2027 $ h: Enable grapheme clustering mode")
					case 9001:
						r = append(r, "CSI ? 9001 $ h: Enable win32 input mode")
					}
				}
			case 'c':
				r = append(r, "CSI c: Request primary device attributes")
			case 'l':
				switch seq.Marker() {
				case '?':
					switch seq.Param(0) {
					case 1:
						r = append(r, "CSI ? 1 l: Disable cursor keys")
					case 25:
						r = append(r, "CSI ? 25 h: Hide cursor")
					case 1000:
						r = append(r, "CSI ? 1000 l: Disable mouse")
					case 1001:
						r = append(r, "CSI ? 1001 l: Disable mouse hilite")
					case 1002:
						r = append(r, "CSI ? 1002 l: Disable mouse cell motion")
					case 1003:
						r = append(r, "CSI ? 1003 $ l: Disable mouse all motion")
					case 1004:
						r = append(r, "CSI ? 1004 $ l: Disable report focus")
					case 1006:
						r = append(r, "CSI ? 1006 $ l: Disable mouse SGR ext")
					case 1049:
						r = append(r, "CSI ? 1049 $ l: Disable altscreen mode")
					case 2004:
						r = append(r, "CSI ? 2004 $ l: Disable bracketed paste mode")
					case 2026:
						r = append(r, "CSI ? 2026 $ l: Disable synchronized output mode")
					case 2027:
						r = append(r, "CSI ? 2027 $ l: Disable grapheme clustering mode")
					case 9001:
						r = append(r, "CSI ? 9001 $ l: Disable win32 input mode")
					}
				}
			case 'm':
				r = append(r, parseSGR(seq)...)
			case 'n':
				switch seq.Param(0) {
				case 6:
					if seq.Marker() > 0 {
						r = append(r, "CSI ? 6 n: Request extended cursor position")
					} else {
						r = append(r, "CSI 6 n: Request cursor position")
					}
				default:
					r = append(r, fmt.Sprintf("CSI %d n: TODO", seq.Param(0)))
				}
			case 'p':
				switch seq.Marker() {
				case '?':
					if seq.Intermediate() == '$' {
						switch seq.Param(0) {
						case 1:
							r = append(r, "CSI ? 1 $ p: Request cursor keys")
						case 25:
							r = append(r, "CSI ? 25 $ p: Request cursor visibility")
						case 1000:
							r = append(r, "CSI ? 1000 $ p: Request mouse")
						case 1001:
							r = append(r, "CSI ? 1001 $ p: Request mouse hilite")
						case 1002:
							r = append(r, "CSI ? 1002 $ p: Request mouse cell motion")
						case 1003:
							r = append(r, "CSI ? 1003 $ p: Request mouse all motion")
						case 1004:
							r = append(r, "CSI ? 1004 $ p: Request report focus")
						case 1006:
							r = append(r, "CSI ? 1006 $ p: Request mouse SGR ext")
						case 1049:
							r = append(r, "CSI ? 1049 $ p: Request altscreen mode")
						case 2004:
							r = append(r, "CSI ? 2004 $ p: Request bracketed paste mode")
						case 2026:
							r = append(r, "CSI ? 2026 $ p: Request synchronized output mode")
						case 2027:
							r = append(r, "CSI ? 2027 $ p: Request grapheme clustering mode")
						case 9001:
							r = append(r, "CSI ? 9001 $ p: Request win32 input mode")
						}
					}
				}
			case 'q':
				if seq.Marker() == '>' && seq.Param(0) == 0 {
					r = append(r, "CSI > 0 q: Request XT version")
				} else {
					cursor := 1
					if seq.Len() > 0 {
						cursor = seq.Param(0)
					}
					r = append(r, fmt.Sprintf("CSI %d q: Set cursor style '%s'", cursor, cursorDesc(cursor)))
				}
			case 'r':
				r = append(r, fmt.Sprintf(
					"CSI %d ; %d r: Set scrolling region to top=%[1]d bottom=%[2]d",
					seq.Param(0),
					seq.Param(1),
				))
			case 's':
				r = append(r, "CSI s: Save cursor position")
			case 'u':
				switch seq.Marker() {
				case 0:
					r = append(r, "CSI u: Restore cursor position")
				case '?':
					r = append(r, "CSI ? u: Request Kitty keyboard")
				case '=':
					r = append(r, fmt.Sprintf(
						"CSI = u: Set Kitty keyboard flags=%q mode=%q",
						kittyFlagsDesc(seq.Param(0)),
						kittyModeDesc(seq.Param(1)),
					))
				case '>':
					if seq.Param(0) == 0 {
						r = append(r, "CSI > 0 u: Disable Kitty keyboard")
					} else {
						r = append(r, fmt.Sprintf(
							"CSI > %d u: Push Kitty keyboard flags=%q",
							seq.Param(0),
							kittyFlagsDesc(seq.Param(0)),
						))
					}
				case '<':
					r = append(r, fmt.Sprintf(
						"CSI < %d u: Pop Kitty keyboard %[1]d times",
						seq.Param(0),
					))
				default:
					r = append(r, fmt.Sprintf(
						"CSI %s: TODO",
						seq.String(),
					))
				}
			default:
				r = append(r, fmt.Sprintf("CSI %s: TODO", seq.String()))
			}
		}
	}

	return r
}

func cursorDesc(i int) string {
	switch i {
	case 2:
		return "Steady block"
	case 3:
		return "Blinking underline"
	case 4:
		return "Steady underline"
	case 5:
		return "Blinking bar (xterm)"
	case 6:
		return "Steady bar (xterm)"
	default:
		return "Blinking block"
	}
}

func kittyFlagsDesc(flag int) string {
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

func kittyModeDesc(mode int) string {
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

func clipboardDesc(s string) string {
	switch s {
	case string(ansi.SystemClipboard):
		return "System"
	case string(ansi.PrimaryClipboard):
		return "Primary"
	}
	return "Unknown"
}

func parseSGR(seq ansi.CsiSequence) []string {
	var r []string
	var done int
	seq.Range(func(i, param int, hasMore bool) bool {
		if done > 0 {
			done--
			return true
		}
		switch param {
		case 0:
			r = append(r, "CSI 0m: Reset all attributes")
		case 1:
			r = append(r, "CSI 1m: Set bold")
		case 2:
			r = append(r, "CSI 2m: Set faint")
		case 3:
			r = append(r, "CSI 3m: Set italic")
		case 4:
			r = append(r, "CSI 4m: Set underline")
		case 5:
			r = append(r, "CSI 5m: Set slow blink")
		case 6:
			r = append(r, "CSI 6m: Set rapid blink")
		case 7:
			r = append(r, "CSI 7m: Set reverse video")
		case 8:
			r = append(r, "CSI 8m: Set concealed")
		case 9:
			r = append(r, "CSI 9m: Set crossed-out")
		case 21:
			r = append(r, "CSI 21m: Set double underline")
		case 22:
			r = append(r, "CSI 22m: Reset bold and faint")
		case 23:
			r = append(r, "CSI 23m: Reset italic")
		case 24:
			r = append(r, "CSI 24m: Reset underline")
		case 25:
			r = append(r, "CSI 25m: Reset blink")
		case 27:
			r = append(r, "CSI 27m: Reset reverse video")
		case 28:
			r = append(r, "CSI 28m: Reset concealed")
		case 29:
			r = append(r, "CSI 29m: Reset crossed-out")
		case 30, 31, 32, 33, 34, 35, 36, 37:
			r = append(r, fmt.Sprintf("CSI %dm: Set foreground color to %s", param, colorName(param-30)))
		case 38:
			fmt.Println("AQUI")
			nextParam := seq.Param(i + 1)
			if nextParam == 5 && i+2 < seq.Len() {
				r = append(r, fmt.Sprintf("CSI 38 ; 5 ; %d m: Set foreground color to 8-bit color %d", seq.Param(i+2), seq.Param(i+2)))
				done += 2
			} else if nextParam == 2 && i+4 < seq.Len() {
				r = append(r, fmt.Sprintf("CSI 38 ; 2 ; %d ; %d ; %d m: Set foreground color to RGB(%d,%d,%d)",
					seq.Param(i+2), seq.Param(i+3), seq.Param(i+4),
					seq.Param(i+2), seq.Param(i+3), seq.Param(i+4)))
				done += 4
			}
		case 39:
			r = append(r, "CSI 39m: Reset foreground color")
		case 40, 41, 42, 43, 44, 45, 46, 47:
			r = append(r, fmt.Sprintf("CSI %d m: Set background color to %s", param, colorName(param-40)))
		case 48:
			nextParam := seq.Param(i + 1)
			if nextParam == 5 && i+2 < seq.Len() {
				r = append(r, fmt.Sprintf("CSI 48 ; 5 ; %d m: Set background color to 8-bit color %d", seq.Param(i+2), seq.Param(i+2)))
				done += 2
			} else if nextParam == 2 && i+4 < seq.Len() {
				r = append(r, fmt.Sprintf("CSI 48 ; 2 ; %d ; %d ; %d m: Set background color to RGB(%d,%d,%d)",
					seq.Param(i+2), seq.Param(i+3), seq.Param(i+4),
					seq.Param(i+2), seq.Param(i+3), seq.Param(i+4)))
				done += 4
			}
		case 49:
			r = append(r, "CSI 49m: Reset background color")
		case 59:
			r = append(r, "CSI 59m: Reset underline color")
		case 90, 91, 92, 93, 94, 95, 96, 97:
			r = append(r, fmt.Sprintf("CSI %dm: Set bright foreground color to %s", param, colorName(param-90)))
		case 100, 101, 102, 103, 104, 105, 106, 107:
			r = append(r, fmt.Sprintf("CSI %dm: Set bright background color to %s", param, colorName(param-100)))
		}
		return true
	})
	return r
}

func colorName(index int) string {
	colors := []string{"Black", "Red", "Green", "Yellow", "Blue", "Magenta", "Cyan", "White"}
	if index >= 0 && index < len(colors) {
		return colors[index]
	}
	return "Unknown"
}
