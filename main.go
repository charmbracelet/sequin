package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
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

	for _, desc := range parse([]byte(in)) {
		fmt.Println(desc)
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
		case ansi.EscSequence:
			switch seq.Command() {
			case '7':
				r = append(r, "ESC 7: Save cursor")
			case '8':
				r = append(r, "ESC 8: Restore cursor")
			}
		case ansi.OscSequence:
			switch seq.Command() {
			case 8:
				var uri, params string
				if len(seq.Params()) > 2 {
					params = seq.Params()[1]
					uri = seq.Params()[2]
				}
				r = append(r, fmt.Sprintf("OSC 8 ; %s ; %s ; Uri ST: Set hyperlink '%[2]s' with params '%[1]s'", params, uri))
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
				r = append(r, fmt.Sprintf("OSC %d: TODO", seq.Command()))
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
			case 'c':
				r = append(r, "CSI c: Request primary device attributes")
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
			case 's':
				r = append(r, "CSI s: Save cursor position")
			case 'u':
				r = append(r, "CSI u: Restore cursor position")
			case 'q':
				if seq.Marker() == '>' && seq.Param(0) == 0 {
					r = append(r, "CSI > 0 q: Request XT version")
				}
				cursor := 1
				if seq.Len() > 0 {
					cursor = seq.Param(0)
				}
				r = append(r, fmt.Sprintf("CSI %d q: Set cursor style '%s'", cursor, cursorDesc(cursor)))
			default:
				r = append(r, fmt.Sprintf("CSI %d: TODO", rune(seq.Command())))
			}
		}
	}

	return r
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
