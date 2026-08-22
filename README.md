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

## 📄 License

MIT
