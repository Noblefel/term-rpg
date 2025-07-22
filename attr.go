package main

import (
	"fmt"
	"math/rand/v2"
)

type entity interface {
	attack(entity)
	// attack with predefined damage.
	// damage effects will be applied in here.
	attackWith(entity, float64)
	// apply effects that happens before player/enemy take action.
	applyEffects()
	// gets data like HP, strength...
	attr() attributes
	// niche case for attacks that ignore both defense & effects
	setHP(float64)
	// takes attack dmg and log it. minimal dmg is 1.
	// defensive effects will be applied here.
	damage(float64)
}

type attributes struct {
	name     string
	hp       float64
	hpcap    float64
	strength float64
	defense  float64
	agility  float64
	effects  map[string]int //name - turn duration
}

func (attr *attributes) attack(target entity) {
	fmt.Printf("%s%s attacked!", success, attr.name)
	attr.attackWith(target, attr.strength)
}

func (attr *attributes) attackWith(target entity, dmg float64) {
	dodge := target.attr().agi()*0.5 - attr.agi()*0.15

	if dodge > rand.Float64()*100 {
		fmt.Println(" \033[38;5;226mmiss\033[0m")
		return
	}

	if attr.has("strengthen") {
		dmg += dmg * 0.1
	}

	if attr.has("vitality") {
		dmg += dmg * 0.05
	}

	if attr.has("weakened") {
		dmg -= dmg * 0.15
	}

	if attr.has("ace") {
		dmg += dmg * 0.28
	}

	crit := attr.agi() - target.attr().agi()*0.25

	if crit > rand.Float64()*100 {
		fmt.Print(" \033[38;5;226mcrit\033[0m")
		dmg *= 1.75
	}

	target.damage(dmg)

	// chain ifs to prevent multiple reflect damage
	if target.attr().has("reflect high") {
		fmt.Print("  reflected")
		attr.damage(dmg * 0.9)
	} else if target.attr().has("reflect") {
		fmt.Print("  reflected")
		attr.damage(dmg * 0.6)
	} else if target.attr().has("reflect low") {
		fmt.Print("  reflected")
		attr.damage(dmg * 0.3)
	}
}

func (attr *attributes) applyEffects() {
	// effects may ignore some defense for balance purposes
	// simply because they're too high for low damage spam like poison

	if attr.has("poisoned") && !attr.has("poison immunity") {
		fmt.Printf("  %s suffer from poison", attr.name)
		attr.damage(attr.defense/2 + attr.hp*0.11 + 10)
	}

	if attr.has("poisoned severe") && !attr.has("poison immunity") {
		fmt.Printf("  %s suffer from severe poison", attr.name)
		attr.damage(attr.defense/2 + attr.hp*0.22 + 20)
	}

	if attr.has("burning") && !attr.has("burning immunity") {
		fmt.Printf("  %s suffer from burning", attr.name)
		attr.damage(attr.defense/2 + attr.hpcap*0.05 + 10)
		delete(attr.effects, "frozen")
	}

	if attr.has("burning severe") && !attr.has("burning immunity") {
		fmt.Printf("  %s suffer from severe burning", attr.name)
		attr.damage(attr.defense/2 + attr.hpcap*0.1 + 20)
		delete(attr.effects, "frozen")
	}

	if attr.has("vitality") {
		heal := (attr.hpcap - attr.hp) * 0.1
		attr.hp = min(attr.hp+heal, attr.hpcap)
		fmt.Printf("  %s recover \033[38;5;83m%.1f\033[0m hp from vitality\n", attr.name, heal)
		attr.effects["bleeding"] -= 10
	}

	if attr.has("heal aura") {
		heal := attr.hpcap * 0.07
		attr.hp = min(attr.hp+heal, attr.hpcap)
		fmt.Printf("  %s recover \033[38;5;83m%.1f\033[0m hp from healing aura\n", attr.name, heal)
		attr.effects["bleeding"] -= 20
	}

	if attr.has("bleeding") {
		// in case of this effect, the "turn" count is changed to severity.
		// and this effect cant go away without healing.
		severity := attr.effects["bleeding"]
		attr.effects["bleeding"] += 8

		if severity > 60 {
			fmt.Printf("  %s suffer from severe bleeding", attr.name)
			attr.damage(attr.defense*0.7 + attr.hp*0.5)
		} else if severity > 30 {
			fmt.Printf("  %s suffer from heavy bleeding", attr.name)
			attr.damage(attr.defense*0.6 + attr.hp*0.24)
		} else if severity > 10 {
			fmt.Printf("  %s suffer from mild bleeding", attr.name)
			attr.damage(attr.defense*0.5 + attr.hp*0.12)
		} else {
			fmt.Printf("  %s suffer from minor bleeding", attr.name)
			attr.damage(attr.defense*0.4 + attr.hp*0.06)
		}
	}
}

func (attr attributes) attr() attributes { return attr }

func (attr attributes) agi() float64 {
	agi := attr.agility

	if attr.has("ace") {
		agi += agi * 0.15
	}

	if attr.has("focus") {
		agi += agi*0.3 + 5
	}

	if attr.has("shiver") {
		agi -= agi*0.3 + 5
	}

	if attr.has("force-field") {
		agi += agi * 0.1
	}

	return agi
}

func (attr *attributes) setHP(hp float64) {
	attr.hp = hp
}

func (attr *attributes) damage(dmg float64) {
	defense := attr.defense

	if attr.has("immunity") {
		fmt.Println(" \033[38;5;226mimmune\033[0m")
		return
	}

	if attr.has("strengthen") {
		defense += defense * 0.1
	}

	if attr.has("vitality") {
		dmg -= dmg * 0.05
	}

	if attr.has("barrier") {
		dmg -= dmg * 0.4
	}

	if attr.has("force-field") {
		dmg -= dmg * 0.15
	}

	if attr.has("ace") {
		dmg -= dmg * 0.28
	}

	// prevent instant shatter
	if attr.effects["frozen"] == 1 && !attr.has("frozen immunity") {
		if roll() < 50 {
			fmt.Print(" \033[38;5;226mshatter\033[0m")
			dmg *= 2
		} else {
			defense += 4 + defense*0.2
		}
	}

	if attr.has("weakened") {
		defense /= 2
	}

	dmg = max(dmg-defense, 1)
	attr.hp = max(attr.hp-dmg, 0)
	fmt.Printf(" \033[38;5;198m%.1f\033[0m\n", dmg)
}

func (attr attributes) has(effect string) bool {
	return attr.effects[effect] > 0
}

// called at the end of every turn
func (attr attributes) decrementEffect() {
	for k := range attr.effects {
		if attr.effects[k] <= 0 {
			delete(attr.effects, k)
		} else {
			attr.effects[k]--
		}
	}
}

func (attr attributes) hpbar() string {
	percentage := attr.hp / attr.hpcap * 100
	bars := bars(40, attr.hp, attr.hpcap)

	if percentage > 60 {
		bars[0] = "\033[38;5;83m" + bars[0] + "\033[0m"
	} else if percentage > 30 {
		bars[0] = "\033[38;5;226m" + bars[0] + "\033[0m"
	} else {
		bars[0] = "\033[38;5;196m" + bars[0] + "\033[0m"
	}

	return bars[0] + bars[1]
}
