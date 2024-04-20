package entity

import (
	"fmt"
	"math/rand"
)

var spawners = []func() Entity{
	newThug,
	newAcolyte,
	newAssasin,
	newSnakes,
	newGolem,
	newVampire,
	newWraith,
}

func SpawnRandom() Entity {
	return spawners[rand.Intn(len(spawners))]()
}

type thug struct{ Base }

func newThug() Entity {
	return &thug{Base{
		Name:     "Thug 🥊",
		Hp:       80,
		Att:      8,
		Def:      6,
		HpCap:    80,
		DmgReduc: 0.05,
		DropRate: 1,
	}}
}

type acolyte struct{ Base }

func newAcolyte() Entity {
	return &acolyte{Base{
		Name:     "Acolyte 🧙",
		Hp:       70,
		Att:      5,
		Def:      1,
		HpCap:    70,
		DmgReduc: 0.4,
		DropRate: 2,
	}}
}

type assasin struct{ Base }

func newAssasin() Entity {
	return &assasin{Base{
		Name:     "Assassin 🗡️",
		Hp:       60,
		Att:      15,
		Def:      3,
		HpCap:    60,
		DropRate: 1.3,
	}}
}

type snakes struct{ Base }

func newSnakes() Entity {
	return &snakes{Base{
		Name:     "Snakes 🐍",
		Hp:       30,
		Att:      13,
		Def:      0,
		HpCap:    30,
		DropRate: 0,
	}}
}

type golem struct{ Base }

func newGolem() Entity {
	return &golem{Base{
		Name:     "Golem 🗿",
		Hp:       151,
		Att:      38,
		Def:      20,
		HpCap:    151,
		DropRate: 3.2,
	}}
}

func (g *golem) Attack(e Entity) (float32, string) {
	if rand.Intn(100) < 30 {
		return g.Base.Attack(e)
	}

	return 0, fmt.Sprintf("%s attacked but missed", g.Name)
}

type vampire struct{ Base }

func newVampire() Entity {
	return &vampire{Base{
		Name:     "Vampire 🧛",
		Hp:       90,
		Att:      10,
		Def:      1,
		HpCap:    90,
		DropRate: 3,
	}}
}

func (v *vampire) Attack(e Entity) (float32, string) {
	dmg, s := v.Base.Attack(e)

	ls := 0.1 + dmg*0.2
	v.RecoverHP(ls)

	return dmg, fmt.Sprintf("%s (+%.1f hp)", s, ls)
}

type wraith struct{ Base }

func newWraith() Entity {
	return &wraith{Base{
		Name:     "Wraith 👻",
		Hp:       78,
		Att:      10,
		Def:      1,
		HpCap:    78,
		DropRate: 1,
	}}
}

func (w *wraith) Attack(e Entity) (float32, string) {
	e.Attr().Hp -= w.Att
	return w.Att, fmt.Sprintf("%s absorbed %.1f of hp", w.Name, w.Att)
}
