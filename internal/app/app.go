package app

import (
	"fmt"

	"github.com/Bloby22/andtls/internal/adb"
	"github.com/Bloby22/andtls/internal/config"
	"github.com/Bloby22/andtls/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Run initializes application configuration, ADB client, and starts the Bubble Tea TUI
func Run() error {
	cfg := config.DefaultConfig()

	// Initialize ADB Client
	adbClient := adb.NewClientWithPath(cfg.ADBPath, cfg.CommandTimeout)

	// Create Bubble Tea Model
	model := ui.NewModel(adbClient, cfg.PollInterval)

	// Launch program in fullscreen alt screen
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("application run error: %w", err)
	}

	return nil
}
