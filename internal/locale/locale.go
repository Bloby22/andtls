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
	HelpToggleHelp         string
	HelpCloseModal         string
	HelpQuit               string
	HelpCloseHint          string

	// Empty state
	EmptyTitle      string
	EmptyTip1       string
	EmptyTip2       string
	EmptyTip3       string
	EmptyTip4       string

	// Status bar
	StatusNavigate  string
	StatusActions   string
	StatusScreenshot string
	StatusRefresh   string
	StatusHelp      string
	StatusQuit      string
	StatusExit      string
	StatusUpdated   string

	// Detail card
	DetailTitle          string
	DetailSerialNumber   string
	DetailADBStatus      string
	DetailModelName      string
	DetailCodename       string
	DetailAndroidOS      string
	DetailIPAddress      string
	DetailBattery        string
	DetailStorage        string
	DetailUSBPort        string
	DetailResolution     string

	// Table headers
	TableStatus    string
	TableSerial    string
	TableModel     string
	TableBattery   string
	TableTransport string

	// Actions
	ActionsTitle            string
	ActionsScreenshotTitle  string
	ActionsScreenshotDesc   string
	ActionsWirelessTitle    string
	ActionsWirelessDesc     string
	ActionsRebootSystem     string
	ActionsRebootSystemDesc string
	ActionsRebootRecovery   string
	ActionsRebootRecDesc    string
	ActionsRebootBoot       string
	ActionsRebootBootDesc   string
	ActionsFooter           string

	// Device states
	StateConnected     string
	StateUnauthorized  string
	StateOffline       string
	StateAuthorizing   string
	StateRecovery      string
	StateSideload      string
	StateBootloader    string

	// Battery statuses
	BatteryCharging     string
	BatteryDischarging  string
	BatteryNotCharging  string
	BatteryFull         string
	BatteryUnknown      string

	// Status messages
	MsgRefreshing         string
	MsgMonitoringNone     string
	MsgMonitoringFound    string // with %d
	MsgUnauthorizedAction string
	MsgCannotScreenshot   string
	MsgCannotWireless     string
	MsgCapturingScreenshot string
	MsgEnablingWireless   string
	MsgRebooting          string
	MsgRebootingRecovery  string
	MsgRebootingBootloader string
	MsgScreenshotSaved    string
	MsgWirelessEnabled    string
	MsgRebootToSystem     string
	MsgRebootToRecovery   string
	MsgRebootToBootloader string
	MsgADBError           string
	MsgScreenshotFailed   string
	MsgWirelessFailed     string
	MsgRebootFailed       string
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
	HelpScreenshot:         "Take screenshot of selected device (saved to ./screenshots/)",
	HelpWirelessADB:        "Enable wireless ADB mode on port 5555 for selected device",
	HelpRebootMenu:         "Open reboot options menu (System, Recovery, Bootloader)",
	HelpActionsMenu:        "Open quick actions menu for selected device",
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
	StatusNavigate:   "Navigate",
	StatusActions:    "Actions",
	StatusScreenshot: "Screenshot",
	StatusRefresh:    "Refresh",
	StatusHelp:       "Help",
	StatusQuit:       "Quit",
	StatusExit:       "Exit",
	StatusUpdated:    "Updated: %s",

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
	DetailUSBPort:       "USB Bus / Port:",
	DetailResolution:    "Display Resolution:",

	// Table headers
	TableStatus:    "STATUS",
	TableSerial:    "SERIAL NUMBER",
	TableModel:     "MODEL / NAME",
	TableBattery:   "BATTERY",
	TableTransport: "TR-ID",

	// Actions
	ActionsTitle:            "⚡ Quick Actions for: %s (%s)",
	ActionsScreenshotTitle:  "Take Screenshot",
	ActionsScreenshotDesc:   "Capture display to ./screenshots/",
	ActionsWirelessTitle:    "Enable Wireless ADB",
	ActionsWirelessDesc:     "Switch daemon to TCP port 5555",
	ActionsRebootSystem:     "Reboot System",
	ActionsRebootSystemDesc: "Normal device restart",
	ActionsRebootRecovery:   "Reboot to Recovery",
	ActionsRebootRecDesc:    "Boot into Android Recovery",
	ActionsRebootBoot:       "Reboot to Bootloader",
	ActionsRebootBootDesc:   "Boot into Fastboot mode",
	ActionsFooter:           "Press [Enter] to run selected, [Key] for shortcut, or [Esc] to cancel",

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
	HelpScreenshot:         "Pořídit snímek obrazovky vybraného zařízení (uloží do ./screenshots/)",
	HelpWirelessADB:        "Zapnout bezdrátový ADB mód na portu 5555 pro vybrané zařízení",
	HelpRebootMenu:         "Otevřít nabídku restartu (Systém, Recovery, Bootloader)",
	HelpActionsMenu:        "Otevřít rychlé akce pro vybrané zařízení",
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
	StatusNavigate:   "Pohyb",
	StatusActions:    "Akce",
	StatusScreenshot: "Snímek",
	StatusRefresh:    "Obnovit",
	StatusHelp:       "Nápověda",
	StatusQuit:       "Konec",
	StatusExit:       "Odejít",
	StatusUpdated:    "Aktualizováno: %s",

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
	DetailUSBPort:       "USB sběrnice / Port:",
	DetailResolution:    "Rozlišení displeje:",

	// Table headers
	TableStatus:    "STAV",
	TableSerial:    "SÉRIOVÉ ČÍSLO",
	TableModel:     "MODEL / NÁZEV",
	TableBattery:   "BATERIE",
	TableTransport: "TR-ID",

	// Actions
	ActionsTitle:            "⚡ Rychlé akce pro: %s (%s)",
	ActionsScreenshotTitle:  "Pořídit snímek",
	ActionsScreenshotDesc:   "Uložit snímek do ./screenshots/",
	ActionsWirelessTitle:    "Zapnout bezdrátový ADB",
	ActionsWirelessDesc:     "Přepnout na TCP port 5555",
	ActionsRebootSystem:     "Restart systému",
	ActionsRebootSystemDesc: "Běžný restart zařízení",
	ActionsRebootRecovery:   "Restart do Recovery",
	ActionsRebootRecDesc:    "Nabootovat do Android Recovery",
	ActionsRebootBoot:       "Restart do Bootloaderu",
	ActionsRebootBootDesc:   "Nabootovat do Fastboot režimu",
	ActionsFooter:           "Stiskněte [Enter] pro spuštění, [Klávesa] pro zkratku, nebo [Esc] pro zrušení",

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

// Set sets the active locale by language code (e.g. "en", "cs")
func Set(code string) {
	code = strings.ToLower(strings.TrimSpace(code))
	if l, ok := registry[code]; ok {
		current = l
	}
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
