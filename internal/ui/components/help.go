package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderHelpModal renders a floating dialog box with the complete keyboard shortcut map
func RenderHelpModal() string {
	var lines []string

	lines = append(lines, styles.ModalTitleStyle.Render("📖 andtls Keyboard Shortcuts & Help"))
	lines = append(lines, "")

	shortcuts := []struct {
		Key  string
		Desc string
	}{
		{"↑ / k", "Move selection cursor up"},
		{"↓ / j", "Move selection cursor down"},
		{"r", "Trigger immediate device poll and refresh"},
		{"s", "Take screenshot of selected device (saved to ./screenshots/)"},
		{"w", "Enable wireless ADB mode on port 5555 for selected device"},
		{"b", "Open reboot options menu (System, Recovery, Bootloader)"},
		{"Enter", "Open quick actions menu for selected device"},
		{"? / h", "Toggle this help dialog overlay"},
		{"Esc", "Close any active modal overlay"},
		{"q / Ctrl+C", "Quit application"},
	}

	for _, s := range shortcuts {
		line := fmt.Sprintf("  %-16s %s", styles.HelpKeyStyle.Render(s.Key), styles.HelpDescStyle.Render(s.Desc))
		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, styles.HelpDescStyle.Render("Press [Esc] or [?] to close this help window"))

	return styles.ModalBoxStyle.Render(strings.Join(lines, "\n"))
}
