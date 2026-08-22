package ui

import (
	"fmt"
	"time"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming UI messages, key events, and asynchronous command responses
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()

		// Global quit keys
		if keyStr == "ctrl+c" {
			return m, tea.Quit
		}

		// Handle key events based on current ViewMode
		switch m.viewMode {
		case ViewModeHelp:
			switch keyStr {
			case "esc", "q", "?", "h", "enter":
				m.viewMode = ViewModeNormal
				return m, nil
			}

		case ViewModeActions:
			actions := components.DefaultActions()
			switch keyStr {
			case "esc", "q":
				m.viewMode = ViewModeNormal
				return m, nil

			case "up", "k":
				if m.selectedActionIndex > 0 {
					m.selectedActionIndex--
				}

			case "down", "j":
				if m.selectedActionIndex < len(actions)-1 {
					m.selectedActionIndex++
				}

			case "enter":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					action := actions[m.selectedActionIndex]
					return m.executeAction(action.ID, sel)
				}

			case "s":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("screenshot", sel)
				}

			case "w":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("wireless", sel)
				}

			case "1":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("reboot_system", sel)
				}

			case "2":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("reboot_recovery", sel)
				}

			case "3":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("reboot_bootloader", sel)
				}
			}

		case ViewModeNormal:
			switch keyStr {
			case "q":
				return m, tea.Quit

			case "up", "k":
				if len(m.devices) > 0 && m.selectedIndex > 0 {
					m.selectedIndex--
				}

			case "down", "j":
				if len(m.devices) > 0 && m.selectedIndex < len(m.devices)-1 {
					m.selectedIndex++
				}

			case "r":
				m.statusMsg = "Refreshing devices..."
				m.isErr = false
				cmds = append(cmds, fetchDevicesCmd(m.adbClient))

			case "?", "h":
				m.viewMode = ViewModeHelp

			case "enter", "a":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.viewMode = ViewModeActions
					m.selectedActionIndex = 0
				} else if sel != nil {
					m.statusMsg = "Device must be authorized to perform actions"
					m.isErr = true
				}

			case "s":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.statusMsg = fmt.Sprintf("Capturing screenshot for %s...", sel.DisplayName())
					m.isErr = false
					cmds = append(cmds, takeScreenshotCmd(m.adbClient, sel.Serial))
				} else {
					m.statusMsg = "Cannot capture screenshot: device not authorized"
					m.isErr = true
				}

			case "w":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.statusMsg = fmt.Sprintf("Enabling wireless ADB on port 5555 for %s...", sel.DisplayName())
					m.isErr = false
					cmds = append(cmds, enableWirelessADBCmd(m.adbClient, sel.Serial, sel.IPAddress))
				} else {
					m.statusMsg = "Cannot enable wireless ADB: device not authorized"
					m.isErr = true
				}

			case "b":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.viewMode = ViewModeActions
					m.selectedActionIndex = 2 // Point to reboot system
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case TickMsg:
		// Execute background device poll and schedule next tick
		cmds = append(cmds, fetchDevicesCmd(m.adbClient))
		cmds = append(cmds, tickCmd(m.pollingInterval))

	case DevicesMsg:
		m.lastChecked = msg.Timestamp

		if msg.Err != nil {
			m.isErr = true
			m.statusMsg = fmt.Sprintf("ADB Error: %v", msg.Err)
			return m, tea.Batch(cmds...)
		}

		m.isErr = false

		// Detect changes against the previous device snapshot
		events := device.Diff(m.devices, msg.Devices)
		if len(events) > 0 {
			latestEvent := events[len(events)-1]
			m.statusMsg = latestEvent.Message
			for _, ev := range events {
				m.lastEvents = append(m.lastEvents, fmt.Sprintf("[%s] %s",
					time.Now().Format("15:04:05"), ev.Message))
			}
			if len(m.lastEvents) > 10 {
				m.lastEvents = m.lastEvents[len(m.lastEvents)-10:]
			}
		} else if m.statusMsg == "" || m.statusMsg == "Monitoring devices..." || m.statusMsg == "Refreshing devices..." {
			if len(msg.Devices) == 0 {
				m.statusMsg = "Monitoring devices (1s polling)... No devices connected"
			} else {
				m.statusMsg = fmt.Sprintf("Monitoring devices (1s polling)... Found %d device(s)", len(msg.Devices))
			}
		}

		// Preserve previously queried telemetry for existing devices
		oldTelemetryMap := make(map[string]device.Device)
		for _, d := range m.devices {
			oldTelemetryMap[d.Serial] = d
		}

		newDevices := make([]device.Device, len(msg.Devices))
		for i, d := range msg.Devices {
			newDevices[i] = d
			if old, ok := oldTelemetryMap[d.Serial]; ok {
				newDevices[i].Battery = old.Battery
				newDevices[i].BatteryStatus = old.BatteryStatus
				newDevices[i].BatteryTemp = old.BatteryTemp
				newDevices[i].BatteryVoltage = old.BatteryVoltage
				newDevices[i].IPAddress = old.IPAddress
				newDevices[i].StorageFree = old.StorageFree
				newDevices[i].StorageTotal = old.StorageTotal
				newDevices[i].AndroidVersion = old.AndroidVersion
				newDevices[i].SDKLevel = old.SDKLevel
				newDevices[i].ScreenResolution = old.ScreenResolution
			}
			// Trigger async telemetry query for authorized devices
			if d.IsAuthorized() {
				cmds = append(cmds, fetchDeviceTelemetryCmd(m.adbClient, d.Serial))
			}
		}

		m.devices = newDevices

		// Adjust selection index within bounds
		if len(m.devices) == 0 {
			m.selectedIndex = 0
		} else if m.selectedIndex >= len(m.devices) {
			m.selectedIndex = len(m.devices) - 1
		}

	case DeviceTelemetryMsg:
		for i := range m.devices {
			if m.devices[i].Serial == msg.Serial {
				if msg.Battery.Level >= 0 {
					m.devices[i].Battery = msg.Battery.Level
					m.devices[i].BatteryStatus = msg.Battery.Status
					m.devices[i].BatteryTemp = msg.Battery.Temperature
					m.devices[i].BatteryVoltage = msg.Battery.Voltage
				}
				if msg.IPAddress != "" {
					m.devices[i].IPAddress = msg.IPAddress
				}
				if msg.Storage.Total != "" {
					m.devices[i].StorageTotal = msg.Storage.Total
					m.devices[i].StorageFree = msg.Storage.Free
				}
				if msg.Resolution != "" {
					m.devices[i].ScreenResolution = msg.Resolution
				}
				if msg.AndroidVersion != "" {
					m.devices[i].AndroidVersion = msg.AndroidVersion
				}
				if msg.SDKLevel != "" {
					m.devices[i].SDKLevel = msg.SDKLevel
				}
				break
			}
		}

	case ActionResultMsg:
		m.isErr = !msg.Success
		m.statusMsg = msg.Message
	}

	return m, tea.Batch(cmds...)
}

// executeAction dispatches the corresponding command for an action ID
func (m Model) executeAction(actionID string, dev *device.Device) (tea.Model, tea.Cmd) {
	switch actionID {
	case "screenshot":
		m.statusMsg = fmt.Sprintf("Capturing screenshot for %s...", dev.DisplayName())
		m.isErr = false
		return m, takeScreenshotCmd(m.adbClient, dev.Serial)

	case "wireless":
		m.statusMsg = fmt.Sprintf("Enabling wireless ADB on port 5555 for %s...", dev.DisplayName())
		m.isErr = false
		return m, enableWirelessADBCmd(m.adbClient, dev.Serial, dev.IPAddress)

	case "reboot_system":
		m.statusMsg = fmt.Sprintf("Rebooting %s...", dev.DisplayName())
		m.isErr = false
		return m, rebootDeviceCmd(m.adbClient, dev.Serial, "")

	case "reboot_recovery":
		m.statusMsg = fmt.Sprintf("Rebooting %s to Recovery...", dev.DisplayName())
		m.isErr = false
		return m, rebootDeviceCmd(m.adbClient, dev.Serial, "recovery")

	case "reboot_bootloader":
		m.statusMsg = fmt.Sprintf("Rebooting %s to Bootloader...", dev.DisplayName())
		m.isErr = false
		return m, rebootDeviceCmd(m.adbClient, dev.Serial, "bootloader")
	}

	return m, nil
}
