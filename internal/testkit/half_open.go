package testkit

import "fmt"

// ValidateHalfOpen checks 0 <= start <= end <= n for a half-open UTF-8 byte interval [start, end)
// into a buffer or string of byte length n (Variant A; see docs/contracts.md).
func ValidateHalfOpen(start, end, n int) error {
	if start < 0 || end < 0 || start > end || end > n {
		return fmt.Errorf("invalid half-open interval [%d,%d) for length %d", start, end, n)
	}
	return nil
}
