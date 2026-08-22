package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/ui/styles"
)

// ActionItem represents a selectable quick action
type ActionItem struct {
	// ID is the internal action identifier
	ID string

	// Title is the user-facing title
	Title string

	// Key is the shortcut key
	Key string

	// Desc is a short summary of the action
	Desc string
}

// DefaultActions returns the list of available actions on an authorized device
func DefaultActions() []ActionItem {
	return []ActionItem{
		{"screenshot", "Take Screenshot", "s", "Capture display to ./screenshots/"},
		{"wireless", "Enable Wireless ADB", "w", "Switch daemon to TCP port 5555"},
		{"reboot_system", "Reboot System", "1", "Normal device restart"},
		{"reboot_recovery", "Reboot to Recovery", "2", "Boot into Android Recovery"},
		{"reboot_bootloader", "Reboot to Bootloader", "3", "Boot into Fastboot mode"},
	}
}

// RenderActionsModal renders an interactive action menu dialog for the selected device
func RenderActionsModal(dev *device.Device, selectedActionIndex int) string {
	if dev == nil {
		return ""
	}

	var lines []string
	title := fmt.Sprintf("⚡ Quick Actions for: %s (%s)", dev.DisplayName(), dev.Serial)
	lines = append(lines, styles.ModalTitleStyle.Render(title))
	lines = append(lines, "")

	actions := DefaultActions()
	for i, act := range actions {
		cursor := "  "
		if i == selectedActionIndex {
			cursor = styles.CursorStyle.Render("▸ ")
		}

		keyBadge := styles.HelpKeyStyle.Render(fmt.Sprintf("[%s]", act.Key))
		itemText := fmt.Sprintf("%s%-6s %-24s %s", cursor, keyBadge, act.Title, styles.HelpDescStyle.Render(act.Desc))

		if i == selectedActionIndex {
			lines = append(lines, styles.TableSelectedRowStyle.Render(itemText))
		} else {
			lines = append(lines, styles.TableRowStyle.Render(itemText))
		}
	}

	lines = append(lines, "")
	lines = append(lines, styles.HelpDescStyle.Render("Press [Enter] to run selected, [Key] for shortcut, or [Esc] to cancel"))

	return styles.ModalBoxStyle.Render(strings.Join(lines, "\n"))
}
