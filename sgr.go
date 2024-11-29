package main

import (
	"fmt"
	"strings"
	"image/color"

	"github.com/charmbracelet/x/ansi"
)

//nolint:mnd
func handleSgr(p *ansi.Parser) (string, error) { //nolint:unparam
	params := p.Params()
	if len(params) == 0 {
		return "Reset style", nil
	}

	var str string
	var comma bool
	for i := 0; i < len(params); i++ {
		param := params[i]
		if comma {
			str += ", "
		}
		comma = true

		// TODO: add more parameters and options
		switch param.Param(0) {
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
				next := params[i+1]
				switch p := next.Param(0); p {
				case 1, 2, 3, 4, 5:
					i++
					switch p {
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
			}
		case 5, 6:
			str += "Blink"
		case 7:
			str += "Inverse"
		case 8:
			str += "Invisible"
		case 9:
			str += "Crossed-out"
		case 21:
			str += "No bold"
		case 22:
			str += "Normal intensity"
		case 23:
			str += "No italic"
		case 24:
			str += "No underline"
		case 25:
			str += "No blink"
		case 27:
			str += "No reverse"
		case 28:
			str += "No conceal"
		case 29:
			str += "No crossed-out"
		case 30, 31, 32, 33, 34, 35, 36, 37:
			str += fmt.Sprintf("ANSI foreground color: %s", basicColors[param.Param(0)-30])
		case 38:
			c := readColor(&i, params)
			str += fmt.Sprintf("%s foreground color: %s", getColorType(c), getColorLabel(c))
		case 39:
			str += "Default foreground color"
		case 40, 41, 42, 43, 44, 45, 46, 47:
			str += fmt.Sprintf("ANSI background color: %s", basicColors[param.Param(0)-40])
		case 48:
			c := readColor(&i, params)
			str += fmt.Sprintf("%s background color: %s", getColorType(c), getColorLabel(c))
		case 49:
			str += "Default background color"
		case 58:
			c := readColor(&i, params)
			str += fmt.Sprintf("%s underline color: %s", getColorType(c), getColorLabel(c))
		case 59:
			str += "Default underline color"
		case 90, 91, 92, 93, 94, 95, 96, 97:
			str += fmt.Sprintf("Bright ANSI foreground color: %s", basicColors[param.Param(0)-90])
		case 100, 101, 102, 103, 104, 105, 106, 107:
			str += fmt.Sprintf("Bright ANSI background color: %s", basicColors[param.Param(0)-100])
		}
	}

	return str, nil
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

func getColorLabel(c ansi.Color) string {
	r, g, b, _ := c.RGBA()
	hexString := fmt.Sprintf("#%.2X%.2X%.2X", r >> 8, g >> 8, b >> 8)
	switch c := c.(type) {
	case ansi.ExtendedColor:
		paletteID := int(c)
		// First 16 colors are the same as ANSI
		if paletteID < 8 {
			return fmt.Sprintf("%d (%s)", c, basicColors[paletteID])
		} else if paletteID >= 8 && paletteID < 16 {
			return fmt.Sprintf("%d (Bright %s)", c, strings.ToLower(basicColors[paletteID-8]))
		}
		return fmt.Sprintf("%d (%s)", c, hexString)
	case ansi.TrueColor, color.Color:
		return hexString
	default:
		return "Unknown color"
	}
}

func getColorType(c ansi.Color) string {
	switch c.(type) {
	case ansi.BasicColor:
		return "ANSI"
	case ansi.ExtendedColor:
		return "ANSI256"
	case ansi.TrueColor, color.Color:
		return "24-bit RGB"
	default:
		return "Unknown"
	}
}

//nolint:mnd
func readColor(idxp *int, params []ansi.Parameter) (c ansi.Color) {
	i := *idxp
	paramsLen := len(params)
	if i > paramsLen-1 {
		return
	}
	// Note: we accept both main and subparams here
	switch param := params[i+1]; param.Param(0) {
	case 2: // 24-bit RGB (truecolor, also called direct color)
		if i > paramsLen-4 {
			return
		}
		c = color.RGBA{
			R: uint8(params[i+2].Param(0)), //nolint:gosec
			G: uint8(params[i+3].Param(0)), //nolint:gosec
			B: uint8(params[i+4].Param(0)), //nolint:gosec
			A: 0xff,
		}
		*idxp += 4
	case 5: // ANSI256: Palette of 256 colors
		if i > paramsLen-2 {
			return
		}
		c = ansi.ExtendedColor(params[i+2].Param(0)) //nolint:gosec
		*idxp += 2
	}
	return
}
