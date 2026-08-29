package device

import (
	"fmt"
)

// ChangeType represents the nature of a device event
type ChangeType int

const (
	// ChangeTypeConnected indicates a newly attached device
	ChangeTypeConnected ChangeType = iota

	// ChangeTypeDisconnected indicates a detached device
	ChangeTypeDisconnected

	// ChangeTypeStateChanged indicates an authorization or mode change
	ChangeTypeStateChanged
)

// ChangeEvent contains details about a detected state transition on a device
type ChangeEvent struct {
	// Type specifies the category of the change event
	Type ChangeType

	// Device is the current snapshot of the device
	Device Device

	// OldState is the previous state before this change
	OldState DeviceState

	// NewState is the updated state after this change
	NewState DeviceState

	// Message is a human-readable notification string
	Message string
}

// Diff compares previous and current device slices and returns a list of detected change events
func Diff(oldDevices, newDevices []Device) []ChangeEvent {
	oldMap := make(map[string]Device, len(oldDevices))
	for _, d := range oldDevices {
		oldMap[d.Serial] = d
	}

	newMap := make(map[string]Device, len(newDevices))
	var events []ChangeEvent

	// Detect new connections and state changes
	for _, newDev := range newDevices {
		newMap[newDev.Serial] = newDev
		if oldDev, exists := oldMap[newDev.Serial]; !exists {
			events = append(events, ChangeEvent{
				Type:     ChangeTypeConnected,
				Device:   newDev,
				NewState: newDev.State,
				Message:  fmt.Sprintf("Device connected: %s (%s)", newDev.DisplayName(), newDev.Serial),
			})
		} else if oldDev.State != newDev.State {
			events = append(events, ChangeEvent{
				Type:     ChangeTypeStateChanged,
				Device:   newDev,
				OldState: oldDev.State,
				NewState: newDev.State,
				Message: fmt.Sprintf("Device %s state changed: %s -> %s",
					newDev.DisplayName(), oldDev.DisplayState(), newDev.DisplayState()),
			})
		}
	}

	// Detect disconnections
	for _, oldDev := range oldDevices {
		if _, exists := newMap[oldDev.Serial]; !exists {
			events = append(events, ChangeEvent{
				Type:     ChangeTypeDisconnected,
				Device:   oldDev,
				OldState: oldDev.State,
				Message:  fmt.Sprintf("Device disconnected: %s (%s)", oldDev.DisplayName(), oldDev.Serial),
			})
		}
	}

	return events
}
