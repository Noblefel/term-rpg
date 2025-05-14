package main

import (
	"fmt"
	"math/rand/v2"
)

func randomEnemy() entity {
	var enemies = []func() entity{
		// newKnight,
		// newWizard,
		// newChangeling,
		// newVampire,
		// newDemon,
		newShardling,
		newGenie,
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
		fmt.Printf("Vampire bites down hard! heal & deal:")
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

		drain := d.strength + target.attr().hp*0.07
		hp := max(target.attr().hp-drain, 0)
		target.setHP(hp)
		fmt.Printf(messages[rand.IntN(len(messages))])
		fmt.Printf(" \033[38;5;198m%.1f\033[0m\n", drain)
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

type shardling struct {
	attributes
}

func newShardling() entity {
	attr := attributes{
		name:     "Shardling",
		hp:       scale(60, 3.2),
		hpcap:    scale(60, 3.2),
		defense:  scale(10, 0.5),
		strength: scale(8, 1.32),
		effects:  make(map[string]int),
	}

	return &shardling{attr}
}

func (s *shardling) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 20 {
		fmt.Printf("Shardling rams itself onto you!")
		target.damage(s.strength)
	} else if roll < 40 {
		fmt.Printf("Shardling strikes with its crystal limbs!")
		target.damage(s.strength * 1.1)
	} else if roll < 80 {
		fmt.Printf("Shardling launched volley of shards!")
		dmg := rand.Float32()*10 + s.strength
		target.damage(dmg)
	} else {
		fmt.Printf("Shardling impales you with crystal spike!")
		target.damage(20 + s.strength)
	}
}

func (s *shardling) damage(n float32) {
	hp := s.hp
	s.attributes.damage(n)
	reflect := (hp - s.hp) * 0.3
	// quick & dirty way by accessing the global var.
	// i dont want to modify the interface just for this one.
	// and this cant self-target.
	fmt.Print("  reflected:")
	player.damage(reflect)
}

type genie struct {
	attributes
}

func newGenie() entity {
	attr := attributes{
		name:     "Evil genie",
		hp:       scale(100, 6),
		hpcap:    scale(100, 6),
		defense:  scale(2, 0.15),
		strength: scale(7, 0.85),
		effects:  make(map[string]int),
	}

	return &genie{attr}
}

func (g *genie) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	var attr *attributes

	if v, ok := target.(*Player); ok {
		attr = &v.attributes
	}

	if v, ok := target.(*genie); ok { //for self target when hit by "trick" player skill
		attr = &v.attributes
	}

	if roll < 5 {
		curse := 0.5 + rand.Float32()*4
		attr.hpcap = max(50, attr.hpcap-curse)

		fmt.Print("Genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("hp cap reduced by %.1f\n", curse)
	} else if roll < 10 {
		curse := 0.1 + rand.Float32()*1
		attr.strength = max(5, attr.strength-curse)

		fmt.Print("Genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("strength reduced by %.1f\n", curse)
	} else if roll < 15 {
		curse := 0.1 + rand.Float32()*0.8
		attr.defense = max(1, attr.defense-curse)

		fmt.Print("Genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("defense reduced by %.1f\n", curse)
	} else if roll < 18 && attr.name == "player" {
		p := target.(*Player)
		p.energycap = max(10, p.energycap-1)

		fmt.Print("Genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Println("energy cap reduced by 1")
	} else if roll < 31 {
		fmt.Println("Genie cast an \033[38;5;226millusion\033[0m!")
		fmt.Print("  self: ")
		target.attack(target)
	} else if roll < 53 {
		fmt.Printf("Genie cast a \033[38;5;226msandstorm\033[0m!")
		target.damage(g.strength*0.8 + scale(25, 1))
	} else if roll < 80 {
		fmt.Printf("Genie blast you with magical energy!")
		target.damage(g.strength*1.4 + 10)
	} else {
		fmt.Printf("Genie reach out for a punch!")
		target.damage(g.strength)
	}
}
