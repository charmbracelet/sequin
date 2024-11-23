package main

import (
	"image/color"
	"strings"

	"github.com/charmbracelet/lipgloss/v2"
)

type theme struct {
	IsRaw bool

	raw         lipgloss.Style
	kind        lipgloss.Style
	sequence    lipgloss.Style
	separator   lipgloss.Style
	text        lipgloss.Style
	error       lipgloss.Style
	explanation lipgloss.Style

	kindColors struct {
		apc, csi, ctrl, dcs, esc, osc, pm, sos, text color.Color
	}
}

func (t theme) kindStyle(kind string) lipgloss.Style {
	kind = strings.ToLower(kind)
	base := t.kind
	if t.IsRaw {
		base = t.raw
	}

	s := map[string]lipgloss.Style{
		"apc":  base.Foreground(t.kindColors.apc),
		"csi":  base.Foreground(t.kindColors.csi),
		"ctrl": base.Foreground(t.kindColors.ctrl),
		"dcs":  base.Foreground(t.kindColors.dcs),
		"esc":  base.Foreground(t.kindColors.esc),
		"osc":  base.Foreground(t.kindColors.osc),
		"pm":   base.Foreground(t.kindColors.pm),
		"sos":  base.Foreground(t.kindColors.sos),
		"text": base.Foreground(t.kindColors.text),
	}[kind]

	if t.IsRaw {
		return s
	}

	switch kind {
	case "csi":
		return s.SetString("CSI")
	case "dcs":
		return s.SetString("DCS")
	case "osc":
		return s.SetString("OSC")
	case "apc":
		return s.SetString("APC")
	case "pm":
		return s.SetString("PM")
	case "sos":
		return s.SetString("SOS")
	case "esc":
		return s.SetString("ESC")
	case "ctrl":
		return s.SetString("Ctrl")
	case "text":
		return s.SetString("Text")
	default:
		return s
	}
}

//nolint:mnd
func charmTheme(hasDarkBG bool) (t theme) {
	lightDark := func(light, dark string) color.Color {
		return lipgloss.LightDark(hasDarkBG)(lipgloss.Color(light), lipgloss.Color(dark))
	}

	t.raw = lipgloss.NewStyle()
	t.kind = lipgloss.NewStyle().
		Width(4).
		Align(lipgloss.Right).
		Bold(true).
		MarginRight(1)
	t.sequence = lipgloss.NewStyle().
		Foreground(lightDark("#917F8B", "#978692"))
	t.separator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#978692")).
		SetString(": ")
	t.text = lipgloss.NewStyle().
		Foreground(lightDark("#D9D9D9", "#D9D9D9"))
	t.error = lipgloss.NewStyle().
		Foreground(lightDark("#EC6A88", "#ff5f87"))
	t.explanation = lipgloss.NewStyle().
		Foreground(lightDark("#3C343A", "#D4CAD1"))

	t.kindColors.apc = lightDark("#F58855", "#FF8383")
	t.kindColors.csi = lightDark("#936EE5", "#8D58FF")
	t.kindColors.ctrl = lightDark("#4DBA94", "#4BD2A3")
	t.kindColors.dcs = lightDark("#86C867", "#CEE88A")
	t.kindColors.esc = lipgloss.Color("#E46FDD")
	t.kindColors.osc = lightDark("#43C7E0", "#1CD4F7")
	t.kindColors.pm = lightDark("#FF8383", "#DC7272")
	t.kindColors.sos = lightDark("#978692", "#6C6068")
	t.kindColors.text = lightDark("#978692", "#6C6068")

	return t
}

func base16Theme(_ bool) theme {
	t := charmTheme(false)

	t.sequence = t.sequence.Foreground(lipgloss.BrightBlack)
	t.separator = t.separator.Foreground(lipgloss.BrightBlack)
	t.text = t.text.Foreground(lipgloss.BrightBlack)
	t.error = t.error.Foreground(lipgloss.BrightRed)
	t.explanation = t.explanation.Foreground(lipgloss.White)

	t.kindColors.apc = lipgloss.Red
	t.kindColors.csi = lipgloss.Blue
	t.kindColors.ctrl = lipgloss.Green
	t.kindColors.dcs = lipgloss.Yellow
	t.kindColors.esc = lipgloss.Magenta
	t.kindColors.osc = lipgloss.Cyan
	t.kindColors.pm = lipgloss.BrightRed
	t.kindColors.sos = lipgloss.BrightBlack
	t.kindColors.text = lipgloss.BrightBlack

	return t
}
