package ui

import (
	"time"

	"github.com/Bloby22/andtls/internal/adb"
	"github.com/Bloby22/andtls/internal/device"
	"github.com/Bloby22/andtls/internal/locale"
	tea "github.com/charmbracelet/bubbletea"
)

// ViewMode represents the current UI overlay display mode
type ViewMode int

const (
	// ViewModeNormal is the default split dashboard view
	ViewModeNormal ViewMode = iota

	// ViewModeHelp displays the keyboard shortcut reference modal
	ViewModeHelp

	// ViewModeActions displays the quick action menu for the active device
	ViewModeActions
)

// Model represents the Bubble Tea application state for the andtls dashboard
type Model struct {
	adbClient           *adb.Client
	devices             []device.Device
	selectedIndex       int
	statusMsg           string
	isErr               bool
	lastChecked         time.Time
	width               int
	height              int
	ready               bool
	pollingInterval     time.Duration
	lastEvents          []string
	viewMode            ViewMode
	selectedActionIndex int
	locale              *locale.Strings
}

// NewModel initializes and returns a new Model instance
func NewModel(client *adb.Client, pollInterval time.Duration) Model {
	if pollInterval <= 0 {
		pollInterval = 1 * time.Second
	}
	locale.Detect()
	return Model{
		adbClient:           client,
		devices:             []device.Device{},
		selectedIndex:       0,
		statusMsg:           "Monitoring devices...",
		isErr:               false,
		lastChecked:         time.Now(),
		pollingInterval:     pollInterval,
		lastEvents:          make([]string, 0),
		viewMode:            ViewModeNormal,
		selectedActionIndex: 0,
		locale:              locale.Get(),
	}
}

// Init sets up initial commands (immediate device fetch and first tick)
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		fetchDevicesCmd(m.adbClient),
		tickCmd(m.pollingInterval),
	)
}

// SelectedDevice returns a pointer to the currently selected device, or nil if none is selected
func (m Model) SelectedDevice() *device.Device {
	if len(m.devices) == 0 || m.selectedIndex < 0 || m.selectedIndex >= len(m.devices) {
		return nil
	}
	return &m.devices[m.selectedIndex]
}
