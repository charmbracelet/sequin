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
		csi, dcs, osc, apc, esc, ctrl, text color.Color
	}
}

func (t theme) kindStyle(kind string) lipgloss.Style {
	kind = strings.ToLower(kind)
	base := t.kind
	if t.IsRaw {
		base = t.raw
	}

	s := map[string]lipgloss.Style{
		"csi":  base.Foreground(t.kindColors.csi),
		"dcs":  base.Foreground(t.kindColors.dcs),
		"osc":  base.Foreground(t.kindColors.osc),
		"apc":  base.Foreground(t.kindColors.apc),
		"esc":  base.Foreground(t.kindColors.esc),
		"ctrl": base.Foreground(t.kindColors.ctrl),
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

func defaultTheme(hasDarkBG bool) (t theme) {
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

	t.kindColors.csi = lightDark("#936EE5", "#8D58FF")
	t.kindColors.dcs = lightDark("#86C867", "#CEE88A")
	t.kindColors.osc = lightDark("#43C7E0", "#1CD4F7")
	t.kindColors.apc = lightDark("#F58855", "#FF8383")
	t.kindColors.esc = lipgloss.Color("#E46FDD")
	t.kindColors.ctrl = lightDark("#4DBA94", "#4BD2A3")
	t.kindColors.text = lightDark("#978692", "#6C6068")

	return t
}
