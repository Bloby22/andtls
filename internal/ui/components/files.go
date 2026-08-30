package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// fileIcon returns an emoji icon based on filename and directory status
func fileIcon(item device.FileInfo) string {
	if item.IsDir {
		return "📁"
	}
	lower := strings.ToLower(item.Name)
	switch {
	case strings.HasSuffix(lower, ".apk"), strings.HasSuffix(lower, ".apkm"), strings.HasSuffix(lower, ".xapk"):
		return "📦"
	case strings.HasSuffix(lower, ".mp4"), strings.HasSuffix(lower, ".mkv"), strings.HasSuffix(lower, ".webm"), strings.HasSuffix(lower, ".avi"):
		return "🎬"
	case strings.HasSuffix(lower, ".png"), strings.HasSuffix(lower, ".jpg"), strings.HasSuffix(lower, ".jpeg"), strings.HasSuffix(lower, ".webp"), strings.HasSuffix(lower, ".gif"):
		return "🖼️"
	case strings.HasSuffix(lower, ".mp3"), strings.HasSuffix(lower, ".flac"), strings.HasSuffix(lower, ".wav"), strings.HasSuffix(lower, ".ogg"), strings.HasSuffix(lower, ".m4a"):
		return "🎵"
	case strings.HasSuffix(lower, ".zip"), strings.HasSuffix(lower, ".tar"), strings.HasSuffix(lower, ".gz"), strings.HasSuffix(lower, ".7z"), strings.HasSuffix(lower, ".rar"):
		return "🗜️"
	case strings.HasSuffix(lower, ".pdf"), strings.HasSuffix(lower, ".txt"), strings.HasSuffix(lower, ".log"), strings.HasSuffix(lower, ".json"), strings.HasSuffix(lower, ".xml"):
		return "📄"
	default:
		return "📄"
	}
}

// RenderFilesView renders the interactive file manager browser
func RenderFilesView(dev *device.Device, currentDir string, files []device.FileInfo, selectedIndex int, l *locale.Strings, termWidth, termHeight int) string {
	if dev == nil {
		return ""
	}

	var sb strings.Builder

	// 1. Title bar
	title := fmt.Sprintf(l.FilesTitle, currentDir, dev.DisplayName())
	sb.WriteString(styles.ModalTitleStyle.Render(title) + "\n\n")

	// 2. Table Header
	headerCols := fmt.Sprintf(
		"  %-36s  %-12s  %-12s  %-16s",
		l.FilesHeaderName, l.FilesHeaderSize, l.FilesHeaderPerm, l.FilesHeaderMod,
	)
	sb.WriteString(styles.TableHeaderStyle.Render(headerCols) + "\n")

	// 3. File list
	if len(files) == 0 {
		emptyMsg := lipgloss.NewStyle().Foreground(styles.ColorMuted).Italic(true).Render("  " + l.FilesEmpty)
		sb.WriteString(emptyMsg + "\n\n")
	} else {
		// Pagination / scroll window calculation
		maxVisible := 14
		if termHeight > 25 {
			maxVisible = termHeight - 11
		}

		startIdx := 0
		if selectedIndex >= maxVisible {
			startIdx = selectedIndex - maxVisible + 1
		}
		endIdx := startIdx + maxVisible
		if endIdx > len(files) {
			endIdx = len(files)
		}

		for i := startIdx; i < endIdx; i++ {
			f := files[i]
			cursor := "  "
			isSelected := i == selectedIndex
			if isSelected {
				cursor = styles.CursorStyle.Render("▸ ")
			}

			icon := fileIcon(f)
			nameDisplay := fmt.Sprintf("%s %s", icon, truncateString(f.Name, 32))
			sizeDisplay := f.Size
			permDisplay := f.Permissions
			modDisplay := f.ModTime

			rowStr := fmt.Sprintf(
				"%s%-36s  %-12s  %-12s  %-16s",
				cursor, nameDisplay, sizeDisplay, permDisplay, modDisplay,
			)

			if isSelected {
				sb.WriteString(styles.TableSelectedRowStyle.Render(rowStr) + "\n")
			} else {
				if f.IsDir {
					sb.WriteString(lipgloss.NewStyle().Foreground(styles.ColorSecondary).Bold(true).Render(rowStr) + "\n")
				} else {
					sb.WriteString(styles.TableRowStyle.Render(rowStr) + "\n")
				}
			}
		}
		sb.WriteString("\n")
	}

	// 4. Footer hint
	sb.WriteString(styles.HelpDescStyle.Render(l.FilesHint))

	return styles.ModalBoxStyle.Render(sb.String())
}
