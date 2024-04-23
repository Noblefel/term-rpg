package entity

import (
	"fmt"
	"math/rand"
)

type Enemy interface {
	// Attack sums the attack with random value. In tests, replaced random value with 10.
	Attack(p *Player) (float32, string)
	// TakeDamage reduces dmg with the defense and decrements the hp
	TakeDamage(p *Player, dmg float32) float32
	// DropLoot multiplies random value with drop rate. In tests, replaced random value with 10.
	DropLoot() float32
	// Heal adds hp to the entity, cannot exceeds hp cap
	Heal(n float32)

	Attr() *EnemyBase
}

type EnemyBase struct {
	Name     string
	Hp       float32
	Att      float32
	Def      float32
	HpCap    float32
	DmgReduc float32
	DropRate float32

	IsDefending bool
	isTesting   bool
}

func (e *EnemyBase) Attack(p *Player) (float32, string) {
	dmg := e.Att

	if e.isTesting {
		dmg += 10
	} else {
		dmg += rand.Float32() * 10
	}

	dmg = p.TakeDamage(dmg)

	return dmg, fmt.Sprintf("%s attacked (%.1f dmg)", e.Name, dmg)
}

func (e *EnemyBase) TakeDamage(p *Player, dmg float32) float32 {
	dmg -= e.Def + (dmg * e.DmgReduc)

	if e.IsDefending {
		dmg -= dmg * 0.2
	}

	if dmg < 0 {
		return 0
	}

	e.Hp -= dmg
	return dmg
}

func (e *EnemyBase) DropLoot() float32 {
	var loot float32

	if e.isTesting {
		loot = 10
	} else {
		loot += 5 + rand.Float32()*30
	}

	loot *= e.DropRate

	if loot <= 0 {
		return 0
	}

	return loot
}

func (e *EnemyBase) Heal(n float32) {
	hp := e.Hp + n

	if hp > e.HpCap {
		hp = e.HpCap
	}

	e.Hp = hp
}

func (e *EnemyBase) Attr() *EnemyBase { return e }

var spawners = []func() Enemy{
	newAcolyte,
	newAssasin,
	newEvilGenie,
	newGolem,
	newSnakes,
	newThug,
	newVampire,
	newWraith,
}

func SpawnRandom() Enemy {
	return spawners[rand.Intn(len(spawners))]()
}

type acolyte struct{ EnemyBase }

func newAcolyte() Enemy {
	return &acolyte{EnemyBase{
		Name:     "Acolyte 🧙",
		Hp:       70,
		Att:      5,
		Def:      1,
		HpCap:    70,
		DmgReduc: 0.4,
		DropRate: 2,
	}}
}

type assasin struct{ EnemyBase }

func newAssasin() Enemy {
	return &assasin{EnemyBase{
		Name:     "Assassin 🗡️",
		Hp:       60,
		Att:      15,
		Def:      3,
		HpCap:    60,
		DropRate: 1.3,
	}}
}

type evilGenie struct{ EnemyBase }

func newEvilGenie() Enemy {
	return &evilGenie{EnemyBase{
		Name:     "Evil Genie 🧞",
		Hp:       75,
		Att:      6,
		Def:      3,
		HpCap:    75,
		DropRate: 1,
		DmgReduc: 0.1,
	}}
}

func (eg *evilGenie) Attack(p *Player) (float32, string) {
	if rand.Intn(100) < 18 {
		switch rand.Intn(4) {
		case 0:
			n := 0.1 + rand.Float32()*4
			p.HpCap -= n
			return 0, fmt.Sprintf("%s casted a curse on your hp cap by %.1f", eg.Name, n)
		case 1:
			n := 0.1 + rand.Float32()*1.5
			p.Att -= n
			return 0, fmt.Sprintf("%s casted a curse, weakening your attack by %.1f", eg.Name, n)
		case 2:
			n := 0.1 + rand.Float32()*1
			p.Def -= n
			return 0, fmt.Sprintf("%s casted a curse, weakening your defense by %.1f", eg.Name, n)
		case 3:
			p.DmgReduc -= 0.01
			return 0, eg.Name + " casts a curse on your dmg reduction by 1%%"
		}
	}

	return eg.EnemyBase.Attack(p)
}

func newGolem() Enemy {
	return &golem{EnemyBase{
		Name:     "Golem 🗿",
		Hp:       151,
		Att:      38,
		Def:      20,
		HpCap:    151,
		DropRate: 3.2,
	}}
}

func (g *golem) Attack(p *Player) (float32, string) {
	if rand.Intn(100) < 30 {
		return g.EnemyBase.Attack(p)
	}

	return 0, fmt.Sprintf("%s attacked but missed", g.Name)
}

type snakes struct{ EnemyBase }

func newSnakes() Enemy {
	return &snakes{EnemyBase{
		Name:     "Snakes 🐍",
		Hp:       30,
		Att:      13,
		Def:      0,
		HpCap:    30,
		DropRate: 0,
	}}
}

type golem struct{ EnemyBase }

type thug struct{ EnemyBase }

func newThug() Enemy {
	return &thug{EnemyBase{
		Name:     "Thug 🥊",
		Hp:       80,
		Att:      8,
		Def:      6,
		HpCap:    80,
		DmgReduc: 0.05,
		DropRate: 1,
	}}
}

type vampire struct{ EnemyBase }

func newVampire() Enemy {
	return &vampire{EnemyBase{
		Name:     "Vampire 🧛",
		Hp:       90,
		Att:      7,
		Def:      1,
		HpCap:    90,
		DropRate: 3,
	}}
}

func (v *vampire) Attack(p *Player) (float32, string) {
	drain := p.Hp * 0.05
	dmg, _ := v.EnemyBase.Attack(p)
	dmg += p.TakeDamage(drain)

	heal := 0.1 + (dmg)*0.2
	v.Heal(heal)

	return dmg, fmt.Sprintf("%s drained %.1f hp and heals by %.1f", v.Name, dmg, heal)
}

type wraith struct{ EnemyBase }

func newWraith() Enemy {
	return &wraith{EnemyBase{
		Name:     "Wraith 👻",
		Hp:       78,
		Att:      10,
		Def:      1,
		HpCap:    78,
		DropRate: 1,
	}}
}

func (w *wraith) Attack(p *Player) (float32, string) {
	p.Hp -= w.Att
	return w.Att, fmt.Sprintf("%s absorbed %.1f of hp", w.Name, w.Att)
}
