package components

import (
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// RenderEmptyState renders a helpful troubleshooting guide when no devices are connected
func RenderEmptyState(l *locale.Strings) string {
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.EmptyTitleStyle.Render(l.EmptyTitle),
		styles.EmptyTipStyle.Render(l.EmptyTip1),
		styles.EmptyTipStyle.Render(l.EmptyTip2),
		styles.EmptyTipStyle.Render(l.EmptyTip3),
		styles.EmptyTipStyle.Render(l.EmptyTip4),
	)

	return styles.EmptyBoxStyle.Render(content)
}
