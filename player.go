package main

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	attributes
	gold      int
	perk      int
	energy    int
	energycap int
}

var skills = []struct {
	name string
	desc string
	cost int
	cd   int
}{
	{"charge", "attack 130% strength", 4, 5},
	{"heal", "recover hp by atleast 8% of hpcap", 5, 5},
	{"frenzy", "sacrifice hp to attack with 250% strength", 6, 6},
	{"vision", "see enemy attributes", 0, 0},
	{"drain", "take 20% of enemy current hp", 4, 4},
	{"absorb", "take 8% of enemy hp cap and ignore defense", 5, 6},
	{"trick", "make the enemy target themselves", 4, 3},
	{"poison", "attack 60% strength and poison enemy for 3 turns", 5, 5},
	{"stun", "attack 30% strength and stun enemy for 2 turns", 6, 4},
	{"fireball", "deal moderate amount of damage", 7, 5},
	{"meteor strike", "deal huge amount of damage", 10, 5},
}

func NewPlayer(perk int) *Player {
	var player Player
	player.hp = 100
	player.hpcap = 100
	player.defense = 1
	player.strength = 20
	player.gold = 30
	player.perk = perk
	player.energy = 20
	player.energycap = 20
	player.effects = make(map[string]int)

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
	dmg := p.strength

	if p.perk == 1 {
		dmg += dmg * 0.2
	}

	if p.perk == 2 {
		percent := p.hp / (p.hpcap * 0.4) //if 0% hp   = 40% extra dmg
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.4 - percent*0.4
		dmg += dmg * mul
	}

	fmt.Printf(success + "You attacked!")
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

func (p *Player) skill(i int, enemy entity) bool {
	skill := skills[i]

	if skill.cost > p.energy {
		fmt.Println("\033[38;5;196mNot enough energy\033[0m")
		return false
	}

	if p.effects["cd_"+skill.name] > 0 {
		fmt.Println("\033[38;5;196mSkill in cooldown\033[0m")
		return false
	}

	fmt.Print(success)
	fmt.Printf("You use \033[38;5;226m%s\033[0m: ", skill.name)

	switch skill.name {
	case "charge":
		temp := p.strength
		p.strength *= 1.3
		p.attack(enemy)
		p.strength = temp
	case "heal":
		heal := 10 + p.hpcap*0.08
		p.hp = min(p.hp+heal, p.hpcap)
		fmt.Printf("recover \033[38;5;83m%.1f\033[0m hp\n", heal)
	case "frenzy":
		sacrifice := 0.20 * p.hp
		sacrifice += 0.05 * p.hpcap
		p.hp = max(p.hp-sacrifice, 0)
		fmt.Printf("\033[38;5;198m-%.1f\033[0m hp and deal", sacrifice)
		enemy.damage(p.strength * 2.5)
	case "vision":
		fmt.Println("you can see they have")
		fmt.Printf("hp cap   : %.1f\n", enemy.attr().hpcap)
		fmt.Printf("strength : %.1f\n", enemy.attr().strength)
		fmt.Printf("defense  : %.1f\n", enemy.attr().defense)
	case "drain":
		drain := enemy.attr().hp * 0.2
		enemy.damage(drain)
	case "absorb":
		absorb := enemy.attr().hpcap * 0.075
		newhp := max(enemy.attr().hp-absorb, 0)
		enemy.setHP(newhp)
		fmt.Printf("take \033[38;5;198m%.1f\033[0m enemy hp\n", absorb)
	case "trick":
		enemy.attack(enemy)
	case "poison":
		temp := p.strength
		p.strength *= 0.6
		p.attack(enemy)
		p.strength = temp
		enemy.attr().effects["poisoned"] += 3
	case "stun":
		temp := p.strength
		p.strength *= 0.3
		p.attack(enemy)
		p.strength = temp
		enemy.attr().effects["stunned"] += 2
	case "fireball":
		dmg := 20 + rand.Float32()*15
		enemy.damage(dmg)
	case "meteor strike":
		dmg := 20 + rand.Float32()*60
		enemy.damage(dmg)
	}

	cooldown := skill.cd

	if p.perk == 3 {
		cooldown--
	}

	if p.perk == 2 && p.hp/p.hpcap <= 0.15 {
		cooldown--
	}

	p.effects["cd_"+skill.name] = cooldown
	p.energy -= skill.cost
	return true
}

func (p *Player) rest() {
	heal := p.hpcap * 0.02
	heal += 15 + rand.Float32()*15
	p.hp = min(p.hp+heal, p.hpcap)
	p.energy = min(p.energy+5, p.energycap)
	p.gold -= 5

	fmt.Print(success)
	fmt.Printf("Recovered \033[38;5;83m%.1f\033[0m hp", heal)
	fmt.Printf(" and \033[38;5;83m5\033[0m energy\n")
}

func (p *Player) train() {
	p.gold -= 10
	roll := roll()

	if roll < 60 {
		fmt.Print(fail)

		fails := []string{
			"You messed up",
			"You feel nothing",
			"You get distracted",
			"You didnt do anything",
			"You only get exhausted",
			"You just stare at the wall",
		}

		fmt.Println(fails[rand.IntN(len(fails))])
		return
	}

	fmt.Print(success)

	if roll < 71 {
		n := 1 + rand.Float32()*5
		p.hpcap += n
		fmt.Printf("HP cap increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 82 {
		n := 0.1 + rand.Float32()*2
		p.strength += n
		fmt.Printf("Strength increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 93 {
		n := 0.1 + rand.Float32()*2
		p.defense += n
		fmt.Printf("Defense increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else {
		p.energycap++
		fmt.Println("Energy cap increased by \033[38;5;83m1\033[0m")
	}
}

func (p *Player) flee(enemy entity) {
	roll := roll()

	if roll < 60 {
		fmt.Print(success)
		fmt.Println("You have fled the battle")
		p.effects["fled"] = 1
		return
	}

	fmt.Print(fail)

	if roll < 68 {
		fmt.Println("Youre too slow and got caught")

		if enemy.attr().effects["stunned"] > 0 {
			fmt.Print(fail)
			fmt.Printf("%s tried to attack but is stunned\n", enemy.attr().name)
		} else {
			enemy.attack(p)
		}
	} else if roll < 76 {
		fmt.Printf("You slipped in the mud,")
		p.damage(2)
	} else if roll < 84 {
		fmt.Printf("You fell into a ditch,")
		p.damage(6)
	} else if roll < 92 {
		dmg := p.hp * 0.05
		p.hp = max(p.hp-dmg, 0)
		fmt.Printf("You walked into a trap, \033[38;5;198m%.1f\033[0m dmg\n", dmg)
	} else {
		fmt.Println("You run around in circle")
	}
}

func (p Player) energybar() string {
	bar := bars(40, float32(p.energy), float32(p.energycap))
	return fmt.Sprintf("\033[38;5;226m" + bar[0] + "\033[0m" + bar[1])
}
