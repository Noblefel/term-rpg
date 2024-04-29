package entity

import "testing"

func TestNewPlayer(t *testing.T) {
	t.Run("if perk is GREED", func(t *testing.T) {
		p := NewPlayer(GREED)
		var want float32 = 50 + 50*0.15

		if p.Money != want {
			t.Errorf("want %.1f money, got %.1f", want, p.Money)
		}
	})

	t.Run("if perk is HAVOC", func(t *testing.T) {
		p := NewPlayer(HAVOC)
		var want float32 = 85

		if p.Hp != want {
			t.Errorf("want %.1f hp, got %.1f", want, p.Hp)
		}
	})

	t.Run("if perk is RESILIENCY", func(t *testing.T) {
		p := NewPlayer(RESILIENCY)
		var want float32 = 5

		if p.Def != want {
			t.Errorf("want %.1f def, got %.1f", want, p.Def)
		}

		want = 0.1
		if p.DmgReduc != want {
			t.Errorf("want %.1f dmg reduc, got %.1f", want, p.DmgReduc)
		}
	})

	t.Run("if perk is TEMPORAL", func(t *testing.T) {
		p := NewPlayer(TEMPORAL)
		if p.ExtraTurnEffect == 0 {
			t.Errorf("want an extra turn effect, got %d", p.ExtraTurnEffect)
		}
	})
}

func TestPlayerTakeAction(t *testing.T) {
	t.Run("guard", func(t *testing.T) {
		var p Player
		if _, ok := p.TakeAction(nil, 2); !ok {
			t.Errorf("should return success")
		}

		if p.GuardTurns != 2 {
			t.Errorf("incorrect effect duration, want %d, got %d", 2, p.GuardTurns)
		}
	})

	t.Run("cannot stack guard", func(t *testing.T) {
		var p Player
		p.GuardTurns = 1
		if _, ok := p.TakeAction(nil, 2); ok {
			t.Errorf("should not return success")
		}
	})

	t.Run("attack", func(t *testing.T) {
		var p Player
		var e EnemyBase
		e.Hp = 100
		p.isTesting = true

		if _, ok := p.TakeAction(&e, 1); !ok {
			t.Errorf("should return success")
		}

		if e.Hp != 90.0 {
			t.Errorf("incorrect dmg taken, want hp to be %.1f, got %.1f", 90.0, e.Hp)
		}
	})

	t.Run("fury", func(t *testing.T) {
		var p Player
		p.Hp = 100
		p.FuryTurns = -10

		if _, ok := p.TakeAction(nil, 3); !ok {
			t.Errorf("should return success")
		}

		if p.FuryTurns != 2 {
			t.Errorf("incorrect effect duration, want %d, got %d", 2, p.FuryTurns)
		}

		if p.Hp == 100 {
			t.Errorf("did not affect player's hp")
		}
	})

	t.Run("can stack fury", func(t *testing.T) {
		var p Player
		p.Hp = 100
		p.TakeAction(nil, 3)

		if _, ok := p.TakeAction(nil, 3); !ok {
			t.Errorf("should return success")
		}

		if p.FuryTurns != 4 {
			t.Errorf("incorrect effect duration, want %d, got %d", 4, p.FuryTurns)
		}
	})

	t.Run("cannot fury below 10 hp", func(t *testing.T) {
		var p Player
		p.Hp = 1

		if _, ok := p.TakeAction(nil, 3); ok {
			t.Errorf("should not return success")
		}
	})

	t.Run("flee", func(t *testing.T) {
		var p Player
		if _, ok := p.TakeAction(nil, 4); !ok {
			t.Errorf("should return success")
		}
	})

	t.Run("invalid", func(t *testing.T) {
		var p Player
		if _, ok := p.TakeAction(nil, -1); ok {
			t.Errorf("should not return success")
		}
	})
}

func TestPlayerAttack(t *testing.T) {
	var p Player
	p.isTesting = true
	p.Perk = HAVOC
	p.FuryTurns = 1

	var e EnemyBase
	got, _ := p.Attack(&e)
	want := float32(10 + 5)
	want += want * 0.25
	if want != got {
		t.Errorf("want %.1f dmg, got %.1f dmg", want, got)
	}
}

func TestAddMoney(t *testing.T) {
	t.Run("Add normally", func(t *testing.T) {
		p := Player{}
		p.AddMoney(100.0)
		var want float32 = 100.0

		if p.Money != want {
			t.Errorf("want %.1f, got %.1f", want, p.Money)
		}
	})

	t.Run("Add if GREED", func(t *testing.T) {
		p := Player{Perk: GREED}
		p.AddMoney(100.0)
		var want float32 = 115.0

		if p.Money != want {
			t.Errorf("want %.1f, got %.1f", want, p.Money)
		}
	})
}

func TestTrain(t *testing.T) {
	t.Run("buff hp cap", func(t *testing.T) {
		p := Player{}
		p.Train(0)

		if p.HpCap < 1 {
			t.Errorf("%.1f should be more than %.1f", p.HpCap, 1.0)
		}
	})

	t.Run("buff attack", func(t *testing.T) {
		p := Player{}
		p.Train(1)

		if p.Att < 0.5 {
			t.Errorf("%.1f should be more than %.1f", p.Att, 0.5)
		}
	})

	t.Run("buff defense", func(t *testing.T) {
		p := Player{}
		p.Train(2)

		if p.Def < 0.5 {
			t.Errorf("%.1f should be more than %.1f", p.Def, 0.5)
		}
	})

	t.Run("buff dmg reduction", func(t *testing.T) {
		p := Player{}
		p.Train(3)

		if p.DmgReduc != 0.01 {
			t.Errorf("%.2f should be %.2f (%.1f%%)", p.DmgReduc, 0.01, 0.01*100)
		}
	})
}
