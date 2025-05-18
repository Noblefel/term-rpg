package main

import "fmt"

type entity interface {
	attack(entity)
	// attack with predefined damage.
	// damage effects will be applied in here.
	attackWith(entity, float32)
	// apply effects that happens before player/enemy take action.
	applyEffects()
	// gets data like HP, strength...
	attr() attributes
	// niche case for demon attack that ignore defense
	setHP(float32)
	// takes attack dmg and log it. minimal dmg is 1.
	// defensive effects will be applied here.
	damage(float32)
}

type attributes struct {
	name     string
	hp       float32
	hpcap    float32
	strength float32
	defense  float32
	effects  map[string]int //name - turn duration
}

func (attr attributes) attack(target entity) {
	fmt.Printf("%s%s attacked!", success, attr.name)
	attr.attackWith(target, attr.strength)
}

func (attr attributes) attackWith(target entity, dmg float32) {
	if attr.effects["strengthen"] > 0 {
		dmg += dmg * 0.1
	}

	if attr.effects["weakened"] > 0 {
		dmg -= dmg * 0.13
	}

	target.damage(dmg)
}

func (attr *attributes) applyEffects() {
	if attr.effects["poisoned"] > 0 {
		dmg := 2 + attr.hp*0.07
		fmt.Printf("  %s took \033[38;5;198m%.1f\033[0m from poison\n", attr.name, dmg)
		attr.hp -= dmg // poison ignore defense
	}

	if attr.effects["poisoned severe"] > 0 {
		dmg := 2 + attr.hp*0.16
		fmt.Printf("  %s took \033[38;5;198m%.1f\033[0m from severe poison\n", attr.name, dmg)
		attr.hp -= dmg
	}

	if attr.effects["burning"] > 0 {
		fmt.Printf("  %s suffer from burning:", attr.name)
		attr.damage(5 + attr.hpcap*0.04)
	}

	if attr.effects["burning severe"] > 0 {
		fmt.Printf("  %s suffer from severe burning:", attr.name)
		attr.damage(5 + attr.hpcap*0.10)
	}

	if attr.effects["heal aura"] > 0 {
		heal := attr.hpcap * 0.07
		attr.hp = min(attr.hp+heal, attr.hpcap)
		fmt.Printf("  %s recover \033[38;5;83m%.1f\033[0m hp from healing aura\n", attr.name, heal)
	}
}

func (attr attributes) attr() attributes {
	return attr
}

func (attr *attributes) setHP(hp float32) {
	attr.hp = hp
}

func (attr *attributes) damage(dmg float32) {
	defense := attr.defense

	if attr.effects["immunity"] > 0 {
		fmt.Println(" \033[38;5;226mimmune\033[0m")
		return
	}

	if attr.effects["barrier"] > 0 {
		dmg -= dmg * 0.4
	}

	if attr.effects["force-field"] > 0 {
		dmg -= dmg * 0.15
	}

	if attr.effects["weakened"] > 0 {
		defense /= 2
	}

	dmg = max(dmg-defense, 1)
	attr.hp = max(attr.hp-dmg, 0)
	fmt.Printf(" \033[38;5;198m%.1f\033[0m\n", dmg)
}

func (attr attributes) decrementEffect() {
	for k := range attr.effects {
		attr.effects[k]--
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
