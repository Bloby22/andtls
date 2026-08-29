package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/adb"
	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderLogcatView renders the interactive Logcat stream view
func RenderLogcatView(dev *device.Device, lines []adb.LogcatLine, filterLevel string, searchFilter string, isPaused bool, l *locale.Strings, termWidth, termHeight int) string {
	if dev == nil {
		return ""
	}

	var sb strings.Builder

	// 1. Title bar
	title := fmt.Sprintf(l.LogcatTitle, dev.DisplayName())
	stateBadge := styles.StatusMsgSuccessStyle.Render(l.LogcatRunning)
	if isPaused {
		stateBadge = styles.StatusMsgErrorStyle.Render(l.LogcatPaused)
	}
	levelBadge := styles.HelpKeyStyle.Render(fmt.Sprintf(l.LogcatFilter, filterLevel))
	searchBadge := ""
	if searchFilter != "" {
		searchBadge = lipgloss.NewStyle().Foreground(styles.ColorWarning).Render(fmt.Sprintf(l.LogcatSearch, searchFilter))
	}

	titleBar := fmt.Sprintf("%s   %s   %s   %s", styles.ModalTitleStyle.Render(title), stateBadge, levelBadge, searchBadge)
	sb.WriteString(titleBar + "\n\n")

	// 2. Filter lines
	var filtered []adb.LogcatLine
	for _, line := range lines {
		if filterLevel != "ALL" && filterLevel != "" {
			switch filterLevel {
			case "D":
				if line.Level == "V" {
					continue
				}
			case "I":
				if line.Level == "V" || line.Level == "D" {
					continue
				}
			case "W":
				if line.Level == "V" || line.Level == "D" || line.Level == "I" {
					continue
				}
			case "E":
				if line.Level != "E" && line.Level != "F" {
					continue
				}
			}
		}

		if searchFilter != "" {
			query := strings.ToLower(searchFilter)
			if !strings.Contains(strings.ToLower(line.Tag), query) &&
				!strings.Contains(strings.ToLower(line.Message), query) &&
				!strings.Contains(strings.ToLower(line.Raw), query) {
				continue
			}
		}

		filtered = append(filtered, line)
	}

	// Calculate visible height
	maxLines := 16
	if termHeight > 25 {
		maxLines = termHeight - 10
	}
	if len(filtered) > maxLines {
		filtered = filtered[len(filtered)-maxLines:]
	}

	// 3. Render lines
	if len(filtered) == 0 {
		emptyMsg := lipgloss.NewStyle().Foreground(styles.ColorMuted).Italic(true).Render("  (No log output matching current filter...)")
		sb.WriteString(emptyMsg + "\n\n")
	} else {
		for _, line := range filtered {
			timeStr := lipgloss.NewStyle().Foreground(styles.ColorMuted).Render(line.Timestamp)
			lvlBadge := styles.RenderLogLevelBadge(line.Level)
			tagStr := lipgloss.NewStyle().Bold(true).Foreground(styles.ColorSecondary).Render(fmt.Sprintf("%-16s", truncateString(line.Tag, 16)))
			msgStr := styles.TableRowStyle.Render(truncateString(line.Message, 90))

			row := fmt.Sprintf("%s %s %s %s", timeStr, lvlBadge, tagStr, msgStr)
			sb.WriteString(row + "\n")
		}
		sb.WriteString("\n")
	}

	// 4. Footer hints
	sb.WriteString(styles.HelpDescStyle.Render(l.LogcatHint))

	return styles.ModalBoxStyle.Render(sb.String())
}

func truncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return s
	}
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
