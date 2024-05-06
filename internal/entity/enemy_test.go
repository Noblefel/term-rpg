package entity

import (
	"reflect"
	"testing"
)

func TestEnemyTakeAction(t *testing.T) {
	t.Run("guard", func(t *testing.T) {
		var e EnemyBase
		EnemyTakeAction(&e, nil, 1)

		if e.GuardTurns != 2 {
			t.Errorf("incorrect effect duration, want %d, got %d", 2, e.GuardTurns)
		}
	})

	t.Run("fury", func(t *testing.T) {
		var e EnemyBase
		e.Hp = 100
		EnemyTakeAction(&e, nil, 20)

		if e.FuryTurns != 2 {
			t.Errorf("incorrect effect duration, want %d, got %d", 2, e.FuryTurns)
		}
	})

	t.Run("attack", func(t *testing.T) {
		var p Player
		var e EnemyBase
		p.Hp = 100
		e.isTesting = true

		EnemyTakeAction(&e, &p, 99)
		if p.Hp != 90.0 {
			t.Errorf("incorrect dmg taken, want hp to be %.1f, got %.1f", 90.0, p.Hp)
		}
	})
}

func TestEnemyTakeDamage(t *testing.T) {
	var e EnemyBase
	e.Hp = 100
	e.takeDamage(nil, 10)
	if e.Hp != 90.0 {
		t.Errorf("incorrect dmg taken, want hp to be %.1f, got %.1f", 90.0, e.Hp)
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		fn           func() Enemy
		expectedType Enemy
	}{
		{"new thug", newThug, &thug{}},
		{"new acolyte", newAcolyte, &acolyte{}},
		{"new assasin", newAssasin, &assasin{}},
		{"new snakes", newSnakes, &snakes{}},
		{"new golem", newGolem, &golem{}},
		{"new vampire", newVampire, &vampire{}},
		{"new wraith", newWraith, &wraith{}},
		{"new evil genie", newEvilGenie, &evilGenie{}},
		{"new spike turtle", newSpikeTurtle, &spikeTurtle{}},
		{"new changeling", newChangeling, &changeling{}},
	}

	for _, tt := range tests {
		got := reflect.TypeOf(tt.fn())
		want := reflect.TypeOf(tt.expectedType)
		if want != got {
			t.Errorf("incorrect type, want %q, got %q", want, got)
		}
	}
}

func TestChangelingMimic(t *testing.T) {
	var p Player
	p.HpCap = 50
	p.DmgReduc = 0.09

	var e changeling
	e.special(&p)

	if e.Hp != p.HpCap || e.DmgReduc != p.DmgReduc {
		t.Errorf("attributes does not match %.1f, %.1f, %.1f, %.1f", e.Hp, p.Hp, e.DmgReduc*100, p.DmgReduc*100)
	}
}

func TestEvilGenieCurse(t *testing.T) {
	t.Run("debuff hp cap", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.HpCap = 100
		e.curse(&p, 0)
		if p.HpCap == 100 {
			t.Errorf("did not affect hp cap")
		}
	})

	t.Run("debuff attack", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.Att = 10
		e.curse(&p, 1)
		if p.Att == 10 {
			t.Errorf("did not affect attack")
		}
	})

	t.Run("debuff defense", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.Def = 10
		e.curse(&p, 2)
		if p.Def == 10 {
			t.Errorf("did not affect defense")
		}
	})

	t.Run("debuff dmg reduction", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.DmgReduc = 1
		e.curse(&p, 3)
		if p.DmgReduc == 1 {
			t.Errorf("did not affect dmg reduction")
		}
	})

}

func TestSpikeTurtleTakeDamage(t *testing.T) {
	var p Player
	var e spikeTurtle
	p.Hp = 100
	e.Hp = 100
	e.takeDamage(&p, 10)
	if p.Hp == 100 {
		t.Errorf("did not affect player's hp, want less than 100")
	}
}

func TestVampireAttack(t *testing.T) {
	t.Run("should heal hp and deal extra dmg", func(t *testing.T) {
		var p Player
		var v vampire
		p.Hp = 100
		v.HpCap = 10
		v.isTesting = true
		v.attack(&p)

		if v.Hp != 3.1 {
			t.Errorf("want recover 3.1 hp, got %.1f", v.Hp)
		}

		if p.Hp != 85 {
			t.Errorf("want player hp to be 85, got %.1f", p.Hp)
		}
	})
}

func TestWraithAttack(t *testing.T) {
	t.Run("always deal fixed damage", func(t *testing.T) {
		var p Player
		var w wraith
		p.Def = 99999
		p.Hp = 100
		w.Att = 10
		w.attack(&p)

		if p.Hp != 90 {
			t.Errorf("want player hp to be 90, got %.1f", p.Hp)
		}
	})
}
