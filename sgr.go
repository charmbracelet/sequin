package main

import (
	"fmt"
	"image/color"

	"github.com/charmbracelet/x/ansi"
)

func handleSgr(parser *ansi.Parser) {
	var str string
	if parser.ParamsLen == 0 {
		str = "Reset style"
	}

	for i := 0; i < parser.ParamsLen; i++ {
		param := ansi.Param(parser.Params[i])
		if i > 0 {
			str += ", "
		}

		// TODO: add more parameters and options
		switch param.Param() {
		case 0:
			str = "Reset style"
		case 1:
			str += "Bold"
		case 2:
			str += "Faint"
		case 3:
			str += "Italic"
		case 4:
			str += "Underline"
			if param.HasMore() {
				// Handle underline styles
				switch next := ansi.Param(parser.Params[i+1]); next.Param() {
				case 0:
					str += " (Disable)"
				case 1:
					str += " (Single)"
				case 2:
					str += " (Double)"
				case 3:
					str += " (Curly)"
				case 4:
					str += " (Dotted)"
				case 5:
					str += " (Dashed)"
				}
			}
		case 5:
			str += "Blink"
		case 7:
			str += "Inverse"
		case 8:
			str += "Invisible"
		case 9:
			str += "Crossed-out"
		case 30, 31, 32, 33, 34, 35, 36, 37:
			str += fmt.Sprintf("Foreground color: %s", basicColors[int(param)-30])
		case 38:
			str += fmt.Sprintf("Foreground color: %d", readColor(&i, parser.Params))
		case 39:
			str += "Default foreground color"
		case 40, 41, 42, 43, 44, 45, 46, 47:
			str += fmt.Sprintf("Background color: %s", basicColors[int(param)-40])
		case 48:
			str += fmt.Sprintf("Background color: %d", readColor(&i, parser.Params))
		case 49:
			str += "Default background color"
		case 90, 91, 92, 93, 94, 95, 96, 97:
			str += fmt.Sprintf("Bright foreground color: %s", basicColors[int(param)-90])
		case 100, 101, 102, 103, 104, 105, 106, 107:
			str += fmt.Sprintf("Bright background color: %s", basicColors[int(param)-100])
		}
	}

	fmt.Print(str)
}

var basicColors = map[int]string{
	0: "Black",
	1: "Red",
	2: "Green",
	3: "Yellow",
	4: "Blue",
	5: "Magenta",
	6: "Cyan",
	7: "White",
}

func readColor(idxp *int, params []int) (c ansi.Color) {
	i := *idxp
	paramsLen := len(params)
	if i > paramsLen-1 {
		return
	}
	// Note: we accept both main and subparams here
	switch param := ansi.Param(params[i+1]); param.Param() {
	case 2: // RGB
		if i > paramsLen-4 {
			return
		}
		c = color.RGBA{
			R: uint8(ansi.Param(params[i+2]).Param()), //nolint:gosec
			G: uint8(ansi.Param(params[i+3]).Param()), //nolint:gosec
			B: uint8(ansi.Param(params[i+4]).Param()), //nolint:gosec
			A: 0xff,
		}
		*idxp += 4
	case 5: // 256 colors
		if i > paramsLen-2 {
			return
		}
		c = ansi.ExtendedColor(ansi.Param(params[i+2]).Param()) //nolint:gosec
		*idxp += 2
	}
	return
}
