package entity

import (
	"reflect"
	"testing"
)

func TestEnemyAttack(t *testing.T) {
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
			b := EnemyBase{Att: tt.att, isTesting: true}
			got, _ := b.Attack(&Player{})

			if tt.want != got {
				t.Errorf("want %.1f, got %.1f", tt.want, got)
			}
		})
	}

	t.Run("with random sum", func(t *testing.T) {
		b := EnemyBase{Att: 5}
		got, _ := b.Attack(&Player{})

		if b.Att > got {
			t.Errorf("want greater than 5, got %.1f", got)
		}
	})
}

func TestEnemyTakeDamage(t *testing.T) {
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
			b := EnemyBase{Hp: 100, Def: tt.def, DmgReduc: tt.dmgReduc, IsDefending: tt.isDefending}
			b.TakeDamage(nil, tt.dmg)

			if b.Hp != tt.expectedHp {
				t.Errorf("expected %.1f hp, got %.1f", tt.expectedHp, b.Hp)
			}
		})
	}
}

func TestDropLoot(t *testing.T) {
	tests := []struct {
		name     string
		dropRate float32
		want     float32
	}{
		{"with 1x drop rate", 1, 10},
		{"with 3.3x drop rate", 3.3, 33},
		{"with 0x drop rate", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := EnemyBase{DropRate: tt.dropRate, isTesting: true}
			got := b.DropLoot()

			if tt.want != got {
				t.Errorf("want %.1f, got %.1f", tt.want, got)
			}
		})
	}

	t.Run("with random value but 0 drop rate", func(t *testing.T) {
		b := EnemyBase{DropRate: 0}
		got := b.DropLoot()

		if got != 0 {
			t.Errorf("want 0, got %.1f", got)
		}
	})
}

func TestEnemyHeal(t *testing.T) {
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
			b := EnemyBase{HpCap: tt.hpCap}
			b.Heal(tt.recover)

			if b.Hp != tt.want {
				t.Errorf("want %.1f, got %.1f", tt.want, b.Hp)
			}
		})
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
	}

	for _, tt := range tests {
		got := reflect.TypeOf(tt.fn())
		want := reflect.TypeOf(tt.expectedType)
		if want != got {
			t.Errorf("incorrect type, want %q, got %q", want, got)
		}
	}
}

func TestVampireAttack(t *testing.T) {
	t.Run("should heal hp and deal extra dmg", func(t *testing.T) {
		p := Player{Hp: 100}
		v := vampire{EnemyBase{HpCap: 10, isTesting: true}}
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
		p := &Player{Def: 99999, Hp: 100}
		w := wraith{EnemyBase{Att: 10}}
		w.Attack(p)

		if p.Hp != 90 {
			t.Errorf("want player hp to be 90, got %.1f", p.Hp)
		}
	})
}
