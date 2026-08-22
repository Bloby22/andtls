# 📱 andtls (Android Tools)

A modern, interactive terminal dashboard (TUI) for real-time monitoring and management of USB and emulator Android devices.

Built with **Go** using [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

---

## 🚀 Features

- ⚡ **Live Device Monitoring**: Real-time background detection of connected/disconnected Android devices (1-second polling) without restarting.
- 📋 **Device Table**: Comprehensive overview of status (Connected, Unauthorized, Offline), serial numbers, model names, product codes, battery levels, and transport IDs.
- 🔍 **Deep Device Inspector**: Detailed hardware, OS, network, battery telemetry (percentage, status, temperature, voltage), and internal storage space.
- 📸 **Quick Screenshot (`s`)**: Capture device display directly and save high-resolution PNG to `./screenshots/`.
- 🌐 **Wireless ADB Toggle (`w`)**: Switch target device to TCP port 5555 for wireless debugging over WiFi.
- 🔄 **Reboot Menu (`b`)**: Fast reboot to normal system, Android Recovery, or Fastboot/Bootloader mode.
- 📖 **Interactive Help (`?`)**: Full modal overlay showing keybindings and shortcuts.
- 🛡️ **Empty State & Troubleshooting**: Step-by-step guidance when no devices are connected or USB debugging is not enabled.

---

## 🛠️ Requirements

- **Go**: version `1.25.0` or higher
- **ADB (Android Debug Bridge)**: part of Android SDK Platform-Tools available in `PATH`
- **Perl**: version `5.20+` (optional, for automation scripts in `scripts/`)

---

## 📦 Installation & Running

### Using build and install scripts

```bash
# Build binary
./build.sh

# Install to system PATH (~/.local/bin or /usr/local/bin)
./install.sh

# Run
andtls
```

### Using Makefile

```bash
# Build binary
make build

# Run directly
make run

# Install to PATH
make install
```

---

## 📜 Automation Scripts (`scripts/`)

The `scripts/` directory contains helper tools in Shell (`.sh`) and Perl (`.pl`):

| Script | Language | Description |
| :--- | :--- | :--- |
| `scripts/dev.sh` | Shell | Interactive dev mode with formatting, vetting, and live launch |
| `scripts/release.sh` | Shell | Cross-compilation packaging for Linux, macOS, and Windows with SHA256 checksums |
| `scripts/check_devices.sh` | Shell | Quick CLI diagnostic for ADB binary, daemon port, and attached devices |
| `scripts/log_parser.pl` | Perl | Colorized logcat stream parser with tag regex and minimum priority filtering |
| `scripts/udev_generator.pl` | Perl | Detects connected Android USB vendor IDs and generates Linux udev rules |
| `scripts/metrics_collector.pl` | Perl | Exports device memory, storage, battery, and CPU metrics to Table/JSON format |

---

## 🎮 Keybindings

| Key | Action |
| :--- | :--- |
| `↑` / `k` | Select previous device |
| `↓` / `j` | Select next device |
| `Enter` / `a` | Open Quick Actions menu for selected device |
| `s` | Take screenshot of selected device (saved to `./screenshots/`) |
| `w` | Enable wireless ADB mode on port 5555 |
| `b` | Open reboot options menu |
| `r` | Instant manual refresh |
| `?` / `h` | Toggle keyboard help modal |
| `Esc` | Close active modal overlay |
| `q` / `Ctrl+C` | Quit application |

---

## 📁 Project Structure

```text
andtls/
├── Makefile                   # Build automation (build, run, install, fmt, vet, clean)
├── build.sh                   # Colored build script with ldflags and checks
├── install.sh                 # Installation script to PATH
├── cmd/
│   └── andtls/
│       └── main.go            # CLI binary entrypoint
├── internal/
│   ├── adb/
│   │   ├── client.go          # ADB command execution (ListDevices, Screenshot, Reboot...)
│   │   └── parser.go          # Parsing engine for adb devices, battery, and storage
│   ├── app/
│   │   └── app.go             # Application runner & lifecycle orchestrator
│   ├── config/
│   │   └── config.go          # Runtime configuration (timeouts, intervals, paths)
│   ├── device/
│   │   ├── device.go          # Device model, states, and telemetry formatters
│   │   └── watcher.go         # Polling / change detection / diff logic
│   └── ui/
│       ├── components/
│       │   ├── actions.go     # Quick action menu modal component
│       │   ├── details.go     # Deep device inspector card component
│       │   ├── empty.go       # Empty state troubleshooting component
│       │   ├── header.go      # Title & header banner component
│       │   ├── help.go        # Keyboard shortcuts help modal component
│       │   ├── statusbar.go   # Status bar & help footer component
│       │   └── table.go       # Interactive device table component
│       ├── styles/
│       │   ├── colors.go      # Color palette definitions (Catppuccin/Dark)
│       │   └── theme.go       # Lip Gloss theme styles & badge helpers
│       ├── model.go           # Bubble Tea application state model
│       ├── msg.go             # UI Messages & async commands
│       ├── update.go          # Event handling & keybindings loop
│       └── view.go            # Main view composition
├── scripts/
│   ├── check_devices.sh       # Fast environment and connection diagnostic tool
│   ├── dev.sh                 # Local development launcher
│   ├── log_parser.pl          # Perl colorized logcat parser and filter
│   ├── metrics_collector.pl   # Perl telemetry exporter (Table & JSON)
│   ├── release.sh             # Multi-platform release packaging
│   └── udev_generator.pl      # Linux Android udev rules generator
├── main.go                    # Root entrypoint alias
├── go.mod
├── go.sum
└── README.md
```

---

## 📄 License

MIT
