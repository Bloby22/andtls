package components

import (
	"fmt"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/styles"
)

// RenderDetailsCard renders an inspection card with detailed properties of the selected device
func RenderDetailsCard(dev *device.Device, l *locale.Strings) string {
	if dev == nil {
		return ""
	}

	var details []string
	details = append(details, styles.DetailTitleStyle.Render(l.DetailTitle))

	// Row 1: Serial & Status
	row1 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render(l.DetailSerialNumber), styles.DetailValueStyle.Render(dev.Serial),
		styles.DetailKeyStyle.Render(l.DetailADBStatus), styles.RenderStatusBadge(dev.State, l),
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
		styles.DetailKeyStyle.Render(l.DetailModelName), styles.DetailValueStyle.Render(modelVal),
		styles.DetailKeyStyle.Render(l.DetailCodename), styles.DetailValueStyle.Render(prodVal),
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
		styles.DetailKeyStyle.Render(l.DetailAndroidOS), styles.DetailValueStyle.Render(osVal),
		styles.DetailKeyStyle.Render(l.DetailIPAddress), styles.DetailValueStyle.Render(ipVal),
	)
	details = append(details, row3)

	// Row 4: Battery Details & Storage
	batVal := dev.BatteryDetailString()
	storageVal := dev.StorageString()
	row4 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render(l.DetailBattery), styles.DetailValueStyle.Render(batVal),
		styles.DetailKeyStyle.Render(l.DetailStorage), styles.DetailValueStyle.Render(storageVal),
	)
	details = append(details, row4)

	// Row 5: RAM & CPU Utilization with Progress Bars
	ramVal := dev.RAMString()
	if dev.RAMPercent > 0 {
		ramVal = fmt.Sprintf("%s %s", styles.RenderProgressBar(dev.RAMPercent, 8), dev.RAMString())
	}
	cpuVal := dev.CPUString()
	if dev.CPUPercent > 0 {
		cpuVal = fmt.Sprintf("%s %s", styles.RenderProgressBar(dev.CPUPercent, 8), dev.CPUString())
	}
	row5 := fmt.Sprintf(
		"%s %s    %s %s",
		styles.DetailKeyStyle.Render(l.DetailRAM), styles.DetailValueStyle.Render(ramVal),
		styles.DetailKeyStyle.Render(l.DetailCPU), styles.DetailValueStyle.Render(cpuVal),
	)
	details = append(details, row5)

	// Row 6: Port, Resolution & Uptime
	portVal := dev.USBPort
	if portVal == "" {
		portVal = "USB / TR-" + dev.TransportID
	}
	screenVal := dev.ScreenResolution
	if screenVal == "" {
		screenVal = "—"
	}
	uptimeVal := dev.Uptime
	if uptimeVal == "" {
		uptimeVal = "—"
	}
	row6 := fmt.Sprintf(
		"%s %s    %s %s    %s %s",
		styles.DetailKeyStyle.Render(l.DetailUSBPort), styles.DetailValueStyle.Render(portVal),
		styles.DetailKeyStyle.Render(l.DetailResolution), styles.DetailValueStyle.Render(screenVal),
		styles.DetailKeyStyle.Render(l.DetailUptime), styles.DetailValueStyle.Render(uptimeVal),
	)
	details = append(details, row6)

	return styles.DetailCardStyle.Render(strings.Join(details, "\n"))
}
