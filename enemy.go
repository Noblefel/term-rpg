package main

import (
	"fmt"
	"math/rand/v2"
)

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
		hp:       scale(90, 7),
		hpcap:    scale(90, 7),
		defense:  scale(6, 0.35),
		strength: scale(9, 1),
		effects:  make(map[string]int),
	}
	return &knight{attr}
}

func (k *knight) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 15 {
		def := 0.1 + rand.Float32()*scale(1, 0.35)
		k.defense += def
		k.hp = min(k.hp+5+def, k.hpcap)
		fmt.Printf("The knight reinforced his armor! ")
		fmt.Printf("\033[38;5;83m+%.1f\033[0m defense\n", def)
		return
	}

	messages := []string{
		"The knight charged at you!",
		"The knight cuts you with his sword!",
		"The knight slashes you with his sword!",
		"The knight smashed you with his shield!",
	}
	fmt.Printf(messages[rand.IntN(len(messages))])
	target.damage(k.strength)
}

type wizard struct {
	attributes
}

func newWizard() entity {
	attr := attributes{
		name:     "Wizard",
		hp:       scale(55, 3),
		hpcap:    scale(55, 3),
		defense:  scale(2, 0.15),
		strength: scale(4, 0.65),
		effects:  make(map[string]int),
	}
	return &wizard{attr}
}

func (w *wizard) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 20 {
		fmt.Printf("Wizard cannot cast, hit you with his staff instead!")
		target.damage(w.strength)
	} else if roll < 30 {
		fmt.Printf("Wizard cast \033[38;5;226menhanced\033[0m attack!")
		target.damage(w.strength * 3)
	} else if roll < 50 {
		fmt.Printf("Wizard cast \033[38;5;226mfireball\033[0m!")
		dmg := rand.Float32() * scale(30, 1.75)
		target.damage(dmg)
	} else if roll < 70 {
		fmt.Printf("Wizard cast \033[38;5;226mlightning\033[0m!")
		dmg := 10 + rand.Float32()*scale(20, 1)
		target.damage(dmg)
	} else if roll < 80 {
		fmt.Printf("Wizard summons \033[38;5;226mmeteor\033[0m!")
		target.damage(35)
	} else {
		heal := w.hpcap * 0.20
		w.hp = min(w.hp+heal, w.hpcap)
		fmt.Printf("Wizard cast healing! ")
		fmt.Printf("recover \033[38;5;83m%.1f\033[0m hp\n", heal)
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
	fmt.Print(success)

	if !c.mimic {
		c.hp = target.attr().hp
		c.hpcap = target.attr().hpcap
		c.strength = target.attr().strength
		c.defense = target.attr().defense
		c.mimic = true
		fmt.Println("Changeling \033[38;5;226mmorph\033[0m itself to look like you!")
		return
	}

	fmt.Printf("Changeling attacked! ")
	target.damage(c.strength)
}

type vampire struct {
	attributes
}

func newVampire() entity {
	attr := attributes{
		name:     "Vampire",
		hp:       scale(80, 5),
		hpcap:    scale(80, 5),
		defense:  scale(4, 0.15),
		strength: scale(10, 1.9),
		effects:  make(map[string]int),
	}
	return &vampire{attr}
}

func (v *vampire) attack(target entity) {
	roll := roll()

	if roll < 15 {
		fmt.Print(fail)
		fmt.Printf("Vampire get exposed to sunlight! getting")
		v.damage(scale(10, 0.9))
		return
	}

	fmt.Print(success)

	if roll < 70 {
		temp := target.attr().hp
		fmt.Printf("Vampire bites down hard! heal/deal")
		target.damage(v.strength)
		v.hp = min(v.hp+temp-target.attr().hp, v.hpcap)
	} else if roll < 85 {
		fmt.Printf("Vampire slashed you with claws!")
		target.damage(v.strength)
	} else {
		fmt.Printf("Vampire summoned then swarm you with bats!")
		target.damage(v.strength * 1.2)
	}
}

type demon struct {
	attributes
}

func newDemon() entity {
	attr := attributes{
		name:     "Demon",
		hp:       scale(112, 8),
		hpcap:    scale(112, 8),
		defense:  scale(4, 0.3),
		strength: scale(10, 1.8),
		effects:  make(map[string]int),
	}
	return &demon{attr}
}

func (d *demon) attack(target entity) {
	fmt.Print(success)

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
		fmt.Printf(messages[rand.IntN(len(messages))])
		fmt.Printf(" \033[38;5;198m%.1f\033[0m damage\n", drain)
		return
	}

	messages := []string{
		"Demon charged at you!",
		"Demon thrusts its trident at you!",
		"Demon swung its trident and slashes you!",
		"Demon lunged before punching you in the gut!",
	}

	fmt.Printf(messages[rand.IntN(len(messages))])
	target.damage(d.strength)
}
