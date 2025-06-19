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

	skills [5]int //indexes
}

var skills = []struct {
	name string
	desc string
	cost int
	cd   int
}{
	{"charge", "attack with 130% strength", 4, 2},
	{"frenzy", "sacrifice hp to attack with 250% strength", 6, 6},
	{"great blow", "sacrifice the next turn to attack with 210% strength", 5, 3},
	{"poison", "attack 85% strength and poison enemy for 3 turns", 5, 5},
	{"stun", "attack 60% strength and stun enemy for 2 turns", 6, 4},
	{"swift strike", "attack 85% strength (doesnt consume turn)", 4, 4},
	{"knives throw", "attack 40 fixed damage (doesnt consume turn, no cd)", 4, 0},
	{"fireball", "deal moderate amount of damage and inflict burning", 6, 5},
	{"meteor strike", "deal huge amount of damage", 9, 5},
	{"strengthen", "attack 100% strength to increase damage by 10% for 3 turns", 4, 7},
	{"barrier", "reduce incoming damage by 40% for 2 turns", 4, 3},
	{"force-field", "reduce incoming damage by 15% for 5 turns", 5, 8},
	{"heal spell", "recover hp by atleast 15% hp cap", 8, 5},
	{"heal aura", "recover hp by atleast 7% of hpcap for 3 turns", 8, 6},
	{"heal potion", "recover hp by 40 (fixed number)", 7, 10},
	{"drain", "take 22% of enemy current hp", 4, 4},
	{"absorb", "take 10% of enemy hp cap, ignore defense and effects", 5, 6},
	{"trick", "make the enemy self-target", 4, 3},
	{"vision", "see enemy attributes (no cost, no cd, doesnt consume turn)", 0, 0},
}

func NewPlayer() *Player {
	var player Player
	player.name = "player"
	player.gold = 30
	player.perk = -1

	player.hp = 250
	player.hpcap = 250
	player.strength = 50
	player.defense = 15
	player.agility = 5
	player.energy = 20
	player.energycap = 20

	// charge, stun, poison, fireball, heal potion
	player.skills = [5]int{0, 3, 4, 7, 14}
	player.effects = make(map[string]int)

	return &player
}

func (p *Player) attack(enemy entity) {
	fmt.Printf(success + "you attacked!")
	p.attackWith(enemy, p.strength)
}

// player perks modifier applied here
func (p *Player) attackWith(enemy entity, dmg float32) {
	if p.energy <= 3 {
		fmt.Print(" exhausted")
		dmg -= dmg * 0.1
	}

	if p.perk == 1 {
		dmg += dmg * 0.2
	}

	if p.perk == 2 {
		percent := p.hp / (p.hpcap * 0.4) //if 1% hp   = 60% extra dmg
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.6 - percent*0.6
		dmg += dmg * mul
	}

	if _, ok := enemy.(*undead); ok && p.perk == 5 {
		dmg += dmg * 0.33
	}

	if p.perk == 7 {
		roll := roll()

		if roll < 25 {
			// multiplier - 30% increase/decrease range
			mul := 0.3 - rand.Float32()*0.6
			dmg += dmg * mul
		} else if roll < 50 {
			// flat val (scaled)
			val := scale(1, 1)
			dmg += val/2 - rand.Float32()*val
		} else if roll < 75 {
			// flat val
			dmg += 10 - 0*20
		}
	}

	p.attributes.attackWith(enemy, dmg)
}

// player perks modifier applied here
func (p *Player) damage(dmg float32) {
	if p.perk == 0 {
		dmg -= dmg * 0.1
		dmg = min(dmg, p.hpcap*0.16)
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
	cost := skill.cost

	if p.effects["confused"] > 0 {
		cost++
	}

	if cost > p.energy {
		fmt.Println("\033[38;5;196mNot enough energy\033[0m")
		return false
	}

	if p.effects["cd"+skill.name] > 0 {
		fmt.Println("\033[38;5;196mSkill in cooldown\033[0m")
		return false
	}

	cooldown := skill.cd

	if p.perk == 3 {
		cooldown -= 2
	}

	if p.perk == 2 && p.hp/p.hpcap <= 0.25 {
		cooldown--
	}

	if p.perk == 7 {
		cooldown = rand.IntN(7)
	}

	p.effects["cd"+skill.name] = cooldown
	p.energy -= cost

	fmt.Print(success)
	fmt.Printf("You use \033[38;5;226m%s\033[0m!", skill.name)

	switch skill.name {
	case "charge":
		p.attackWith(enemy, p.strength*1.3)
	case "frenzy":
		sacrifice := 0.15 * p.hp
		sacrifice += 0.05 * p.hpcap
		p.hp = max(p.hp-sacrifice, 0)
		fmt.Printf(" \033[38;5;198m-%.1f\033[0m hp and deal", sacrifice)
		p.attackWith(enemy, p.strength*2.5)
	case "great blow":
		p.effects["stunned"] = 2
		p.attackWith(enemy, p.strength*2.1)
	case "poison":
		enemy.attr().effects["poisoned"] = 3
		p.attackWith(enemy, p.strength*0.85)
	case "stun":
		enemy.attr().effects["stunned"] = 2
		p.attackWith(enemy, p.strength*0.6)
	case "swift strike":
		p.attackWith(enemy, p.strength*0.85)
		return false
	case "knives throw":
		p.attackWith(enemy, 40)
		return false
	case "fireball":
		enemy.attr().effects["burning"] = 2
		p.attackWith(enemy, 80)
	case "meteor strike":
		dmg := 50 + rand.Float32()*170
		p.attackWith(enemy, dmg)
	case "strengthen":
		p.attackWith(enemy, p.strength) // attack first so it wont get the bonus yet
		p.effects["strengthen"] = 4     // +1 for 3 attacks
	case "barrier":
		p.effects["barrier"] = 2
		fmt.Println(" reducing damage for 2 turn")
	case "force-field":
		p.effects["force-field"] = 5
		fmt.Println(" reducing damage for 5 turn")
	case "heal spell":
		heal := p.hpcap * 0.15
		p.hp = min(p.hp+heal, p.hpcap)
		fmt.Printf(" recover \033[38;5;83m%.1f\033[0m hp\n", heal)
	case "heal aura":
		p.effects["heal aura"] = 4 // +1 because it start in next turn
		fmt.Println(" recover \033[38;5;83m7%\033[0m hp for 3 turns")
	case "heal potion":
		p.hp = min(p.hp+40, p.hpcap)
		fmt.Println(" recover \033[38;5;83m40\033[0m hp")
	case "drain":
		drain := enemy.attr().hp * 0.22
		enemy.damage(drain)
	case "absorb":
		absorb := enemy.attr().hpcap * 0.1
		newhp := max(enemy.attr().hp-absorb, 0)
		enemy.setHP(newhp)
		fmt.Printf(" take \033[38;5;198m%.1f\033[0m enemy hp\n", absorb)
	case "vision":
		fmt.Println(" you can see they have")
		fmt.Printf("  hp cap: %.1f |", enemy.attr().hpcap)
		fmt.Printf(" strength: %.1f |", enemy.attr().strength)
		fmt.Printf(" defense: %.1f |", enemy.attr().defense)
		fmt.Printf(" agility: %.1f\n", enemy.attr().agility)
		return false
	case "trick":
		fmt.Print("\n  self: ")
		enemy.attack(enemy)
	}

	return true
}

func (p *Player) rest() {
	heal := p.hpcap * 0.02
	heal += 20 + rand.Float32()*20
	p.hp = min(p.hp+heal, p.hpcap)
	p.energy = min(p.energy+5, p.energycap)

	fmt.Print(success)
	fmt.Printf("recovered \033[38;5;83m%.1f\033[0m hp", heal)
	fmt.Printf(" and \033[38;5;83m5\033[0m energy\n")
}

func (p *Player) train() {
	roll := roll()

	if roll < 51 {
		fmt.Print(fail)

		fails := []string{
			"you messed up",
			"you feel nothing",
			"you get distracted",
			"you didnt do anything",
			"you only get exhausted",
			"you just stare at the wall",
		}

		fmt.Println(fails[rand.IntN(len(fails))])
		return
	}

	fmt.Print(success)

	if roll < 62 {
		n := 2.5 + rand.Float32()*5
		p.hpcap += n
		fmt.Printf("hp cap increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 73 {
		n := 0.5 + rand.Float32()*2
		p.strength += n
		fmt.Printf("strength increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 84 {
		n := 0.5 + rand.Float32()*2
		p.defense += n
		fmt.Printf("defense increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 95 {
		n := 0.1 + rand.Float32()*1
		p.agility += n
		fmt.Printf("agility increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else {
		p.energycap++
		fmt.Println("energy cap increased by \033[38;5;83m1\033[0m")
	}
}

func (p *Player) flee(enemy entity) {
	roll := roll()
	flee := float32(roll) - p.agility*0.5 + enemy.attr().agility*0.2

	if flee < 20 || (flee < 85 && p.perk == 6) {
		fmt.Print(success)
		fmt.Println("you have fled the battle")
		p.effects["fled"] = 1
		return
	}

	fmt.Print(fail)

	if roll < 68 {
		fmt.Println("youre too slow and got caught")

		if enemy.attr().effects["stunned"] > 0 {
			fmt.Print(fail)
			fmt.Printf("%s tried to attack but is stunned\n", enemy.attr().name)
		} else {
			enemy.attack(p)
		}
	} else if roll < 76 {
		fmt.Printf("you slipped in the mud,")
		p.damage(18)
	} else if roll < 84 {
		fmt.Printf("you fell into a ditch,")
		p.damage(36)
	} else if roll < 92 {
		dmg := p.hp * 0.05
		p.hp = max(p.hp-dmg, 0)
		fmt.Printf("you walked into a trap, \033[38;5;198m%.1f\033[0m dmg\n", dmg)
	} else {
		fmt.Println("you run around in circle")
	}
}

func (p *Player) setPerk(newperk int) {
	if newperk == 0 {
		p.hp += 20
		p.hpcap += 20
		p.defense += 5
	}

	if newperk == 1 {
		p.hp = max(1, p.hp-50)
		p.hpcap = max(1, p.hpcap-50)
		p.energy = max(1, p.energy-4)
		p.energycap = max(1, p.energycap-4)
	}

	if newperk == 3 {
		p.energycap += 2
	}

	if newperk == 6 {
		p.agility += 5
	}

	// adjustment when switching from old perks

	if p.perk == 0 {
		p.hp = max(1, p.hp-20)
		p.hpcap = max(1, p.hpcap-20)
		p.defense = max(1, p.defense-5)
	}

	if p.perk == 1 {
		p.hp += 50
		p.hpcap += 50
		p.energy += 4
		p.energycap += 4
	}

	if p.perk == 3 {
		p.energycap = max(1, p.energycap-2)
	}

	if p.perk == 6 {
		p.agility = max(1, p.agility-5)
	}

	p.perk = newperk
}

func (p Player) getperk() string {
	perk := []string{
		"Resilient",
		"Havoc",
		"Berserk",
		"Ingenious",
		"Poisoner",
		"Deadman",
		"Survivor",
		"Insanity",
	}
	return perk[p.perk]
}

func (p Player) energybar() string {
	bar := bars(40, float32(p.energy), float32(p.energycap))
	return fmt.Sprintf("\033[38;5;226m" + bar[0] + "\033[0m" + bar[1])
}
