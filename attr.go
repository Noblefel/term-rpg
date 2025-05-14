package main

import "fmt"

type entity interface {
	attack(entity)
	// gets data like HP, strength...
	attr() attributes
	// to modify HP immediately
	setHP(float32)
	// takes attack dmg and log it. Minimal dmg is 1.
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
	target.damage(attr.strength)
}

func (attr attributes) attr() attributes {
	return attr
}

func (attr *attributes) setHP(hp float32) {
	attr.hp = hp
}

func (attr *attributes) damage(dmg float32) {
	dmg = max(dmg-attr.defense, 1)
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
