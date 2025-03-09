package main

import (
	"bytes"
	"cmp"
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
	markerShift   = parser.PrefixShift
	intermedShift = parser.IntermedShift
)

var (
	buf bytes.Buffer
	raw bool
	shouldPrintMnemonics bool
	// Version as provided by goreleaser.
	Version = ""
)

func main() {
	if err := cmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func cmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "sequin",
		Short: "Human-readable ANSI sequences",
		Args:  cobra.ArbitraryArgs,
		Example: `
printf '\x1b[m' | sequin
sequin <file
sequin -- some command to execute
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			w := colorprofile.NewWriter(cmd.OutOrStdout(), os.Environ())
			var in []byte
			var err error
			if len(args) == 0 {
				in, err = io.ReadAll(cmd.InOrStdin())
			} else {
				in, err = executeCommand(cmd.Context(), args)
			}
			if err != nil {
				return err
			}
			return process(w, in)
		},
	}
	root.Flags().BoolVarP(&raw, "raw", "r", false, "raw mode (no explanation)")
	root.Flags().BoolVarP(
		&shouldPrintMnemonics,
		"mnemonics", "m",
		false,
		`print mnemonics for escape sequence types (ignored in raw mode)`,
	)

	// make sure to disable any additional info if we enter the raw mode
	if raw {
		shouldPrintMnemonics = false
	}

	root.Version = cmp.Or(Version, "unknown (built from source)")
	return root
}

func process(w *colorprofile.Writer, in []byte) error {
	var t theme
	switch strings.ToLower(os.Getenv("SEQUIN_THEME")) {
	case "ansi", "carlos", "secret_carlos", "matchy":
		t = base16Theme(false)
	default:
		hasDarkBG := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
		t = charmTheme(hasDarkBG)
	}

	t.IsRaw = raw

	seqPrint := func(kind string, seq []byte) {
		s := fmt.Sprintf("%q", seq)
		s = strings.TrimPrefix(s, `"`)
		s = strings.TrimSuffix(s, `"`)
		if raw {
			_, _ = fmt.Fprint(w, t.kindStyle(kind).Render(s))
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
		// BEL
		if !bytes.Equal(seq, []byte{ansi.BEL}) {
			// Remove only if not a literal bell
			s = strings.TrimSuffix(s, "\\a")
		}
		// SOS
		s = strings.TrimPrefix(s, "\\x98")
		s = strings.TrimPrefix(s, "\\x1bX")
		// PM
		s = strings.TrimPrefix(s, "\\x9e")
		s = strings.TrimPrefix(s, "\\x1b^")
		// APC
		s = strings.TrimPrefix(s, "\\x9f")
		s = strings.TrimPrefix(s, "\\x1b_")
		// ESC
		if !bytes.Equal(seq, []byte{ansi.ESC}) {
			// Remove only if not a standalone ESC
			s = strings.TrimPrefix(s, "\\x1b")
		}
		// ST
		if !bytes.Equal(seq, []byte{ansi.ST}) {
			// Remove only if accompanied by a sequence introducer
			s = strings.TrimSuffix(s, "\\x9c")
		}
		s = strings.TrimSuffix(s, "\\x1b\\\\")

		_, _ = fmt.Fprintf(w, "%s", t.kindStyle(kind))

		switch kind {
		case "Ctrl":
			additionalInfo := ""

			if shouldPrintMnemonics {
				additionalInfo = t.mnemonic.Render("<"+ctrlCodes[seq[0]].mnemonic+"> ")
			}

			_, _ = fmt.Fprintf(
				w,
				"%s%s%s%s\n",
				t.sequence.Render(s),
				t.separator,
				additionalInfo,
				t.explanation.Render(ctrlCodes[seq[0]].explanation),
			)

		case "PM":
			_, _ = fmt.Fprintf(
				w,
				"%s%s\n",
				t.separator,
				t.explanation.Render(fmt.Sprintf("Privacy message %q", s)),
			)

		case "SOS":
			_, _ = fmt.Fprintf(
				w,
				"%s%s\n",
				t.separator,
				t.explanation.Render(fmt.Sprintf("Control string %q", s)),
			)

		case "":
			_, _ = fmt.Fprintf(
				w,
				"%s%sUnknown %q\n",
				t.sequence.Render(s),
				t.separator,
				seq,
			)

		default:
			// For sequences with own handlers print only the kind and the sequence.
			// Explanation will be provided by handlers themselves:
			_, _ = fmt.Fprintf(
				w,
				"%s%s",
				t.sequence.Render(s),
				t.separator,
			)
		}
	}

	flushPrint := func() {
		if buf.Len() == 0 {
			return
		}
		if raw {
			_, _ = fmt.Fprint(w, t.kindStyle("Text").Render(buf.String()))
		} else {
			_, _ = fmt.Fprintf(w, "%s%s\n", t.kindStyle("text"), t.text.Render(buf.String()))
		}

		buf.Reset()
	}

	handle := func(reg map[int]handlerFn, p *ansi.Parser) {
		if raw {
			return
		}

		handler, ok := reg[p.Command()]
		if !ok {
			_, _ = fmt.Fprintln(w, t.error.Render(errUnhandled.Error()))
			return
		}
		out, err := handler(p)
		if err != nil {
			_, _ = fmt.Fprintln(w, t.error.Render(err.Error()))
			return
		}
		if out.mnemonic != "" && shouldPrintMnemonics {
			_, _ = fmt.Fprintln(
				w,
				t.mnemonic.Render("<" + out.mnemonic + ">"),
				t.explanation.Render(out.explanation),
			)
			return
		}
		_, _ = fmt.Fprintln(w, t.explanation.Render(out.explanation))
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

		case ansi.HasPmPrefix(seq):
			flushPrint()
			seqPrint("PM", seq)

		case ansi.HasSosPrefix(seq):
			flushPrint()
			seqPrint("SOS", seq)

		case ansi.HasApcPrefix(seq):
			flushPrint()
			seqPrint("APC", seq)

			switch {
			case ansi.HasPrefix(p.Data(), []byte("G")):
				// TODO: Kitty graphics
			}

			_, _ = fmt.Fprintln(w)

		case ansi.HasEscPrefix(seq):
			flushPrint()

			if len(seq) == 1 && !raw {
				// just an ESC
				_, _ = fmt.Fprintf(
					w,
					"%s%s%s%s\n",
					t.kindStyle("Ctrl"),
					t.sequence.Render("ESC"),
					t.separator,
					t.explanation.Render("Escape"),
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
			buf.WriteString(t.explanation.Render(string(seq)))

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

var ctrlCodes = map[byte]seqInfo{
	// C0
	0:  seqInfo{"NUL", "Null"},
	1:  seqInfo{"SOH", "Start of heading"},
	2:  seqInfo{"STX", "Start of text"},
	3:  seqInfo{"ETX", "End of text"},
	4:  seqInfo{"EOT", "End of transmission"},
	5:  seqInfo{"ENQ", "Enquiry"},
	6:  seqInfo{"ACK", "Acknowledge"},
	7:  seqInfo{"BEL", "Bell"},
	8:  seqInfo{"BS", "Backspace"},
	9:  seqInfo{"HT", "Horizontal tab"},
	10: seqInfo{"LF", "Line feed"},
	11: seqInfo{"VT", "Vertical tab"},
	12: seqInfo{"FF", "Form feed"},
	13: seqInfo{"CR", "Carriage return"},
	14: seqInfo{"SO", "Shift out"},
	15: seqInfo{"SI", "Shift in"},
	16: seqInfo{"DLE", "Data link escape"},
	17: seqInfo{"DC1", "Device control 1"},
	18: seqInfo{"DC2", "Device control 2"},
	19: seqInfo{"DC3", "Device control 3"},
	20: seqInfo{"DC4", "Device control 4"},
	21: seqInfo{"NAK", "Negative acknowledge"},
	22: seqInfo{"SYN", "Synchronous idle"},
	23: seqInfo{"ETB", "End of transmission block"},
	24: seqInfo{"CAN", "Cancel"},
	25: seqInfo{"EM", "End of medium"},
	26: seqInfo{"SUB", "Substitute"},
	27: seqInfo{"ESC", "Escape"},
	28: seqInfo{"FS", "File separator"},
	29: seqInfo{"GS", "Group separator"},
	30: seqInfo{"RS", "Record separator"},
	31: seqInfo{"US", "Unit separator"},

	// RFC 20, section 4.1 "Control Characters" includes DEL with the note:
	// "In the strict sense, DEL is not a control character."
	127: seqInfo{"DEL", "Delete"},

	// C1
	0x80: seqInfo{"PAD", "Padding character"},
	0x81: seqInfo{"HOP", "High octet preset"},
	0x82: seqInfo{"BPH", "Break permitted here"},
	0x83: seqInfo{"NBH", "No break here"},
	0x84: seqInfo{"IND", "Index"},
	0x85: seqInfo{"NEL", "Next line"},
	0x86: seqInfo{"SSA", "Start of selected area"},
	0x87: seqInfo{"ESA", "End of selected area"},
	0x88: seqInfo{"HTS", "Character tabulation set"},
	0x89: seqInfo{"HTJ", "Character tabulation with justification"},
	0x8a: seqInfo{"VTS", "Line tabulation set"},
	0x8b: seqInfo{"PLD", "Partial line forward"},
	0x8c: seqInfo{"PLU", "Partial line backward"},
	0x8d: seqInfo{"RI", "Reverse line feed"},
	0x8e: seqInfo{"SS2", "Single shift 2"},
	0x8f: seqInfo{"SS3", "Single shift 3"},
	0x90: seqInfo{"DCS", "Device control string"},
	0x91: seqInfo{"PU1", "Private use 1"},
	0x92: seqInfo{"PU2", "Private use 2"},
	0x93: seqInfo{"STS", "Set transmit state"},
	0x94: seqInfo{"CCH", "Cancel character"},
	0x95: seqInfo{"MW", "Message waiting"},
	0x96: seqInfo{"SPA", "Start of guarded area"},
	0x97: seqInfo{"EPA", "End of guarded area"},
	0x98: seqInfo{"SOS", "Start of string"},
	0x99: seqInfo{"SGCI", "Single graphic character introducer"},
	0x9a: seqInfo{"SCI", "Single character introducer"},
	0x9b: seqInfo{"CSI", "Control sequence introducer"},
	0x9c: seqInfo{"ST", "String terminator"},
	0x9d: seqInfo{"OSC", "Operating system command"},
	0x9e: seqInfo{"PM", "Privacy message"},
	0x9f: seqInfo{"APC", "Application program command"},
}
