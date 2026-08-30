package ui

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	"github.com/Bloby22/andtls/internal/ui/components"
	"github.com/Bloby22/andtls/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming UI messages, key events, and asynchronous command responses
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	l := m.locale

	switch msg := msg.(type) {
	case tea.KeyMsg:
		keyStr := msg.String()

		// Global quit
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

		case ViewModeLogcat:
			switch keyStr {
			case "esc", "q":
				m.viewMode = ViewModeNormal
				return m, nil
			case " ":
				m.logcatPaused = !m.logcatPaused
				return m, nil
			case "c":
				if sel := m.SelectedDevice(); sel != nil {
					cmds = append(cmds, clearLogcatCmd(m.adbClient, sel.Serial))
				}
				return m, tea.Batch(cmds...)
			case "1":
				m.logcatFilterLevel = "ALL"
				return m, nil
			case "2":
				m.logcatFilterLevel = "D"
				return m, nil
			case "3":
				m.logcatFilterLevel = "I"
				return m, nil
			case "4":
				m.logcatFilterLevel = "W"
				return m, nil
			case "5":
				m.logcatFilterLevel = "E"
				return m, nil
			}

		case ViewModeFiles:
			switch keyStr {
			case "esc", "q":
				m.viewMode = ViewModeNormal
				return m, nil
			case "up", "k":
				if m.filesSelectedIndex > 0 {
					m.filesSelectedIndex--
				}
			case "down", "j":
				if m.filesSelectedIndex < len(m.filesList)-1 {
					m.filesSelectedIndex++
				}
			case "enter":
				if len(m.filesList) > 0 && m.filesSelectedIndex < len(m.filesList) {
					item := m.filesList[m.filesSelectedIndex]
					sel := m.SelectedDevice()
					if sel != nil {
						if item.IsDir {
							m.filesCurrentDir = item.Path
							m.statusMsg = fmt.Sprintf("Navigating to: %s", item.Path)
							cmds = append(cmds, fetchFilesCmd(m.adbClient, sel.Serial, item.Path))
						} else {
							m.statusMsg = fmt.Sprintf("Downloading: %s", item.Name)
							cmds = append(cmds, pullFileCmd(m.adbClient, sel.Serial, item.Path, "downloads"))
						}
					}
				}
			case "backspace", "b", "left":
				if m.filesCurrentDir != "/" && m.filesCurrentDir != "" {
					parentDir := path.Dir(strings.TrimSuffix(m.filesCurrentDir, "/"))
					if parentDir == "" {
						parentDir = "/"
					}
					m.filesCurrentDir = parentDir
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Navigating to: %s", parentDir)
						cmds = append(cmds, fetchFilesCmd(m.adbClient, sel.Serial, parentDir))
					}
				}
			case "d":
				if len(m.filesList) > 0 && m.filesSelectedIndex < len(m.filesList) {
					item := m.filesList[m.filesSelectedIndex]
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Downloading: %s", item.Name)
						cmds = append(cmds, pullFileCmd(m.adbClient, sel.Serial, item.Path, "downloads"))
					}
				}
			case "x":
				if len(m.filesList) > 0 && m.filesSelectedIndex < len(m.filesList) {
					item := m.filesList[m.filesSelectedIndex]
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Deleting: %s", item.Path)
						cmds = append(cmds, deleteFileCmd(m.adbClient, sel.Serial, item.Path, m.filesCurrentDir))
					}
				}
			}

		case ViewModePackages:
			switch keyStr {
			case "esc", "q":
				m.viewMode = ViewModeNormal
				return m, nil
			case "up", "k":
				if m.packagesSelectedIndex > 0 {
					m.packagesSelectedIndex--
				}
			case "down", "j":
				if m.packagesSelectedIndex < len(m.packagesList)-1 {
					m.packagesSelectedIndex++
				}
			case "tab":
				m.packagesThirdPartyOnly = !m.packagesThirdPartyOnly
				if sel := m.SelectedDevice(); sel != nil {
					cmds = append(cmds, fetchPackagesCmd(m.adbClient, sel.Serial, m.packagesThirdPartyOnly))
				}
			case "enter":
				if len(m.packagesList) > 0 && m.packagesSelectedIndex < len(m.packagesList) {
					pkg := m.packagesList[m.packagesSelectedIndex]
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Launching: %s", pkg)
						cmds = append(cmds, launchAppCmd(m.adbClient, sel.Serial, pkg))
					}
				}
			case "x":
				if len(m.packagesList) > 0 && m.packagesSelectedIndex < len(m.packagesList) {
					pkg := m.packagesList[m.packagesSelectedIndex]
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Stopping: %s", pkg)
						cmds = append(cmds, forceStopAppCmd(m.adbClient, sel.Serial, pkg))
					}
				}
			case "c":
				if len(m.packagesList) > 0 && m.packagesSelectedIndex < len(m.packagesList) {
					pkg := m.packagesList[m.packagesSelectedIndex]
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Clearing data: %s", pkg)
						cmds = append(cmds, clearAppDataCmd(m.adbClient, sel.Serial, pkg))
					}
				}
			case "u":
				if len(m.packagesList) > 0 && m.packagesSelectedIndex < len(m.packagesList) {
					pkg := m.packagesList[m.packagesSelectedIndex]
					if sel := m.SelectedDevice(); sel != nil {
						m.statusMsg = fmt.Sprintf("Uninstalling: %s", pkg)
						cmds = append(cmds, uninstallAppCmd(m.adbClient, sel.Serial, pkg))
					}
				}
			}

		case ViewModeRemote:
			sel := m.SelectedDevice()
			if sel == nil || !sel.IsAuthorized() {
				m.viewMode = ViewModeNormal
				return m, nil
			}

			switch keyStr {
			case "esc", "q":
				m.viewMode = ViewModeNormal
				return m, nil
			case "p":
				m.remoteLastAction = "Power Button"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 26, "Power"))
			case "+", "=":
				m.remoteLastAction = "Volume Up"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 24, "Vol+"))
			case "-", "_":
				m.remoteLastAction = "Volume Down"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 25, "Vol-"))
			case "h":
				m.remoteLastAction = "Home Button"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 3, "Home"))
			case "b":
				m.remoteLastAction = "Back Button"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 4, "Back"))
			case "r":
				m.remoteLastAction = "Recent Apps"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 187, "Recent Apps"))
			case "w":
				m.remoteLastAction = "Wake & Unlock"
				cmds = append(cmds, sendRemoteKeyCmd(m.adbClient, sel.Serial, 224, "Wake"), sendRemoteKeyCmd(m.adbClient, sel.Serial, 82, "Unlock"))
			case "t":
				m.remoteLastAction = "Input: Hello from andtls!"
				cmds = append(cmds, sendRemoteTextCmd(m.adbClient, sel.Serial, "Hello from andtls!"))
			}

		case ViewModeActions:
			actions := components.DefaultActions(l)
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

			case "m":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("mirror", sel)
				}

			case "v":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("record", sel)
				}

			case "f":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("files", sel)
				}

			case "l":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("logcat", sel)
				}

			case "p":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("packages", sel)
				}

			case "o":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("remote", sel)
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

			case "e":
				m.viewMode = ViewModeNormal
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("export", sel)
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
				m.statusMsg = l.MsgRefreshing
				m.isErr = false
				cmds = append(cmds, fetchDevicesCmd(m.adbClient))

			case "t":
				newTheme := styles.CycleTheme()
				m.statusMsg = fmt.Sprintf(l.MsgThemeChanged, newTheme)
				m.isErr = false

			case "L":
				_ = locale.Cycle()
				m.locale = locale.Get()
				m.statusMsg = m.locale.MsgLanguageChanged
				m.isErr = false

			case "?", "h":
				m.viewMode = ViewModeHelp

			case "enter", "a":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.viewMode = ViewModeActions
					m.selectedActionIndex = 0
				} else if sel != nil {
					m.statusMsg = l.MsgUnauthorizedAction
					m.isErr = true
				}

			case "m":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("mirror", sel)
				}

			case "v":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("record", sel)
				}

			case "f":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("files", sel)
				}

			case "l":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("logcat", sel)
				}

			case "p":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("packages", sel)
				}

			case "o":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("remote", sel)
				}

			case "e":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					return m.executeAction("export", sel)
				}

			case "s":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.statusMsg = fmt.Sprintf(l.MsgCapturingScreenshot, sel.DisplayName())
					m.isErr = false
					cmds = append(cmds, takeScreenshotCmd(m.adbClient, sel.Serial))
				} else {
					m.statusMsg = l.MsgCannotScreenshot
					m.isErr = true
				}

			case "w":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.statusMsg = fmt.Sprintf(l.MsgEnablingWireless, sel.DisplayName())
					m.isErr = false
					cmds = append(cmds, enableWirelessADBCmd(m.adbClient, sel.Serial, sel.IPAddress))
				} else {
					m.statusMsg = l.MsgCannotWireless
					m.isErr = true
				}

			case "b":
				if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
					m.viewMode = ViewModeActions
					m.selectedActionIndex = 9 // Point to reboot system
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case TickMsg:
		cmds = append(cmds, fetchDevicesCmd(m.adbClient))
		cmds = append(cmds, tickCmd(m.pollingInterval))

		// If viewing logcat and not paused, stream logs
		if m.viewMode == ViewModeLogcat && !m.logcatPaused {
			if sel := m.SelectedDevice(); sel != nil && sel.IsAuthorized() {
				cmds = append(cmds, fetchLogcatCmd(m.adbClient, sel.Serial, 80))
			}
		}

	case DevicesMsg:
		m.lastChecked = msg.Timestamp

		if msg.Err != nil {
			m.isErr = true
			m.statusMsg = fmt.Sprintf(l.MsgADBError, msg.Err)
			return m, tea.Batch(cmds...)
		}

		m.isErr = false

		// Detect changes
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
		} else if m.statusMsg == "" || m.statusMsg == l.MsgMonitoringNone || m.statusMsg == l.MsgRefreshing {
			if len(msg.Devices) == 0 {
				m.statusMsg = l.MsgMonitoringNone
			} else {
				m.statusMsg = fmt.Sprintf(l.MsgMonitoringFound, len(msg.Devices))
			}
		}

		// Preserve previously queried telemetry
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
				newDevices[i].RAMUsedMB = old.RAMUsedMB
				newDevices[i].RAMTotalMB = old.RAMTotalMB
				newDevices[i].RAMPercent = old.RAMPercent
				newDevices[i].CPUPercent = old.CPUPercent
				newDevices[i].Uptime = old.Uptime
				newDevices[i].AndroidVersion = old.AndroidVersion
				newDevices[i].SDKLevel = old.SDKLevel
				newDevices[i].ScreenResolution = old.ScreenResolution
			}
			if d.IsAuthorized() {
				cmds = append(cmds, fetchDeviceTelemetryCmd(m.adbClient, d.Serial))
			}
		}

		m.devices = newDevices

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
				if msg.RAMTotalMB > 0 {
					m.devices[i].RAMUsedMB = msg.RAMUsedMB
					m.devices[i].RAMTotalMB = msg.RAMTotalMB
					m.devices[i].RAMPercent = msg.RAMPercent
				}
				if msg.CPUPercent > 0 {
					m.devices[i].CPUPercent = msg.CPUPercent
				}
				if msg.Uptime != "" {
					m.devices[i].Uptime = msg.Uptime
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

	case LogcatMsg:
		if msg.Err == nil {
			m.logcatLines = msg.Lines
		}

	case FilesMsg:
		if msg.Err == nil {
			m.filesCurrentDir = msg.Dir
			m.filesList = msg.Files
			m.filesSelectedIndex = 0
		} else {
			m.isErr = true
			m.statusMsg = fmt.Sprintf("Files error: %v", msg.Err)
		}

	case PackagesMsg:
		if msg.Err == nil {
			m.packagesList = msg.Packages
			m.packagesSelectedIndex = 0
		} else {
			m.isErr = true
			m.statusMsg = fmt.Sprintf("Packages error: %v", msg.Err)
		}

	case ActionResultMsg:
		m.isErr = !msg.Success
		m.statusMsg = msg.Message
	}

	return m, tea.Batch(cmds...)
}

// executeAction dispatches the corresponding command for an action ID
func (m Model) executeAction(actionID string, dev *device.Device) (tea.Model, tea.Cmd) {
	l := m.locale
	var cmds []tea.Cmd

	switch actionID {
	case "mirror":
		cmds = append(cmds, launchMirrorCmd(m.adbClient, dev.Serial))

	case "record":
		if m.isRecording {
			m.isRecording = false
			m.statusMsg = "Stopping recording & saving video..."
			cmds = append(cmds, stopRecordCmd(m.adbClient, dev.Serial))
		} else {
			m.isRecording = true
			m.recordStartTime = time.Now()
			m.statusMsg = l.MsgRecordStarted
			cmds = append(cmds, startRecordCmd(m.adbClient, dev.Serial))
		}

	case "files":
		m.viewMode = ViewModeFiles
		m.filesCurrentDir = "/sdcard"
		cmds = append(cmds, fetchFilesCmd(m.adbClient, dev.Serial, "/sdcard"))

	case "logcat":
		m.viewMode = ViewModeLogcat
		m.logcatPaused = false
		cmds = append(cmds, fetchLogcatCmd(m.adbClient, dev.Serial, 80))

	case "packages":
		m.viewMode = ViewModePackages
		cmds = append(cmds, fetchPackagesCmd(m.adbClient, dev.Serial, m.packagesThirdPartyOnly))

	case "remote":
		m.viewMode = ViewModeRemote
		m.remoteLastAction = ""

	case "screenshot":
		m.statusMsg = fmt.Sprintf(l.MsgCapturingScreenshot, dev.DisplayName())
		m.isErr = false
		cmds = append(cmds, takeScreenshotCmd(m.adbClient, dev.Serial))

	case "wireless":
		m.statusMsg = fmt.Sprintf(l.MsgEnablingWireless, dev.DisplayName())
		m.isErr = false
		cmds = append(cmds, enableWirelessADBCmd(m.adbClient, dev.Serial, dev.IPAddress))

	case "export":
		m.statusMsg = "Exporting diagnostics report..."
		m.isErr = false
		cmds = append(cmds, exportReportCmd(m.adbClient, dev.Serial, *dev))

	case "reboot_system":
		m.statusMsg = fmt.Sprintf(l.MsgRebooting, dev.DisplayName())
		m.isErr = false
		cmds = append(cmds, rebootDeviceCmd(m.adbClient, dev.Serial, ""))

	case "reboot_recovery":
		m.statusMsg = fmt.Sprintf(l.MsgRebootingRecovery, dev.DisplayName())
		m.isErr = false
		cmds = append(cmds, rebootDeviceCmd(m.adbClient, dev.Serial, "recovery"))

	case "reboot_bootloader":
		m.statusMsg = fmt.Sprintf(l.MsgRebootingBootloader, dev.DisplayName())
		m.isErr = false
		cmds = append(cmds, rebootDeviceCmd(m.adbClient, dev.Serial, "bootloader"))
	}

	return m, tea.Batch(cmds...)
}

