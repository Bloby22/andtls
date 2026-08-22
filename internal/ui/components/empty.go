package components

import (
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderEmptyState renders a helpful troubleshooting guide when no devices are connected
func RenderEmptyState() string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.EmptyTitleStyle.Render("⚡ No Android Devices Detected"),
		styles.EmptyTipStyle.Render("1. Connect your Android phone or tablet via USB cable."),
		styles.EmptyTipStyle.Render("2. In Android Settings -> Developer Options, enable 'USB Debugging'."),
		styles.EmptyTipStyle.Render("3. On your device screen, accept the 'Allow USB debugging?' prompt."),
		styles.EmptyTipStyle.Render("4. Press [r] to immediately refresh the connection list."),
	)

	return styles.EmptyBoxStyle.Render(content)
}
