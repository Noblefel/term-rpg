package entity

import (
	"fmt"
	"math/rand"
)

const (
	GREED int = iota + 1
	RESILIENCY
	HAVOC
)

var Perks = map[int]string{
	GREED:      "💰 Greed (Gain 15% more loot)",
	RESILIENCY: "🛡️  Resiliency (+1 Def point and 10% dmg reduction)",
	HAVOC:      "⚔️   Havoc (+25% Attack, but -15 HP cap)",
}

type Player struct {
	Base
	Money float32
	Perk  int
}

func NewPlayer(perk int) *Player {
	p := &Player{Perk: perk}
	p.HpCap = 100
	p.Att = 10
	p.Def = 4

	if perk == HAVOC {
		p.HpCap = 85
	}

	if perk == RESILIENCY {
		p.Def += 1
		p.DmgReduc += 0.1
	}

	p.AddMoney(50)
	p.RecoverHP(100)

	return p
}

// Attack would further modify the value from the base Attack
func (p *Player) Attack() float32 {
	dmg := p.Base.Attack()

	if p.Perk == HAVOC {
		dmg += dmg * 0.25
	}

	return dmg
}

func (p *Player) RecoverHP(min float32) {
	health := p.Hp + min

	if health > p.HpCap {
		health = p.HpCap
	}

	p.Hp = health
}

func (p *Player) AddMoney(n float32) float32 {
	if p.Perk == GREED {
		n += n * 0.15
	}

	p.Money += n
	return n
}

// Train buffs player's attribute with random value
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
	case 3:
		p.DmgReduc += 0.01
		return "Dmg reduction increased by 1%"
	}

	return ""
}
