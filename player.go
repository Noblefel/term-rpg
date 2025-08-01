package main

import (
	"fmt"
	"math/rand/v2"
)

type Player struct {
	attributes
	gold      int
	energy    int
	energycap int

	perk   int    //perk index
	weapon int    //weapon index
	armor  int    //armor index
	skills [5]int //skill indexes
}

var perks = []string{
	"No perk",
	"Resilient", "Havoc", "Berserk", "Wizardry",
	"Poisoner", "Deadman", "Survivor", "Insanity",
	"Shock", "Frigid", "Ranger", "Fencer", "Smith",
}

var skills = []struct {
	name string
	desc string
	cost int
	cd   int
}{
	{"charge", "attack with 130% strength", 4, 2},
	{"frenzy", "sacrifice hp to attack with 250% strength", 8, 6},
	{"great blow", "sacrifice the next turn to attack with 210% strength", 5, 3},
	{"poison", "attack 85% strength and poison enemy for 3 turns", 5, 5},
	{"stun", "attack 60% strength and stun enemy for 2 turns", 6, 4},
	{"icy blast", "attack 60% strength and 30% chance to inflict freeze", 5, 3},
	{"swift strike", "attack 85% strength (doesnt use turn)", 4, 4},
	{"knives throw", "attack 40 fixed damage (doesnt use turn, no cd, no wep effects)", 4, 0},
	{"fireball", "deal moderate amount of damage and inflict burning (no wep effects)", 7, 5},
	{"meteor strike", "deal huge amount of damage (no wep effects)", 11, 5},
	{"strengthen", "attack 100% strength to increase damage by 10% for 3 turns", 4, 7},
	{"focus attack", "attack 100% strength to get focus state for 3 turns", 4, 5},
	{"devour", "attack 150% strength and heal by 5% hp cap", 8, 5},
	{"vitality", "get vitality effect for 5 turns", 4, 5},
	{"barrier", "reduce incoming damage by 40% for 2 turns", 4, 3},
	{"force-field", "reduce incoming damage by 15% for 5 turns", 5, 8},
	{"heal spell", "recover hp by atleast 15% hp cap", 8, 5},
	{"heal aura", "recover hp by atleast 7% of hpcap for 3 turns", 8, 6},
	{"heal potion", "recover hp by 40 (fixed number)", 7, 10},
	{"drain", "take 22% of enemy current hp as damage", 4, 4},
	{"absorb", "take 10% of enemy hp cap, ignore defense and effects", 5, 6},
	{"trick", "make the enemy self-target", 4, 3},
	{"vision", "see enemy attributes (no cost, no cd, doesnt use turn)", 0, 0},
}

var weapons = []struct {
	name string
	desc string
	cost int
}{
	{"no weapon", "", 0},
	//common
	{"sword", "+15 strength", 40},
	{"needle", "+10 strength, ignore 10% defense", 30},
	{"club", "+12 strength, +5% multiplier", 38},
	{"daggers", "+9 strength, +2 agility, -2 defense", 41},
	{"staff", "+8 strength, +2 energy cap", 45},
	{"gloves", "+6 strength, +4 defense", 40},
	//rare
	{"greatsword", "+40 strength", 300},
	{"flaming sword", "+26 strength, 15% chance to inflict burning", 247},
	{"rapier", "+20 strength, ignore 25% defense", 258},
	{"warhammer", "+30 strength, +10% multiplier", 330},
	{"chain daggers", "+15 strength, +4 agility, -2 defense", 285},
	{"enchanted staff", "+16 strength, +4 energy cap", 262},
	{"gauntlets", "+12 strength, +8 defense", 240},
	//more rare
	{"demonic blade", "+30 strength, +5% enemy current hp as damage", 690},
	{"daunting mace", "+20 strength, 30 hp, 7% self hp cap as damage", 1398},
	{"crimson blade", "+30 strength, heal by 5 hp (fixed)", 780},
	{"dragonscale blade", "+50 strength, 30% chance to inflict burning or 10% severe burning", 1437},
	{"astral rapier", "+40 strength, ignore 40% defense", 1125},
	{"lance", "+20 strength, +50% multiplier on first attack", 1308},
	{"obsidian warhammer", "+60 strength, +20% multiplier increase", 1500},
	{"holy staff", "+30 strength, +6 energy cap", 1200},
	{"king's gauntlets", "+24 strength, +16 defense", 1125},
	//exceptional
	{"voidforged rapier", "+40 strength, ignore defense", 3500},
	{"soulreaper", "+25 strength, +15% enemy current hp as damage", 3732},
	{"celestial staff", "+40 strength, +8 energy cap, -1 cooldown, heal 2% hp cap", 4200},
	{"vanguard lance", "+40 strength, +100% multiplier on first attack", 4200},
	{"earthbreaker", "+100 strength, +25% multiplier increase", 4704},
}

var armory = []struct {
	name string
	desc string
	cost int
}{
	{"no armor", "", 0},
	//common
	{"leather armor", "40 hp, 2 defense", 50},
	{"wooden shield", "20 hp, 4 defense", 63},
	{"basic plate", "80 hp", 80},
	{"spike helmet", "10 hp, 2 defense, 5 strength", 65},
	{"cloak", "10 hp, 3 agility", 51},
	//rare
	{"iron shield", "20 hp, 12 defense", 440},
	{"standard armor", "150 hp, 6 defense, -1 agility", 450},
	{"crystal armor", "75 hp, 6 defense, reflect (small)", 500},
	{"enchanted plate", "120 hp, 8% reduction", 475},
	{"heavy plate", "240 hp", 437},
	{"wolfshead", "75 hp, 5 defense, 10 strength", 480},
	{"arcane vest", "30 hp, 4 energy cap", 420},
	//more rare
	{"deepsea mantle", "24 hp, initial damage cannot exceed 20% hp cap", 1400},
	{"obsidian shield", "20 hp, 33 defense", 1310},
	{"obsitian armor", "325 hp, 15 defense, -2 agility", 1500},
	{"mythril plate", "200 hp, 16% reduction", 1540},
	{"reinforced plate", "500 hp, -2 agility", 1410},
	{"chainmail", "150 hp, 6 defense, +15% defense value", 1450},
	{"king's helmet", "120 hp, 10 defense, 20 strength", 1390},
	{"misty cloak", "50 hp, 12 agility", 1200},
	//exceptional
	{"void mantle", "24 hp, immune to the first damage", 3000},
	{"energy shield", "20 hp, 2 energy cap, get 35% reduction every 4 turn", 3000},
	{"amethyst armor", "75 hp, 6 defense, reflect (high)", 3100},
	{"conqueror's armor", "400 hp, 20 defense, 8% reduction", 3200},
}

func NewPlayer() *Player {
	var player Player
	player.name = "player"
	player.gold = 100

	player.hp = 200
	player.hpcap = 200
	player.strength = 50
	player.defense = 15
	player.agility = 10
	player.energy = 20
	player.energycap = 20

	// charge, stun, poison, fireball, heal potion
	player.skills = [5]int{0, 3, 4, 8, 14}
	player.effects = make(map[string]int)

	return &player
}

func (p *Player) attack(enemy entity) {
	fmt.Printf(success + "you attacked!")

	if p.is("Fencer") {
		p.attackWith(enemy, p.strength*0.55)
		fmt.Printf(success + "you attacked!")
		p.attackWith(enemy, p.strength*0.55)
		return
	}

	p.attackWith(enemy, p.strength)
}

// player perks modifier applied here
func (p *Player) attackWith(enemy entity, dmg float64) {
	dmg = p.useWeapon(dmg, enemy)

	if p.energy <= 3 {
		fmt.Print(" exhausted")
		dmg -= dmg * 0.2
	}

	if p.is("Havoc") {
		dmg += dmg * 0.2
	}

	if p.is("Berserk") {
		percent := p.hp / (p.hpcap * 0.4) //if 1% hp   = 60% extra dmg
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.6 - percent*0.6
		dmg += dmg * mul
	}

	if p.is("Insanity") {
		roll := roll()

		if roll < 25 {
			// multiplier - 30% increase/decrease range
			mul := 0.3 - rand.Float64()*0.6
			dmg += dmg * mul
		} else if roll < 50 {
			// flat val (scaled)
			val := scale(1, 1)
			dmg += val/2 - rand.Float64()*val
		} else if roll < 75 {
			// flat val
			dmg += 10 - rand.Float64()*20
		}
	}

	if p.is("Frigid") && roll() < 15 {
		fmt.Print(" \033[38;5;226mfreeze\033[0m")
		enemy.attr().effects["frozen"] = 2
	}

	p.attributes.attackWith(enemy, dmg)
}

// player perks modifier applied here
func (p *Player) damage(dmg float64) {
	dmg = p.useArmor(dmg)

	if p.is("Resilient") {
		dmg -= dmg * 0.1
	}

	if p.is("Berserk") {
		percent := p.hp / (p.hpcap * 0.4) //if 0% hp   = 40% reduction
		percent = min(percent, 1)         //if 40%+ hp = none
		mul := 0.4 - percent*0.4
		dmg -= dmg * mul
	}

	p.attributes.damage(dmg)
}

func (p *Player) skill(i int, enemy entity) bool {
	if p.has("disoriented") {
		fmt.Println("\033[38;5;196mcannot use skill when disoriented\033[0m")
		return false
	}

	skill := skills[i]
	cost := skill.cost

	if p.has("confused") {
		cost++
	}

	if cost > p.energy {
		fmt.Println("\033[38;5;196mnot enough energy\033[0m")
		return false
	}

	if p.has("cd" + skill.name) {
		fmt.Println("\033[38;5;196mskill in cooldown\033[0m")
		return false
	}

	cooldown := skill.cd

	if p.is("Wizardry") {
		cooldown -= 2
	}

	if p.is("Berserk") && p.hp/p.hpcap <= 0.25 {
		cooldown--
	}

	if p.is("Insanity") {
		cooldown = rand.IntN(8)
	}

	if p.weapon > 0 && weapons[p.weapon].name == "celestial staff" {
		cooldown--
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
	case "icy blast":
		if roll() < 30 {
			enemy.attr().effects["frozen"] = 2
		}
		p.attackWith(enemy, p.strength*0.6)
	case "swift strike":
		p.attackWith(enemy, p.strength*0.85)
		return false
	case "knives throw":
		p.effects["no wep effects"] = 1
		p.attackWith(enemy, 40)
		return false
	case "fireball":
		enemy.attr().effects["burning"] = 2
		p.effects["no wep effects"] = 1
		p.attackWith(enemy, 80)
	case "meteor strike":
		dmg := 50 + rand.Float64()*170
		p.effects["no wep effects"] = 1
		p.attackWith(enemy, dmg)
	case "strengthen":
		p.attackWith(enemy, p.strength) // attack first so it wont get the bonus yet
		p.effects["strengthen"] = 4     // +1 for 3 attacks'
	case "focus attack":
		p.attackWith(enemy, p.strength)
		p.effects["focus"] = 4 // +1 for 3 attacks
	case "devour":
		p.attackWith(enemy, p.strength*1.5)
		heal := p.hpcap * 0.05
		p.effects["bleeding"] -= 10
		p.hp = min(p.hp+heal, p.hpcap)
		fmt.Printf("  recover \033[38;5;83m%.1f\033[0m hp\n", heal)
	case "vitality":
		p.effects["vitality"] = 5
		fmt.Println(" applying bonuses for 5 turn")
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
		delete(p.effects, "bleeding")
	case "heal aura":
		p.effects["heal aura"] = 4 // +1 because it start in next turn
		fmt.Println(" recover \033[38;5;83m7%\033[0m hp for 3 turns")
	case "heal potion":
		p.hp = min(p.hp+40, p.hpcap)
		fmt.Println(" recover \033[38;5;83m40\033[0m hp")
		delete(p.effects, "bleeding")
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
	heal := p.hpcap * 0.1
	heal += 20 + rand.Float64()*20
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
		n := 2.5 + rand.Float64()*5
		p.hpcap += n
		fmt.Printf("hp cap increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 73 {
		n := 0.5 + rand.Float64()*2
		p.strength += n
		fmt.Printf("strength increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 84 {
		n := 0.5 + rand.Float64()*2
		p.defense += n
		fmt.Printf("defense increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else if roll < 95 {
		n := 0.1 + rand.Float64()*1
		p.agility += n
		fmt.Printf("agility increased by \033[38;5;83m%.1f\033[0m\n", n)
	} else {
		p.energycap++
		fmt.Println("energy cap increased by \033[38;5;83m1\033[0m")
	}
}

func (p *Player) flee(enemy entity) {
	roll := roll()
	flee := float64(roll) - p.agi()*0.5 + enemy.attr().agi()*0.2

	if flee < 25 || (roll < 95 && p.is("Survivor")) {
		fmt.Print(success)
		fmt.Println("you have fled the battle")
		p.effects["fled"] = 1
		return
	}

	fmt.Print(fail)

	if roll < 68 {
		fmt.Println("youre too slow and got caught")

		if enemy.attr().has("frozen") && !enemy.attr().has("frozen immunity") {
			fmt.Print(fail)
			fmt.Printf("%s tried to attack but is frozen\n", enemy.attr().name)
		} else if enemy.attr().has("stunned") {
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

func (p *Player) setWeapon(index int) {
	switch weapons[index].name {
	case "sword":
		p.strength += 15
	case "needle":
		p.strength += 10
	case "club":
		p.strength += 12
	case "daggers":
		p.strength += 9
		p.agility += 2
		p.defense -= 2
	case "staff":
		p.strength += 8
		p.energycap += 2
	case "gloves":
		p.strength += 6
		p.defense += 4
	case "greatsword":
		p.strength += 40
	case "flaming sword":
		p.strength += 26
	case "rapier":
		p.strength += 20
	case "warhammer":
		p.strength += 30
	case "chain daggers":
		p.strength += 15
		p.agility += 4
	case "enchanted staff":
		p.strength += 16
		p.energycap += 4
	case "gauntlets":
		p.strength += 12
		p.defense += 8
	case "demonic blade":
		p.strength += 30
	case "crimson blade":
		p.strength += 30
	case "dragonscale blade":
		p.strength += 50
	case "astral rapier":
		p.strength += 40
	case "lance":
		p.strength += 20
	case "obsidian warhammer":
		p.strength += 60
	case "holy staff":
		p.strength += 30
		p.energycap += 6
	case "king's gauntlets":
		p.strength += 24
		p.defense += 16
	case "daunting mace":
		p.strength += 20
		p.hpcap += 30
	case "voidforged rapier":
		p.strength += 40
	case "soulreaper":
		p.strength += 25
	case "celestial staff":
		p.strength += 40
		p.energycap += 8
	case "vanguard lance":
		p.strength += 40
	case "earthbreaker":
		p.strength += 100
	}

	switch weapons[p.weapon].name {
	case "sword":
		p.strength -= 15
	case "needle":
		p.strength -= 10
	case "club":
		p.strength -= 12
	case "daggers":
		p.strength -= 9
		p.agility -= 2
		p.defense += 2
	case "staff":
		p.strength -= 8
		p.energycap -= 2
	case "gloves":
		p.strength -= 6
		p.defense -= 4
	case "greatsword":
		p.strength -= 40
	case "flaming sword":
		p.strength -= 26
	case "rapier":
		p.strength -= 20
	case "warhammer":
		p.strength -= 30
	case "chain daggers":
		p.strength -= 15
		p.agility -= 4
		p.defense += 2
	case "enchanted staff":
		p.strength -= 16
	case "gauntlets":
		p.strength -= 12
		p.defense -= 8
	case "demonic blade":
		p.strength -= 30
	case "crimson blade":
		p.strength -= 30
	case "dragonscale blade":
		p.strength -= 50
	case "astral rapier":
		p.strength -= 40
	case "lance":
		p.strength -= 20
	case "obsidian warhammer":
		p.strength -= 60
	case "holy staff":
		p.strength -= 30
		p.energycap -= 6
	case "king's gauntlets":
		p.strength -= 24
		p.defense -= 16
	case "daunting mace":
		p.strength -= 20
		p.hpcap -= 30
	case "voidforged rapier":
		p.strength -= 40
	case "soulreaper":
		p.strength -= 25
	case "celestial staff":
		p.strength -= 40
		p.energycap -= 8
	case "vanguard lance":
		p.strength -= 40
	case "earthbreaker":
		p.strength -= 100
	}

	p.hpcap = max(p.hpcap, 50)
	p.hp = min(p.hp, p.hpcap)
	p.strength = max(p.strength, 5)
	p.defense = max(p.defense, 1)
	p.agility = max(p.agility, 1)
	p.energycap = max(p.energycap, 10)
	p.energy = min(p.energy, p.energycap)
	p.weapon = index
}

// specific weapon effects
func (p *Player) useWeapon(dmg float64, enemy entity) float64 {
	mul := 1.0

	if perks[p.perk] == "Fencer" {
		mul = 0.55
	}

	if perks[p.perk] == "Smith" {
		mul = 1.14
	}

	// for knives throw, fireball, meteor skill
	if p.has("no wep effects") {
		delete(p.effects, "no wep effects")
		return dmg
	}

	switch weapons[p.weapon].name {
	case "needle":
		dmg += enemy.attr().defense * 0.1 * mul
	case "club":
		dmg += dmg * 0.05 * mul
	case "flaming sword":
		if roll() < 15 {
			enemy.attr().effects["burning"] = 2
		}
	case "rapier":
		dmg += enemy.attr().defense * 0.25 * mul
	case "warhammer":
		dmg += dmg * 0.1 * mul
	case "demonic blade":
		dmg += enemy.attr().hp * 0.05 * mul
	case "crimson blade":
		heal := 5 * mul
		p.effects["bleeding"] -= 5
		p.hp = min(p.hp+heal, p.hpcap)
		fmt.Printf(" heal by %.1f", heal)
	case "dragonscale blade":
		roll := roll()

		if roll < 10 {
			enemy.attr().effects["burning severe"] = 2
		} else if roll < 40 {
			enemy.attr().effects["burning"] = 2
		}
	case "astral rapier":
		dmg += enemy.attr().defense * 0.4 * mul
	case "lance":
		if !p.has("lance") {
			dmg += dmg * 0.5 * mul
			p.effects["lance"] = 99
		}
	case "obsidian warhammer":
		dmg += dmg * 0.2 * mul
	case "daunting mace":
		dmg += p.hpcap * 0.07 * mul
	case "voidforged rapier":
		dmg += enemy.attr().defense * mul
	case "soulreaper":
		dmg += enemy.attr().hp * 0.15 * mul
	case "celestial staff":
		heal := p.hpcap * 0.02 * mul
		p.effects["bleeding"] -= 10
		p.hp = min(p.hp+heal, p.hpcap)
		fmt.Printf(" heal by %.1f", heal)
	case "vanguard lance":
		if !p.has("lance") {
			dmg += dmg * mul
			p.effects["lance"] = 99
		}
	case "earthbreaker":
		dmg += dmg * 0.25 * mul
	}

	return dmg
}

func (p *Player) setArmor(index int) {
	switch armory[index].name {
	case "leather armor":
		p.hpcap += 40
		p.defense += 2
	case "wooden shield":
		p.hpcap += 20
		p.defense += 4
	case "basic plate":
		p.hpcap += 80
	case "spike helmet":
		p.hpcap += 10
		p.defense += 2
		p.strength += 5
	case "cloak":
		p.hpcap += 10
		p.agility += 3
	case "iron shield":
		p.hpcap += 20
		p.defense += 12
	case "standard armor":
		p.hpcap += 150
		p.defense += 6
		p.agility -= 1
	case "crystal armor":
		p.hpcap += 75
		p.defense += 6
	case "enchanted plate":
		p.hpcap += 120
	case "heavy plate":
		p.hpcap += 240
	case "wolfshead":
		p.hpcap += 75
		p.defense += 5
		p.strength += 10
	case "arcane vest":
		p.hpcap += 30
		p.energycap += 4
	case "deepsea mantle":
		p.hpcap += 24
	case "obsidian shield":
		p.hpcap += 20
		p.defense += 33
	case "obsidian armor":
		p.hpcap += 325
		p.defense += 15
		p.agility -= 2
	case "mythril plate":
		p.hpcap += 200
	case "reinforced plate":
		p.hpcap += 500
		p.agility -= 2
	case "chainmail":
		p.hpcap += 150
		p.defense += 6
	case "king's helmet":
		p.hpcap += 120
		p.defense += 10
		p.strength += 20
	case "misty cloak":
		p.hpcap += 50
		p.agility += 12
	case "void mantle":
		p.hpcap += 24
	case "energy shield":
		p.hpcap += 20
		p.energycap += 2
	case "amethyst armor":
		p.hpcap += 75
		p.defense += 6
	case "conqueror's armor":
		p.hpcap += 400
		p.defense += 20
	}

	switch armory[p.armor].name {
	case "leather armor":
		p.hpcap -= 40
		p.defense -= 2
	case "wooden shield":
		p.hpcap -= 20
		p.defense -= 4
	case "basic plate":
		p.hpcap -= 80
	case "spike helmet":
		p.hpcap -= 10
		p.defense -= 2
		p.strength -= 5
	case "cloak":
		p.hpcap -= 10
		p.agility -= 3
	case "iron shield":
		p.hpcap -= 20
		p.defense -= 12
	case "standard armor":
		p.hpcap -= 150
		p.defense -= 6
		p.agility += 1
	case "crystal armor":
		p.hpcap -= 75
		p.defense -= 6
	case "enchanted plate":
		p.hpcap -= 120
	case "heavy plate":
		p.hpcap -= 240
	case "wolfshead":
		p.hpcap -= 75
		p.defense -= 5
		p.strength -= 10
	case "arcane vest":
		p.hpcap -= 30
		p.energycap -= 4
	case "deepsea mantle":
		p.hpcap -= 24
	case "obsidian shield":
		p.hpcap -= 20
		p.defense -= 33
	case "obsidian armor":
		p.hpcap -= 325
		p.defense -= 15
		p.agility += 2
	case "mythril plate":
		p.hpcap -= 200
	case "reinforced plate":
		p.hpcap -= 500
		p.agility += 2
	case "chainmail":
		p.hpcap -= 150
		p.defense -= 6
	case "king's helmet":
		p.hpcap -= 120
		p.defense -= 10
		p.strength -= 20
	case "misty cloak":
		p.hpcap -= 50
		p.agility -= 12
	case "void mantle":
		p.hpcap -= 24
	case "energy shield":
		p.hpcap -= 20
		p.energycap -= 2
	case "amethyst armor":
		p.hpcap -= 75
		p.defense -= 6
	case "conqueror's armor":
		p.hpcap -= 400
		p.defense -= 20
	}

	p.hpcap = max(p.hpcap, 50)
	p.hp = min(p.hp, p.hpcap)
	p.strength = max(p.strength, 5)
	p.defense = max(p.defense, 1)
	p.agility = max(p.agility, 1)
	p.energycap = max(p.energycap, 10)
	p.energy = min(p.energy, p.energycap)
	p.armor = index
}

// specific armor effects
func (p *Player) useArmor(dmg float64) float64 {
	switch armory[p.armor].name {
	case "crystal armor":
		p.effects["reflect low"] = 5
	case "enchanted plate", "conqueror's armor":
		dmg -= dmg * 0.08
	case "deepsea mantle":
		dmg = min(dmg, p.hpcap*0.20+p.defense)
	case "mythril plate":
		dmg -= dmg * 0.16
	case "chainmail":
		dmg -= p.defense * 0.15
	case "void mantle":
		if !p.has("void mantle") {
			p.effects["void mantle"] = 99
			p.effects["immunity"] = 1
		}
	case "energy shield":
		if !p.has("energy shield") {
			dmg -= dmg * 0.35
			p.effects["energy shield"] = 4
		}
	case "amethyst armor":
		p.effects["reflect high"] = 5
	}

	return dmg
}

func (p *Player) setPerk(index int) {
	switch perks[index] {
	case "Resilient":
		p.hpcap += 20
		p.defense += 5
	case "Havoc":
		p.hpcap -= 50
		p.energycap -= 4
	case "Wizardry":
		p.energycap += 2
	case "Survivor":
		p.agility += 5
	case "Shock":
		p.strength += 5
	case "Frigid":
		p.agility -= 5
	}

	switch perks[p.perk] {
	case "Resilient":
		p.hpcap -= 20
		p.defense -= 5
	case "Havoc":
		p.hpcap += 50
		p.energycap += 4
	case "Wizardry":
		p.energycap -= 2
	case "Survivor":
		p.agility -= 5
	case "Shock":
		p.strength -= 5
	case "Frigid":
		p.agility += 5
	}

	p.hpcap = max(p.hpcap, 50)
	p.hp = min(p.hp, p.hpcap)
	p.strength = max(p.strength, 5)
	p.defense = max(p.defense, 1)
	p.agility = max(p.agility, 1)
	p.energycap = max(p.energycap, 10)
	p.energy = min(p.energy, p.energycap)
	p.perk = index
}

func (p Player) is(perk string) bool { return perks[p.perk] == perk }

func (p Player) energybar() string {
	bar := bars(40, float64(p.energy), float64(p.energycap))
	return fmt.Sprintf("\033[38;5;226m" + bar[0] + "\033[0m" + bar[1])
}
