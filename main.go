package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

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

func exec(cmd *cobra.Command, _ []string) error {
	in, err := io.ReadAll(cmd.InOrStdin())
	if err != nil {
		return err
	}

	flushPrint := func() {
		if buf.Len() == 0 {
			return
		}
		cmd.Printf("Text: %q\n", buf.String())
		buf.Reset()
	}

	handle := func(reg map[int]handlerFn, p *ansi.Parser) {
		handler, ok := reg[p.Cmd]
		if !ok {
			cmd.PrintErrln(errUnhandled)
			return
		}
		out, err := handler(p)
		if err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
		cmd.Println(out)
	}

	var state byte
	p := ansi.GetParser()
	defer ansi.PutParser(p)

	for len(in) > 0 {
		seq, width, n, newState := ansi.DecodeSequence(in, state, p)

		switch {
		case ansi.HasCsiPrefix(seq):
			flushPrint()
			cmd.Printf("CSI %q: ", seq)
			handle(csiHandlers, p)

		case ansi.HasDcsPrefix(seq):
			flushPrint()
			cmd.Printf("DCS %q: ", seq)
			handle(dcsHandlers, p)

		case ansi.HasOscPrefix(seq):
			flushPrint()
			cmd.Printf("OSC %q: ", seq)
			handle(oscHandlers, p)

		case ansi.HasApcPrefix(seq):
			flushPrint()
			cmd.Printf("APC %q", seq)

			switch {
			case ansi.HasPrefix(p.Data, []byte("G")):
				// Kitty graphics
			}

			cmd.Println()

		case ansi.HasEscPrefix(seq):
			flushPrint()

			if len(seq) == 1 {
				// just an ESC
				cmd.Println("Control code ESC")
				break
			}

			cmd.Printf("ESC: %q: ", seq)
			handle(escHandler, p)

		case width == 0 && len(seq) == 1:
			flushPrint()
			// control code
			cmd.Printf("Control code %q: %s\n", seq, ctrlCodes[seq[0]])

		case width > 0:
			// Text
			buf.Write(seq)

		default:
			flushPrint()
			cmd.Printf("Unknown %q\n", seq)
		}

		in = in[n:]
		state = newState
	}

	flushPrint()
	return nil
}

var ctrlCodes = map[byte]string{
	// C0
	0:  "null",
	1:  "start of heading",
	2:  "start of text",
	3:  "end of text",
	4:  "end of transmission",
	5:  "enquiry",
	6:  "acknowledge",
	7:  "bell",
	8:  "backspace",
	9:  "horizontal tab",
	10: "line feed",
	11: "vertical tab",
	12: "form feed",
	13: "carriage return",
	14: "shift out",
	15: "shift in",
	16: "data link escape",
	17: "device control 1",
	18: "device control 2",
	19: "device control 3",
	20: "device control 4",
	21: "negative acknowledge",
	22: "synchronous idle",
	23: "end of transmission block",
	24: "cancel",
	25: "end of medium",
	26: "substitute",
	27: "escape",
	28: "file separator",
	29: "group separator",
	30: "record separator",
	31: "unit separator",

	// C1
	0x80: "padding character",
	0x81: "high octet preset",
	0x82: "break permitted here",
	0x83: "no break here",
	0x84: "index",
	0x85: "next line",
	0x86: "start of selected area",
	0x87: "end of selected area",
	0x88: "character tabulation set",
	0x89: "character tabulation with justification",
	0x8a: "line tabulation set",
	0x8b: "partial line forward",
	0x8c: "partial line backward",
	0x8d: "reverse line feed",
	0x8e: "single shift 2",
	0x8f: "single shift 3",
	0x90: "device control string",
	0x91: "private use 1",
	0x92: "private use 2",
	0x93: "set transmit state",
	0x94: "cancel character",
	0x95: "message waiting",
	0x96: "start of guarded area",
	0x97: "end of guarded area",
	0x98: "start of string",
	0x99: "single graphic character introducer",
	0x9a: "single character introducer",
	0x9b: "control sequence introducer",
	0x9c: "string terminator",
	0x9d: "operating system command",
	0x9e: "privacy message",
	0x9f: "application program command",
}
