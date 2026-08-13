package testkit

import (
	"os"
	"strings"
	"testing"
)

// EnvDifferentialPython controls Go-vs-Python differential tests.
//
//	unset  — run when python3 + submodule are available, otherwise skip
//	"0"    — always skip
//	"1" or "require" — fail the test if Python razdel is not available (CI)
const EnvDifferentialPython = "RAZDEL_DIFFERENTIAL_PYTHON"

// DifferentialPythonDisabled reports RAZDEL_DIFFERENTIAL_PYTHON=0.
func DifferentialPythonDisabled() bool {
	return os.Getenv(EnvDifferentialPython) == "0"
}

// DifferentialPythonRequired reports that missing Python razdel must fail the test.
func DifferentialPythonRequired() bool {
	v := os.Getenv(EnvDifferentialPython)
	return v == "1" || strings.EqualFold(v, "require")
}

// SkipOrFailPython skips tb, or fatals when differential Python is required.
func SkipOrFailPython(tb testing.TB, reason string) {
	tb.Helper()
	if DifferentialPythonRequired() {
		tb.Fatalf("python differential required (%s=1): %s", EnvDifferentialPython, reason)
	}
	tb.Skip(reason)
}
