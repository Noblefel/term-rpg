package entity

import "math/rand"

func SpawnRandom() Entity {
	switch rand.Intn(4) {
	case 0:
		return newThug()
	case 1:
		return newAcolyte()
	case 2:
		return newAssasin()
	case 3:
		return newSnakes()
	}

	return nil
}

type thug struct{ Base }

func newThug() Entity {
	return &thug{Base{
		Name:     "Thug 🥊",
		Hp:       80,
		Att:      8,
		Def:      6,
		HpCap:    80,
		DmgReduc: 0.04,
		DropRate: 0.8,
	}}
}

type acolyte struct{ Base }

func newAcolyte() Entity {
	return &acolyte{Base{
		Name:     "Acolyte 🧙",
		Hp:       70,
		Att:      4,
		Def:      1,
		HpCap:    70,
		DmgReduc: 0.3,
		DropRate: 1,
	}}
}

type assasin struct{ Base }

func newAssasin() Entity {
	return &assasin{Base{
		Name:     "Assassin 🗡️",
		Hp:       60,
		Att:      14,
		Def:      3,
		HpCap:    60,
		DropRate: 1.2,
	}}
}

type snakes struct{ Base }

func newSnakes() Entity {
	return &snakes{Base{
		Name:     "Snakes 🐍",
		Hp:       30,
		Att:      12.5,
		Def:      1,
		HpCap:    30,
		DropRate: 0,
	}}
}
