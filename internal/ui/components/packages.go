package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderPackagesView renders the application and package manager
func RenderPackagesView(dev *device.Device, packages []string, selectedIndex int, searchFilter string, thirdPartyOnly bool, l *locale.Strings, termWidth, termHeight int) string {
	if dev == nil {
		return ""
	}

	var sb strings.Builder

	// 1. Filter packages by search query
	var filtered []string
	query := strings.ToLower(strings.TrimSpace(searchFilter))
	for _, pkg := range packages {
		if query == "" || strings.Contains(strings.ToLower(pkg), query) {
			filtered = append(filtered, pkg)
		}
	}

	// 2. Title & Header
	title := fmt.Sprintf(l.PackagesTitle, len(filtered), dev.DisplayName())
	modeBadge := styles.HelpKeyStyle.Render("[User Apps (3rd Party)]")
	if !thirdPartyOnly {
		modeBadge = lipgloss.NewStyle().Foreground(styles.ColorWarning).Render("[All Packages (System + User)]")
	}
	searchBadge := ""
	if searchFilter != "" {
		searchBadge = lipgloss.NewStyle().Foreground(styles.ColorSecondary).Render(fmt.Sprintf(l.PackagesFilter, searchFilter))
	}

	titleBar := fmt.Sprintf("%s   %s   %s", styles.ModalTitleStyle.Render(title), modeBadge, searchBadge)
	sb.WriteString(titleBar + "\n\n")

	// 3. Package List
	if len(filtered) == 0 {
		emptyMsg := lipgloss.NewStyle().Foreground(styles.ColorMuted).Italic(true).Render("  (No packages matching search filter...)")
		sb.WriteString(emptyMsg + "\n\n")
	} else {
		maxVisible := 14
		if termHeight > 25 {
			maxVisible = termHeight - 11
		}

		if selectedIndex >= len(filtered) {
			selectedIndex = len(filtered) - 1
		}
		if selectedIndex < 0 {
			selectedIndex = 0
		}

		startIdx := 0
		if selectedIndex >= maxVisible {
			startIdx = selectedIndex - maxVisible + 1
		}
		endIdx := startIdx + maxVisible
		if endIdx > len(filtered) {
			endIdx = len(filtered)
		}

		for i := startIdx; i < endIdx; i++ {
			pkg := filtered[i]
			cursor := "  "
			isSelected := i == selectedIndex
			if isSelected {
				cursor = styles.CursorStyle.Render("▸ ")
			}

			icon := "📦"
			if strings.HasPrefix(pkg, "com.google.") || strings.HasPrefix(pkg, "com.android.") {
				icon = "⚙️ "
			}

			rowStr := fmt.Sprintf("%s%s %-54s", cursor, icon, pkg)
			if isSelected {
				sb.WriteString(styles.TableSelectedRowStyle.Render(rowStr) + "\n")
			} else {
				sb.WriteString(styles.TableRowStyle.Render(rowStr) + "\n")
			}
		}
		sb.WriteString("\n")
	}

	// 4. Footer hints
	footer := fmt.Sprintf("%s   %s", styles.HelpDescStyle.Render(l.PackagesHint), styles.HelpKeyStyle.Render(l.PackagesToggle))
	sb.WriteString(footer)

	return styles.ModalBoxStyle.Render(sb.String())
}
