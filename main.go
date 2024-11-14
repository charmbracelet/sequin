package main

import (
	"bytes"
	"fmt"
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

var buf bytes.Buffer

func main() {
	if err := cmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sequin",
		Short: "Human-readable ANSI sequences",
		Args:  cobra.NoArgs,
		Example: `
printf '\x1b[m' | sequin
sequin <file
	`,
		RunE: exec,
	}
}

var (
	kindStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#C070D4")).
			Padding(0, 1).
			MaxWidth(5).
			Width(5).
			AlignHorizontal(lipgloss.Center).
			Bold(true)
	seqStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04D7D7")).
			Italic(true)
	txtStyle = lipgloss.NewStyle().
			Faint(true).
			Italic(true)
	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("204"))
)

func exec(cmd *cobra.Command, _ []string) error {
	w := colorprofile.NewWriter(cmd.OutOrStdout(), os.Environ())
	in, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return err
	}

	pseq := func(kind string, seq []byte) {
		s := seqStyle.Render(strings.ReplaceAll(fmt.Sprintf("%q", seq), `"`, ""))
		fmt.Fprintf(w, "%s %s: ", kindStyle.Render(kind), s)
	}

	flushPrint := func() {
		if buf.Len() == 0 {
			return
		}
		fmt.Fprintf(w, "%s %s\n", kindStyle.Render("TXT"), txtStyle.Render(buf.String()))
		buf.Reset()
	}

	handle := func(reg map[int]handlerFn, p *ansi.Parser) {
		handler, ok := reg[p.Cmd]
		if !ok {
			fmt.Fprintln(w, errStyle.Render(errUnhandled.Error()))
			return
		}
		out, err := handler(p)
		if err != nil {
			fmt.Fprintln(w, errStyle.Render(err.Error()))
			return
		}
		fmt.Fprintln(w, out)
	}

	var state byte
	p := ansi.GetParser()
	defer ansi.PutParser(p)

	for len(in) > 0 {
		seq, width, n, newState := ansi.DecodeSequence(in, state, p)

		switch {
		case ansi.HasCsiPrefix(seq):
			flushPrint()
			pseq("CSI", seq)
			handle(csiHandlers, p)

		case ansi.HasDcsPrefix(seq):
			flushPrint()
			pseq("DCS", seq)
			handle(dcsHandlers, p)

		case ansi.HasOscPrefix(seq):
			flushPrint()
			pseq("OSC", seq)
			handle(oscHandlers, p)

		case ansi.HasApcPrefix(seq):
			flushPrint()
			pseq("APC", seq)

			switch {
			case ansi.HasPrefix(p.Data, []byte("G")):
				// TODO: Kitty graphics
			}

			fmt.Fprintln(w)

		case ansi.HasEscPrefix(seq):
			flushPrint()

			if len(seq) == 1 {
				// just an ESC
				pseq("ESC", seq)
				fmt.Fprintln(w, "Control code ESC")
				break
			}

			pseq("ESC", seq)
			handle(escHandler, p)

		case width == 0 && len(seq) == 1:
			flushPrint()
			// control code
			pseq("CTR", seq)
			fmt.Fprintln(w, ctrlCodes[seq[0]])

		case width > 0:
			// Text
			buf.Write(seq)

		default:
			flushPrint()
			fmt.Fprintf(w, "Unknown %q\n", seq)
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
