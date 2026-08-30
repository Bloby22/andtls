package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderHelpModal renders a floating dialog box with the complete keyboard shortcut map
func RenderHelpModal(l *locale.Strings) string {
	var lines []string

	lines = append(lines, styles.ModalTitleStyle.Render(l.HelpTitle))
	lines = append(lines, "")

	shortcuts := []struct {
		Key  string
		Desc string
	}{
		{"↑ / ↓, k / j", l.HelpMoveUp + " / " + l.HelpMoveDown},
		{"Enter / a", l.HelpActionsMenu},
		{"m", l.HelpMirror},
		{"v", l.HelpRecording},
		{"f", l.HelpFileManager},
		{"l", l.HelpLogcat},
		{"p", l.HelpPackages},
		{"k", l.HelpRemote},
		{"s", l.HelpScreenshot},
		{"w", l.HelpWirelessADB},
		{"b", l.HelpRebootMenu},
		{"e", l.HelpExportReport},
		{"r", l.HelpRefresh},
		{"t", l.HelpTheme},
		{"L", l.HelpLanguage},
		{"? / h", l.HelpToggleHelp},
		{"Esc", l.HelpCloseModal},
		{"q / Ctrl+C", l.HelpQuit},
	}

	for _, s := range shortcuts {
		line := fmt.Sprintf("  %-18s %s", styles.HelpKeyStyle.Render(s.Key), styles.HelpDescStyle.Render(s.Desc))
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, styles.HelpDescStyle.Render(l.HelpCloseHint))

	return styles.ModalBoxStyle.Render(strings.Join(lines, "\n"))
}
