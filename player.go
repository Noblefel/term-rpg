package main

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	attributes
	gold int
	perk int
}

func NewPlayer(perk int) *Player {
	var player Player
	player.hp = 100
	player.hpcap = 100
	player.defense = 1
	player.strength = 20
	player.gold = 30
	player.perk = perk

	if perk == 0 {
		player.hp += 5
		player.hpcap += 5
		player.defense += 2.5
	}

	if perk == 1 {
		player.hp -= 25
		player.hpcap -= 25
		player.gold = 0
	}

	return &player
}

func (p *Player) attack(enemy entity) {
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

	fmt.Fprintf(out, success+"You attacked!")
	enemy.damage(dmg)
}

func (p *Player) damage(dmg float32) {
	if p.perk == 0 {
		dmg -= dmg * 0.1
	}

	if p.perk == 2 {
		percent := p.hp / (p.hpcap * 0.4) //if 0% hp   = 40% reduction
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.4 - percent*0.4
		dmg -= dmg * mul
	}

	p.attributes.damage(dmg)
}

func (p *Player) rest() {
	n := 15 + rand.Float32()*15
	p.hp = min(p.hpcap, n+p.hp)
	p.gold -= 5

	fmt.Fprint(out, success)
	fmt.Fprintf(out, "Recovered \033[38;5;83m%.1f\033[0m hp\n", n)
}

func (p *Player) train() {
	p.gold -= 10

	if rand.IntN(10) < 6 {
		fmt.Fprint(out, fail)

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

	fmt.Fprint(out, success)

	switch rand.IntN(3) {
	case 0:
		n := 1 + rand.Float32()*5
		p.hpcap += n
		fmt.Fprintf(out, "HP cap increased by \033[38;5;83m%.1f\033[0m\n", n)
	case 1:
		n := 0.1 + rand.Float32()*2
		p.strength += n
		fmt.Fprintf(out, "Strength increased by \033[38;5;83m%.1f\033[0m\n", n)
	case 2:
		n := 0.1 + rand.Float32()*2
		p.defense += n
		fmt.Fprintf(out, "Defense increased by \033[38;5;83m%.1f\033[0m\n", n)
	}
}

func (p *Player) flee(enemy entity) bool {
	if rand.IntN(10) < 6 {
		fmt.Fprint(out, success)
		fmt.Fprintln(out, "You have fled the battle")
		return true
	}

	fmt.Fprint(out, fail)

	switch rand.IntN(5) {
	case 0:
		fmt.Fprintln(out, "Youre too slow and got caught")
		enemy.attack(p)
	case 1:
		fmt.Fprintf(out, "You slipped in the mud,")
		p.damage(2)
	case 2:
		fmt.Fprintf(out, "You fell into a ditch,")
		p.damage(6)
	case 3:
		fmt.Fprintln(out, "You run around in circle")
	case 4:
		dmg := p.hp * 0.05
		p.hp = max(p.hp-dmg, 0)
		fmt.Fprintf(out, "You walked into a trap, -5%% hp (%.1f)\n", dmg)
	}

	return false
}
