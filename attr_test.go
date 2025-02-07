package main

import (
	"strings"
	"testing"
)

func TestAttributes_Attack(t *testing.T) {
	target := attributes{hp: 100}
	attacker := attributes{strength: 50}
	attacker.attack(&target)

	if target.hp != 50 {
		t.Errorf("target hp should be 50, got %.1f", target.hp)
	}
}

func TestAttributes_Damage(t *testing.T) {
	target := attributes{hp: 100, defense: 10}
	target.damage(50)

	if target.hp != 60 {
		t.Errorf("target hp should be 60, got %.1f", target.hp)
	}

	target.hp = 100
	target.damage(9999)

	if target.hp != 0 {
		t.Errorf("target hp should be 0, got %.1f", target.hp)
	}
}

func TestAttributes_Hpbar(t *testing.T) {
	t.Run("green bar (above 60%)", func(t *testing.T) {
		attr := attributes{hp: 10, hpcap: 10}
		got := attr.hpbar()

		if !strings.Contains(got, "38;5;83m") {
			t.Errorf("missing color code '38;5;83m'\ngot:%q", got)
		}
	})

	t.Run("orange bar (above 30%)", func(t *testing.T) {
		attr := attributes{hp: 5, hpcap: 10}
		got := attr.hpbar()

		if !strings.Contains(got, "38;5;226m") {
			t.Errorf("missing color code '38;5;226m'\ngot:%q", got)
		}
	})

	t.Run("red bar (below 30%)", func(t *testing.T) {
		attr := attributes{hp: 2, hpcap: 10}
		got := attr.hpbar()

		if !strings.Contains(got, "38;5;196m") {
			t.Errorf("missing green color code '38;5;196m'\ngot:%q", got)
		}
	})
}
