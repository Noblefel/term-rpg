package main

import (
	"strings"
	"testing"
)

func TestBars(t *testing.T) {
	t.Run("10 length with max value", func(t *testing.T) {
		bars := bars(10, 10, 10)

		got := strings.Count(bars[0], "━")
		if got != 10 {
			t.Errorf("1st bar should be 10, got: %d", got)
		}

		got = strings.Count(bars[1], "━")
		if got != 0 {
			t.Errorf("2nd bar should be 0, got: %d", got)
		}
	})

	t.Run("10 length with empty value", func(t *testing.T) {
		bars := bars(10, 0, 10)

		got := strings.Count(bars[0], "━")
		if got != 0 {
			t.Errorf("1st bar should be 0, got: %d", got)
		}

		got = strings.Count(bars[1], "━")
		if got != 10 {
			t.Errorf("2nd bar should be 10, got: %d", got)
		}
	})

	t.Run("edge case if value > cap", func(t *testing.T) {
		bars := bars(10, 99999, 10)

		got := strings.Count(bars[0], "━")
		if got != 10 {
			t.Errorf("1st bar should be capped at 10, got: %d", got)
		}

		got = strings.Count(bars[1], "━")
		if got != 0 {
			t.Errorf("2nd bar should be 0, got: %d", got)
		}
	})
}
