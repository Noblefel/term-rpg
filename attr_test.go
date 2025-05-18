package main

import (
	"strings"
	"testing"
)

func TestAttributes_AttackWith(t *testing.T) {
	target := attributes{hp: 100}
	attacker := attributes{strength: 50}
	attacker.effects = make(map[string]int)
	attacker.attackWith(&target, attacker.strength)

	if target.hp != 50 {
		t.Errorf("target hp should be 50, got %.1f", target.hp)
	}

	t.Run("strengthen effect", func(t *testing.T) {
		target.hp = 100
		attacker.effects["strengthen"] = 1
		attacker.attackWith(&target, 10)
		clear(attacker.effects)

		if 100-target.hp != 11 {
			t.Errorf("damage should be 11 (10%% bonus), got %.1f", 100-target.hp)
		}
	})

	t.Run("weakened effect", func(t *testing.T) {
		target.hp = 100
		attacker.effects["weakened"] = 1
		attacker.attackWith(&target, 10)

		if !equal(100-target.hp, 8.7) {
			t.Errorf("damage should be 8.7 (13%% reduction), got %.1f", 100-target.hp)
		}
	})
}

func TestAttributes_ApplyEffects(t *testing.T) {
	var target attributes
	target.effects = make(map[string]int)

	t.Run("poisoned", func(t *testing.T) {
		target.hpcap = 200
		target.hp = 100
		target.effects["poisoned"] = 1
		target.applyEffects()
		clear(target.effects)

		dmg := float32(100)*0.07 + 2
		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (7%% of hp + 2), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("poisoned severe", func(t *testing.T) {
		target.hpcap = 200
		target.hp = 100
		target.effects["poisoned severe"] = 1
		target.applyEffects()
		clear(target.effects)

		dmg := float32(100)*0.16 + 2
		if !equal(dmg, 100-target.hp) {
			t.Errorf("damage should be %.1f (16%% of hp + 2), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("burning", func(t *testing.T) {
		target.hpcap = 200
		target.hp = 100
		target.effects["burning"] = 1
		target.applyEffects()
		clear(target.effects)

		dmg := float32(200)*0.04 + 5
		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (4%% of hpcap + 5), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("burning severe", func(t *testing.T) {
		target.hpcap = 200
		target.hp = 100
		target.effects["burning severe"] = 1
		target.applyEffects()
		clear(target.effects)

		dmg := float32(200)*0.10 + 5
		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (10%% of hpcap + 5), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("heal aura", func(t *testing.T) {
		target.hpcap = 100
		target.hp = 0
		target.effects["heal aura"] = 1
		target.applyEffects()
		clear(target.effects)

		if target.hp != 7 {
			t.Errorf("should heal by 7 (7%% of hpcap), got %.1f", target.hp)
		}
	})
}

func TestAttributes_Damage(t *testing.T) {
	target := attributes{hp: 100, defense: 10}
	target.effects = make(map[string]int)
	target.damage(50)

	if target.hp != 60 {
		t.Errorf("target hp should be 60, got %.1f", target.hp)
	}

	target.hp = 100
	target.defense = 0
	target.damage(9999)

	if target.hp != 0 {
		t.Errorf("target hp should be 0, got %.1f", target.hp)
	}

	t.Run("immunity effect", func(t *testing.T) {
		target.hp = 100
		target.effects["immunity"] = 1
		target.damage(9999999)
		clear(target.effects)

		if 100-target.hp != 0 {
			t.Errorf("damage should be 0 (immune), got %.1f", 100-target.hp)
		}
	})

	t.Run("barrier effect", func(t *testing.T) {
		target.hp = 100
		target.effects["barrier"] = 1
		target.damage(10)
		clear(target.effects)

		if 100-target.hp != 6 {
			t.Errorf("damage should be 6 (40%% reduction), got %.1f", 100-target.hp)
		}
	})

	t.Run("force-field effect", func(t *testing.T) {
		target.hp = 100
		target.effects["force-field"] = 1
		target.damage(10)
		clear(target.effects)

		if 100-target.hp != 8.5 {
			t.Errorf("damage should be 8.5 (15%% reduction), got %.1f", 100-target.hp)
		}
	})

	t.Run("weakened effect", func(t *testing.T) {
		target.hp = 100
		target.defense = 10
		target.effects["weakened"] = 1
		target.damage(10)

		if 100-target.hp != 5 {
			t.Errorf("damage should be 5 (50%% defense reduction), got %.1f", 100-target.hp)
		}
	})
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
