package adb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Bloby22/andtls/internal/device"
)

var (
	// ErrADBNotFound indicates that the adb executable was not found in PATH
	ErrADBNotFound = errors.New("adb was not found in system PATH (make sure Android SDK platform-tools are installed)")

	// ErrTimeout indicates that an ADB operation exceeded its execution timeout
	ErrTimeout = errors.New("adb operation timed out")
)

// Client provides methods to execute ADB commands against connected devices
type Client struct {
	adbPath        string
	defaultTimeout time.Duration
}

// NewClient returns a new ADB client using default settings ("adb" binary, 5s timeout)
func NewClient() *Client {
	return &Client{
		adbPath:        "adb",
		defaultTimeout: 5 * time.Second,
	}
}

// NewClientWithPath returns a new ADB client using a specified binary path and timeout
func NewClientWithPath(adbPath string, timeout time.Duration) *Client {
	if adbPath == "" {
		adbPath = "adb"
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Client{
		adbPath:        adbPath,
		defaultTimeout: timeout,
	}
}

// CheckADB verifies whether the ADB executable is available and operational
func (c *Client) CheckADB(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "version")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return ErrTimeout
		}
		return fmt.Errorf("%w: %v (%s)", ErrADBNotFound, err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// ListDevices executes `adb devices -l` and returns the parsed list of connected devices
func (c *Client) ListDevices(ctx context.Context) ([]device.Device, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "devices", "-l")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, ErrTimeout
		}
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return nil, fmt.Errorf("failed to run 'adb devices -l': %s", errMsg)
	}

	devices := ParseDevicesOutput(stdout.String())
	return devices, nil
}

// GetBattery queries battery level and status via `dumpsys battery`
func (c *Client) GetBattery(ctx context.Context, serial string) (int, string, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "dumpsys", "battery")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return -1, "", fmt.Errorf("dumpsys battery failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return ParseBatteryOutput(stdout.String())
}

// GetBatteryInfo queries detailed battery metrics (level, status, temperature, voltage)
func (c *Client) GetBatteryInfo(ctx context.Context, serial string) (BatteryInfo, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "dumpsys", "battery")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return BatteryInfo{Level: -1, Status: "Unknown"}, fmt.Errorf("dumpsys battery failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return ParseBatteryInfo(stdout.String())
}

// GetProp queries an Android system property via `getprop <prop>`
func (c *Client) GetProp(ctx context.Context, serial, prop string) (string, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "getprop", prop)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("getprop failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

// GetIPAddress discovers the primary network IP address of the target device
func (c *Client) GetIPAddress(ctx context.Context, serial string) (string, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "ip", "addr", "show", "wlan0")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err == nil {
		if ip := ParseIPAddress(stdout.String()); ip != "" {
			return ip, nil
		}
	}

	// Fallback to ip route
	cmdRoute := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "ip", "route")
	var stdoutRoute bytes.Buffer
	cmdRoute.Stdout = &stdoutRoute

	if err := cmdRoute.Run(); err == nil {
		if ip := ParseIPAddress(stdoutRoute.String()); ip != "" {
			return ip, nil
		}
	}

	return "", errors.New("no IP address found")
}

// GetStorageInfo queries storage capacity and available space on the device /data partition
func (c *Client) GetStorageInfo(ctx context.Context, serial string) (StorageInfo, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "df", "-h", "/data")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return StorageInfo{}, fmt.Errorf("df failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return ParseStorageOutput(stdout.String())
}

// GetDisplayResolution queries the physical screen resolution of the device
func (c *Client) GetDisplayResolution(ctx context.Context, serial string) (string, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "wm", "size")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Output format: "Physical size: 1080x2400"
	out := strings.TrimSpace(stdout.String())
	parts := strings.Split(out, ":")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[1]), nil
	}
	return out, nil
}

// TakeScreenshot captures a device screenshot and saves it locally into a destination directory
func (c *Client) TakeScreenshot(ctx context.Context, serial, destDir string) (string, error) {
	if destDir == "" {
		destDir = "screenshots"
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create screenshot directory: %w", err)
	}

	fileName := fmt.Sprintf("screenshot_%s_%s.png", serial, time.Now().Format("20060102_150405"))
	destPath := filepath.Join(destDir, fileName)

	// Stream screencap directly from adb exec-out
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "exec-out", "screencap", "-p")
	outFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination image file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		_ = os.Remove(destPath)
		return "", fmt.Errorf("screencap failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return destPath, nil
}

// EnableWirelessADB sets the device ADB daemon to listen on TCP port (default: 5555)
func (c *Client) EnableWirelessADB(ctx context.Context, serial string, port int) error {
	if port <= 0 {
		port = 5555
	}
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "tcpip", fmt.Sprintf("%d", port))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to enable TCP mode: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// ConnectWirelessADB connects the local ADB client to a wireless device via IP:port
func (c *Client) ConnectWirelessADB(ctx context.Context, ip string, port int) error {
	if port <= 0 {
		port = 5555
	}
	target := fmt.Sprintf("%s:%d", ip, port)
	cmd := exec.CommandContext(ctx, c.adbPath, "connect", target)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("adb connect failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	if strings.Contains(stdout.String(), "unable") || strings.Contains(stdout.String(), "failed") {
		return fmt.Errorf("connection refused: %s", strings.TrimSpace(stdout.String()))
	}
	return nil
}

// Shell runs an arbitrary shell command on the target device
func (c *Client) Shell(ctx context.Context, serial, command string) (string, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("adb shell failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}

// Reboot restarts the device into the specified mode (e.g. "", "recovery", "bootloader")
func (c *Client) Reboot(ctx context.Context, serial, mode string) error {
	args := []string{"-s", serial, "reboot"}
	if mode != "" {
		args = append(args, mode)
	}
	cmd := exec.CommandContext(ctx, c.adbPath, args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("reboot failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// Install installs an APK file onto the target device
func (c *Client) Install(ctx context.Context, serial, apkPath string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "install", "-r", apkPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// Uninstall removes an installed application package by package name
func (c *Client) Uninstall(ctx context.Context, serial, pkgName string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "uninstall", pkgName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("uninstall failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// LaunchScrcpy starts the scrcpy screen mirroring tool for the specified device in a detached process
func (c *Client) LaunchScrcpy(serial string) error {
	scrcpyPath, err := exec.LookPath("scrcpy")
	if err != nil {
		return fmt.Errorf("scrcpy not found in PATH (Install via: 'winget install Genymobile.scrcpy' on Windows or 'brew install scrcpy' on macOS)")
	}

	cmd := exec.Command(scrcpyPath, "-s", serial, "--stay-awake", "--window-title", fmt.Sprintf("andtls: %s", serial))
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start scrcpy: %w", err)
	}

	// Release process in background
	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

// ListFiles returns entries for a directory on the device
func (c *Client) ListFiles(ctx context.Context, serial, remoteDir string) ([]device.FileInfo, error) {
	if remoteDir == "" {
		remoteDir = "/sdcard"
	}
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "ls", "-la", remoteDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ls failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return ParseFilesOutput(stdout.String(), remoteDir), nil
}

// PullFile downloads a file or directory from the device to local storage
func (c *Client) PullFile(ctx context.Context, serial, remotePath, localPath string) error {
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "pull", remotePath, localPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pull failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// PushFile uploads a local file to the device filesystem
func (c *Client) PushFile(ctx context.Context, serial, localPath, remotePath string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "push", localPath, remotePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("push failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// RemoveFile deletes a file or directory on the device
func (c *Client) RemoveFile(ctx context.Context, serial, remotePath string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "rm", "-rf", remotePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rm failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// GetLogcat dumps recent logcat lines from the device
func (c *Client) GetLogcat(ctx context.Context, serial string, linesCount int) ([]LogcatLine, error) {
	if linesCount <= 0 {
		linesCount = 100
	}
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "logcat", "-d", "-v", "threadtime", "-t", fmt.Sprintf("%d", linesCount))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("logcat failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return ParseLogcatLines(stdout.String()), nil
}

// ClearLogcat flushes the device log buffers
func (c *Client) ClearLogcat(ctx context.Context, serial string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "logcat", "-c")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clear logcat failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// ListPackages queries installed application packages on the device
func (c *Client) ListPackages(ctx context.Context, serial string, thirdPartyOnly bool) ([]string, error) {
	args := []string{"-s", serial, "shell", "pm", "list", "packages"}
	if thirdPartyOnly {
		args = append(args, "-3")
	}
	cmd := exec.CommandContext(ctx, c.adbPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("pm list packages failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}

	return ParsePackagesOutput(stdout.String()), nil
}

// LaunchApp launches an application on the device by package name
func (c *Client) LaunchApp(ctx context.Context, serial, pkgName string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "monkey", "-p", pkgName, "-c", "android.intent.category.LAUNCHER", "1")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("launch app failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// ForceStopApp terminates all processes associated with a package
func (c *Client) ForceStopApp(ctx context.Context, serial, pkgName string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "am", "force-stop", pkgName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("force-stop failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// ClearAppData resets application data and caches
func (c *Client) ClearAppData(ctx context.Context, serial, pkgName string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "pm", "clear", pkgName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clear data failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// SendKeyEvent injects a hardware key event into the device
func (c *Client) SendKeyEvent(ctx context.Context, serial string, keyCode int) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "input", "keyevent", fmt.Sprintf("%d", keyCode))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("input keyevent failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// SendInputText injects text characters into the currently focused input field
func (c *Client) SendInputText(ctx context.Context, serial, text string) error {
	// ADB input text requires spaces to be %s and special chars escaped
	escaped := strings.ReplaceAll(text, " ", "%s")
	escaped = strings.ReplaceAll(escaped, "&", "\\&")
	escaped = strings.ReplaceAll(escaped, "<", "\\<")
	escaped = strings.ReplaceAll(escaped, ">", "\\>")
	escaped = strings.ReplaceAll(escaped, ";", "\\;")
	escaped = strings.ReplaceAll(escaped, "\"", "\\\"")

	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "input", "text", escaped)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("input text failed: %v (%s)", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}

// GetMemoryInfo queries physical RAM consumption
func (c *Client) GetMemoryInfo(ctx context.Context, serial string) (usedMB, totalMB int, percent float64, err error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "cat", "/proc/meminfo")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return 0, 0, 0, err
	}

	usedMB, totalMB, percent = ParseMemInfo(stdout.String())
	return usedMB, totalMB, percent, nil
}

// GetCPUUsage queries total CPU utilization percentage
func (c *Client) GetCPUUsage(ctx context.Context, serial string) (float64, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "dumpsys", "cpuinfo")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return 0.0, err
	}

	return ParseCPUInfo(stdout.String()), nil
}

// GetUptime queries device operating system uptime
func (c *Client) GetUptime(ctx context.Context, serial string) (string, error) {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "cat", "/proc/uptime")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return ParseUptime(stdout.String()), nil
}

// StartScreenRecord starts a background video recording on device storage
func (c *Client) StartScreenRecord(ctx context.Context, serial, remotePath string) error {
	if remotePath == "" {
		remotePath = "/sdcard/andtls_record.mp4"
	}
	// Start with 3 minute limit
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "screenrecord", "--time-limit", "180", remotePath)
	return cmd.Start()
}

// StopScreenRecord cleanly terminates the active screen recording process
func (c *Client) StopScreenRecord(ctx context.Context, serial string) error {
	cmd := exec.CommandContext(ctx, c.adbPath, "-s", serial, "shell", "pkill -2 -f screenrecord || killall -2 screenrecord || true")
	return cmd.Run()
}

// ExportReport generates a complete diagnostic report for the specified device and saves it to destDir
func (c *Client) ExportReport(ctx context.Context, serial string, dev device.Device, destDir string) (string, error) {
	if destDir == "" {
		destDir = "reports"
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}

	fileName := fmt.Sprintf("report_%s_%s.md", serial, time.Now().Format("20060102_150405"))
	destPath := filepath.Join(destDir, fileName)

	// Fetch properties
	propsRaw, _ := c.Shell(ctx, serial, "getprop")
	dfRaw, _ := c.Shell(ctx, serial, "df -h")
	packages, _ := c.ListPackages(ctx, serial, true)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# 📱 Android Device Diagnostics Report: %s\n\n", dev.DisplayName()))
	sb.WriteString(fmt.Sprintf("- **Generated At**: %s\n", time.Now().Format(time.RFC1123)))
	sb.WriteString(fmt.Sprintf("- **Serial Number**: `%s`\n", dev.Serial))
	sb.WriteString(fmt.Sprintf("- **Status**: %s\n", dev.State))
	sb.WriteString(fmt.Sprintf("- **Model**: %s\n", dev.Model))
	sb.WriteString(fmt.Sprintf("- **Product / Codename**: %s\n", dev.Product))
	sb.WriteString(fmt.Sprintf("- **Android Version**: %s\n", dev.OSVersionString()))
	sb.WriteString(fmt.Sprintf("- **Screen Resolution**: %s\n", dev.ScreenResolution))
	sb.WriteString(fmt.Sprintf("- **IP Address**: %s\n", dev.IPAddress))
	sb.WriteString(fmt.Sprintf("- **Battery**: %s\n", dev.BatteryDetailString()))
	sb.WriteString(fmt.Sprintf("- **Storage**: %s\n", dev.StorageString()))
	sb.WriteString(fmt.Sprintf("- **RAM Utilization**: %s\n", dev.RAMString()))
	sb.WriteString(fmt.Sprintf("- **System Uptime**: %s\n", dev.Uptime))
	sb.WriteString(fmt.Sprintf("- **User Installed Apps**: %d packages\n\n", len(packages)))

	sb.WriteString("## 💾 Filesystem Partitions (`df -h`)\n\n```\n")
	sb.WriteString(dfRaw)
	sb.WriteString("\n```\n\n")

	sb.WriteString("## 📦 Installed Third-Party Packages\n\n")
	for _, pkg := range packages {
		sb.WriteString(fmt.Sprintf("- `%s`\n", pkg))
	}
	sb.WriteString("\n")

	sb.WriteString("## ⚙️ System Properties (`getprop`)\n\n```ini\n")
	sb.WriteString(propsRaw)
	sb.WriteString("\n```\n")

	if err := os.WriteFile(destPath, []byte(sb.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to save report file: %w", err)
	}

	return destPath, nil
}
