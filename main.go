package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/ansi/parser"
)

const (
	markerShift   = parser.MarkerShift
	intermedShift = parser.IntermedShift
)

var buf bytes.Buffer

func flushPrint() {
	if buf.Len() == 0 {
		return
	}
	fmt.Printf("Text: %q\n", buf.String())
	buf.Reset()
}

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
			flushPrint()
			fmt.Printf("CSI %q: ", seq)

			handler, ok := csiHandlers[p.Cmd]
			if ok {
				handler(p)
			}

			fmt.Println()

		case ansi.HasDcsPrefix(seq):
			flushPrint()

		case ansi.HasOscPrefix(seq):
			flushPrint()
			fmt.Printf("OSC %q: ", seq)
			handler, ok := oscHandlers[p.Cmd]
			if ok {
				handler(p)
			}

			fmt.Println()

		case ansi.HasApcPrefix(seq):
			flushPrint()
			fmt.Printf("APC %q", seq)

			switch {
			case ansi.HasPrefix(p.Data, []byte("G")):
				// Kitty graphics
			}

			fmt.Println()

		case ansi.HasEscPrefix(seq):
			flushPrint()

			if len(seq) == 1 {
				// just an ESC
				fmt.Println("Control code ESC")
				break
			}

			fmt.Printf("ESC: %q", seq)
			switch p.Cmd {
			case 7:
				// save cursor
			case 8:
				// restore cursor
			}

			fmt.Println()

		default:
			if width == 0 && len(seq) == 1 {
				flushPrint()
				// control code
				fmt.Printf("Control code %q\n", seq)
				break
			} else if width == 0 {
				flushPrint()
				fmt.Printf("Unknown %q\n", seq)
				break
			}

			buf.WriteString(seq)
		}

		in = in[n:]
		state = newState
	}

	flushPrint()
}
