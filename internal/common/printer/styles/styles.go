package styles

import "github.com/charmbracelet/lipgloss"

// Colours adapt to the terminal background so output
// reads well on both light and dark themes.
var (
	ColorSuccess = lipgloss.AdaptiveColor{Light: "#00A344", Dark: "#00C853"} // green
	ColorError   = lipgloss.AdaptiveColor{Light: "#D32F2F", Dark: "#FF5252"} // red
	ColorInfo    = lipgloss.AdaptiveColor{Light: "#0288D1", Dark: "#40C4FF"} // blue
	ColorWarn    = lipgloss.AdaptiveColor{Light: "#B28704", Dark: "#FFD740"} // yellow
	ColorAccent  = lipgloss.AdaptiveColor{Light: "#A8641B", Dark: "#E0A458"} // savanna amber
	ColorText    = lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#E4E4E4"} // near black/white
	ColorMuted   = lipgloss.AdaptiveColor{Light: "#8A8A8A", Dark: "#6C6C6C"} // grey
	ColorSubtle  = lipgloss.AdaptiveColor{Light: "#5C5C5C", Dark: "#A8A8A8"} // light grey
)

var (
	Base = lipgloss.NewStyle()

	Success = Base.Foreground(ColorSuccess).Bold(true)
	Error   = Base.Foreground(ColorError).Bold(true)
	Info    = Base.Foreground(ColorInfo).Bold(true)
	Warn    = Base.Foreground(ColorWarn).Bold(true)

	// Muted is used for hints, footers and empty states.
	Muted = Base.Foreground(ColorMuted)

	// Title is the heading printed at the top of a command's output.
	Title = lipgloss.NewStyle().
		Foreground(ColorAccent).
		Bold(true)

	// Label is the key style used in key/value views.
	Label = lipgloss.NewStyle().
		Foreground(ColorSubtle).
		Bold(true)

	// Value is the value style used in key/value views.
	Value = lipgloss.NewStyle().
		Foreground(ColorText)

	Header = lipgloss.NewStyle().
		Foreground(ColorText).
		Padding(0, 1).
		Bold(true)

	// Box is a rounded card used for summaries and detail headers.
	Box = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(0, 1)
)

// Icons
const (
	IconSuccess = "✔"
	IconError   = "✖"
	IconInfo    = "ℹ"
	IconWarn    = "⚠"
	IconTitle   = "▍"
	IconArrow   = "→"
	IconDot     = "·"
)
