package testkit

import (
	"reflect"
	"testing"
)

func TestParsePartition(t *testing.T) {
	t.Parallel()
	tests := []struct {
		in   string
		want []string
	}{
		{in: "", want: nil},
		{in: "a", want: []string{"a"}},
		{in: "a|b|c", want: []string{"a", "b", "c"}},
		{in: "a||b", want: []string{"a", "", "b"}},
		{in: "|x", want: []string{"", "x"}},
		{in: "x|", want: []string{"x", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			got := ParsePartition(tt.in)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParsePartition(%q): got %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}
