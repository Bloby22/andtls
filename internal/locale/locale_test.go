package locale

import (
	"testing"
)

func TestLocaleCycleAndSet(t *testing.T) {
	Set("en")
	if CurrentCode() != "en" {
		t.Errorf("expected code en, got %s", CurrentCode())
	}

	code := Cycle()
	if code != "cs" || CurrentCode() != "cs" {
		t.Errorf("expected code cs after cycle, got %s", code)
	}

	l := Get()
	if l.AppTitle != "📱 ANDTLS" || l.StatusScreenshot != "Snímek" {
		t.Errorf("Czech locale strings mismatch: %+v", l)
	}

	Cycle()
	if CurrentCode() != "en" {
		t.Errorf("expected en after second cycle")
	}
}
