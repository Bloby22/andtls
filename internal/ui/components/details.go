package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderDetailsCard renders an inspection card with detailed properties of the selected device
func RenderDetailsCard(dev *device.Device) string {
	if dev == nil {
		return ""
	}

	var details []string
	details = append(details, styles.DetailTitleStyle.Render("🔍 Device Inspector"))

	// Row 1: Serial & Status
	row1 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render("Serial Number:"), styles.DetailValueStyle.Render(dev.Serial),
		styles.DetailKeyStyle.Render("ADB Status:"), styles.RenderStatusBadge(dev.State),
	)
	details = append(details, row1)

	// Row 2: Model & Codename
	modelVal := dev.Model
	if modelVal == "" {
		modelVal = "—"
	}
	prodVal := dev.Product
	if prodVal == "" {
		prodVal = "—"
	}
	row2 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render("Model Name:"), styles.DetailValueStyle.Render(modelVal),
		styles.DetailKeyStyle.Render("Codename / Product:"), styles.DetailValueStyle.Render(prodVal),
	)
	details = append(details, row2)

	// Row 3: OS Version & IP Address
	osVal := dev.OSVersionString()
	ipVal := dev.IPAddress
	if ipVal == "" {
		ipVal = "— (USB / Unknown)"
	}
	row3 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render("Android OS:"), styles.DetailValueStyle.Render(osVal),
		styles.DetailKeyStyle.Render("IP Address:"), styles.DetailValueStyle.Render(ipVal),
	)
	details = append(details, row3)

	// Row 4: Battery Details & Storage
	batVal := dev.BatteryDetailString()
	storageVal := dev.StorageString()
	row4 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render("Battery Telemetry:"), styles.DetailValueStyle.Render(batVal),
		styles.DetailKeyStyle.Render("Storage Space:"), styles.DetailValueStyle.Render(storageVal),
	)
	details = append(details, row4)

	// Row 5: Port & Screen Resolution
	portVal := dev.USBPort
	if portVal == "" {
		portVal = "USB / Transport " + dev.TransportID
	}
	screenVal := dev.ScreenResolution
	if screenVal == "" {
		screenVal = "—"
	}
	row5 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render("USB Bus / Port:"), styles.DetailValueStyle.Render(portVal),
		styles.DetailKeyStyle.Render("Display Resolution:"), styles.DetailValueStyle.Render(screenVal),
	)
	details = append(details, row5)

	return styles.DetailCardStyle.Render(strings.Join(details, "\n"))
}
