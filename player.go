package main

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	hp       float32
	hpcap    float32
	strength float32
	defense  float32
	gold     int
	perk     int
}

func NewPlayer(perk int) *Player {
	player := Player{
		hp:       100,
		hpcap:    100,
		defense:  1,
		strength: 20,
		gold:     50,
		perk:     perk,
	}

	if perk == 0 {
		player.hp += 5
		player.hpcap += 5
		player.defense += 2.5
	}

	if perk == 1 {
		player.hp -= 25
		player.hpcap -= 25
		player.gold -= 40
	}

	return &player
}

func (p *Player) Attack(e Enemy) {
	dmg := p.strength + rand.Float32()*5

	if p.perk == 1 {
		dmg += dmg * 0.2
	}

	if p.perk == 2 {
		percent := p.hp / (p.hpcap * 0.4) //if 0% hp   = 40% extra dmg
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.4 - percent*0.4
		dmg += dmg * mul
	}

	fmt.Fprint(out, "\033[38;5;83m✔\033[0m You attacked!")
	e.TakeDamage(dmg)
}

func (p *Player) TakeDamage(dmg float32) {
	if p.perk == 0 {
		dmg -= dmg * 0.1
	}

	if p.perk == 2 {
		percent := p.hp / (p.hpcap * 0.4) //if 0% hp   = 40% reduction
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.4 - percent*0.4
		dmg -= dmg * mul
	}

	dmg = max(dmg-p.defense, 1)
	p.hp = max(p.hp-dmg, 0)
	fmt.Fprintf(out, " \033[38;5;198m%.1f\033[0m damage\n", dmg)
}

func (p *Player) Rest() {
	n := 15 + rand.Float32()*15
	p.hp = min(p.hpcap, n+p.hp)
	p.gold -= 5

	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")
	fmt.Fprintf(out, "Recovered \033[38;5;83m%.1f\033[0m hp\n", n)
}

func (p *Player) Train() {
	p.gold -= 10

	if rand.IntN(10) < 6 {
		fmt.Fprint(out, "\033[38;5;196m✘\033[0m ")

		fails := []string{
			"You messed up",
			"You feel nothing",
			"You get distracted",
			"You didnt do anything",
			"You only get exhausted",
			"You just stare at the wall",
		}

		fmt.Fprintln(out, fails[rand.IntN(len(fails))])
		return
	}

	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")

	switch rand.IntN(3) {
	case 0:
		n := 1 + rand.Float32()*5
		p.hpcap += n
		fmt.Fprintf(out, "HP cap increased by \033[38;5;83m%.1f\033[0m\n", n)
	case 1:
		n := 0.5 + rand.Float32()*2
		p.strength += n
		fmt.Fprintf(out, "Strength increased by \033[38;5;83m%.1f\033[0m\n", n)
	case 2:
		n := 0.5 + rand.Float32()*2
		p.defense += n
		fmt.Fprintf(out, "Defense increased by \033[38;5;83m%.1f\033[0m\n", n)
	}
}

func (p *Player) Flee(enemy Enemy) bool {
	if rand.IntN(10) < 6 {
		fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")
		fmt.Fprintln(out, "You have fled the battle")
		return true
	}

	fmt.Fprint(out, "\033[38;5;196m✘\033[0m ")

	switch rand.IntN(5) {
	case 0:
		fmt.Fprintln(out, "Youre too slow and got caught")
		enemy.Attack()
	case 1:
		fmt.Fprint(out, "You slipped in the mud,")
		p.TakeDamage(2)
	case 2:
		fmt.Fprint(out, "You fell into a ditch,")
		p.TakeDamage(6)
	case 3:
		fmt.Fprintln(out, "You run around in circle")
	case 4:
		dmg := p.hp * 0.05
		p.hp = max(p.hp-dmg, 0)
		fmt.Fprintf(out, "You walked into a trap, -5%% hp (%.1f)\n", dmg)
	}

	return false
}
