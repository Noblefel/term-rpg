package main

import (
	"fmt"
	"math/rand/v2"
)

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
	fmt.Fprintf(out, "%s%s attacked!", success, attr.name)
	target.damage(attr.strength)
}

func (attr attributes) attr() attributes { return attr }

func (attr *attributes) setHP(hp float32) { attr.hp = hp }

func (attr *attributes) damage(dmg float32) {
	dmg = max(dmg-attr.defense, 1)
	attr.hp = max(attr.hp-dmg, 0)
	fmt.Fprintf(out, " \033[38;5;198m%.1f\033[0m damage\n", dmg)
}

func (attr attributes) decrementEffect() {
	for k := range attr.effects {
		attr.effects[k]--
	}
}

func randomEnemy() entity {
	var enemies = []func() entity{
		newKnight,
		newWizard,
		newChangeling,
		newVampire,
		newDemon,
	}

	spawn := enemies[rand.IntN(len(enemies))]
	return spawn()
}

type knight struct {
	attributes
}

func newKnight() entity {
	attr := attributes{
		name:     "Knight",
		hp:       scale(90, 4),
		hpcap:    scale(90, 4),
		defense:  scale(6, 0.4),
		strength: scale(9, 1),
		effects:  make(map[string]int),
	}
	return &knight{attr}
}

func (k *knight) attack(target entity) {
	fmt.Fprint(out, success)
	roll := roll()

	if roll < 15 {
		def := rand.Float32() * scale(5, 0.4)
		k.defense += def
		k.hp = min(k.hp+5+def, k.hpcap)
		fmt.Fprintf(out, "The knight reinforced his armor! ")
		fmt.Fprintf(out, "\033[38;5;83m+%.1f\033[0m defense\n", def)
		return
	}

	messages := []string{
		"The knight charged at you!",
		"The knight cuts you with his sword!",
		"The knight slashes you with his sword!",
		"The knight smashed you with his shield!",
	}
	fmt.Fprintf(out, messages[rand.IntN(len(messages))])
	target.damage(k.strength)
}

type wizard struct {
	attributes
}

func newWizard() entity {
	attr := attributes{
		name:     "Wizard",
		hp:       scale(55, 1.5),
		hpcap:    scale(55, 1.5),
		defense:  scale(2, 0.25),
		strength: scale(4, 0.65),
		effects:  make(map[string]int),
	}
	return &wizard{attr}
}

func (w *wizard) attack(target entity) {
	fmt.Fprint(out, success)
	roll := roll()

	if roll < 20 {
		fmt.Fprintf(out, "Wizard cannot cast, hit you with his staff instead!")
		target.damage(w.strength)
	} else if roll < 30 {
		fmt.Fprintf(out, "Wizard cast \033[38;5;226menhanced\033[0m attack!")
		target.damage(w.strength * 3)
	} else if roll < 50 {
		fmt.Fprintf(out, "Wizard cast \033[38;5;226mfireball\033[0m!")
		dmg := rand.Float32() * scale(30, 1.75)
		target.damage(dmg)
	} else if roll < 70 {
		fmt.Fprintf(out, "Wizard cast \033[38;5;226mlightning\033[0m!")
		dmg := 10 + rand.Float32()*scale(20, 1)
		target.damage(dmg)
	} else if roll < 80 {
		fmt.Fprintf(out, "Wizard summons \033[38;5;226mmeteor\033[0m!")
		target.damage(35)
	} else {
		heal := w.hpcap * 0.20
		w.hp = min(w.hp+heal, w.hpcap)
		fmt.Fprintf(out, "Wizard cast healing! ")
		fmt.Fprintf(out, "recover \033[38;5;83m%.1f\033[0m hp\n", heal)
	}
}

type changeling struct {
	attributes
	mimic bool
}

func newChangeling() entity {
	attr := attributes{
		name:     "Changeling",
		hp:       80,
		hpcap:    80,
		defense:  1,
		strength: 1,
		effects:  make(map[string]int),
	}
	return &changeling{attr, false}
}

func (c *changeling) attack(target entity) {
	fmt.Fprint(out, success)

	if !c.mimic {
		c.hp = target.attr().hp
		c.hpcap = target.attr().hpcap
		c.strength = target.attr().strength
		c.defense = target.attr().defense
		c.mimic = true
		fmt.Fprintln(out, "Changeling \033[38;5;226mmorph\033[0m itself to look like you!")
		return
	}

	fmt.Fprintf(out, "Changeling attacked! ")
	target.damage(c.strength)
}

type vampire struct {
	attributes
}

func newVampire() entity {
	attr := attributes{
		name:     "Vampire",
		hp:       scale(80, 2.6),
		hpcap:    scale(80, 2.6),
		defense:  scale(4, 0.3),
		strength: scale(10, 1.9),
		effects:  make(map[string]int),
	}
	return &vampire{attr}
}

func (v *vampire) attack(target entity) {
	roll := roll()

	if roll < 15 {
		fmt.Fprint(out, fail)
		fmt.Fprintf(out, "Vampire get exposed to sunlight! getting")
		v.damage(scale(10, 0.9))
		return
	}

	fmt.Fprint(out, success)

	if roll < 70 {
		temp := target.attr().hp
		fmt.Fprintf(out, "Vampire bites down hard! heal/deal")
		target.damage(v.strength)
		v.hp = min(v.hp+temp-target.attr().hp, v.hpcap)
	} else if roll < 85 {
		fmt.Fprintf(out, "Vampire slashed you with claws!")
		target.damage(v.strength)
	} else {
		fmt.Fprintf(out, "Vampire summoned then swarm you with bats!")
		target.damage(v.strength * 1.2)
	}
}

type demon struct {
	attributes
}

func newDemon() entity {
	attr := attributes{
		name:     "Demon",
		hp:       scale(112, 3),
		hpcap:    scale(112, 3),
		defense:  scale(4, 0.5),
		strength: scale(10, 1.8),
		effects:  make(map[string]int),
	}
	return &demon{attr}
}

func (d *demon) attack(target entity) {
	fmt.Fprint(out, success)

	if roll() < 60 {
		messages := []string{
			"Demon drains your life force!",
			"Demon absors your life force!",
			"Demon reach out and take good portion of your hp!",
			"Demon conjure dark energy that weakens you whole!",
		}

		drain := d.strength + target.attr().hp*0.05
		hp := max(target.attr().hp-drain, 0)
		target.setHP(hp)
		fmt.Fprintf(out, messages[rand.IntN(len(messages))])
		fmt.Fprintf(out, " \033[38;5;198m%.1f\033[0m damage\n", drain)
		return
	}

	messages := []string{
		"Demon charged at you!",
		"Demon thrusts its trident at you!",
		"Demon swung its trident and slashes you!",
		"Demon lunged before punching you in the gut!",
	}

	fmt.Fprintf(out, messages[rand.IntN(len(messages))])
	target.damage(d.strength)
}
