package locale

import (
	"os"
	"strings"
)

// Strings holds all translatable UI text for the application
type Strings struct {
	// Header
	AppTitle    string
	AppSubtitle string

	// Help modal
	HelpTitle              string
	HelpMoveUp             string
	HelpMoveDown           string
	HelpRefresh            string
	HelpScreenshot         string
	HelpWirelessADB        string
	HelpRebootMenu         string
	HelpActionsMenu        string
	HelpMirror             string
	HelpRecording          string
	HelpFileManager        string
	HelpLogcat             string
	HelpPackages           string
	HelpRemote             string
	HelpExportReport       string
	HelpTheme              string
	HelpLanguage           string
	HelpToggleHelp         string
	HelpCloseModal         string
	HelpQuit               string
	HelpCloseHint          string

	// Empty state
	EmptyTitle string
	EmptyTip1  string
	EmptyTip2  string
	EmptyTip3  string
	EmptyTip4  string

	// Status bar
	StatusNavigate        string
	StatusActions         string
	StatusScreenshot      string
	StatusMirror          string
	StatusFiles           string
	StatusLogcat          string
	StatusPackages        string
	StatusRemote          string
	StatusRefresh         string
	StatusHelp            string
	StatusQuit            string
	StatusExit            string
	StatusUpdated         string
	StatusRecordingActive string

	// Detail card
	DetailTitle         string
	DetailSerialNumber  string
	DetailADBStatus     string
	DetailModelName     string
	DetailCodename      string
	DetailAndroidOS     string
	DetailIPAddress     string
	DetailBattery       string
	DetailStorage       string
	DetailRAM           string
	DetailCPU           string
	DetailUptime        string
	DetailUSBPort       string
	DetailResolution    string
	DetailInstalledApps string

	// Table headers
	TableStatus    string
	TableSerial    string
	TableModel     string
	TableBattery   string
	TableTransport string

	// Actions
	ActionsTitle            string
	ActionsMirrorTitle      string
	ActionsMirrorDesc       string
	ActionsRecordTitle      string
	ActionsRecordDesc       string
	ActionsFilesTitle       string
	ActionsFilesDesc        string
	ActionsLogcatTitle      string
	ActionsLogcatDesc       string
	ActionsPackagesTitle    string
	ActionsPackagesDesc     string
	ActionsRemoteTitle      string
	ActionsRemoteDesc       string
	ActionsScreenshotTitle  string
	ActionsScreenshotDesc   string
	ActionsWirelessTitle    string
	ActionsWirelessDesc     string
	ActionsExportTitle      string
	ActionsExportDesc       string
	ActionsRebootSystem     string
	ActionsRebootSystemDesc string
	ActionsRebootRecovery   string
	ActionsRebootRecDesc    string
	ActionsRebootBoot       string
	ActionsRebootBootDesc   string
	ActionsFooter           string

	// Logcat View
	LogcatTitle   string
	LogcatFilter  string
	LogcatSearch  string
	LogcatPaused  string
	LogcatRunning string
	LogcatHint    string

	// File Manager View
	FilesTitle      string
	FilesHeaderName string
	FilesHeaderSize string
	FilesHeaderPerm string
	FilesHeaderMod  string
	FilesHint       string
	FilesEmpty      string

	// Package Manager View
	PackagesTitle  string
	PackagesFilter string
	PackagesToggle string
	PackagesHint   string

	// Remote Control View
	RemoteTitle     string
	RemotePower     string
	RemoteVolUp     string
	RemoteVolDown   string
	RemoteHome      string
	RemoteBack      string
	RemoteRecent    string
	RemoteWake      string
	RemoteTextInput string
	RemoteFooter    string

	// Device states
	StateConnected    string
	StateUnauthorized string
	StateOffline      string
	StateAuthorizing  string
	StateRecovery     string
	StateSideload     string
	StateBootloader   string

	// Battery statuses
	BatteryCharging    string
	BatteryDischarging string
	BatteryNotCharging string
	BatteryFull        string
	BatteryUnknown     string

	// Status messages
	MsgRefreshing          string
	MsgMonitoringNone      string
	MsgMonitoringFound     string
	MsgUnauthorizedAction  string
	MsgCannotScreenshot    string
	MsgCannotWireless      string
	MsgCapturingScreenshot string
	MsgEnablingWireless    string
	MsgRebooting           string
	MsgRebootingRecovery   string
	MsgRebootingBootloader string
	MsgScreenshotSaved     string
	MsgWirelessEnabled     string
	MsgRebootToSystem      string
	MsgRebootToRecovery    string
	MsgRebootToBootloader  string
	MsgADBError            string
	MsgScreenshotFailed    string
	MsgWirelessFailed      string
	MsgRebootFailed        string
	MsgMirrorStarted       string
	MsgMirrorFailed        string
	MsgRecordStarted       string
	MsgRecordSaved         string
	MsgRecordFailed        string
	MsgFileDownloaded      string
	MsgFileDeleted         string
	MsgFileOpFailed        string
	MsgAppLaunched         string
	MsgAppStopped          string
	MsgAppCleared          string
	MsgAppUninstalled      string
	MsgKeySent             string
	MsgReportExported      string
	MsgLanguageChanged     string
	MsgThemeChanged        string
}

// en is the default English locale
var en = &Strings{
	// Header
	AppTitle:    "📱 ANDTLS",
	AppSubtitle: "Android Tools Dashboard",

	// Help modal
	HelpTitle:              "📖 andtls Keyboard Shortcuts & Help",
	HelpMoveUp:             "Move selection cursor up",
	HelpMoveDown:           "Move selection cursor down",
	HelpRefresh:            "Trigger immediate device poll and refresh",
	HelpScreenshot:         "Take screenshot (saved to ./screenshots/)",
	HelpWirelessADB:        "Enable wireless ADB mode on port 5555",
	HelpRebootMenu:         "Open reboot options menu (System, Recovery, Bootloader)",
	HelpActionsMenu:        "Open quick actions menu for selected device",
	HelpMirror:             "Mirror device screen via scrcpy (m)",
	HelpRecording:          "Toggle screen video recording (v)",
	HelpFileManager:        "Open device file manager (f)",
	HelpLogcat:             "Open live Logcat viewer (l)",
	HelpPackages:           "Open package & app manager (p)",
	HelpRemote:             "Open virtual remote control (o)",
	HelpExportReport:       "Export full diagnostics report (e)",
	HelpTheme:              "Cycle color themes (t)",
	HelpLanguage:           "Toggle UI language CS / EN (L)",
	HelpToggleHelp:         "Toggle this help dialog overlay",
	HelpCloseModal:         "Close any active modal overlay",
	HelpQuit:               "Quit application",
	HelpCloseHint:          "Press [Esc] or [?] to close this help window",

	// Empty state
	EmptyTitle: "⚡ No Android Devices Detected",
	EmptyTip1:  "1. Connect your Android phone or tablet via USB cable.",
	EmptyTip2:  "2. In Android Settings -> Developer Options, enable 'USB Debugging'.",
	EmptyTip3:  "3. On your device screen, accept the 'Allow USB debugging?' prompt.",
	EmptyTip4:  "4. Press [r] to immediately refresh the connection list.",

	// Status bar
	StatusNavigate:        "Navigate",
	StatusActions:         "Actions",
	StatusScreenshot:      "Screenshot",
	StatusMirror:          "Mirror",
	StatusFiles:           "Files",
	StatusLogcat:          "Logcat",
	StatusPackages:        "Apps",
	StatusRemote:          "Remote",
	StatusRefresh:         "Refresh",
	StatusHelp:            "Help",
	StatusQuit:            "Quit",
	StatusExit:            "Exit",
	StatusUpdated:         "Updated: %s",
	StatusRecordingActive: "● REC [%s]",

	// Detail card
	DetailTitle:         "🔍 Device Inspector",
	DetailSerialNumber:  "Serial Number:",
	DetailADBStatus:     "ADB Status:",
	DetailModelName:     "Model Name:",
	DetailCodename:      "Codename / Product:",
	DetailAndroidOS:     "Android OS:",
	DetailIPAddress:     "IP Address:",
	DetailBattery:       "Battery Telemetry:",
	DetailStorage:       "Storage Space:",
	DetailRAM:           "RAM Memory:",
	DetailCPU:           "CPU Utilization:",
	DetailUptime:        "System Uptime:",
	DetailUSBPort:       "USB Bus / Port:",
	DetailResolution:    "Display Resolution:",
	DetailInstalledApps: "Installed Apps:",

	// Table headers
	TableStatus:    "STATUS",
	TableSerial:    "SERIAL NUMBER",
	TableModel:     "MODEL / NAME",
	TableBattery:   "BATTERY",
	TableTransport: "TR-ID",

	// Actions
	ActionsTitle:            "⚡ Quick Actions for: %s (%s)",
	ActionsMirrorTitle:      "Screen Mirroring (scrcpy)",
	ActionsMirrorDesc:       "Mirror and control screen in window",
	ActionsRecordTitle:      "Record Screen Video",
	ActionsRecordDesc:       "Capture video to ./recordings/",
	ActionsFilesTitle:       "Device File Manager",
	ActionsFilesDesc:        "Browse, download, and delete files",
	ActionsLogcatTitle:      "Live Logcat Viewer",
	ActionsLogcatDesc:       "Stream real-time device logs",
	ActionsPackagesTitle:    "Package & App Manager",
	ActionsPackagesDesc:     "Launch, stop, clear data, uninstall",
	ActionsRemoteTitle:      "Remote Control & Keypad",
	ActionsRemoteDesc:       "Power, Volume, Home, Back, Input",
	ActionsScreenshotTitle:  "Take Screenshot",
	ActionsScreenshotDesc:   "Capture display to ./screenshots/",
	ActionsWirelessTitle:    "Enable Wireless ADB",
	ActionsWirelessDesc:     "Switch daemon to TCP port 5555",
	ActionsExportTitle:      "Export Diagnostics Report",
	ActionsExportDesc:       "Save full Markdown spec to ./reports/",
	ActionsRebootSystem:     "Reboot System",
	ActionsRebootSystemDesc: "Normal device restart",
	ActionsRebootRecovery:   "Reboot to Recovery",
	ActionsRebootRecDesc:    "Boot into Android Recovery",
	ActionsRebootBoot:       "Reboot to Bootloader",
	ActionsRebootBootDesc:   "Boot into Fastboot mode",
	ActionsFooter:           "Press [Enter] to run selected, [Key] for shortcut, or [Esc] to cancel",

	// Logcat View
	LogcatTitle:   "📋 Live Logcat Stream — %s",
	LogcatFilter:  "Level: [%s]",
	LogcatSearch:  "Search: %s",
	LogcatPaused:  "[PAUSED]",
	LogcatRunning: "[STREAMING]",
	LogcatHint:    "[Space] Pause/Resume • [c] Clear • [1-5] Level • [/] Search • [Esc] Exit",

	// File Manager View
	FilesTitle:      "📁 Device File Manager — %s (%s)",
	FilesHeaderName: "NAME",
	FilesHeaderSize: "SIZE",
	FilesHeaderPerm: "PERMISSIONS",
	FilesHeaderMod:  "MODIFIED",
	FilesHint:       "[Enter] Open • [Backspace/b] Parent dir • [d] Download/Pull • [x] Delete • [Esc] Exit",
	FilesEmpty:      "Directory is empty or permission denied",

	// Package Manager View
	PackagesTitle:  "📦 Installed Applications (%d apps) — %s",
	PackagesFilter: "Filter: %s",
	PackagesToggle: "[Tab] Toggle 3rd Party / All",
	PackagesHint:   "[Enter] Launch • [x] Force Stop • [c] Clear Data • [u] Uninstall • [Esc] Exit",

	// Remote Control View
	RemoteTitle:     "🎮 Virtual Remote Control & Input — %s",
	RemotePower:     "Power Button",
	RemoteVolUp:     "Volume Up",
	RemoteVolDown:   "Volume Down",
	RemoteHome:      "Home Screen",
	RemoteBack:      "Back Button",
	RemoteRecent:    "Recent Apps",
	RemoteWake:      "Wake / Unlock Display",
	RemoteTextInput: "Send Text Input",
	RemoteFooter:    "Press [Key] to trigger hardware button, or [Esc] to return",

	// Device states
	StateConnected:    "Connected",
	StateUnauthorized: "Unauthorized",
	StateOffline:      "Offline",
	StateAuthorizing:  "Authorizing...",
	StateRecovery:     "Recovery",
	StateSideload:     "Sideload",
	StateBootloader:   "Bootloader",

	// Battery statuses
	BatteryCharging:    "Charging",
	BatteryDischarging: "Discharging",
	BatteryNotCharging: "Not charging",
	BatteryFull:        "Full",
	BatteryUnknown:     "Unknown",

	// Status messages
	MsgRefreshing:          "Refreshing devices...",
	MsgMonitoringNone:      "Monitoring devices (1s polling)... No devices connected",
	MsgMonitoringFound:     "Monitoring devices (1s polling)... Found %d device(s)",
	MsgUnauthorizedAction:  "Device must be authorized to perform actions",
	MsgCannotScreenshot:    "Cannot capture screenshot: device not authorized",
	MsgCannotWireless:      "Cannot enable wireless ADB: device not authorized",
	MsgCapturingScreenshot: "Capturing screenshot for %s...",
	MsgEnablingWireless:    "Enabling wireless ADB on port 5555 for %s...",
	MsgRebooting:           "Rebooting %s...",
	MsgRebootingRecovery:   "Rebooting %s to Recovery...",
	MsgRebootingBootloader: "Rebooting %s to Bootloader...",
	MsgScreenshotSaved:     "Screenshot saved to: %s",
	MsgWirelessEnabled:     "Wireless ADB enabled on port 5555 (Connect with: adb connect %s:5555)",
	MsgRebootToSystem:      "Device %s rebooting to system...",
	MsgRebootToRecovery:    "Device %s rebooting to recovery...",
	MsgRebootToBootloader:  "Device %s rebooting to bootloader...",
	MsgADBError:            "ADB Error: %v",
	MsgScreenshotFailed:    "Screenshot failed: %v",
	MsgWirelessFailed:      "Wireless ADB failed: %v",
	MsgRebootFailed:        "Reboot failed: %v",
	MsgMirrorStarted:       "Screen mirroring active in external window (scrcpy)",
	MsgMirrorFailed:        "Screen mirroring failed: %v",
	MsgRecordStarted:       "Screen recording started! Press [v] again to stop and save video.",
	MsgRecordSaved:         "Screen recording saved to: %s",
	MsgRecordFailed:        "Screen recording failed: %v",
	MsgFileDownloaded:      "Downloaded file to: %s",
	MsgFileDeleted:         "Deleted on device: %s",
	MsgFileOpFailed:        "File operation failed: %v",
	MsgAppLaunched:         "Launched application: %s",
	MsgAppStopped:          "Force-stopped application: %s",
	MsgAppCleared:          "Cleared application data for: %s",
	MsgAppUninstalled:      "Uninstalled application: %s",
	MsgKeySent:             "Key command dispatched: %s",
	MsgReportExported:      "Diagnostics report saved to: %s",
	MsgLanguageChanged:     "Language switched to English",
	MsgThemeChanged:        "Theme switched to: %s",
}

// cs is the Czech locale
var cs = &Strings{
	// Header
	AppTitle:    "📱 ANDTLS",
	AppSubtitle: "Android nástroje — přehled",

	// Help modal
	HelpTitle:              "📖 andtls Klávesové zkratky a nápověda",
	HelpMoveUp:             "Posunout kurzor výběru nahoru",
	HelpMoveDown:           "Posunout kurzor výběru dolů",
	HelpRefresh:            "Okamžitě obnovit seznam zařízení",
	HelpScreenshot:         "Pořídit snímek obrazovky (uloží do ./screenshots/)",
	HelpWirelessADB:        "Zapnout bezdrátový ADB mód na portu 5555",
	HelpRebootMenu:         "Otevřít nabídku restartu (Systém, Recovery, Bootloader)",
	HelpActionsMenu:        "Otevřít rychlé akce pro vybrané zařízení",
	HelpMirror:             "Zrcadlit obrazovku přes scrcpy (m)",
	HelpRecording:          "Spustit/Zastavit nahrávání obrazovky (v)",
	HelpFileManager:        "Otevřít správce souborů zařízení (f)",
	HelpLogcat:             "Otevřít živý prohlížeč Logcat (l)",
	HelpPackages:           "Otevřít správce aplikací a balíčků (p)",
	HelpRemote:             "Otevřít virtuální dálkový ovladač (o)",
	HelpExportReport:       "Exportovat diagnostický report (e)",
	HelpTheme:              "Přepnout barevné téma (t)",
	HelpLanguage:           "Přepnout jazyk rozhraní CS / EN (L)",
	HelpToggleHelp:         "Přepnout zobrazení této nápovědy",
	HelpCloseModal:         "Zavřít aktivní překryvné okno",
	HelpQuit:               "Ukončit aplikaci",
	HelpCloseHint:          "Stiskněte [Esc] nebo [?] pro zavření této nápovědy",

	// Empty state
	EmptyTitle: "⚡ Žádná zařízení Android nebyla nalezena",
	EmptyTip1:  "1. Připojte telefon nebo tablet Android přes USB kabel.",
	EmptyTip2:  "2. V Nastavení Androidu -> Možnosti vývojáře zapněte 'Ladění přes USB'.",
	EmptyTip3:  "3. Na obrazovce zařízení přijměte výzvu 'Povolit ladění přes USB?'.",
	EmptyTip4:  "4. Stiskněte [r] pro okamžité obnovení seznamu připojení.",

	// Status bar
	StatusNavigate:        "Pohyb",
	StatusActions:         "Akce",
	StatusScreenshot:      "Snímek",
	StatusMirror:          "Zrcadlit",
	StatusFiles:           "Soubory",
	StatusLogcat:          "Logcat",
	StatusPackages:        "Aplikace",
	StatusRemote:          "Ovladač",
	StatusRefresh:         "Obnovit",
	StatusHelp:            "Nápověda",
	StatusQuit:            "Konec",
	StatusExit:            "Odejít",
	StatusUpdated:         "Aktualizováno: %s",
	StatusRecordingActive: "● ZÁZNAM [%s]",

	// Detail card
	DetailTitle:         "🔍 Inspektor zařízení",
	DetailSerialNumber:  "Sériové číslo:",
	DetailADBStatus:     "Stav ADB:",
	DetailModelName:     "Název modelu:",
	DetailCodename:      "Kódové jméno / Produkt:",
	DetailAndroidOS:     "Android OS:",
	DetailIPAddress:     "IP adresa:",
	DetailBattery:       "Baterie:",
	DetailStorage:       "Úložiště:",
	DetailRAM:           "Operační paměť RAM:",
	DetailCPU:           "Vytížení CPU:",
	DetailUptime:        "Doba běhu systému:",
	DetailUSBPort:       "USB sběrnice / Port:",
	DetailResolution:    "Rozlišení displeje:",
	DetailInstalledApps: "Nainstalované aplikace:",

	// Table headers
	TableStatus:    "STAV",
	TableSerial:    "SÉRIOVÉ ČÍSLO",
	TableModel:     "MODEL / NÁZEV",
	TableBattery:   "BATERIE",
	TableTransport: "TR-ID",

	// Actions
	ActionsTitle:            "⚡ Rychlé akce pro: %s (%s)",
	ActionsMirrorTitle:      "Zrcadlení obrazovky (scrcpy)",
	ActionsMirrorDesc:       "Zrcadlit a ovládat displej v okně",
	ActionsRecordTitle:      "Nahrávání obrazovky (video)",
	ActionsRecordDesc:       "Uložit záznam do ./recordings/",
	ActionsFilesTitle:       "Správce souborů zařízení",
	ActionsFilesDesc:        "Procházet, stahovat a mazat soubory",
	ActionsLogcatTitle:      "Živý prohlížeč Logcat",
	ActionsLogcatDesc:       "Sledovat systémové logy v reálném čase",
	ActionsPackagesTitle:    "Správce aplikací a balíčků",
	ActionsPackagesDesc:     "Spouštět, zastavovat, mazat data, odinstalovat",
	ActionsRemoteTitle:      "Dálkové ovládání a klávesy",
	ActionsRemoteDesc:       "Napájení, hlasitost, domů, zpět, text",
	ActionsScreenshotTitle:  "Pořídit snímek",
	ActionsScreenshotDesc:   "Uložit snímek do ./screenshots/",
	ActionsWirelessTitle:    "Zapnout bezdrátový ADB",
	ActionsWirelessDesc:     "Přepnout na TCP port 5555",
	ActionsExportTitle:      "Exportovat diagnostický report",
	ActionsExportDesc:       "Uložit specifikace v Markdown do ./reports/",
	ActionsRebootSystem:     "Restart systému",
	ActionsRebootSystemDesc: "Běžný restart zařízení",
	ActionsRebootRecovery:   "Restart do Recovery",
	ActionsRebootRecDesc:    "Nabootovat do Android Recovery",
	ActionsRebootBoot:       "Restart do Bootloaderu",
	ActionsRebootBootDesc:   "Nabootovat do Fastboot režimu",
	ActionsFooter:           "Stiskněte [Enter] pro spuštění, [Klávesa] pro zkratku, nebo [Esc] pro zrušení",

	// Logcat View
	LogcatTitle:   "📋 Živý Logcat proud — %s",
	LogcatFilter:  "Úroveň: [%s]",
	LogcatSearch:  "Hledat: %s",
	LogcatPaused:  "[POZASTAVENO]",
	LogcatRunning: "[BĚŽÍ]",
	LogcatHint:    "[Space] Pozastavit/Pokračovat • [c] Vymazat • [1-5] Úroveň • [/] Hledat • [Esc] Zpět",

	// File Manager View
	FilesTitle:      "📁 Správce souborů zařízení — %s (%s)",
	FilesHeaderName: "NÁZEV",
	FilesHeaderSize: "VELIKOST",
	FilesHeaderPerm: "OPRÁVNĚNÍ",
	FilesHeaderMod:  "ZMĚNĚNO",
	FilesHint:       "[Enter] Otevřít • [Backspace/b] Nadřazená složka • [d] Stáhnout • [x] Smazat • [Esc] Zpět",
	FilesEmpty:      "Složka je prázdná nebo k ní nemáte přístup",

	// Package Manager View
	PackagesTitle:  "📦 Nainstalované aplikace (%d aplikací) — %s",
	PackagesFilter: "Filtr: %s",
	PackagesToggle: "[Tab] Přepnout 3rd Party / Vše",
	PackagesHint:   "[Enter] Spustit • [x] Zastavit • [c] Smazat data • [u] Odinstalovat • [Esc] Zpět",

	// Remote Control View
	RemoteTitle:     "🎮 Virtuální dálkový ovladač a vstup — %s",
	RemotePower:     "Tlačítko napájení (Power)",
	RemoteVolUp:     "Zvýšit hlasitost (Vol+)",
	RemoteVolDown:   "Snížit hlasitost (Vol-)",
	RemoteHome:      "Domovská obrazovka (Home)",
	RemoteBack:      "Tlačítko Zpět (Back)",
	RemoteRecent:    "Poslední aplikace (Recent)",
	RemoteWake:      "Probudit / Odemknout displej",
	RemoteTextInput: "Odeslat textový vstup",
	RemoteFooter:    "Stiskněte [Klávesu] pro odeslání příkazu, nebo [Esc] pro návrat",

	// Device states
	StateConnected:    "Připojeno",
	StateUnauthorized: "Neautorizováno",
	StateOffline:      "Offline",
	StateAuthorizing:  "Autorizace...",
	StateRecovery:     "Recovery",
	StateSideload:     "Sideload",
	StateBootloader:   "Bootloader",

	// Battery statuses
	BatteryCharging:    "Nabíjí se",
	BatteryDischarging: "Vybíjí se",
	BatteryNotCharging: "Nenabíjí se",
	BatteryFull:        "Plná",
	BatteryUnknown:     "Neznámý",

	// Status messages
	MsgRefreshing:          "Obnovuji zařízení...",
	MsgMonitoringNone:      "Sleduji zařízení (1s)... Žádná zařízení nejsou připojena",
	MsgMonitoringFound:     "Sleduji zařízení (1s)... Nalezeno %d zařízení",
	MsgUnauthorizedAction:  "Zařízení musí být autorizováno pro provádění akcí",
	MsgCannotScreenshot:    "Nelze pořídit snímek: zařízení není autorizováno",
	MsgCannotWireless:      "Nelze zapnout bezdrátový ADB: zařízení není autorizováno",
	MsgCapturingScreenshot: "Pořizuji snímek pro %s...",
	MsgEnablingWireless:    "Zapínám bezdrátový ADB na portu 5555 pro %s...",
	MsgRebooting:           "Restartuji %s...",
	MsgRebootingRecovery:   "Restartuji %s do Recovery...",
	MsgRebootingBootloader: "Restartuji %s do Bootloaderu...",
	MsgScreenshotSaved:     "Snímek uložen do: %s",
	MsgWirelessEnabled:     "Bezdrátový ADB zapnut na portu 5555 (Připojte pomocí: adb connect %s:5555)",
	MsgRebootToSystem:      "Zařízení %s se restartuje do systému...",
	MsgRebootToRecovery:    "Zařízení %s se restartuje do Recovery...",
	MsgRebootToBootloader:  "Zařízení %s se restartuje do Bootloaderu...",
	MsgADBError:            "ADB chyba: %v",
	MsgScreenshotFailed:    "Snímek se nepodařilo: %v",
	MsgWirelessFailed:      "Bezdrátový ADB se nepodařilo: %v",
	MsgRebootFailed:        "Restart se nepodařil: %v",
	MsgMirrorStarted:       "Zrcadlení obrazovky spuštěno v samostatném okně (scrcpy)",
	MsgMirrorFailed:        "Zrcadlení obrazovky selhalo: %v",
	MsgRecordStarted:       "Nahrávání obrazovky spuštěno! Dalším stiskem [v] záznam ukončíte a uložíte.",
	MsgRecordSaved:         "Záznam obrazovky uložen do: %s",
	MsgRecordFailed:        "Nahrávání obrazovky selhalo: %v",
	MsgFileDownloaded:      "Soubor stažen do: %s",
	MsgFileDeleted:         "Smazáno na zařízení: %s",
	MsgFileOpFailed:        "Operace se souborem selhala: %v",
	MsgAppLaunched:         "Aplikace spuštěna: %s",
	MsgAppStopped:          "Aplikace vynuceně zastavena: %s",
	MsgAppCleared:          "Data aplikace vymazána: %s",
	MsgAppUninstalled:      "Aplikace odinstalována: %s",
	MsgKeySent:             "Klávesový příkaz odeslán: %s",
	MsgReportExported:      "Diagnostický report uložen do: %s",
	MsgLanguageChanged:     "Jazyk přepnut na češtinu",
	MsgThemeChanged:        "Téma změněno na: %s",
}

// registry maps language code to its Strings
var registry = map[string]*Strings{
	"en": en,
	"cs": cs,
}

// supportedCodes lists all available locale codes
var supportedCodes = []string{"en", "cs"}

// current holds the active locale (default: en)
var current = en
var currentCode = "en"

// Set sets the active locale by language code (e.g. "en", "cs")
func Set(code string) {
	code = strings.ToLower(strings.TrimSpace(code))
	if l, ok := registry[code]; ok {
		current = l
		currentCode = code
	}
}

// Cycle toggles between supported languages and returns the newly active code
func Cycle() string {
	if currentCode == "en" {
		Set("cs")
		return "cs"
	}
	Set("en")
	return "en"
}

// CurrentCode returns the active locale code
func CurrentCode() string {
	return currentCode
}

// Get returns the current active locale strings
func Get() *Strings {
	return current
}

// Detect reads the system locale from LANG / LC_ALL / LANGUAGE env vars
// and sets the active locale accordingly. Falls back to English.
func Detect() {
	// Explicit override
	if lang := os.Getenv("ANDTLS_LANG"); lang != "" {
		Set(lang)
		return
	}

	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if val := os.Getenv(env); val != "" {
			val = strings.ToLower(val)
			for _, code := range supportedCodes {
				if strings.HasPrefix(val, code) {
					Set(code)
					return
				}
			}
		}
	}

	// Default to English
	Set("en")
}

// SupportedLanguages returns a list of supported locale codes
func SupportedLanguages() []string {
	return append([]string{}, supportedCodes...)
}

