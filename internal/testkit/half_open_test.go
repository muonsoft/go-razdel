package testkit

import (
	"fmt"
	"testing"
)

func TestValidateHalfOpen_ok(t *testing.T) {
	t.Parallel()
	if err := ValidateHalfOpen(0, 0, 0); err != nil {
		t.Fatal(err)
	}
	if err := ValidateHalfOpen(0, 3, 5); err != nil {
		t.Fatal(err)
	}
	if err := ValidateHalfOpen(2, 2, 4); err != nil {
		t.Fatal(err)
	}
}

func TestValidateHalfOpen_rejects(t *testing.T) {
	t.Parallel()
	tests := []struct {
		start, end, n int
	}{
		{-1, 0, 1},
		{0, -1, 1},
		{2, 1, 3},
		{0, 4, 3},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d_%d_%d", tt.start, tt.end, tt.n), func(t *testing.T) {
			t.Parallel()
			if err := ValidateHalfOpen(tt.start, tt.end, tt.n); err == nil {
				t.Errorf("ValidateHalfOpen(%d,%d,%d): want error", tt.start, tt.end, tt.n)
			}
		})
	}
}
