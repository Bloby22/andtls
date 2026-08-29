package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderRemoteView renders the virtual keypad and remote control interface
func RenderRemoteView(dev *device.Device, l *locale.Strings, lastAction string) string {
	if dev == nil {
		return ""
	}

	var sb strings.Builder

	title := fmt.Sprintf(l.RemoteTitle, dev.DisplayName())
	sb.WriteString(styles.ModalTitleStyle.Render(title) + "\n\n")

	// Render buttons grid
	btnStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorSecondary).
		Padding(0, 2).
		MarginRight(2)

	powerBtn := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorDanger).
		Foreground(styles.ColorDanger).
		Bold(true).
		Padding(0, 2).
		Render("[p] " + l.RemotePower)

	wakeBtn := btnStyle.Render("[w] " + l.RemoteWake)
	row1 := lipgloss.JoinHorizontal(lipgloss.Center, powerBtn, "  ", wakeBtn)
	sb.WriteString(row1 + "\n\n")

	volUpBtn := btnStyle.Render("[+] " + l.RemoteVolUp)
	volDownBtn := btnStyle.Render("[-] " + l.RemoteVolDown)
	row2 := lipgloss.JoinHorizontal(lipgloss.Center, volUpBtn, "  ", volDownBtn)
	sb.WriteString(row2 + "\n\n")

	homeBtn := btnStyle.Render("[h] " + l.RemoteHome)
	backBtn := btnStyle.Render("[b] " + l.RemoteBack)
	recentBtn := btnStyle.Render("[r] " + l.RemoteRecent)
	row3 := lipgloss.JoinHorizontal(lipgloss.Center, backBtn, "  ", homeBtn, "  ", recentBtn)
	sb.WriteString(row3 + "\n\n")

	textBtn := btnStyle.Render("[t] " + l.RemoteTextInput)
	sb.WriteString(textBtn + "\n\n")

	if lastAction != "" {
		lastActStr := lipgloss.NewStyle().Foreground(styles.ColorSuccess).Bold(true).Render("✔ " + lastAction)
		sb.WriteString(lastActStr + "\n\n")
	}

	sb.WriteString(styles.HelpDescStyle.Render(l.RemoteFooter))

	return styles.ModalBoxStyle.Render(sb.String())
}
