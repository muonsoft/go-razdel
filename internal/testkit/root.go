package testkit

import (
	"os"
	"path/filepath"
	"testing"
)

// ModuleRoot walks up from the working directory to find the module root (directory containing go.mod).
func ModuleRoot(tb testing.TB) string {
	tb.Helper()
	dir, err := os.Getwd()
	if err != nil {
		tb.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			tb.Fatal("go.mod not found when walking up from working directory")
		}
		dir = parent
	}
}
