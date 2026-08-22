package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
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
func DefaultActions(l *locale.Strings) []ActionItem {
	return []ActionItem{
		{"screenshot", l.ActionsScreenshotTitle, "s", l.ActionsScreenshotDesc},
		{"wireless", l.ActionsWirelessTitle, "w", l.ActionsWirelessDesc},
		{"reboot_system", l.ActionsRebootSystem, "1", l.ActionsRebootSystemDesc},
		{"reboot_recovery", l.ActionsRebootRecovery, "2", l.ActionsRebootRecDesc},
		{"reboot_bootloader", l.ActionsRebootBoot, "3", l.ActionsRebootBootDesc},
	}
}

// RenderActionsModal renders an interactive action menu dialog for the selected device
func RenderActionsModal(dev *device.Device, selectedActionIndex int, l *locale.Strings) string {
	if dev == nil {
		return ""
	}

	var lines []string
	title := fmt.Sprintf(l.ActionsTitle, dev.DisplayName(), dev.Serial)
	lines = append(lines, styles.ModalTitleStyle.Render(title))
	lines = append(lines, "")

	actions := DefaultActions(l)
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
	lines = append(lines, styles.HelpDescStyle.Render(l.ActionsFooter))

	return styles.ModalBoxStyle.Render(strings.Join(lines, "\n"))
}
