package torrents

import (
	"testing"
)

func TestWorseThan(t *testing.T) {

	if High.WorseThan(High) ||
		High.WorseThan(Medium) ||
		Medium.WorseThan(Medium) ||
		Medium.WorseThan(Low) ||
		Low.WorseThan(Low) ||
		Low.WorseThan(Default) ||
		Default.WorseThan(Default) {
		t.Errorf("Error with VideoQuality numeric conversion")
	}
}

func TestBetterThan(t *testing.T) {

	if High.BetterThan(High) ||
		Medium.BetterThan(High) ||
		Medium.BetterThan(Medium) ||
		Low.BetterThan(Medium) ||
		Low.BetterThan(Low) ||
		Default.BetterThan(Low) ||
		Default.BetterThan(Default) {
		t.Errorf("Error with VideoQuality numeric conversion")
	}
}

func TestNumeric(t *testing.T) {

	if High.numeric() < Medium.numeric() ||
		Medium.numeric() < Low.numeric() ||
		Low.numeric() < Default.numeric() {
		t.Errorf("Error with VideoQuality numeric conversion")
	}
}
