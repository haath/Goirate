package torrents

import (
	"testing"
)

func TestSizeKB(t *testing.T) {
	table := []struct {
		in  string
		out int64
	}{
		{"64 KB", 64},
	}

	for _, tt := range table {
		t.Run(tt.in, func(t *testing.T) {

			f := SearchFilters{MinSize: tt.in, MaxSize: tt.in}

			s1, err := f.MinSizeKB()

			if err != nil {
				t.Error(err)
			}

			if s1 != tt.out {
				t.Errorf("got %v, want %v", s1, tt.out)
			}

			s2, err := f.MaxSizeKB()

			if err != nil {
				t.Error(err)
			}

			if s2 != tt.out {
				t.Errorf("got %v, want %v", s2, tt.out)
			}
		})
	}
}
