package entity

import (
	"fmt"
	"math/rand"
)

const (
	GREED int = iota + 1
	RESILIENCY
	HAVOC
	TEMPORAL
)

var Perks = map[int]string{
	GREED:      "💰 Greed (Gain 15% more loot)",
	RESILIENCY: "🛡️  Resiliency (+1 Def point and 10% dmg reduction)",
	HAVOC:      "⚔️   Havoc (+25% Attack, but -15 HP cap)",
	TEMPORAL:   "⏰  Temporal (+8 seconds to actions bonus modifier)",
}

type Player struct {
	Hp       float32
	Att      float32
	Def      float32
	HpCap    float32
	DmgReduc float32
	Money    float32
	Perk     int

	IsDefending bool
	isTesting   bool
}

func NewPlayer(perk int) *Player {
	p := &Player{Perk: perk}
	p.HpCap = 100
	p.Att = 10
	p.Def = 4

	switch perk {
	case HAVOC:
		p.HpCap = 85
	case RESILIENCY:
		p.Def += 1
		p.DmgReduc += 0.1
	}

	p.AddMoney(50)
	p.Heal(100)

	return p
}

func (p *Player) Attack(e Enemy) (float32, string) {
	dmg := p.Att

	if p.isTesting {
		dmg += 10
	} else {
		dmg += rand.Float32() * 10
	}

	if p.Perk == HAVOC {
		dmg += dmg * 0.25
	}

	dmg = e.TakeDamage(p, dmg)
	return dmg, fmt.Sprintf("You attacked (%.1f dmg)", dmg)
}

func (p *Player) TakeDamage(dmg float32) float32 {
	dmg -= p.Def + (dmg * p.DmgReduc)

	if p.IsDefending {
		dmg -= dmg * 0.2
	}

	if dmg < 0 {
		return 0
	}

	p.Hp -= dmg
	return dmg
}

func (p *Player) Heal(n float32) {
	hp := p.Hp + n

	if hp > p.HpCap {
		hp = p.HpCap
	}

	p.Hp = hp
}

func (p *Player) AddMoney(n float32) float32 {
	if p.Perk == GREED {
		n += n * 0.15
	}

	p.Money += n
	return n
}

func (p *Player) Train(n int) string {
	switch n {
	case 0:
		n := 1 + rand.Float32()*5
		p.HpCap += n
		return fmt.Sprintf("Hp cap increased by %.1f", n)
	case 1:
		n := 0.5 + rand.Float32()*2
		p.Att += n
		return fmt.Sprintf("Attack increased by %.1f", n)
	case 2:
		n := 0.5 + rand.Float32()*2
		p.Def += n
		return fmt.Sprintf("Defense increased by %.1f", n)
	default:
		p.DmgReduc += 0.01
		return "Dmg reduction increased by 1%%"
	}
}
