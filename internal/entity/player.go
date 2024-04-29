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
	TEMPORAL:   "⌛ Temporal (+1 extra turn for bonus effects)",
}

type Player struct {
	base

	FuryTurns       int
	Perk            int
	Money           float32
	HasFled         bool
	ExtraTurnEffect int
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
		p.Def++
		p.DmgReduc += 0.1
	case TEMPORAL:
		p.ExtraTurnEffect++
	}

	p.AddMoney(50)
	p.Heal(100)

	return p
}

func (p *Player) TakeAction(e Enemy, n int) (log string, ok bool) {
	defer func() {
		if ok {
			p.GuardTurns--
			p.FuryTurns--
		}
	}()

	switch n {
	case 1:
		_, s := p.Attack(e)
		return s, true
	case 2:
		if p.GuardTurns > 0 {
			return "Action already in effect\n", false
		}

		p.GuardTurns = 3 + p.ExtraTurnEffect
		return "You brace yourself 🛡️", true
	case 3:
		if p.Hp <= 10 {
			return "Not enough hp to perform this action\n", false
		}

		if p.FuryTurns < 0 {
			p.FuryTurns = 3 + p.ExtraTurnEffect
		} else {
			p.FuryTurns += 3 + p.ExtraTurnEffect
		}

		sacrifice := 1 + (p.Hp * 0.1) + (rand.Float32() * 4)
		p.Hp -= sacrifice
		return fmt.Sprintf("You descent into fury 🔥 (-%.1f hp)", sacrifice), true
	case 4:
		p.HasFled = true
		return "You decided to fight another day 🏃", true
	}

	return "Invalid action\n", false
}

func (p *Player) Attack(e Enemy) (float32, string) {
	dmg := p.base.attack()

	if p.FuryTurns > 0 {
		dmg += 5
	}

	if p.Perk == HAVOC {
		dmg += dmg * 0.25
	}

	dmg = e.TakeDamage(p, dmg)
	return dmg, fmt.Sprintf("You attacked (%.1f dmg)", dmg)
}

func (p *Player) TakeDamage(dmg float32) float32 { return p.base.takeDamage(dmg) }

func (p *Player) Heal(n float32) { p.base.heal(n) }

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
