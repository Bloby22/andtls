package components

import (
	"fmt"
	"time"

	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders the bottom status bar and keybinding help
func RenderStatusBar(statusMsg string, isErr bool, lastChecked time.Time, hasSelectedDevice bool) string {
	var statusLeft string
	if isErr {
		statusLeft = styles.StatusMsgErrorStyle.Render("✖ " + statusMsg)
	} else {
		statusLeft = styles.StatusMsgSuccessStyle.Render("✔ " + statusMsg)
	}

	lastUpdated := fmt.Sprintf("Updated: %s", lastChecked.Format("15:04:05"))
	statusRight := lipgloss.NewStyle().Foreground(styles.ColorMuted).Render(lastUpdated)

	statusBar := styles.StatusBarStyle.Render(
		fmt.Sprintf("%s   │   %s", statusLeft, statusRight),
	)

	var helpKeys string
	if hasSelectedDevice {
		helpKeys = fmt.Sprintf(
			"%s %s   %s %s   %s %s   %s %s   %s %s   %s %s",
			styles.HelpKeyStyle.Render("[↑/↓]"), styles.HelpDescStyle.Render("Navigate"),
			styles.HelpKeyStyle.Render("[Enter]"), styles.HelpDescStyle.Render("Actions"),
			styles.HelpKeyStyle.Render("[s]"), styles.HelpDescStyle.Render("Screenshot"),
			styles.HelpKeyStyle.Render("[r]"), styles.HelpDescStyle.Render("Refresh"),
			styles.HelpKeyStyle.Render("[?]"), styles.HelpDescStyle.Render("Help"),
			styles.HelpKeyStyle.Render("[q]"), styles.HelpDescStyle.Render("Quit"),
		)
	} else {
		helpKeys = fmt.Sprintf(
			"%s %s   %s %s   %s %s   %s %s",
			styles.HelpKeyStyle.Render("[r]"), styles.HelpDescStyle.Render("Refresh"),
			styles.HelpKeyStyle.Render("[?]"), styles.HelpDescStyle.Render("Help"),
			styles.HelpKeyStyle.Render("[q]"), styles.HelpDescStyle.Render("Quit"),
			styles.HelpKeyStyle.Render("[Ctrl+C]"), styles.HelpDescStyle.Render("Exit"),
		)
	}

	return fmt.Sprintf("%s\n\n%s\n", statusBar, helpKeys)
}
