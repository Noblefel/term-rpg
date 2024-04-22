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
	newEvilGenie,
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
		Att:      7,
		Def:      1,
		HpCap:    90,
		DropRate: 3,
	}}
}

func (v *vampire) Attack(e Entity) (float32, string) {
	extra := e.Attr().Hp * 0.05
	dmg, _ := v.Base.Attack(e)
	dmg += extra
	e.TakeDamage(extra)

	heal := 0.1 + (dmg)*0.2
	v.RecoverHP(heal)

	return dmg, fmt.Sprintf("%s drained %.1f hp and heals by %.1f", v.Name, dmg, heal)
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

type evilGenie struct{ Base }

func newEvilGenie() Entity {
	return &evilGenie{Base{
		Name:     "Evil Genie 🧞",
		Hp:       75,
		Att:      6,
		Def:      3,
		HpCap:    75,
		DropRate: 1,
		DmgReduc: 0.1,
	}}
}

func (eg *evilGenie) Attack(e Entity) (float32, string) {
	if rand.Intn(100) < 18 {
		switch rand.Intn(4) {
		case 0:
			n := 0.1 + rand.Float32()*4
			e.Attr().HpCap -= n
			return 0, fmt.Sprintf("%s casted a curse on your hp cap by %.1f", eg.Name, n)
		case 1:
			n := 0.1 + rand.Float32()*1.5
			e.Attr().Att -= n
			return 0, fmt.Sprintf("%s casted a curse, weakening your attack by %.1f", eg.Name, n)
		case 2:
			n := 0.1 + rand.Float32()*1
			e.Attr().Def -= n
			return 0, fmt.Sprintf("%s casted a curse, weakening your defense by %.1f", eg.Name, n)
		case 3:
			e.Attr().DmgReduc -= 0.01
			return 0, eg.Name + " casts a curse on your dmg reduction by 1%%"
		}
	}

	return eg.Base.Attack(e)
}
