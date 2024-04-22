package entity

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		fn           func() Entity
		expectedType Entity
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
		e := Base{Hp: 100}
		v := vampire{Base{HpCap: 10, isTesting: true}}
		v.Attack(&e)

		if v.Hp != 3.1 {
			t.Errorf("want recover 3.1 hp, got %.1f", v.Hp)
		}

		if e.Hp != 85 {
			t.Errorf("want enemy hp to be 85, got %.1f", e.Hp)
		}
	})
}

func TestWraithAttack(t *testing.T) {
	t.Run("always deal fixed damage", func(t *testing.T) {
		e := &Base{Def: 99999, Hp: 100}
		w := wraith{Base{Att: 10}}
		w.Attack(e)

		if e.Hp != 90 {
			t.Errorf("want hp to be 90, got %.1f", e.Hp)
		}
	})
}
