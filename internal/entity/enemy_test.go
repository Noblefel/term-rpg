package entity

import (
	"reflect"
	"testing"
)

func TestEnemyTakeAction(t *testing.T) {
	t.Run("guard", func(t *testing.T) {
		var e EnemyBase
		e.TakeAction(nil, nil, 1)

		if e.GuardTurns != 2 {
			t.Errorf("incorrect effect duration, want %d, got %d", 2, e.GuardTurns)
		}
	})

	t.Run("attack", func(t *testing.T) {
		var p Player
		var e EnemyBase
		p.Hp = 100
		e.isTesting = true

		e.TakeAction(&e, &p, 99)
		if p.Hp != 90.0 {
			t.Errorf("incorrect dmg taken, want hp to be %.1f, got %.1f", 90.0, p.Hp)
		}
	})
}

func TestEnemyTakeDamage(t *testing.T) {
	var e EnemyBase
	e.Hp = 100
	e.TakeDamage(nil, 10)
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
	}

	for _, tt := range tests {
		got := reflect.TypeOf(tt.fn())
		want := reflect.TypeOf(tt.expectedType)
		if want != got {
			t.Errorf("incorrect type, want %q, got %q", want, got)
		}
	}
}

func TestEvilGenieCurse(t *testing.T) {
	t.Run("debuff hp cap", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.HpCap = 100
		e.Curse(&p, 0)
		if p.HpCap == 100 {
			t.Errorf("did not affect hp cap")
		}
	})

	t.Run("debuff attack", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.Att = 10
		e.Curse(&p, 1)
		if p.Att == 10 {
			t.Errorf("did not affect attack")
		}
	})

	t.Run("debuff defense", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.Def = 10
		e.Curse(&p, 2)
		if p.Def == 10 {
			t.Errorf("did not affect defense")
		}
	})

	t.Run("debuff dmg reduction", func(t *testing.T) {
		var p Player
		var e evilGenie
		p.DmgReduc = 1
		e.Curse(&p, 3)
		if p.DmgReduc == 1 {
			t.Errorf("did not affect dmg reduction")
		}
	})

}

func TestSpikeTurtleTakeDamage(t *testing.T) {
	t.Run("reflect damage", func(t *testing.T) {
		var p Player
		var e spikeTurtle
		p.Hp = 100
		e.Hp = 100
		e.TakeDamage(&p, 10)
		if p.Hp == 100 {
			t.Errorf("did not affect player's hp, want less than 100")
		}
	})

	t.Run("should not kill player", func(t *testing.T) {
		var p Player
		var e spikeTurtle
		p.Hp = 0
		e.Hp = 100
		e.TakeDamage(&p, 10)
		if p.Hp <= 0 {
			t.Errorf("player hp is lte 0, want %.1f, got %.1f", 0.1, p.Hp)
		}
	})
}

func TestVampireAttack(t *testing.T) {
	t.Run("should heal hp and deal extra dmg", func(t *testing.T) {
		var p Player
		var v vampire
		p.Hp = 100
		v.HpCap = 10
		v.isTesting = true
		v.Attack(&p)

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
		w.Attack(&p)

		if p.Hp != 90 {
			t.Errorf("want player hp to be 90, got %.1f", p.Hp)
		}
	})
}
