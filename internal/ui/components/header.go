package components

import (
	"fmt"

	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderHeader renders the top application title bar
func RenderHeader(version string, l *locale.Strings) string {
	title := styles.HeaderStyle.Render(" " + l.AppTitle + " ")
	subtitle := styles.HeaderSubStyle.Render(" " + l.AppSubtitle + " ")
	ver := styles.HeaderVersionStyle.Render("v" + version)

	return fmt.Sprintf("%s%s %s\n\n", title, subtitle, ver)
}
