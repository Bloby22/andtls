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
