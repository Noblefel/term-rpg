package entity

import "testing"

func TestBaseAttack(t *testing.T) {
	tests := []struct {
		name string
		att  float32
		want float32
	}{
		{"with 5 attack stat", 5, 15},
		{"with 0 attack stat", 0, 10},
		{"with negative attack stat", -50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := base{Att: tt.att, isTesting: true}
			got := b.attack()

			if tt.want != got {
				t.Errorf("want %.1f, got %.1f", tt.want, got)
			}
		})
	}

	t.Run("with random sum", func(t *testing.T) {
		b := base{Att: 5}
		got := b.attack()

		if b.Att > got {
			t.Errorf("want greater than 5, got %.1f", got)
		}
	})
}

func TestBaseTakeDamage(t *testing.T) {
	tests := []struct {
		name       string
		def        float32
		dmg        float32
		dmgReduc   float32
		expectedHp float32
		onGuard    bool
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
			b := base{Hp: 100, Def: tt.def, DmgReduc: tt.dmgReduc}
			if tt.onGuard {
				b.GuardTurns = 1
			}
			b.takeDamage(tt.dmg)

			if b.Hp != tt.expectedHp {
				t.Errorf("expected %.1f hp, got %.1f", tt.expectedHp, b.Hp)
			}
		})
	}
}

func TestBaseHeal(t *testing.T) {
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
			b := base{HpCap: tt.hpCap}
			b.heal(tt.recover)

			if b.Hp != tt.want {
				t.Errorf("want %.1f, got %.1f", tt.want, b.Hp)
			}
		})
	}
}
