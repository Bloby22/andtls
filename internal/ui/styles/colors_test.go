package styles

import (
	"testing"
)

func TestThemeCycle(t *testing.T) {
	initial := CurrentThemeName()
	next := CycleTheme()
	if next == initial {
		t.Errorf("expected new theme name, got %s", next)
	}

	ApplyTheme(ThemeCyberpunk)
	if CurrentThemeName() != "Cyberpunk" {
		t.Errorf("expected Cyberpunk theme")
	}
}
