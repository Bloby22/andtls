package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bloby22/andtls/internal/ui/components"
)

// AppVersion is the current semver release version string
const AppVersion = "0.2.0"

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

	case ViewModeLogcat:
		if sel := m.SelectedDevice(); sel != nil {
			sb.WriteString(components.RenderLogcatView(sel, m.logcatLines, m.logcatFilterLevel, m.logcatSearchFilter, m.logcatPaused, l, m.width, m.height))
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(components.RenderEmptyState(l))
			sb.WriteString("\n")
		}

	case ViewModeFiles:
		if sel := m.SelectedDevice(); sel != nil {
			sb.WriteString(components.RenderFilesView(sel, m.filesCurrentDir, m.filesList, m.filesSelectedIndex, l, m.width, m.height))
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(components.RenderEmptyState(l))
			sb.WriteString("\n")
		}

	case ViewModePackages:
		if sel := m.SelectedDevice(); sel != nil {
			sb.WriteString(components.RenderPackagesView(sel, m.packagesList, m.packagesSelectedIndex, m.packagesSearchFilter, m.packagesThirdPartyOnly, l, m.width, m.height))
			sb.WriteString("\n\n")
		} else {
			sb.WriteString(components.RenderEmptyState(l))
			sb.WriteString("\n")
		}

	case ViewModeRemote:
		if sel := m.SelectedDevice(); sel != nil {
			sb.WriteString(components.RenderRemoteView(sel, l, m.remoteLastAction))
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
	var recordDurStr string
	if m.isRecording {
		dur := time.Since(m.recordStartTime)
		recordDurStr = fmt.Sprintf("%02d:%02d", int(dur.Minutes()), int(dur.Seconds())%60)
	}
	sb.WriteString(components.RenderStatusBar(m.statusMsg, m.isErr, m.lastChecked, hasSelection, m.isRecording, recordDurStr, l))

	return sb.String()
}
