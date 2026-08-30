package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderTable renders the interactive device list table with highlighting for the active selection
func RenderTable(devices []device.Device, selectedIndex int, l *locale.Strings) string {
	var tableRows []string

	// Table column header line
	headerCols := fmt.Sprintf(
		"  %-18s  %-22s  %-24s  %-20s  %-8s",
		l.TableStatus, l.TableSerial, l.TableModel, l.TableBattery, l.TableTransport,
	)
	tableRows = append(tableRows, styles.TableHeaderStyle.Render(headerCols))

	for i, dev := range devices {
		cursor := "  "
		isSelected := i == selectedIndex
		if isSelected {
			cursor = styles.CursorStyle.Render("▸ ")
		}

		statusBadge := styles.RenderStatusBadge(dev.State, l)
		serial := dev.Serial
		if len(serial) > 20 {
			serial = serial[:17] + "..."
		}

		modelName := dev.DisplayName()
		if len(modelName) > 22 {
			modelName = modelName[:19] + "..."
		}

		batteryBadge := styles.RenderBatteryBadge(dev.Battery, dev.BatteryStatus, l)
		trID := dev.TransportID
		if trID == "" {
			trID = "—"
		}

		rowStr := fmt.Sprintf(
			"%s%-18s  %-22s  %-24s  %-20s  %-8s",
			cursor,
			statusBadge,
			serial,
			modelName,
			batteryBadge,
			trID,
		)

		if isSelected {
			tableRows = append(tableRows, styles.TableSelectedRowStyle.Render(rowStr))
		} else {
			tableRows = append(tableRows, styles.TableRowStyle.Render(rowStr))
		}
	}

	tableContent := strings.Join(tableRows, "\n")
	return styles.TableBoxStyle.Render(tableContent)
}
