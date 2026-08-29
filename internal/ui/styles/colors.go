package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines a complete UI color scheme
type Theme struct {
	Name      string
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Danger    lipgloss.Color
	Info      lipgloss.Color
	Muted     lipgloss.Color
	Subtle    lipgloss.Color
	Bg        lipgloss.Color
	BgAlt     lipgloss.Color
	Highlight lipgloss.Color
	Text      lipgloss.Color
	TextBold  lipgloss.Color
}

// Preset themes
var (
	ThemeCyberpunk = Theme{
		Name:      "Cyberpunk",
		Primary:   lipgloss.Color("#7D56F4"),
		Secondary: lipgloss.Color("#00D2FF"),
		Success:   lipgloss.Color("#10B981"),
		Warning:   lipgloss.Color("#F59E0B"),
		Danger:    lipgloss.Color("#EF4444"),
		Info:      lipgloss.Color("#3B82F6"),
		Muted:     lipgloss.Color("#6B7280"),
		Subtle:    lipgloss.Color("#374151"),
		Bg:        lipgloss.Color("#1E1E2E"),
		BgAlt:     lipgloss.Color("#181825"),
		Highlight: lipgloss.Color("#2E2E3E"),
		Text:      lipgloss.Color("#CDD6F4"),
		TextBold:  lipgloss.Color("#FFFFFF"),
	}

	ThemeCatppuccin = Theme{
		Name:      "Catppuccin",
		Primary:   lipgloss.Color("#CBA6F7"), // Mauve
		Secondary: lipgloss.Color("#89DCEB"), // Sky
		Success:   lipgloss.Color("#A6E3A1"), // Green
		Warning:   lipgloss.Color("#F9E2AF"), // Yellow
		Danger:    lipgloss.Color("#F38BA8"), // Red
		Info:      lipgloss.Color("#89B4FA"), // Blue
		Muted:     lipgloss.Color("#7F849C"),
		Subtle:    lipgloss.Color("#45475A"),
		Bg:        lipgloss.Color("#1E1E2E"),
		BgAlt:     lipgloss.Color("#11111B"),
		Highlight: lipgloss.Color("#313244"),
		Text:      lipgloss.Color("#CDD6F4"),
		TextBold:  lipgloss.Color("#FFFFFF"),
	}

	ThemeDracula = Theme{
		Name:      "Dracula",
		Primary:   lipgloss.Color("#BD93F9"), // Purple
		Secondary: lipgloss.Color("#FF79C6"), // Pink
		Success:   lipgloss.Color("#50FA7B"), // Green
		Warning:   lipgloss.Color("#FFB86C"), // Orange
		Danger:    lipgloss.Color("#FF5555"), // Red
		Info:      lipgloss.Color("#8BE9FD"), // Cyan
		Muted:     lipgloss.Color("#6272A4"),
		Subtle:    lipgloss.Color("#44475A"),
		Bg:        lipgloss.Color("#282A36"),
		BgAlt:     lipgloss.Color("#1E1F29"),
		Highlight: lipgloss.Color("#44475A"),
		Text:      lipgloss.Color("#F8F8F2"),
		TextBold:  lipgloss.Color("#FFFFFF"),
	}

	ThemeNord = Theme{
		Name:      "Nord",
		Primary:   lipgloss.Color("#88C0D0"), // Frost
		Secondary: lipgloss.Color("#81A1C1"), // Blue
		Success:   lipgloss.Color("#A3BE8C"), // Green
		Warning:   lipgloss.Color("#EBCB8B"), // Yellow
		Danger:    lipgloss.Color("#BF616A"), // Red
		Info:      lipgloss.Color("#5E81AC"), // Deep Blue
		Muted:     lipgloss.Color("#616E88"),
		Subtle:    lipgloss.Color("#4C566A"),
		Bg:        lipgloss.Color("#2E3440"),
		BgAlt:     lipgloss.Color("#242933"),
		Highlight: lipgloss.Color("#3B4252"),
		Text:      lipgloss.Color("#ECEFF4"),
		TextBold:  lipgloss.Color("#FFFFFF"),
	}

	ThemeMatrix = Theme{
		Name:      "Matrix",
		Primary:   lipgloss.Color("#00FF66"), // Neon Green
		Secondary: lipgloss.Color("#00CC44"), // Mint
		Success:   lipgloss.Color("#33FF33"),
		Warning:   lipgloss.Color("#FFFF00"),
		Danger:    lipgloss.Color("#FF3333"),
		Info:      lipgloss.Color("#00FFFF"),
		Muted:     lipgloss.Color("#4A7C59"),
		Subtle:    lipgloss.Color("#1B3B22"),
		Bg:        lipgloss.Color("#0D140F"),
		BgAlt:     lipgloss.Color("#080D0A"),
		Highlight: lipgloss.Color("#1A2E20"),
		Text:      lipgloss.Color("#D0F0D0"),
		TextBold:  lipgloss.Color("#FFFFFF"),
	}

	Themes = []Theme{ThemeCyberpunk, ThemeCatppuccin, ThemeDracula, ThemeNord, ThemeMatrix}
	activeThemeIndex = 0
)

// Active color variables referencing the currently active theme
var (
	ColorPrimary   = ThemeCyberpunk.Primary
	ColorSecondary = ThemeCyberpunk.Secondary
	ColorSuccess   = ThemeCyberpunk.Success
	ColorWarning   = ThemeCyberpunk.Warning
	ColorDanger    = ThemeCyberpunk.Danger
	ColorInfo      = ThemeCyberpunk.Info
	ColorMuted     = ThemeCyberpunk.Muted
	ColorSubtle    = ThemeCyberpunk.Subtle
	ColorBg        = ThemeCyberpunk.Bg
	ColorBgAlt     = ThemeCyberpunk.BgAlt
	ColorHighlight = ThemeCyberpunk.Highlight
	ColorText      = ThemeCyberpunk.Text
	ColorTextBold  = ThemeCyberpunk.TextBold
)

// CycleTheme switches to the next color theme and refreshes all styles
func CycleTheme() string {
	activeThemeIndex = (activeThemeIndex + 1) % len(Themes)
	ApplyTheme(Themes[activeThemeIndex])
	return Themes[activeThemeIndex].Name
}

// CurrentThemeName returns the name of the active theme
func CurrentThemeName() string {
	return Themes[activeThemeIndex].Name
}

// ApplyTheme applies a theme palette and rebuilds styles
func ApplyTheme(t Theme) {
	for i, theme := range Themes {
		if theme.Name == t.Name {
			activeThemeIndex = i
			break
		}
	}

	ColorPrimary = t.Primary
	ColorSecondary = t.Secondary
	ColorSuccess = t.Success
	ColorWarning = t.Warning
	ColorDanger = t.Danger
	ColorInfo = t.Info
	ColorMuted = t.Muted
	ColorSubtle = t.Subtle
	ColorBg = t.Bg
	ColorBgAlt = t.BgAlt
	ColorHighlight = t.Highlight
	ColorText = t.Text
	ColorTextBold = t.TextBold

	RebuildStyles()
}
