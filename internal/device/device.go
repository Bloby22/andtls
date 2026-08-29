package device

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bloby22/andtls/internal/locale"
)

// DeviceState represents the operational state of an Android device according to ADB
type DeviceState string

const (
	// StateDevice indicates that the device is fully booted, connected, and authorized
	StateDevice DeviceState = "device"

	// StateUnauthorized indicates that USB debugging is not yet authorized on the device screen
	StateUnauthorized DeviceState = "unauthorized"

	// StateOffline indicates that the device is recognized by the daemon but not communicating
	StateOffline DeviceState = "offline"

	// StateAuthorizing indicates that the authorization handshake is in progress
	StateAuthorizing DeviceState = "authorizing"

	// StateRecovery indicates that the device is booted into Android Recovery mode
	StateRecovery DeviceState = "recovery"

	// StateSideload indicates that the device is in ADB sideload mode
	StateSideload DeviceState = "sideload"

	// StateBootloader indicates that the device is in fastboot / bootloader mode
	StateBootloader DeviceState = "bootloader"

	// StateUnknown indicates an unrecognized device state
	StateUnknown DeviceState = "unknown"
)

// Device represents an attached Android device with its hardware, network, and operational metadata
type Device struct {
	// Serial is the unique identifier or USB serial of the device
	Serial string `json:"serial"`

	// State is the current ADB connection state
	State DeviceState `json:"state"`

	// Product is the codename of the product (ro.product.name)
	Product string `json:"product,omitempty"`

	// Model is the commercial model name (ro.product.model)
	Model string `json:"model,omitempty"`

	// DeviceName is the device hardware name (ro.product.device)
	DeviceName string `json:"device_name,omitempty"`

	// TransportID is the ADB internal transport identifier
	TransportID string `json:"transport_id,omitempty"`

	// USBPort is the USB bus port string if available
	USBPort string `json:"usb_port,omitempty"`

	// Battery is the percentage level (0-100), or -1 if unknown
	Battery int `json:"battery"`

	// BatteryStatus describes charging state (e.g. Charging, Discharging, Full)
	BatteryStatus string `json:"battery_status,omitempty"`

	// BatteryTemp is the battery temperature in Celsius
	BatteryTemp float64 `json:"battery_temp,omitempty"`

	// BatteryVoltage is the battery voltage in millivolts
	BatteryVoltage int `json:"battery_voltage,omitempty"`

	// AndroidVersion is the OS version (e.g. "14", "13.0")
	AndroidVersion string `json:"android_version,omitempty"`

	// SDKLevel is the Android API SDK level (e.g. "34", "33")
	SDKLevel string `json:"sdk_level,omitempty"`

	// IPAddress is the local WiFi/Ethernet IP address of the device
	IPAddress string `json:"ip_address,omitempty"`

	// StorageFree is the available user storage string (e.g. "45.2 GB")
	StorageFree string `json:"storage_free,omitempty"`

	// StorageTotal is the total user storage string (e.g. "128 GB")
	StorageTotal string `json:"storage_total,omitempty"`

	// RAMUsedMB is the memory currently in use in megabytes
	RAMUsedMB int `json:"ram_used_mb,omitempty"`

	// RAMTotalMB is the total physical memory in megabytes
	RAMTotalMB int `json:"ram_total_mb,omitempty"`

	// RAMPercent is the memory consumption percentage (0-100)
	RAMPercent float64 `json:"ram_percent,omitempty"`

	// CPUPercent is the current CPU utilization percentage (0-100)
	CPUPercent float64 `json:"cpu_percent,omitempty"`

	// Uptime is the system uptime duration string
	Uptime string `json:"uptime,omitempty"`

	// InstalledAppsCount is the total number of installed application packages
	InstalledAppsCount int `json:"installed_apps_count,omitempty"`

	// ScreenResolution is the display resolution string (e.g. "1080x2400")
	ScreenResolution string `json:"screen_resolution,omitempty"`

	// LastSeen is the timestamp when the device was last observed
	LastSeen time.Time `json:"last_seen"`
}

// FileInfo represents a file or directory entry on the Android device filesystem
type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDir       bool   `json:"is_dir"`
	Size        string `json:"size"`
	Permissions string `json:"permissions"`
	ModTime     string `json:"mod_time"`
}

// NewDevice creates a new Device instance initialized with default values
func NewDevice(serial string, state DeviceState) Device {
	return Device{
		Serial:   serial,
		State:    state,
		Battery:  -1,
		LastSeen: time.Now(),
	}
}

// IsAuthorized returns true if the device is connected and authorized for ADB commands
func (d Device) IsAuthorized() bool {
	return d.State == StateDevice
}

// IsOnline returns true if the device responds to interactive ADB commands
func (d Device) IsOnline() bool {
	return d.State == StateDevice || d.State == StateRecovery
}

// DisplayName returns a formatted human-readable name for the device
func (d Device) DisplayName() string {
	if d.Model != "" {
		return strings.ReplaceAll(d.Model, "_", " ")
	}
	if d.Product != "" {
		return strings.ReplaceAll(d.Product, "_", " ")
	}
	if d.DeviceName != "" {
		return d.DeviceName
	}
	return d.Serial
}

// DisplayState returns a formatted string representation of the device state
func (d Device) DisplayState() string {
	l := locale.Get()
	switch d.State {
	case StateDevice:
		return l.StateConnected
	case StateUnauthorized:
		return l.StateUnauthorized
	case StateOffline:
		return l.StateOffline
	case StateAuthorizing:
		return l.StateAuthorizing
	case StateRecovery:
		return l.StateRecovery
	case StateSideload:
		return l.StateSideload
	case StateBootloader:
		return l.StateBootloader
	default:
		return string(d.State)
	}
}

// BatteryString returns a formatted battery status string
func (d Device) BatteryString() string {
	if d.Battery < 0 {
		return "—"
	}
	if d.BatteryStatus != "" && d.BatteryStatus != "Unknown" {
		return fmt.Sprintf("%d%% (%s)", d.Battery, d.BatteryStatus)
	}
	return fmt.Sprintf("%d%%", d.Battery)
}

// BatteryDetailString returns detailed battery metrics including temperature and voltage
func (d Device) BatteryDetailString() string {
	if d.Battery < 0 {
		return "—"
	}
	var parts []string
	parts = append(parts, fmt.Sprintf("%d%%", d.Battery))
	if d.BatteryStatus != "" && d.BatteryStatus != "Unknown" {
		parts = append(parts, d.BatteryStatus)
	}
	if d.BatteryTemp > 0 {
		parts = append(parts, fmt.Sprintf("%.1f°C", d.BatteryTemp))
	}
	if d.BatteryVoltage > 0 {
		parts = append(parts, fmt.Sprintf("%.2fV", float64(d.BatteryVoltage)/1000.0))
	}
	return strings.Join(parts, " • ")
}

// RAMString returns a formatted RAM memory utilization string
func (d Device) RAMString() string {
	if d.RAMTotalMB <= 0 {
		return "—"
	}
	usedGB := float64(d.RAMUsedMB) / 1024.0
	totalGB := float64(d.RAMTotalMB) / 1024.0
	return fmt.Sprintf("%.1f / %.1f GB (%.0f%%)", usedGB, totalGB, d.RAMPercent)
}

// CPUString returns a formatted CPU utilization percentage
func (d Device) CPUString() string {
	if d.CPUPercent <= 0 {
		return "—"
	}
	return fmt.Sprintf("%.1f%%", d.CPUPercent)
}

// StorageString returns a formatted storage usage string
func (d Device) StorageString() string {
	if d.StorageFree == "" && d.StorageTotal == "" {
		return "—"
	}
	if d.StorageFree != "" && d.StorageTotal != "" {
		return fmt.Sprintf("%s free / %s total", d.StorageFree, d.StorageTotal)
	}
	if d.StorageFree != "" {
		return fmt.Sprintf("%s free", d.StorageFree)
	}
	return d.StorageTotal
}

// OSVersionString returns the formatted Android version and SDK level
func (d Device) OSVersionString() string {
	if d.AndroidVersion == "" && d.SDKLevel == "" {
		return "—"
	}
	if d.AndroidVersion != "" && d.SDKLevel != "" {
		return fmt.Sprintf("Android %s (API %s)", d.AndroidVersion, d.SDKLevel)
	}
	if d.AndroidVersion != "" {
		return fmt.Sprintf("Android %s", d.AndroidVersion)
	}
	return fmt.Sprintf("API %s", d.SDKLevel)
}
