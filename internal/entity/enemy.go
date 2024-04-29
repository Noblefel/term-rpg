package entity

import (
	"fmt"
	"math/rand"
)

type Enemy interface {
	TakeAction(self Enemy, p *Player, n int) string
	Attack(p *Player) (float32, string)
	TakeDamage(p *Player, dmg float32) float32
	Heal(n float32)
	Attr() EnemyBase
}

type EnemyBase struct {
	base
	Name string
}

func (e *EnemyBase) TakeAction(self Enemy, p *Player, n int) string {
	defer func() {
		e.GuardTurns--
	}()

	if e.GuardTurns <= 0 && n < 10 {
		e.GuardTurns = 3
		return fmt.Sprintf("%s braces themselves 🛡️", e.Name)
	}

	_, log := self.Attack(p)
	return log
}

func (e *EnemyBase) Attack(p *Player) (float32, string) {
	dmg := e.base.attack()
	dmg = p.TakeDamage(dmg)
	return dmg, fmt.Sprintf("%s attacked (%.1f dmg)", e.Name, dmg)
}

func (e *EnemyBase) TakeDamage(p *Player, dmg float32) float32 {
	return e.base.takeDamage(dmg)
}

func (e *EnemyBase) Heal(n float32) { e.base.heal(n) }

func (e EnemyBase) Attr() EnemyBase { return e }

var spawners = []func() Enemy{
	newAcolyte,
	newAssasin,
	newEvilGenie,
	newGolem,
	newSnakes,
	newSpikeTurtle,
	newThug,
	newVampire,
	newWraith,
}

func SpawnRandom() Enemy {
	return spawners[rand.Intn(len(spawners))]()
}

type acolyte struct{ EnemyBase }

func newAcolyte() Enemy {
	var e acolyte
	e.Name = "Acolyte 🧙"
	e.Hp = 70
	e.Att = 5
	e.Def = 1
	e.HpCap = 70
	e.DmgReduc = 0.4
	return &e
}

type assasin struct{ EnemyBase }

func newAssasin() Enemy {
	var e assasin
	e.Name = "Assassin 🗡️ "
	e.Hp = 60
	e.Att = 15
	e.Def = 3
	e.HpCap = 60
	return &e
}

type evilGenie struct{ EnemyBase }

func newEvilGenie() Enemy {
	var e evilGenie
	e.Name = "Evil Genie 🧞"
	e.Hp = 75
	e.Att = 6
	e.Def = 3
	e.HpCap = 75
	e.DmgReduc = 0.1
	return &e
}

func (eg *evilGenie) TakeAction(self Enemy, p *Player, n int) string {
	if n < 18 {
		return eg.Curse(p, rand.Intn(4))
	}

	return eg.EnemyBase.TakeAction(self, p, n)
}

func (eg *evilGenie) Curse(p *Player, n int) string {
	switch n {
	case 0:
		n := 0.1 + rand.Float32()*4
		p.HpCap -= n
		return fmt.Sprintf("%s casted a curse on your hp cap by %.1f", eg.Name, n)
	case 1:
		n := 0.1 + rand.Float32()*1.5
		p.Att -= n
		return fmt.Sprintf("%s casted a curse, weakening your attack by %.1f", eg.Name, n)
	case 2:
		n := 0.1 + rand.Float32()*1
		p.Def -= n
		return fmt.Sprintf("%s casted a curse, weakening your defense by %.1f", eg.Name, n)
	default:
		p.DmgReduc -= 0.01
		return eg.Name + " casts a curse on your dmg reduction by 1%%"
	}
}

type golem struct{ EnemyBase }

func newGolem() Enemy {
	var e golem
	e.Name = "Golem 🗿"
	e.Hp = 151
	e.Att = 38
	e.Def = 20
	e.HpCap = 151
	return &e
}

func (g *golem) Attack(p *Player) (float32, string) {
	if rand.Intn(100) < 30 {
		return g.EnemyBase.Attack(p)
	}

	return 0, fmt.Sprintf("%s attacked but missed", g.Name)
}

type snakes struct{ EnemyBase }

func newSnakes() Enemy {
	var e snakes
	e.Name = "Snakes 🐍"
	e.Hp = 30
	e.Att = 13
	e.Def = 0
	e.HpCap = 30
	return &e
}

type spikeTurtle struct{ EnemyBase }

func newSpikeTurtle() Enemy {
	var e spikeTurtle
	e.Name = "Spike Turtle 🐢"
	e.Hp = 120
	e.Att = 6
	e.Def = 10
	e.HpCap = 120
	return &e
}

func (st *spikeTurtle) TakeDamage(p *Player, dmg float32) float32 {
	reflect := dmg * (0.1 + rand.Float32()*0.2)
	p.TakeDamage(reflect)
	if p.Hp <= 0 {
		p.Hp = 0.1 // prevent negative hp
	}

	return st.EnemyBase.TakeDamage(nil, dmg)
}

type thug struct{ EnemyBase }

func newThug() Enemy {
	var e thug
	e.Name = "Thug 🥊"
	e.Hp = 80
	e.Att = 10
	e.Def = 6
	e.HpCap = 80
	e.DmgReduc = 0.05
	return &e
}

type vampire struct{ EnemyBase }

func newVampire() Enemy {
	var e vampire
	e.Name = "Vampire 🧛"
	e.Hp = 90
	e.Att = 7
	e.Def = 1
	e.HpCap = 90
	return &e
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
	var e wraith
	e.Name = "Wraith 👻"
	e.Hp = 78
	e.Att = 10
	e.Def = 1
	e.HpCap = 78
	return &e
}

func (w *wraith) Attack(p *Player) (float32, string) {
	p.Hp -= w.Att
	return w.Att, fmt.Sprintf("%s absorbed %.1f of hp", w.Name, w.Att)
}
