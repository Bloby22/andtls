package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/Bloby22/andtls/internal/adb"
	"github.com/Bloby22/andtls/internal/device"
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

	// Resolution is the display resolution
	Resolution string

	// AndroidVersion is the OS version
	AndroidVersion string

	// SDKLevel is the API SDK level
	SDKLevel string
}

// ActionResultMsg contains the result of an interactive action (screenshot, reboot, wireless ADB)
type ActionResultMsg struct {
	// Success indicates if the action completed without error
	Success bool

	// Message is the outcome message to display in the status bar
	Message string
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
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		msg := DeviceTelemetryMsg{Serial: serial}

		// Query battery info
		if bInfo, err := client.GetBatteryInfo(ctx, serial); err == nil {
			msg.Battery = bInfo
		}

		// Query IP address
		if ip, err := client.GetIPAddress(ctx, serial); err == nil {
			msg.IPAddress = ip
		}

		// Query storage info
		if storage, err := client.GetStorageInfo(ctx, serial); err == nil {
			msg.Storage = storage
		}

		// Query resolution
		if res, err := client.GetDisplayResolution(ctx, serial); err == nil {
			msg.Resolution = res
		}

		// Query Android version & SDK
		if ver, err := client.GetProp(ctx, serial, "ro.build.version.release"); err == nil {
			msg.AndroidVersion = ver
		}
		if sdk, err := client.GetProp(ctx, serial, "ro.build.version.sdk"); err == nil {
			msg.SDKLevel = sdk
		}

		return msg
	}
}

// takeScreenshotCmd executes an asynchronous screenshot capture command
func takeScreenshotCmd(client *adb.Client, serial string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		filePath, err := client.TakeScreenshot(ctx, serial, "screenshots")
		if err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf("Screenshot failed: %v", err),
			}
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf("Screenshot saved to: %s", filePath),
		}
	}
}

// enableWirelessADBCmd sets device to TCP mode and optionally connects
func enableWirelessADBCmd(client *adb.Client, serial string, ip string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.EnableWirelessADB(ctx, serial, 5555); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf("Wireless ADB failed: %v", err),
			}
		}

		if ip != "" {
			_ = client.ConnectWirelessADB(ctx, ip, 5555)
		}

		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf("Wireless ADB enabled on port 5555 (Connect with: adb connect %s:5555)", ip),
		}
	}
}

// rebootDeviceCmd triggers an asynchronous device reboot
func rebootDeviceCmd(client *adb.Client, serial, mode string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Reboot(ctx, serial, mode); err != nil {
			return ActionResultMsg{
				Success: false,
				Message: fmt.Sprintf("Reboot failed: %v", err),
			}
		}

		modeName := mode
		if modeName == "" {
			modeName = "system"
		}
		return ActionResultMsg{
			Success: true,
			Message: fmt.Sprintf("Device %s rebooting to %s...", serial, modeName),
		}
	}
}
