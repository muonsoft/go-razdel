package testkit

import "testing"

func TestDifferentialPythonRequired(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv(EnvDifferentialPython, "")
		if DifferentialPythonRequired() || DifferentialPythonDisabled() {
			t.Fatal("empty env should be optional")
		}
	})
	t.Run("disabled", func(t *testing.T) {
		t.Setenv(EnvDifferentialPython, "0")
		if !DifferentialPythonDisabled() || DifferentialPythonRequired() {
			t.Fatal("0 should disable and not require")
		}
	})
	t.Run("required_1", func(t *testing.T) {
		t.Setenv(EnvDifferentialPython, "1")
		if !DifferentialPythonRequired() || DifferentialPythonDisabled() {
			t.Fatal("1 should require")
		}
	})
	t.Run("required_require", func(t *testing.T) {
		t.Setenv(EnvDifferentialPython, "require")
		if !DifferentialPythonRequired() {
			t.Fatal("require should require")
		}
	})
}
