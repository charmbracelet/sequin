package main

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/parser"
	"github.com/spf13/cobra"
)

const (
	markerShift   = parser.MarkerShift
	intermedShift = parser.IntermedShift
)

var (
	buf bytes.Buffer
	raw bool
)

func main() {
	if err := cmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "sequin",
		Short: "Human-readable ANSI sequences",
		Args:  cobra.NoArgs,
		Example: `
printf '\x1b[m' | sequin
sequin <file
	`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			w := colorprofile.NewWriter(cmd.OutOrStdout(), os.Environ())
			in, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				//nolint:wrapcheck
				return err
			}
			return exec(w, in)
		},
	}
	root.Flags().BoolVarP(&raw, "raw", "r", false, "raw mode (no explanation)")
	return root
}

//nolint:mnd
func exec(w *colorprofile.Writer, in []byte) error {
	hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	lightDark := func(light, dark string) color.Color {
		return lipgloss.LightDark(hasDarkBG)(lipgloss.Color(light), lipgloss.Color(dark))
	}

	rawModeStyle := lipgloss.NewStyle()
	rawKindStyle := lipgloss.NewStyle().
		Width(4).
		Align(lipgloss.Right).
		Bold(true).
		MarginRight(1)
	seqStyle := lipgloss.NewStyle().
		Foreground(lightDark("#917F8B", "#978692"))
	separator := lipgloss.NewStyle().
		Foreground(lightDark("", "#978692")).
		SetString(": ")
	textStyle := lipgloss.NewStyle().
		Foreground(lightDark("#D9D9D9", "#D9D9D9"))
	errStyle := lipgloss.NewStyle().
		Foreground(lightDark("#EC6A88", "#ff5f87"))
	explanationStyle := lipgloss.NewStyle().
		Foreground(lightDark("#3C343A", "#D4CAD1"))

	kindColors := map[string]color.Color{
		"CSI":  lightDark("#936EE5", "#8D58FF"),
		"DCS":  lightDark("#86C867", "#CEE88A"),
		"OSC":  lightDark("#43C7E0", "#1CD4F7"),
		"APC":  lightDark("#F58855", "#FF8383"),
		"ESC":  lightDark("#E46FDD", "#E46FDD"),
		"Ctrl": lightDark("#4DBA94", "#4BD2A3"),
		"Text": lightDark("#978692", "#6C6068"),
	}

	kindStyle := func(kind string) lipgloss.Style {
		if raw {
			return rawModeStyle.Foreground(kindColors[kind])
		}
		return rawKindStyle.Foreground(kindColors[kind])
	}

	textKindStyle := kindStyle("Text").SetString("Text")

	seqPrint := func(kind string, seq []byte) {
		s := fmt.Sprintf("%q", seq)
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		if raw {
			_, _ = fmt.Fprint(w, kindStyle(kind).Render(s))
			return
		}

		// Trim introducers and terminators
		// CSI
		s = strings.TrimPrefix(s, "\\x9b")
		s = strings.TrimPrefix(s, "\\x1b[")
		// DCS
		s = strings.TrimPrefix(s, "\\x90")
		s = strings.TrimPrefix(s, "\\x1bP")
		// OSC
		s = strings.TrimPrefix(s, "\\x9d")
		s = strings.TrimPrefix(s, "\\x1b]")
		s = strings.TrimSuffix(s, "\\a")
		// APC
		s = strings.TrimPrefix(s, "\\x9f")
		s = strings.TrimPrefix(s, "\\x1b_")
		// ESC
		if !bytes.Equal(seq, []byte{ansi.ESC}) {
			s = strings.TrimPrefix(s, "\\x1b")
		}
		// ST
		s = strings.TrimSuffix(s, "\\x9c")
		s = strings.TrimSuffix(s, "\\x1b\\\\")

		_, _ = fmt.Fprintf(
			w,
			"%s%s%s",
			kindStyle(kind).Render(kind),
			seqStyle.Render(s),
			separator,
		)

		switch kind {
		case "Ctrl":
			_, _ = fmt.Fprintln(w, explanationStyle.Render(ctrlCodes[seq[0]]))
		case "":
			_, _ = fmt.Fprintf(w, "Unknown %q\n", seq)
		}
	}

	flushPrint := func() {
		if buf.Len() == 0 {
			return
		}
		if raw {
			_, _ = fmt.Fprint(w, kindStyle("Text").Render(buf.String()))
		} else {
			_, _ = fmt.Fprintf(w, "%s%s\n", textKindStyle, textStyle.Render(buf.String()))
		}

		buf.Reset()
	}

	handle := func(reg map[int]handlerFn, p *ansi.Parser) {
		if raw {
			return
		}

		handler, ok := reg[p.Cmd]
		if !ok {
			_, _ = fmt.Fprintln(w, errStyle.Render(errUnhandled.Error()))
			return
		}
		out, err := handler(p)
		if err != nil {
			_, _ = fmt.Fprintln(w, errStyle.Render(err.Error()))
			return
		}
		_, _ = fmt.Fprintln(w, explanationStyle.Render(out))
	}

	var state byte
	p := ansi.GetParser()
	defer ansi.PutParser(p)

	for len(in) > 0 {
		seq, width, n, newState := ansi.DecodeSequence(in, state, p)

		switch {
		case ansi.HasCsiPrefix(seq):
			flushPrint()
			seqPrint("CSI", seq)
			handle(csiHandlers, p)

		case ansi.HasDcsPrefix(seq):
			flushPrint()
			seqPrint("DCS", seq)
			handle(dcsHandlers, p)

		case ansi.HasOscPrefix(seq):
			flushPrint()
			seqPrint("OSC", seq)
			handle(oscHandlers, p)

		case ansi.HasApcPrefix(seq):
			flushPrint()
			seqPrint("APC", seq)

			switch {
			case ansi.HasPrefix(p.Data, []byte("G")):
				// TODO: Kitty graphics
			}

			_, _ = fmt.Fprintln(w)

		case ansi.HasEscPrefix(seq):
			flushPrint()

			if len(seq) == 1 {
				// just an ESC
				_, _ = fmt.Fprintf(
					w,
					"%s%s%s%s\n",
					kindStyle("Ctrl").SetString("Ctrl"),
					seqStyle.Render("ESC"),
					separator,
					explanationStyle.Render("Escape"),
				)
				break
			}

			seqPrint("ESC", seq)
			handle(escHandler, p)

		case width == 0 && len(seq) == 1:
			flushPrint()
			// control code
			seqPrint("Ctrl", seq)

		case width > 0:
			// Text
			buf.WriteString(explanationStyle.Render(string(seq)))

		default:
			flushPrint()
			seqPrint("", seq)
		}

		in = in[n:]
		state = newState
	}

	flushPrint()
	return nil
}

var ctrlCodes = map[byte]string{
	// C0
	0:  "Null",
	1:  "Start of heading",
	2:  "Start of text",
	3:  "End of text",
	4:  "End of transmission",
	5:  "Enquiry",
	6:  "Acknowledge",
	7:  "Bell",
	8:  "Backspace",
	9:  "Horizontal tab",
	10: "Line feed",
	11: "Vertical tab",
	12: "Form feed",
	13: "Carriage return",
	14: "Shift out",
	15: "Shift in",
	16: "Data link escape",
	17: "Device control 1",
	18: "Device control 2",
	19: "Device control 3",
	20: "Device control 4",
	21: "Negative acknowledge",
	22: "Synchronous idle",
	23: "End of transmission block",
	24: "Cancel",
	25: "End of medium",
	26: "Substitute",
	27: "Escape",
	28: "File separator",
	29: "Group separator",
	30: "Record separator",
	31: "Unit separator",

	// C1
	0x80: "Padding character",
	0x81: "High octet preset",
	0x82: "Break permitted here",
	0x83: "No break here",
	0x84: "Index",
	0x85: "Next line",
	0x86: "Start of selected area",
	0x87: "End of selected area",
	0x88: "Character tabulation set",
	0x89: "Character tabulation with justification",
	0x8a: "Line tabulation set",
	0x8b: "Partial line forward",
	0x8c: "Partial line backward",
	0x8d: "Reverse line feed",
	0x8e: "Single shift 2",
	0x8f: "Single shift 3",
	0x90: "Device control string",
	0x91: "Private use 1",
	0x92: "Private use 2",
	0x93: "Set transmit state",
	0x94: "Cancel character",
	0x95: "Message waiting",
	0x96: "Start of guarded area",
	0x97: "End of guarded area",
	0x98: "Start of string",
	0x99: "Single graphic character introducer",
	0x9a: "Single character introducer",
	0x9b: "Control sequence introducer",
	0x9c: "String terminator",
	0x9d: "Operating system command",
	0x9e: "Privacy message",
	0x9f: "Application program command",
}
