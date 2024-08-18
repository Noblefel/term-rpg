package main

import (
	"strings"
	"testing"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

func TestChoice(t *testing.T) {
	t.Run("should pick 'c' with 2 arrow down keypress", func(t *testing.T) {
		var sb strings.Builder
		choice := NewChoice(&sb, "a", "b", "c")

		go func() {
			time.Sleep(10 * time.Millisecond)
			keyboard.SimulateKeyPress(keys.Down)
			keyboard.SimulateKeyPress(keys.Down)
			keyboard.SimulateKeyPress(keys.Enter)
		}()

		if got := choice.Select(); got != 2 {
			t.Fatalf("should return index of 2 from the options, got %d", got)
		}
	})

	t.Run("should pick 'a' by default", func(t *testing.T) {
		var sb strings.Builder
		choice := NewChoice(&sb, "a", "b", "c")

		go func() {
			time.Sleep(10 * time.Millisecond)
			keyboard.SimulateKeyPress(keys.Enter)
		}()

		if got := choice.Select(); got != 0 {
			t.Fatalf("should return index of 0 from the options, got %d", got)
		}
	})

	t.Run("should pick 'b' with 2 arrow up keypress", func(t *testing.T) {
		var sb strings.Builder
		choice := NewChoice(&sb, "a", "b", "c")

		go func() {
			time.Sleep(10 * time.Millisecond)
			keyboard.SimulateKeyPress(keys.Up)
			keyboard.SimulateKeyPress(keys.Up)
			keyboard.SimulateKeyPress(keys.Enter)
		}()

		if got := choice.Select(); got != 1 {
			t.Fatalf("should return index of 1 from the options, got %d", got)
		}
	})
}

func TestChoice_PrintOptions(t *testing.T) {
	var sb strings.Builder
	choice := NewChoice(&sb, "a", "b", "c")
	choice.current = 2
	choice.printOptions()

	if !strings.Contains(sb.String(), "\033[1m> c") {
		t.Errorf("should highlight the selected option\ngot: %q", sb.String())
	}
}
