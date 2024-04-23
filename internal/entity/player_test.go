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
}

func TestPlayerAttack(t *testing.T) {
	tests := []struct {
		name string
		att  float32
		perk int
		want float32
	}{
		{"with 5 attck stat", 5, 0, 15},
		{"with 0 attack stat", 0, 0, 10},
		{"with negative attack stat", -50, 0, 0},
		{"with HAVOC", 0, HAVOC, 12.5},
		{"without HAVOC", 0, -1, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Player{Perk: tt.perk, Att: tt.att}
			p.isTesting = true
			got, _ := p.Attack(&EnemyBase{})

			if got != tt.want {
				t.Errorf("want %.1f, got %.1f", tt.want, got)
			}
		})
	}

	t.Run("with random sum", func(t *testing.T) {
		b := Player{Att: 5}
		got, _ := b.Attack(&EnemyBase{})

		if b.Att > got {
			t.Errorf("want greater than 5, got %.1f", got)
		}
	})
}

func TestPlayerTakeDamage(t *testing.T) {
	tests := []struct {
		name        string
		def         float32
		dmg         float32
		dmgReduc    float32
		expectedHp  float32
		isDefending bool
	}{
		{"10 damage with 4 def", 4, 10, 0, 94, false},
		{"10 damage with 100 def", 100, 10, 0, 100, false},
		{"50 damage with 25 def", 25, 50, 0, 75, false},
		{"8.8 damage with 1.3 def", 1.3, 8.8, 0, 92.5, false},
		{"30 damage with 30% dmg reduction", 0, 30, 0.3, 79, false},
		{"10 damage while defending", 0, 10, 0, 92, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Player{Hp: 100, Def: tt.def, DmgReduc: tt.dmgReduc, IsDefending: tt.isDefending}
			b.TakeDamage(tt.dmg)

			if b.Hp != tt.expectedHp {
				t.Errorf("expected %.1f hp, got %.1f", tt.expectedHp, b.Hp)
			}
		})
	}
}

func TestPlayerHeal(t *testing.T) {
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
			b := Player{HpCap: tt.hpCap}
			b.Heal(tt.recover)

			if b.Hp != tt.want {
				t.Errorf("want %.1f, got %.1f", tt.want, b.Hp)
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

func TestTrain(t *testing.T) {
	t.Run("buff hp cap", func(t *testing.T) {
		p := Player{}
		p.Train(0)

		if p.HpCap < 1 {
			t.Errorf("%.1f should be more than %.1f", p.HpCap, 1.0)
		}
	})

	t.Run("buff att", func(t *testing.T) {
		p := Player{}
		p.Train(1)

		if p.Att < 0.5 {
			t.Errorf("%.1f should be more than %.1f", p.Att, 0.5)
		}
	})

	t.Run("buff def", func(t *testing.T) {
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
