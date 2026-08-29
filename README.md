# 📱 andtls (Android Tools TUI Dashboard)

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Bubble Tea](https://img.shields.io/badge/UI-Bubble%20Tea-FA5252?style=flat)](https://github.com/charmbracelet/bubbletea)
[![Lip Gloss](https://img.shields.io/badge/Styling-Lip%20Gloss-7950F2?style=flat)](https://github.com/charmbracelet/lipgloss)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-blue)](https://github.com/Bloby22/andtls)
[![Locales](https://img.shields.io/badge/Locales-English%20🇬🇧%20%7C%20Čeština%20🇨🇿-orange)](#-bilingual-support)

**andtls** is a high-performance, keyboard-driven terminal dashboard (TUI) for real-time monitoring, debugging, remote control, and device management of Android devices (physical USB devices and emulators).

```
 ╭──────────────────────────────────────────────────────────────────────────────────────────╮
 │  📱 ANDTLS — Android Device Manager v0.2.0                                               │
 ╰──────────────────────────────────────────────────────────────────────────────────────────╯

  STATUS        SERIAL NUMBER             MODEL / NAME         BATTERY      TR-ID
 ▸ ● Connected   emulator-5554             Pixel 8 Pro          🔋 92% (USB)  1
   ● Connected   192.168.1.120:5555        Samsung Galaxy S24   🔋 78% (WiFi) 2

 ╭──────────────────────────── 🔍 Device Inspector ─────────────────────────────────────────╮
 │  Serial Number:  emulator-5554           ADB Status:     ● Connected (Authorized)       │
 │  Model Name:     Pixel 8 Pro             Codename:       husky                          │
 │  Android OS:     Android 14 (API 34)     IP Address:     192.168.1.120                  │
 │  Battery:        92% (Charging, 29.4°C)  Storage Space:  65G free / 110G total (41%)    │
 │  RAM Memory:     [████░░░░] 4.0 / 8.0 GB CPU Load:       [██░░░░░░] 14.2%               │
 │  USB / Port:     USB / TR-1              Resolution:     1344x2992      Uptime: 2d 4h   │
 ╰──────────────────────────────────────────────────────────────────────────────────────────╯

 ✔ Monitoring 2 connected device(s)... │ 🎨 Cyberpunk │ 🌐 CS │ Updated: 21:38:00
 [↑/↓] Navigate  [Enter] Actions  [m] Mirror  [f] Files  [l] Logcat  [p] Apps  [s] Shot  [L] Lang  [t] Theme  [?] Help
```

---

## 🚀 Key Subsystems & Features

### 🖥️ Display & Control
- 🪞 **Instant Screen Mirroring (`m` / `scrcpy`)**: Launch high-fps display mirroring and keyboard/mouse passthrough using `scrcpy` in the background with a single keypress.
- 📹 **Screen Video Recording (`v`)**: Start/stop on-device video recording (`screenrecord`) with a live recording timer indicator; output `.mp4` is automatically pulled to `./recordings/`.
- 📸 **Display Screenshot (`s`)**: Grab high-resolution display framebuffers directly to `./screenshots/` with automatic timestamped filenames.
- 🎮 **Virtual Remote Control & Keypad (`o`)**: Send hardware key events (Power, Volume +/-, Home, Back, Recent Apps, Wake/Unlock) and inject text strings into active input fields.

### 🗂️ File & Application Management
- 📁 **Interactive File Manager (`f`)**: Browse the device filesystem (e.g. `/sdcard/`), navigate subdirectories, download/pull files (`d`) to `./downloads/`, and delete remote files (`x`).
- 📦 **Package & Application Manager (`p`)**: List all installed packages, toggle between user (3rd party) and system apps (`Tab`), launch applications (`Enter`), force-stop apps (`x`), clear app data (`c`), and uninstall packages (`u`).

### 🩺 Real-Time Telemetry & Monitoring
- 📋 **Live Device Table**: Instant background detection of connected/disconnected USB and TCP/IP devices with status badges (Connected, Unauthorized, Offline, Recovery).
- 📊 **Hardware Telemetry Gauges**: Real-time progress bars for **RAM utilization** (used/total MB and %) and **CPU load**, alongside system **Uptime**, storage metrics, and battery health (temperature, voltage, charging state).
- 📜 **Live Logcat Stream Viewer (`l`)**: Stream device logs in real time with log-level filtering (`1-5`: All, Debug, Info, Warn, Error), pause/resume (`Space`), and buffer flush (`c`).
- 📑 **Diagnostics Report Export (`e`)**: Generate detailed Markdown hardware and software diagnostic reports saved to `./reports/report_<serial>_<timestamp>.md`.

### 🎨 Customization & Localization
- 🎨 **Dynamic Color Themes (`t`)**: Instant runtime theme switching:
  - **Cyberpunk** (Neon Pink, Cyan, Yellow)
  - **Catppuccin** (Mauve, Peach, Lavender)
  - **Dracula** (Purple, Green, Pink)
  - **Nord** (Frost Blue, Polar Night)
  - **Matrix** (Green phosphor monochrome)
- 🌐 **Bilingual Support (`L`)**: Switch the entire UI between **English** 🇬🇧 and **Czech** 🇨🇿 on the fly (also auto-detects system locale).

---

## 🎮 Keybindings & Shortcuts Reference

### 🌐 Global / Main Dashboard Mode
| Key | Action |
| :--- | :--- |
| `↑` / `k` | Select previous device |
| `↓` / `j` | Select next device |
| `Enter` / `a` | Open Quick Actions menu for selected device |
| `m` | **Screen Mirroring** (launch `scrcpy`) |
| `v` | **Record Screen Video** (start / stop .mp4 recording) |
| `f` | **Device File Manager** (browse `/sdcard/`, download, delete) |
| `l` | **Live Logcat Stream** (real-time logcat viewer) |
| `p` | **Package & App Manager** (launch, stop, clear data, uninstall) |
| `o` | **Virtual Remote Control** (hardware buttons & text input) |
| `s` | **Take Screenshot** (saved to `./screenshots/`) |
| `w` | **Enable Wireless ADB** (switch daemon to TCP port 5555) |
| `e` | **Export Diagnostics Report** (saved to `./reports/`) |
| `b` | Open Reboot menu (System / Recovery / Bootloader) |
| `r` | Instant manual device refresh |
| `t` | **Cycle Color Themes** (Cyberpunk, Catppuccin, Dracula, Nord, Matrix) |
| `L` | **Toggle Language** (Čeština 🇨🇿 / English 🇬🇧) |
| `?` / `h` | Toggle keyboard shortcuts help modal |
| `q` / `Ctrl+C` | Quit application |

### 📁 File Manager Mode (`f`)
| Key | Action |
| :--- | :--- |
| `↑` / `k`, `↓` / `j` | Navigate file list |
| `Enter` | Open directory / Download file |
| `Backspace` / `b` / `←` | Navigate to parent directory (`..`) |
| `d` | Download / Pull selected file to `./downloads/` |
| `x` | Delete selected file on device |
| `Esc` / `q` | Return to main dashboard |

### 📦 Package Manager Mode (`p`)
| Key | Action |
| :--- | :--- |
| `↑` / `k`, `↓` / `j` | Navigate package list |
| `Tab` | Toggle between User Apps (3rd-party) and All Packages (System) |
| `Enter` | Launch application |
| `x` | Force Stop application process (`am force-stop`) |
| `c` | Clear application data & cache (`pm clear`) |
| `u` | Uninstall application package (`pm uninstall`) |
| `Esc` / `q` | Return to main dashboard |

### 📜 Live Logcat Mode (`l`)
| Key | Action |
| :--- | :--- |
| `Space` | Pause / Resume live streaming |
| `c` | Clear logcat buffer |
| `1` | Show all log levels |
| `2` | Filter **Debug** (`D`) and higher |
| `3` | Filter **Info** (`I`) and higher |
| `4` | Filter **Warn** (`W`) and higher |
| `5` | Filter **Error** (`E`) and higher |
| `Esc` / `q` | Return to main dashboard |

### 🎮 Virtual Remote Mode (`o`)
| Key | Action | Keycode |
| :--- | :--- | :--- |
| `p` | Power / Screen Lock | `KEYCODE_POWER (26)` |
| `+` / `=` | Volume Up | `KEYCODE_VOLUME_UP (24)` |
| `-` / `_` | Volume Down | `KEYCODE_VOLUME_DOWN (25)` |
| `h` | Home Button | `KEYCODE_HOME (3)` |
| `b` | Back Button | `KEYCODE_BACK (4)` |
| `r` | Recent Apps / Overview | `KEYCODE_APP_SWITCH (187)` |
| `w` | Wake & Unlock Screen | `KEYCODE_WAKEUP (224)` + `KEYCODE_MENU (82)` |
| `t` | Inject Text String | Dispatches `input text` |
| `Esc` / `q` | Return to main dashboard | — |

---

## 🛠️ Prerequisites

1. **Go** `1.25+` (for compiling from source)
2. **Android SDK Platform-Tools (`adb`)**: must be available in your system `PATH`.
3. **scrcpy** *(optional, recommended)*: for display screen mirroring (`m`).

---

## 📦 Building & Installation

### Build Standalone Binary

```bash
# Clone repository
git clone https://github.com/Bloby22/andtls.git
cd andtls

# Build executable
go build -o andtls.exe .

# Run
./andtls.exe
```

### Run Directly via Go
```bash
go run ./cmd/andtls
```

### Run Tests
```bash
go test ./...
```

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.