package components

import (
	"fmt"
	"time"

	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders the bottom status bar and keybinding help
func RenderStatusBar(statusMsg string, isErr bool, lastChecked time.Time, hasSelectedDevice bool, l *locale.Strings) string {
	var statusLeft string
	if isErr {
		statusLeft = styles.StatusMsgErrorStyle.Render("✖ " + statusMsg)
	} else {
		statusLeft = styles.StatusMsgSuccessStyle.Render("✔ " + statusMsg)
	}

	lastUpdated := fmt.Sprintf(l.StatusUpdated, lastChecked.Format("15:04:05"))
	statusRight := lipgloss.NewStyle().Foreground(styles.ColorMuted).Render(lastUpdated)

	statusBar := styles.StatusBarStyle.Render(
		fmt.Sprintf("%s   │   %s", statusLeft, statusRight),
	)

	var helpKeys string
	if hasSelectedDevice {
		helpKeys = fmt.Sprintf(
			"%s %s   %s %s   %s %s   %s %s   %s %s   %s %s",
			styles.HelpKeyStyle.Render("[↑/↓]"), styles.HelpDescStyle.Render(l.StatusNavigate),
			styles.HelpKeyStyle.Render("[Enter]"), styles.HelpDescStyle.Render(l.StatusActions),
			styles.HelpKeyStyle.Render("[s]"), styles.HelpDescStyle.Render(l.StatusScreenshot),
			styles.HelpKeyStyle.Render("[r]"), styles.HelpDescStyle.Render(l.StatusRefresh),
			styles.HelpKeyStyle.Render("[?]"), styles.HelpDescStyle.Render(l.StatusHelp),
			styles.HelpKeyStyle.Render("[q]"), styles.HelpDescStyle.Render(l.StatusQuit),
		)
	} else {
		helpKeys = fmt.Sprintf(
			"%s %s   %s %s   %s %s   %s %s",
			styles.HelpKeyStyle.Render("[r]"), styles.HelpDescStyle.Render(l.StatusRefresh),
			styles.HelpKeyStyle.Render("[?]"), styles.HelpDescStyle.Render(l.StatusHelp),
			styles.HelpKeyStyle.Render("[q]"), styles.HelpDescStyle.Render(l.StatusQuit),
			styles.HelpKeyStyle.Render("[Ctrl+C]"), styles.HelpDescStyle.Render(l.StatusExit),
		)
	}

	return fmt.Sprintf("%s\n\n%s\n", statusBar, helpKeys)
}
