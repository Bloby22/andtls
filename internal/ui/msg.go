package ui

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Bloby22/andtls/internal/adb"
	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	tea "github.com/charmbracelet/bubbletea"
)

// TickMsg is fired periodically to poll for changes in connected devices
type TickMsg time.Time

// DevicesMsg contains the result of an asynchronous device query
type DevicesMsg struct {
	// Devices is the list of currently connected devices
	Devices []device.Device

	// Err is any error encountered while executing adb devices
	Err error

	// Timestamp is the time when the query completed
	Timestamp time.Time
}

// DeviceTelemetryMsg contains detailed telemetry results for a specific device
type DeviceTelemetryMsg struct {
	// Serial is the device serial identifier
	Serial string

	// Battery is the detailed battery info
	Battery adb.BatteryInfo

	// IPAddress is the discovered network IP
	IPAddress string

	// Storage is the internal storage metric
	Storage adb.StorageInfo

	// RAM metrics
	RAMUsedMB  int
	RAMTotalMB int
	RAMPercent float64

	// CPU utilization
	CPUPercent float64

	// System uptime
	Uptime string

	// Installed user apps count
	InstalledAppsCount int

	// Resolution is the display resolution
	Resolution string

	// AndroidVersion is the OS version
	AndroidVersion string

	// SDKLevel is the API SDK level
	SDKLevel string
}

// ActionResultMsg contains the result of an interactive action
type ActionResultMsg struct {
	// Success indicates if the action completed without error
	Success bool

	// Message is the outcome message to display in the status bar
	Message string
}

// LogcatMsg contains recent logcat lines for the active device
type LogcatMsg struct {
	Serial string
	Lines  []adb.LogcatLine
	Err    error
}

// FilesMsg contains directory listing results for file manager
type FilesMsg struct {
	Serial string
	Dir    string
	Files  []device.FileInfo
	Err    error
}

// PackagesMsg contains list of installed applications
type PackagesMsg struct {
	Serial   string
	Packages []string
	Err      error
}

// tickCmd returns a Bubble Tea command that fires a TickMsg after the specified interval
func tickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// fetchDevicesCmd executes an asynchronous query to list connected devices via ADB
func fetchDevicesCmd(client *adb.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		devices, err := client.ListDevices(ctx)
		return DevicesMsg{
			Devices:   devices,
			Err:       err,
			Timestamp: time.Now(),
		}
	}
}

// fetchDeviceTelemetryCmd executes asynchronous telemetry queries for a specific authorized device
func fetchDeviceTelemetryCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		msg := DeviceTelemetryMsg{Serial: serial}

		// Battery
		if bInfo, err := client.GetBatteryInfo(ctx, serial); err == nil {
			msg.Battery = bInfo
		}

		// IP Address
		if ip, err := client.GetIPAddress(ctx, serial); err == nil {
			msg.IPAddress = ip
		}

		// Storage
		if storage, err := client.GetStorageInfo(ctx, serial); err == nil {
			msg.Storage = storage
		}

		// RAM
		if usedMB, totalMB, percent, err := client.GetMemoryInfo(ctx, serial); err == nil {
			msg.RAMUsedMB = usedMB
			msg.RAMTotalMB = totalMB
			msg.RAMPercent = percent
		}

		// CPU
		if cpu, err := client.GetCPUUsage(ctx, serial); err == nil {
			msg.CPUPercent = cpu
		}

		// Uptime
		if uptime, err := client.GetUptime(ctx, serial); err == nil {
			msg.Uptime = uptime
		}

		// Resolution
		if res, err := client.GetDisplayResolution(ctx, serial); err == nil {
			msg.Resolution = res
		}

		// OS version & SDK
		if ver, err := client.GetProp(ctx, serial, "ro.build.version.release"); err == nil {
			msg.AndroidVersion = ver
		}
		if sdk, err := client.GetProp(ctx, serial, "ro.build.version.sdk"); err == nil {
			msg.SDKLevel = sdk
		}

		return msg
	}
}

// fetchLogcatCmd streams recent logcat lines
func fetchLogcatCmd(client *adb.Client, serial string, linesCount int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		lines, err := client.GetLogcat(ctx, serial, linesCount)
		return LogcatMsg{
			Serial: serial,
			Lines:  lines,
			Err:    err,
		}
	}
}

// clearLogcatCmd flushes device logcat buffer
func clearLogcatCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_ = client.ClearLogcat(ctx, serial)
		lines, err := client.GetLogcat(ctx, serial, 10)
		return LogcatMsg{
			Serial: serial,
			Lines:  lines,
			Err:    err,
		}
	}
}

// fetchFilesCmd queries file listing for device directory
func fetchFilesCmd(client *adb.Client, serial, dir string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		files, err := client.ListFiles(ctx, serial, dir)
		return FilesMsg{
			Serial: serial,
			Dir:    dir,
			Files:  files,
			Err:    err,
		}
	}
}

// pullFileCmd downloads a file from device
func pullFileCmd(client *adb.Client, serial, remotePath, localDestDir string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		l := locale.Get()
		fileName := filepath.Base(remotePath)
		localPath := filepath.Join(localDestDir, fileName)

		if err := client.PullFile(ctx, serial, remotePath, localPath); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgFileOpFailed, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgFileDownloaded, localPath),
		}
	}
}

// deleteFileCmd deletes a file on device
func deleteFileCmd(client *adb.Client, serial, remotePath, currentDir string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.RemoveFile(ctx, serial, remotePath); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgFileOpFailed, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgFileDeleted, remotePath),
		}
	}
}

// fetchPackagesCmd queries packages list
func fetchPackagesCmd(client *adb.Client, serial string, thirdPartyOnly bool) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()

		pkgs, err := client.ListPackages(ctx, serial, thirdPartyOnly)
		return PackagesMsg{
			Serial:   serial,
			Packages: pkgs,
			Err:      err,
		}
	}
}

// launchAppCmd starts application on device
func launchAppCmd(client *adb.Client, serial, pkgName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.LaunchApp(ctx, serial, pkgName); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgAppLaunched, pkgName),
		}
	}
}

// forceStopAppCmd stops an app on device
func forceStopAppCmd(client *adb.Client, serial, pkgName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.ForceStopApp(ctx, serial, pkgName); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgAppStopped, pkgName),
		}
	}
}

// clearAppDataCmd clears app data
func clearAppDataCmd(client *adb.Client, serial, pkgName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.ClearAppData(ctx, serial, pkgName); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgAppCleared, pkgName),
		}
	}
}

// uninstallAppCmd uninstalls an app
func uninstallAppCmd(client *adb.Client, serial, pkgName string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.Uninstall(ctx, serial, pkgName); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgAppUninstalled, pkgName),
		}
	}
}

// sendRemoteKeyCmd dispatches a hardware key press
func sendRemoteKeyCmd(client *adb.Client, serial string, keyCode int, label string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.SendKeyEvent(ctx, serial, keyCode); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgKeySent, label),
		}
	}
}

// sendRemoteTextCmd types text into focused input on device
func sendRemoteTextCmd(client *adb.Client, serial, text string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.SendInputText(ctx, serial, text); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgKeySent, fmt.Sprintf("Text: %s", text)),
		}
	}
}

// launchMirrorCmd starts scrcpy
func launchMirrorCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		l := locale.Get()
		if err := client.LaunchScrcpy(serial); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgMirrorFailed, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: l.MsgMirrorStarted,
		}
	}
}

// startRecordCmd begins video recording
func startRecordCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.StartScreenRecord(ctx, serial, "/sdcard/andtls_record.mp4"); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgRecordFailed, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: l.MsgRecordStarted,
		}
	}
}

// stopRecordCmd stops video recording and pulls MP4
func stopRecordCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		l := locale.Get()
		_ = client.StopScreenRecord(ctx, serial)
		time.Sleep(1 * time.Second) // allow MP4 container header to flush

		localPath := filepath.Join("recordings", fmt.Sprintf("recording_%s_%s.mp4", serial, time.Now().Format("20060102_150405")))
		if err := client.PullFile(ctx, serial, "/sdcard/andtls_record.mp4", localPath); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgRecordFailed, err),
			}
		}
		_ = client.RemoveFile(ctx, serial, "/sdcard/andtls_record.mp4")

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgRecordSaved, localPath),
		}
	}
}

// exportReportCmd exports diagnostics markdown report
func exportReportCmd(client *adb.Client, serial string, dev device.Device) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		l := locale.Get()
		filePath, err := client.ExportReport(ctx, serial, dev, "reports")
		if err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgADBError, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgReportExported, filePath),
		}
	}
}

// takeScreenshotCmd executes an asynchronous screenshot capture command
func takeScreenshotCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		l := locale.Get()
		filePath, err := client.TakeScreenshot(ctx, serial, "screenshots")
		if err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgScreenshotFailed, err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgScreenshotSaved, filePath),
		}
	}
}

// enableWirelessADBCmd sets device to TCP mode and optionally connects
func enableWirelessADBCmd(client *adb.Client, serial string, ip string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.EnableWirelessADB(ctx, serial, 5555); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgWirelessFailed, err),
			}
		}

		if ip != "" {
			_ = client.ConnectWirelessADB(ctx, ip, 5555)
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf(l.MsgWirelessEnabled, ip),
		}
	}
}

// rebootDeviceCmd triggers an asynchronous device reboot
func rebootDeviceCmd(client *adb.Client, serial, mode string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		l := locale.Get()
		if err := client.Reboot(ctx, serial, mode); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf(l.MsgRebootFailed, err),
			}
		}

		switch mode {
		case "recovery":
			return ActionResultMsg{
				Success: true,
				Message: fmt.Sprintf(l.MsgRebootToRecovery, serial),
			}
		case "bootloader":
			return ActionResultMsg{
				Success: true,
				Message: fmt.Sprintf(l.MsgRebootToBootloader, serial),
			}
		default:
			return ActionResultMsg{
				Success: true,
				Message: fmt.Sprintf(l.MsgRebootToSystem, serial),
			}
		}
	}
}

