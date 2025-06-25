package main

import (
	"strings"
	"testing"
)

func TestAttributes_AttackWith(t *testing.T) {
	target := &attributes{hp: 100}
	target.effects = make(map[string]int)

	attacker := &attributes{strength: 50}
	attacker.effects = make(map[string]int)
	attacker.attackWith(target, attacker.strength)

	if target.hp != 50 {
		t.Errorf("target hp should be 50, got %.1f", target.hp)
	}

	t.Run("miss", func(t *testing.T) {
		target.hp = 100
		target.agility = 9999
		attacker.attack(target)
		target.agility = 0

		if target.hp != 100 {
			t.Errorf("target hp should be 100 because it missed, got %.1f", target.hp)
		}
	})

	t.Run("strengthen effect", func(t *testing.T) {
		target.hp = 100
		attacker.effects["strengthen"] = 1
		attacker.attackWith(target, 10)
		clear(attacker.effects)

		if 100-target.hp != 11 {
			t.Errorf("damage should be 11 (10%% bonus), got %.1f", 100-target.hp)
		}
	})

	t.Run("weakened effect", func(t *testing.T) {
		target.hp = 100
		attacker.effects["weakened"] = 1
		attacker.attackWith(target, 10)
		clear(attacker.effects)

		if !equal(100-target.hp, 8.7) {
			t.Errorf("damage should be 8.7 (13%% reduction), got %.1f", 100-target.hp)
		}
	})

	t.Run("ace", func(t *testing.T) {
		target.hp = 100
		attacker.effects["ace"] = 1
		attacker.attackWith(target, 10)
		clear(attacker.effects)

		if !equal(100-target.hp, 12.8) {
			t.Errorf("damage should be 12.8 (28%% bonus), got %.1f", 100-target.hp)
		}
	})

	t.Run("crit", func(t *testing.T) {
		target.hp = 100
		attacker.strength = 10
		attacker.agility = 9999
		attacker.attack(target)
		attacker.agility = 0

		if !equal(100-target.hp, 17.5) {
			t.Errorf("damage should be 17.5 (x175%%), got %.1f", 100-target.hp)
		}
	})

	t.Run("reflect", func(t *testing.T) {
		target.effects["reflect"] = 1
		attacker.hp = 10
		attacker.hpcap = 10
		attacker.strength = 10
		attacker.attack(target)

		if !equal(10-attacker.hp, 3) {
			t.Errorf("should inflict 3 damage (30%% reflection), got %.1f", 10-attacker.hp)
		}
	})
}

func TestAttributes_ApplyEffects(t *testing.T) {
	var target attributes
	target.hpcap = 200
	target.hp = 100
	target.defense = 10
	target.effects = make(map[string]int)

	t.Run("poisoned", func(t *testing.T) {
		target.effects["poisoned"] = 1
		target.applyEffects()
		clear(target.effects)

		dmg := 100*0.11 + target.defense*0.5 + 10
		dmg -= target.defense

		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (10%% of hp + 10 + 50%% def), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("poisoned severe", func(t *testing.T) {
		target.hp = 100
		target.effects["poisoned severe"] = 1
		target.applyEffects()
		clear(target.effects)

		dmg := 100*0.22 + target.defense*0.5 + 20
		dmg -= target.defense

		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (20%% of hp + 20 + 50%% def), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("burning", func(t *testing.T) {
		target.hp = 100
		target.effects["burning"] = 1
		target.effects["frozen"] = 5
		target.applyEffects()
		clear(target.effects)

		dmg := 200*0.05 + target.defense*0.5 + 10
		dmg -= target.defense

		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (6%% hp cap + 10 + 50%% def), got %.1f", dmg, 100-target.hp)
		}

		if target.effects["frozen"] != 0 {
			t.Error("should remove frozen effect")
		}
	})

	t.Run("burning severe", func(t *testing.T) {
		target.hp = 100
		target.effects["burning severe"] = 1
		target.effects["frozen"] = 5
		target.applyEffects()
		clear(target.effects)

		dmg := 200*0.1 + target.defense*0.5 + 20
		dmg -= target.defense

		if dmg != 100-target.hp {
			t.Errorf("damage should be %.1f (12%% hp cap + 20 + 50%% def), got %.1f", dmg, 100-target.hp)
		}

		if target.effects["frozen"] != 0 {
			t.Error("should remove frozen effect")
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

	t.Run("frozen effect", func(t *testing.T) {
		rolltest = 1
		target.hp = 100
		target.effects["frozen"] = 1
		target.damage(10)

		if 100-target.hp != 20 {
			t.Errorf("damage should be 20 (shatter x2 bonus), got %.1f", 100-target.hp)
		}

		rolltest = 75
		target.hp = 100
		target.effects["frozen"] = 1
		target.defense = 10
		target.damage(20)
		clear(target.effects)

		dmg := 20 - target.defense*1.25
		if 100-target.hp != dmg {
			t.Errorf("damage should be %.1f (25%% defense increase), got %.1f", dmg, 100-target.hp)
		}
	})

	t.Run("weakened effect", func(t *testing.T) {
		target.hp = 100
		target.defense = 10
		target.effects["weakened"] = 1
		target.damage(10)
		clear(target.effects)

		if 100-target.hp != 5 {
			t.Errorf("damage should be 5 (50%% defense reduction), got %.1f", 100-target.hp)
		}
	})

	t.Run("ace effect", func(t *testing.T) {
		target.hp = 100
		target.defense = 0
		target.effects["ace"] = 1
		target.damage(10)

		if !equal(100-target.hp, 7.2) {
			t.Errorf("damage should be 7.2 (28%% reduction), got %.1f", 100-target.hp)
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
