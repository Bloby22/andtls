package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderStatusBar renders the bottom status bar and keybinding help
func RenderStatusBar(statusMsg string, isErr bool, lastChecked time.Time, hasSelectedDevice bool, isRecording bool, recordDur string, l *locale.Strings) string {
	var statusLeft string
	if isRecording {
		statusLeft = lipgloss.NewStyle().Bold(true).Foreground(styles.ColorDanger).Blink(true).Render(fmt.Sprintf(l.StatusRecordingActive, recordDur)) + "  "
	}

	if isErr {
		statusLeft += styles.StatusMsgErrorStyle.Render("✖ " + statusMsg)
	} else {
		statusLeft += styles.StatusMsgSuccessStyle.Render("✔ " + statusMsg)
	}

	themeBadge := lipgloss.NewStyle().Foreground(styles.ColorPrimary).Bold(true).Render("🎨 " + styles.CurrentThemeName())
	langBadge := lipgloss.NewStyle().Foreground(styles.ColorSecondary).Bold(true).Render("🌐 " + strings.ToUpper(locale.CurrentCode()))
	lastUpdated := fmt.Sprintf(l.StatusUpdated, lastChecked.Format("15:04:05"))
	statusRight := lipgloss.NewStyle().Foreground(styles.ColorMuted).Render(lastUpdated)

	statusBar := styles.StatusBarStyle.Render(
		fmt.Sprintf("%s   │   %s   │   %s   │   %s", statusLeft, themeBadge, langBadge, statusRight),
	)

	var helpKeys string
	if hasSelectedDevice {
		helpKeys = fmt.Sprintf(
			"%s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s",
			styles.HelpKeyStyle.Render("[↑/↓]"), styles.HelpDescStyle.Render(l.StatusNavigate),
			styles.HelpKeyStyle.Render("[Enter]"), styles.HelpDescStyle.Render(l.StatusActions),
			styles.HelpKeyStyle.Render("[m]"), styles.HelpDescStyle.Render(l.StatusMirror),
			styles.HelpKeyStyle.Render("[f]"), styles.HelpDescStyle.Render(l.StatusFiles),
			styles.HelpKeyStyle.Render("[l]"), styles.HelpDescStyle.Render(l.StatusLogcat),
			styles.HelpKeyStyle.Render("[p]"), styles.HelpDescStyle.Render(l.StatusPackages),
			styles.HelpKeyStyle.Render("[s]"), styles.HelpDescStyle.Render(l.StatusScreenshot),
			styles.HelpKeyStyle.Render("[L]"), styles.HelpDescStyle.Render("Lang"),
			styles.HelpKeyStyle.Render("[t]"), styles.HelpDescStyle.Render("Theme"),
			styles.HelpKeyStyle.Render("[?]"), styles.HelpDescStyle.Render(l.StatusHelp),
		)
	} else {
		helpKeys = fmt.Sprintf(
			"%s %s   %s %s   %s %s   %s %s   %s %s",
			styles.HelpKeyStyle.Render("[r]"), styles.HelpDescStyle.Render(l.StatusRefresh),
			styles.HelpKeyStyle.Render("[L]"), styles.HelpDescStyle.Render("Lang"),
			styles.HelpKeyStyle.Render("[t]"), styles.HelpDescStyle.Render("Theme"),
			styles.HelpKeyStyle.Render("[?]"), styles.HelpDescStyle.Render(l.StatusHelp),
			styles.HelpKeyStyle.Render("[q]"), styles.HelpDescStyle.Render(l.StatusQuit),
		)
	}

	return fmt.Sprintf("%s\n\n%s\n", statusBar, helpKeys)
}
