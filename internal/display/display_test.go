package display

import (
	"reflect"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func assertStringW(t testing.TB, w1, w2 strings.Builder) {
	t.Helper()
	s1 := w1.String()
	s2 := w2.String()

	if s1 != s2 {
		t.Errorf("\nwant: \n%sgot: \n%s", s1, s2)
	}
}

func TestNew(t *testing.T) {
	d := New(nil)
	got := reflect.TypeOf(d)
	want := reflect.TypeOf(&Display{})

	if want != got {
		t.Errorf("incorrect type, want %q, got %q", want, got)
	}
}

func TestPrint(t *testing.T) {
	f := color.New(color.FgWhite).FprintfFunc()

	var w strings.Builder
	d := New(&w)
	d.Printf(f, "test")

	var w2 strings.Builder
	f(&w2, "test")

	assertStringW(t, w2, w)
}

func TestCenter(t *testing.T) {
	t.Run("with color func", func(t *testing.T) {
		f := color.New(color.FgWhite).FprintfFunc()
		var w strings.Builder
		d := New(&w)
		d.Center(f, "test")

		var w2 strings.Builder
		w2.WriteString("\n")
		w2.WriteString("                     test")
		w2.WriteString("\n")

		assertStringW(t, w2, w)
	})

	t.Run("without color func", func(t *testing.T) {
		var w strings.Builder
		d := New(&w)
		d.Center(nil, "test")

		var w2 strings.Builder
		w2.WriteString("\n")
		w2.WriteString("                     test")
		w2.WriteString("\n")

		assertStringW(t, w2, w)
	})
}

func TestBar(t *testing.T) {
	tests := []struct {
		name   string
		currHp float32
		maxHp  float32
		want   string
	}{
		{"at full hp", 100, 100, "██████████████████████████████████████████████████"},
		{"at full hp 2", 44, 44, "██████████████████████████████████████████████████"},
		{"at half hp", 50, 100, "█████████████████████████"},
		{"at half hp 2", 9, 18, "█████████████████████████"},
		{"at 10% hp", 10, 100, "█████"},
		{"at 10% hp 2", 2.5, 22.5, "█████"},
		{"at 0 hp", 0, 100, ""},
		{"if somehow negative", -5, 100, ""},
		{"if currHp exceeds maxHp", 100, 80, "██████████████████████████████████████████████████"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := color.New(color.FgWhite).FprintfFunc()

			var w strings.Builder
			d := New(&w)
			d.isTesting = true
			d.Bar(tt.currHp, tt.maxHp)

			var w2 strings.Builder
			f(&w2, "%s", tt.want)

			assertStringW(t, w2, w)
		})
	}

	t.Run("with isTesting false", func(t *testing.T) {
		f := color.New(color.FgWhite).FprintfFunc()

		var w strings.Builder
		d := New(&w)
		d.Bar(0, 1)

		var w2 strings.Builder
		f(&w2, "%s", "██████████████████████████████████████████████████\n")

		assertStringW(t, w2, w)
	})
}
