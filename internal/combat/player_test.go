package combat

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
}

func TestPlayerAttack(t *testing.T) {
	tests := []struct {
		name string
		perk int
		want float32
	}{
		{"with HAVOC", HAVOC, 12.5},
		{"without HAVOC", -1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{Perk: tt.perk}
			p.isTesting = true
			got := p.Attack()

			if got != tt.want {
				t.Errorf("want %.1f, got %.1f", tt.want, got)
			}
		})
	}
}

func TestRecoverHP(t *testing.T) {
	tests := []struct {
		name    string
		recover float32
		hpCap   float32
		want    float32
	}{
		{"Recover normally", 50.0, 50.0, 50.0},
		{"Recover normally 2", 20.0, 125.0, 20.0},
		{"Recover exceeds cap", 200.0, 100.0, 100.0},
		{"Recover exceeds cap", 1.1, 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{}
			p.HpCap = tt.hpCap
			p.RecoverHP(tt.recover)

			if p.Hp != tt.want {
				t.Errorf("want %.1f, got %.1f", tt.want, p.Hp)
			}
		})
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
