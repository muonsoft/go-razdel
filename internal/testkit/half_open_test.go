package testkit

import "testing"

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
		if err := ValidateHalfOpen(tt.start, tt.end, tt.n); err == nil {
			t.Fatalf("ValidateHalfOpen(%d,%d,%d): want error", tt.start, tt.end, tt.n)
		}
	}
}
