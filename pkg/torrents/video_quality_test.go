package torrents

import (
	"testing"
)

func TestNumeric(t *testing.T) {

	if numeric(High) < numeric(Medium) ||
		numeric(Medium) < numeric(Low) ||
		numeric(Low) < numeric(Default) {
		t.Errorf("Error with VideoQuality numeric conversion")
	}
}
