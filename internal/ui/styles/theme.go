package styles

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/charmbracelet/lipgloss"
)

var (
	HeaderStyle            lipgloss.Style
	HeaderSubStyle         lipgloss.Style
	HeaderVersionStyle     lipgloss.Style
	TableBoxStyle          lipgloss.Style
	TableHeaderStyle       lipgloss.Style
	TableRowStyle          lipgloss.Style
	TableSelectedRowStyle  lipgloss.Style
	CursorStyle            lipgloss.Style
	StatusDeviceStyle      lipgloss.Style
	StatusUnauthorizedStyle lipgloss.Style
	StatusOfflineStyle     lipgloss.Style
	StatusOtherStyle       lipgloss.Style
	DetailCardStyle        lipgloss.Style
	DetailTitleStyle       lipgloss.Style
	DetailKeyStyle         lipgloss.Style
	DetailValueStyle       lipgloss.Style
	EmptyBoxStyle          lipgloss.Style
	EmptyTitleStyle        lipgloss.Style
	EmptyTipStyle          lipgloss.Style
	StatusBarStyle         lipgloss.Style
	StatusMsgSuccessStyle  lipgloss.Style
	StatusMsgErrorStyle    lipgloss.Style
	HelpKeyStyle           lipgloss.Style
	HelpDescStyle          lipgloss.Style
	ModalBoxStyle          lipgloss.Style
	ModalTitleStyle        lipgloss.Style
)

func init() {
	RebuildStyles()
}

// RebuildStyles refreshes all Lip Gloss styles with current theme colors
func RebuildStyles() {
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorTextBold).
		Background(ColorPrimary).
		Padding(0, 1)

	HeaderSubStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	HeaderVersionStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorMuted)

	TableBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 1).
		MarginBottom(1)

	TableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottomForeground(ColorSubtle).
		PaddingBottom(0)

	TableRowStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Padding(0, 0)

	TableSelectedRowStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorTextBold).
		Background(ColorHighlight).
		Padding(0, 0)

	CursorStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	StatusDeviceStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSuccess)

	StatusUnauthorizedStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWarning)

	StatusOfflineStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorDanger)

	StatusOtherStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	DetailCardStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Padding(0, 1).
		MarginBottom(1)

	DetailTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		MarginBottom(0)

	DetailKeyStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorMuted).
		Width(20)

	DetailValueStyle = lipgloss.NewStyle().
		Foreground(ColorTextBold)

	EmptyBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorMuted).
		Padding(1, 2).
		Align(lipgloss.Center).
		MarginBottom(1)

	EmptyTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorWarning).
		MarginBottom(1)

	EmptyTipStyle = lipgloss.NewStyle().
		Foreground(ColorMuted)

	StatusBarStyle = lipgloss.NewStyle().
		Background(ColorBgAlt).
		Foreground(ColorText).
		Padding(0, 1)

	StatusMsgSuccessStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)

	StatusMsgErrorStyle = lipgloss.NewStyle().
		Foreground(ColorDanger).
		Bold(true)

	HelpKeyStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary)

	HelpDescStyle = lipgloss.NewStyle().
		Foreground(ColorMuted)

	ModalBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		MarginBottom(1)

	ModalTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorSecondary).
		MarginBottom(1)
}

// RenderStatusBadge formats and styles the device state indicator badge
func RenderStatusBadge(state device.DeviceState, l *locale.Strings) string {
	var label string
	switch state {
	case device.StateDevice:
		label = l.StateConnected
	case device.StateUnauthorized:
		label = l.StateUnauthorized
	case device.StateOffline:
		label = l.StateOffline
	case device.StateAuthorizing:
		label = l.StateAuthorizing
	case device.StateRecovery:
		label = l.StateRecovery
	case device.StateSideload:
		label = l.StateSideload
	case device.StateBootloader:
		label = l.StateBootloader
	default:
		label = string(state)
	}
	upper := strings.ToUpper(label)

	switch state {
	case device.StateDevice:
		return StatusDeviceStyle.Render("● " + upper)
	case device.StateUnauthorized, device.StateAuthorizing:
		return StatusUnauthorizedStyle.Render("● " + upper)
	case device.StateOffline:
		return StatusOfflineStyle.Render("● " + upper)
	default:
		return StatusOtherStyle.Render("● " + upper)
	}
}

// RenderBatteryBadge formats and colors the battery percentage and charging indicator
func RenderBatteryBadge(battery int, status string, l *locale.Strings) string {
	if battery < 0 {
		return lipgloss.NewStyle().Foreground(ColorMuted).Render("—")
	}

	var color lipgloss.Color
	if battery >= 50 {
		color = ColorSuccess
	} else if battery >= 20 {
		color = ColorWarning
	} else {
		color = ColorDanger
	}

	style := lipgloss.NewStyle().Foreground(color).Bold(true)
	icon := "🔋"
	if status == l.BatteryCharging {
		icon = "⚡"
	}

	if status != "" && status != l.BatteryUnknown {
		return style.Render(fmt.Sprintf("%s %d%% (%s)", icon, battery, status))
	}
	return style.Render(fmt.Sprintf("%s %d%%", icon, battery))
}

// RenderProgressBar returns a visual ASCII progress bar
func RenderProgressBar(percent float64, width int) string {
	if width <= 4 {
		width = 10
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int((percent / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	var color lipgloss.Color
	if percent >= 85 {
		color = ColorDanger
	} else if percent >= 65 {
		color = ColorWarning
	} else {
		color = ColorSuccess
	}

	barFilled := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled))
	barEmpty := lipgloss.NewStyle().Foreground(ColorSubtle).Render(strings.Repeat("░", empty))

	return fmt.Sprintf("[%s%s]", barFilled, barEmpty)
}

// RenderLogLevelBadge returns a colored log level badge (V, D, I, W, E, F)
func RenderLogLevelBadge(level string) string {
	switch level {
	case "E", "F":
		return lipgloss.NewStyle().Bold(true).Foreground(ColorDanger).Render(" " + level + " ")
	case "W":
		return lipgloss.NewStyle().Bold(true).Foreground(ColorWarning).Render(" " + level + " ")
	case "I":
		return lipgloss.NewStyle().Bold(true).Foreground(ColorSuccess).Render(" " + level + " ")
	case "D":
		return lipgloss.NewStyle().Bold(true).Foreground(ColorSecondary).Render(" " + level + " ")
	default:
		return lipgloss.NewStyle().Foreground(ColorMuted).Render(" " + level + " ")
	}
}

