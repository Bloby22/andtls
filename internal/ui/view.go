package ui

import (
	"strings"

	"github.com/Bloby22/andtls/internal/ui/components"
)

// AppVersion is the current semver release version string
const AppVersion = "0.1.0"

// View renders the complete TUI dashboard layout or active modal
func (m Model) View() string {
	var sb strings.Builder
	l := m.locale

	// 1. Header component
	sb.WriteString(components.RenderHeader(AppVersion, l))

	// 2. Modal Overlay or Main Dashboard
	switch m.viewMode {
	case ViewModeHelp:
		sb.WriteString(components.RenderHelpModal(l))
		sb.WriteString("\n\n")

	case ViewModeActions:
		if sel := m.SelectedDevice(); sel != nil {
			sb.WriteString(components.RenderActionsModal(sel, m.selectedActionIndex, l))
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(components.RenderEmptyState(l))
			sb.WriteString("\n")
		}

	default:
		// Normal Dashboard View
		if len(m.devices) == 0 {
			sb.WriteString(components.RenderEmptyState(l))
			sb.WriteString("\n")
		} else {
			sb.WriteString(components.RenderTable(m.devices, m.selectedIndex, l))
			sb.WriteString("\n")

			// 3. Selected Device Inspector Card
			if selected := m.SelectedDevice(); selected != nil {
				sb.WriteString(components.RenderDetailsCard(selected, l))
				sb.WriteString("\n")
			}
		}
	}

	// 4. Status bar and Help footer component
	hasSelection := m.SelectedDevice() != nil && m.SelectedDevice().IsAuthorized()
	sb.WriteString(components.RenderStatusBar(m.statusMsg, m.isErr, m.lastChecked, hasSelection, l))

	return sb.String()
}
