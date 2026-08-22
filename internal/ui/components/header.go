package components

import (
	"fmt"

	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderHeader renders the top application title bar
func RenderHeader(version string) string {
	title := styles.HeaderStyle.Render(" 📱 ANDTLS ")
	subtitle := styles.HeaderSubStyle.Render(" Android Tools Dashboard ")
	ver := styles.HeaderVersionStyle.Render("v" + version)

	return fmt.Sprintf("%s%s %s\n\n", title, subtitle, ver)
}
