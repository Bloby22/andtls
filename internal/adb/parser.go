package adb

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
)

// BatteryInfo holds structured battery telemetry data
type BatteryInfo struct {
	// Level is the battery percentage (0-100)
	Level int

	// Status is the charging status description
	Status string

	// Temperature is the temperature in degrees Celsius
	Temperature float64

	// Voltage is the battery voltage in millivolts
	Voltage int
}

// StorageInfo holds storage metrics for a filesystem path
type StorageInfo struct {
	// Total is the total capacity string
	Total string

	// Free is the available free storage string
	Free string

	// Used is the consumed storage string
	Used string
}

// ParseDevicesOutput parses the raw text output of `adb devices -l` into a slice of Device objects
func ParseDevicesOutput(rawOutput string) []device.Device {
	var devices []device.Device
	lines := strings.Split(rawOutput, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Ignore headers, daemon initialization messages and warnings
		if strings.HasPrefix(line, "List of devices attached") ||
			strings.HasPrefix(line, "* daemon") ||
			strings.HasPrefix(line, "--- adb") ||
			strings.HasPrefix(line, "adb:") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		serial := fields[0]
		state := device.DeviceState(fields[1])

		dev := device.NewDevice(serial, state)

		// Parse key:value metadata (product:..., model:..., device:..., transport_id:..., usb:...)
		for _, field := range fields[2:] {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) != 2 {
				continue
			}
			key, val := parts[0], parts[1]
			switch key {
			case "product":
				dev.Product = val
			case "model":
				dev.Model = val
			case "device":
				dev.DeviceName = val
			case "transport_id":
				dev.TransportID = val
			case "usb":
				dev.USBPort = val
			}
		}

		devices = append(devices, dev)
	}

	return devices
}

// ParseBatteryOutput parses raw `dumpsys battery` text output into percentage level and charging status string
func ParseBatteryOutput(rawOutput string) (int, string, error) {
	info, err := ParseBatteryInfo(rawOutput)
	if err != nil {
		return -1, "", err
	}
	return info.Level, info.Status, nil
}

// ParseBatteryInfo parses raw `dumpsys battery` text output into a structured BatteryInfo struct
func ParseBatteryInfo(rawOutput string) (BatteryInfo, error) {
	info := BatteryInfo{
		Level:  -1,
		Status: locale.Get().BatteryUnknown,
	}

	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "level:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					info.Level = val
				}
			}
		} else if strings.HasPrefix(line, "status:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Android BatteryManager status codes:
				// 1: BATTERY_STATUS_UNKNOWN
				// 2: BATTERY_STATUS_CHARGING
				// 3: BATTERY_STATUS_DISCHARGING
				// 4: BATTERY_STATUS_NOT_CHARGING
				// 5: BATTERY_STATUS_FULL
				l := locale.Get()
				switch parts[1] {
				case "2":
					info.Status = l.BatteryCharging
				case "3":
					info.Status = l.BatteryDischarging
				case "4":
					info.Status = l.BatteryNotCharging
				case "5":
					info.Status = l.BatteryFull
				default:
					info.Status = l.BatteryUnknown
				}
			}
		} else if strings.HasPrefix(line, "temperature:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					info.Temperature = float64(val) / 10.0
				}
			}
		} else if strings.HasPrefix(line, "voltage:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if val, err := strconv.Atoi(parts[1]); err == nil {
					info.Voltage = val
				}
			}
		}
	}

	if info.Level < 0 {
		return info, errors.New("failed to parse battery level from output")
	}

	return info, nil
}

// ParseIPAddress parses raw network output (e.g. `ip route` or `ip addr`) into the primary IPv4 address
func ParseIPAddress(rawOutput string) string {
	// Check for wlan0 / wlan inet match
	reInet := regexp.MustCompile(`inet\s+([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)/`)
	matches := reInet.FindAllStringSubmatch(rawOutput, -1)
	for _, m := range matches {
		if len(m) >= 2 && !strings.HasPrefix(m[1], "127.") {
			return m[1]
		}
	}

	// Check for generic IP format (e.g. from ip route src ...)
	reSrc := regexp.MustCompile(`src\s+([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
	srcMatch := reSrc.FindStringSubmatch(rawOutput)
	if len(srcMatch) >= 2 && !strings.HasPrefix(srcMatch[1], "127.") {
		return srcMatch[1]
	}

	return ""
}

// ParseStorageOutput parses `df -h /data` or `df /data` output into StorageInfo
func ParseStorageOutput(rawOutput string) (StorageInfo, error) {
	info := StorageInfo{}
	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Filesystem") || line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 5 {
			// Typical df -h output: Filesystem Size Used Avail Use% Mounted on
			info.Total = fields[1]
			info.Used = fields[2]
			info.Free = fields[3]
			return info, nil
		}
	}
	return info, errors.New("failed to parse storage output")
}
