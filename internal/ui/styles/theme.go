package styles

import (
	"fmt"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/charmbracelet/lipgloss"
)

// HeaderStyle styles the main application banner block
var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorTextBold).
	Background(ColorPrimary).
	Padding(0, 1)

// HeaderSubStyle styles the application subtitle
var HeaderSubStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary)

// HeaderVersionStyle styles the version text in the header
var HeaderVersionStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorMuted)

// TableBoxStyle styles the border and frame surrounding the device table
var TableBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorPrimary).
	Padding(0, 1).
	MarginBottom(1)

// TableHeaderStyle styles the column headers of the device table
var TableHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary).
	BorderBottom(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottomForeground(ColorSubtle).
	PaddingBottom(0)

// TableRowStyle styles unselected rows in the device table
var TableRowStyle = lipgloss.NewStyle().
	Foreground(ColorText).
	Padding(0, 0)

// TableSelectedRowStyle styles the currently highlighted row in the device table
var TableSelectedRowStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorTextBold).
	Background(ColorHighlight).
	Padding(0, 0)

// CursorStyle styles the active selection cursor indicator
var CursorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary)

// StatusDeviceStyle styles the connected device state badge
var StatusDeviceStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSuccess)

// StatusUnauthorizedStyle styles the unauthorized device state badge
var StatusUnauthorizedStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorWarning)

// StatusOfflineStyle styles the offline device state badge
var StatusOfflineStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorDanger)

// StatusOtherStyle styles other operational mode state badges
var StatusOtherStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary)

// DetailCardStyle styles the container frame for the device inspector card
var DetailCardStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorSecondary).
	Padding(0, 1).
	MarginBottom(1)

// DetailTitleStyle styles the title bar of the device inspector card
var DetailTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary).
	MarginBottom(0)

// DetailKeyStyle styles metadata label keys in the inspector card
var DetailKeyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorMuted).
	Width(18)

// DetailValueStyle styles metadata values in the inspector card
var DetailValueStyle = lipgloss.NewStyle().
	Foreground(ColorTextBold)

// EmptyBoxStyle styles the troubleshooting container when no devices are connected
var EmptyBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ColorMuted).
	Padding(1, 2).
	Align(lipgloss.Center).
	MarginBottom(1)

// EmptyTitleStyle styles the main heading inside the empty state box
var EmptyTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorWarning).
	MarginBottom(1)

// EmptyTipStyle styles the individual instruction lines in the empty state box
var EmptyTipStyle = lipgloss.NewStyle().
	Foreground(ColorMuted)

// StatusBarStyle styles the bottom status bar strip
var StatusBarStyle = lipgloss.NewStyle().
	Background(ColorBgAlt).
	Foreground(ColorText).
	Padding(0, 1)

// StatusMsgSuccessStyle styles standard status messages
var StatusMsgSuccessStyle = lipgloss.NewStyle().
	Foreground(ColorSuccess).
	Bold(true)

// StatusMsgErrorStyle styles error notifications in the status bar
var StatusMsgErrorStyle = lipgloss.NewStyle().
	Foreground(ColorDanger).
	Bold(true)

// HelpKeyStyle styles keyboard shortcut keys in help displays
var HelpKeyStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary)

// HelpDescStyle styles description labels for keyboard shortcuts
var HelpDescStyle = lipgloss.NewStyle().
	Foreground(ColorMuted)

// ModalBoxStyle styles floating modal dialog boxes (help, actions)
var ModalBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.DoubleBorder()).
	BorderForeground(ColorPrimary).
	Padding(1, 2).
	MarginBottom(1)

// ModalTitleStyle styles headings inside modal dialogs
var ModalTitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ColorSecondary).
	MarginBottom(1)

// RenderStatusBadge formats and styles the device state indicator badge
func RenderStatusBadge(state device.DeviceState) string {
	switch state {
	case device.StateDevice:
		return StatusDeviceStyle.Render("● CONNECTED")
	case device.StateUnauthorized:
		return StatusUnauthorizedStyle.Render("● UNAUTHORIZED")
	case device.StateOffline:
		return StatusOfflineStyle.Render("● OFFLINE")
	case device.StateAuthorizing:
		return StatusUnauthorizedStyle.Render("● AUTHORIZING...")
	case device.StateRecovery:
		return StatusOtherStyle.Render("● RECOVERY")
	case device.StateSideload:
		return StatusOtherStyle.Render("● SIDELOAD")
	case device.StateBootloader:
		return StatusOtherStyle.Render("● BOOTLOADER")
	default:
		return StatusOtherStyle.Render("● " + string(state))
	}
}

// RenderBatteryBadge formats and colors the battery percentage and charging indicator
func RenderBatteryBadge(battery int, status string) string {
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
	if status == "Charging" || status == "Nabíjí se" {
		icon = "⚡"
	}

	if status != "" && status != "Unknown" {
		return style.Render(fmt.Sprintf("%s %d%% (%s)", icon, battery, status))
	}
	return style.Render(fmt.Sprintf("%s %d%%", icon, battery))
}
