package combat

import "testing"

func TestAttack(t *testing.T) {
	tests := []struct {
		name string
		att  float32
		want float32
	}{
		{"with 5 attck stat", 5, 15},
		{"with 0 attack stat", 0, 10},
		{"with negative attack stat", -50, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Base{Att: tt.att, isTesting: true}
			got := b.Attack()

			if tt.want != got {
				t.Errorf("want %.1f, got %.1f", tt.want, got)
			}
		})
	}

	t.Run("with random sum", func(t *testing.T) {
		b := Base{Att: 5}
		got := b.Attack()
		var want float32

		if b.Att > got {
			t.Errorf("want greater than %.1f, got %.1f", want, got)
		}
	})
}

func TestTakeDamage(t *testing.T) {
	tests := []struct {
		name       string
		def        float32
		dmg        float32
		dmgReduc   float32
		expectedHp float32
	}{
		{"10 damage with 4 def", 4, 10, 0, 94},
		{"10 damage with 100 def", 100, 10, 0, 100},
		{"50 damage with 25 def", 25, 50, 0, 75},
		{"8.8 damage with 1.3 def", 1.3, 8.8, 0, 92.5},
		{"30 damage with 30% dmg reduction", 0, 30, 0.3, 79},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Base{Hp: 100, Def: tt.def, DmgReduc: tt.dmgReduc}
			b.TakeDamage(tt.dmg)

			if b.Hp != tt.expectedHp {
				t.Errorf("expected %.1f hp, got %.1f", tt.expectedHp, b.Hp)
			}
		})
	}
}
