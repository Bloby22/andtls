package adb

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
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

// LogcatLine represents a structured log entry from Android logcat
type LogcatLine struct {
	// Timestamp is the log entry time (e.g. "08-29 19:15:30.123")
	Timestamp string

	// Level is the severity level: V (Verbose), D (Debug), I (Info), W (Warn), E (Error), F (Fatal)
	Level string

	// Tag is the logging component/class tag
	Tag string

	// PID is the originating process identifier
	PID string

	// TID is the originating thread identifier
	TID string

	// Message is the actual log message content
	Message string

	// Raw is the original unprocessed line
	Raw string
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

// ParsePackagesOutput parses `pm list packages` output into a sorted list of package names
func ParsePackagesOutput(rawOutput string) []string {
	var packages []string
	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "package:") {
			pkg := strings.TrimPrefix(line, "package:")
			pkg = strings.TrimSpace(pkg)
			if pkg != "" {
				packages = append(packages, pkg)
			}
		}
	}
	sort.Strings(packages)
	return packages
}

// ParseMemInfo parses `/proc/meminfo` into used MB, total MB, and memory usage percentage
func ParseMemInfo(rawOutput string) (usedMB, totalMB int, percent float64) {
	var totalKB, freeKB, buffersKB, cachedKB, availKB int
	hasAvail := false

	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			val, _ := strconv.Atoi(parts[1])
			switch parts[0] {
			case "MemTotal:":
				totalKB = val
			case "MemFree:":
				freeKB = val
			case "MemAvailable:":
				availKB = val
				hasAvail = true
			case "Buffers:":
				buffersKB = val
			case "Cached:":
				cachedKB = val
			}
		}
	}

	if totalKB <= 0 {
		return 0, 0, 0
	}

	totalMB = totalKB / 1024
	var usedKB int
	if hasAvail && availKB > 0 {
		usedKB = totalKB - availKB
	} else {
		usedKB = totalKB - (freeKB + buffersKB + cachedKB)
	}

	if usedKB < 0 {
		usedKB = 0
	}
	usedMB = usedKB / 1024
	percent = (float64(usedKB) / float64(totalKB)) * 100.0
	if percent > 100.0 {
		percent = 100.0
	}

	return usedMB, totalMB, percent
}

// ParseCPUInfo parses `dumpsys cpuinfo` output for total CPU utilization percentage
func ParseCPUInfo(rawOutput string) float64 {
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)%\s+TOTAL`)
	match := re.FindStringSubmatch(rawOutput)
	if len(match) >= 2 {
		if val, err := strconv.ParseFloat(match[1], 64); err == nil {
			return val
		}
	}
	return 0.0
}

// ParseUptime parses `/proc/uptime` or `uptime` output into a human-readable duration
func ParseUptime(rawOutput string) string {
	rawOutput = strings.TrimSpace(rawOutput)
	if rawOutput == "" {
		return ""
	}

	// Try /proc/uptime format: "123456.78 987654.32"
	fields := strings.Fields(rawOutput)
	if len(fields) >= 1 {
		if secFloat, err := strconv.ParseFloat(fields[0], 64); err == nil {
			totalSeconds := int(secFloat)
			days := totalSeconds / 86400
			hours := (totalSeconds % 86400) / 3600
			minutes := (totalSeconds % 3600) / 60

			if days > 0 {
				return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
			}
			if hours > 0 {
				return fmt.Sprintf("%dh %dm", hours, minutes)
			}
			return fmt.Sprintf("%dm", minutes)
		}
	}

	return rawOutput
}

var (
	reLogcatThreadTime = regexp.MustCompile(`^(\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\.\d{3})\s+(\d+)\s+(\d+)\s+([VDIWEF])\s+([^:]+):\s*(.*)$`)
	reLogcatTimeBrief  = regexp.MustCompile(`^(\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}\.\d{3})\s+([VDIWEF])/([^\(]+)\(\s*(\d+)\):\s*(.*)$`)
	reLogcatBrief      = regexp.MustCompile(`^([VDIWEF])/([^\(]+)\(\s*(\d+)\):\s*(.*)$`)
)

// ParseLogcatLine parses a single logcat line into a structured LogcatLine
func ParseLogcatLine(line string) LogcatLine {
	line = strings.TrimRight(line, "\r\n")

	// 1. ThreadTime format: "08-29 19:15:30.123  1234  5678 I TagName: Message"
	if m := reLogcatThreadTime.FindStringSubmatch(line); len(m) == 7 {
		return LogcatLine{
			Timestamp: m[1],
			PID:       m[2],
			TID:       m[3],
			Level:     m[4],
			Tag:       strings.TrimSpace(m[5]),
			Message:   m[6],
			Raw:       line,
		}
	}

	// 2. TimeBrief format: "08-29 19:15:30.123 I/TagName( 1234): Message"
	if m := reLogcatTimeBrief.FindStringSubmatch(line); len(m) == 6 {
		return LogcatLine{
			Timestamp: m[1],
			Level:     m[2],
			Tag:       strings.TrimSpace(m[3]),
			PID:       m[4],
			Message:   m[5],
			Raw:       line,
		}
	}

	// 3. Brief format: "I/TagName( 1234): Message"
	if m := reLogcatBrief.FindStringSubmatch(line); len(m) == 5 {
		return LogcatLine{
			Level:   m[1],
			Tag:     strings.TrimSpace(m[2]),
			PID:     m[3],
			Message: m[4],
			Raw:     line,
		}
	}

	// Fallback for unformatted lines
	level := "I"
	if strings.Contains(line, " E ") || strings.HasPrefix(line, "E/") {
		level = "E"
	} else if strings.Contains(line, " W ") || strings.HasPrefix(line, "W/") {
		level = "W"
	} else if strings.Contains(line, " D ") || strings.HasPrefix(line, "D/") {
		level = "D"
	}

	return LogcatLine{
		Level:   level,
		Tag:     "system",
		Message: line,
		Raw:     line,
	}
}

// ParseLogcatLines parses a multiline logcat dump into structured LogcatLine entries
func ParseLogcatLines(rawOutput string) []LogcatLine {
	var results []LogcatLine
	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--------- beginning of") {
			continue
		}
		results = append(results, ParseLogcatLine(line))
	}
	return results
}

// formatFileSize formats a raw byte count integer or string into human-readable size
func formatFileSize(sizeStr string) string {
	bytes, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return sizeStr
	}
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ParseFilesOutput parses raw `ls -la` or `ls -l` shell output into a slice of FileInfo structs
func ParseFilesOutput(rawOutput, currentDir string) []device.FileInfo {
	var dirs []device.FileInfo
	var files []device.FileInfo

	lines := strings.Split(rawOutput, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "total ") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		perms := fields[0]
		isDir := strings.HasPrefix(perms, "d")

		// Extract filename (might contain spaces at the end)
		// Standard Android toybox ls -la:
		// [0]drwxrwx--x [1]4 [2]root [3]sdcard_rw [4]4096 [5]2026-08-29 [6]12:00 [7...]Name
		// Or [0]drwxrwx--x [1]root [2]root [3]4096 [4]Aug [5]29 [6]12:00 [7...]Name
		var sizeStr string
		var name string
		var modTime string

		if len(fields) >= 8 && (strings.Contains(fields[5], "-") || strings.Contains(fields[6], ":")) {
			sizeStr = fields[4]
			modTime = fmt.Sprintf("%s %s", fields[5], fields[6])
			name = strings.Join(fields[7:], " ")
		} else if len(fields) >= 7 {
			sizeStr = fields[3]
			modTime = fmt.Sprintf("%s %s %s", fields[4], fields[5], fields[6])
			if len(fields) >= 8 {
				name = strings.Join(fields[7:], " ")
			} else {
				name = fields[len(fields)-1]
			}
		} else {
			name = fields[len(fields)-1]
			sizeStr = "-"
		}

		// Clean up symlinks format "name -> target"
		if idx := strings.Index(name, " -> "); idx != -1 {
			name = name[:idx]
		}
		name = strings.TrimSpace(name)

		if name == "." || name == ".." || name == "" {
			continue
		}

		filePath := currentDir
		if !strings.HasSuffix(filePath, "/") {
			filePath += "/"
		}
		filePath += name

		displaySize := "-"
		if !isDir {
			displaySize = formatFileSize(sizeStr)
		}

		item := device.FileInfo{
			Name:        name,
			Path:        filePath,
			IsDir:       isDir,
			Size:        displaySize,
			Permissions: perms,
			ModTime:     modTime,
		}

		if isDir {
			dirs = append(dirs, item)
		} else {
			files = append(files, item)
		}
	}

	// Sort directories and files alphabetically
	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return append(dirs, files...)
}
