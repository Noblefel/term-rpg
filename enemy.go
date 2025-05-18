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
		newShardling,
		newGenie,
		newCelestial,
		newShapeshift,
		newUndead,
	}

	spawn := enemies[rand.IntN(len(enemies))]
	return spawn()
}

type knight struct {
	attributes
}

func newKnight() entity {
	attr := attributes{
		name:     "knight",
		hp:       scale(90, 14),
		hpcap:    scale(90, 14),
		defense:  scale(6, 0.3),
		strength: scale(9, 1),
		effects:  make(map[string]int),
	}
	return &knight{attr}
}

func (k *knight) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 15 {
		def := 0.1 + rand.Float32()*(k.defense/10)
		k.defense += def
		k.hp = min(k.hp+5+def, k.hpcap)

		fmt.Print("the knight reinforced his armor! ")
		fmt.Printf("\033[38;5;83m+%.1f\033[0m defense\n", def)
	} else if roll < 30 {
		fmt.Println("the knight \033[38;5;226mstrengthen\033[0m himself!")
		k.attackWith(target, k.strength)
		k.effects["strengthen"] = 4
	} else {
		messages := []string{
			"the knight charged at you!",
			"the knight cuts you with his sword!",
			"the knight slashes you with his sword!",
			"the knight smashed you with his shield!",
		}

		fmt.Printf(messages[rand.IntN(len(messages))])
		k.attackWith(target, k.strength)
	}
}

type wizard struct {
	attributes
}

func newWizard() entity {
	attr := attributes{
		name:     "wizard",
		hp:       scale(55, 6),
		hpcap:    scale(55, 6),
		defense:  scale(1, 0.1),
		strength: scale(4, 0.65),
		effects:  make(map[string]int),
	}
	return &wizard{attr}
}

func (w *wizard) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 15 {
		heal := w.hpcap * 0.20
		w.hp = min(w.hp+heal, w.hpcap)
		fmt.Printf("wizard cast healing! ")
		fmt.Printf("recover \033[38;5;83m%.1f\033[0m hp\n", heal)
	} else if roll < 20 {
		w.effects["immunity"] = 2
		fmt.Println("wizard cast \033[38;5;226mimmunity\033[0m!")
	} else if roll < 30 {
		w.effects["barrier"] = 4
		fmt.Println("wizard cast magical \033[38;5;226mbarrier\033[0m!")
	} else if roll < 35 {
		target.attr().effects["confused"] = 5
		fmt.Println("wizard cast \033[38;5;226mconfusion\033[0m!")
	} else if roll < 50 {
		fmt.Printf("wizard cast \033[38;5;226menhanced\033[0m attack!")
		w.attackWith(target, w.strength*3)
	} else if roll < 70 {
		fmt.Printf("wizard cast \033[38;5;226mfireball\033[0m!")
		dmg := rand.Float32() * scale(30, 1.75)
		w.attackWith(target, dmg)
	} else if roll < 85 {
		fmt.Printf("wizard cast \033[38;5;226mlightning\033[0m!")
		dmg := 10 + rand.Float32()*scale(20, 1)
		w.attackWith(target, dmg)
	} else if roll < 90 {
		fmt.Printf("wizard summons \033[38;5;226mmeteor\033[0m!")
		w.attackWith(target, 45)
	} else {
		fmt.Printf("wizard cannot cast, he use his staff instead!")
		w.attackWith(target, w.strength)
	}
}

type changeling struct {
	attributes
	mimic bool
}

func newChangeling() entity {
	attr := attributes{
		name:     "changeling",
		hp:       80,
		hpcap:    80,
		defense:  1,
		strength: 1,
		effects:  make(map[string]int),
	}
	return &changeling{attr, false}
}

func (c *changeling) attack(target entity) {
	if !c.mimic {
		fmt.Print(success)
		c.hp = target.attr().hp * 0.75
		c.hpcap = target.attr().hpcap * 0.75
		c.strength = target.attr().strength * 0.75
		c.defense = target.attr().defense * 0.75
		c.mimic = true
		fmt.Println("changeling attempt to \033[38;5;226mmorph\033[0m itself!")
		return
	}

	c.attributes.attack(target)
}

type vampire struct {
	attributes
}

func newVampire() entity {
	attr := attributes{
		name:     "vampire",
		hp:       scale(80, 10),
		hpcap:    scale(80, 10),
		defense:  scale(3, 0.13),
		strength: scale(10, 1.9),
		effects:  make(map[string]int),
	}
	return &vampire{attr}
}

func (v *vampire) attack(target entity) {
	roll := roll()

	if roll < 8 {
		fmt.Print(fail)
		fmt.Printf("vampire get exposed to sunlight! getting")
		v.damage(v.hp * 0.13)
		return
	}

	fmt.Print(success)

	if roll < 60 {
		oldhp := target.attr().hp
		fmt.Printf("vampire bites down hard!")
		dmg := v.strength + oldhp*0.01
		v.attackWith(target, dmg)

		heal := (oldhp-target.attr().hp)/2 + 4
		v.hp = min(v.hp+heal, v.hpcap)
	} else if roll < 70 {
		fmt.Printf("vampire bites down, giving \033[38;5;226mpoison\033[0m!")
		target.attr().effects["poisoned"] = 3
		dmg := v.strength + target.attr().hp*0.01
		v.attackWith(target, dmg)
	} else if roll < 85 {
		fmt.Printf("vampire summoned then swarm you with bats!")
		v.attackWith(target, v.strength*1.2)
	} else {
		fmt.Printf("vampire slashed you with claws!")
		v.attackWith(target, v.strength)
	}
}

type demon struct {
	attributes
}

func newDemon() entity {
	attr := attributes{
		name:     "demon",
		hp:       scale(112, 16),
		hpcap:    scale(112, 16),
		defense:  scale(4, 0.28),
		strength: scale(10, 1.8),
		effects:  make(map[string]int),
	}
	return &demon{attr}
}

func (d *demon) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 60 {
		messages := []string{
			"demon \033[38;5;226mdrains\033[0m the target life force!",
			"demon \033[38;5;226mabsorbs\033[0m the target life force!",
			"demon conjure \033[38;5;226mlife draining\033[0m magic!",
		}

		drain := d.strength + target.attr().hp*0.07
		hp := max(target.attr().hp-drain, 0)
		target.setHP(hp)
		fmt.Printf(messages[rand.IntN(len(messages))])
		fmt.Printf(" \033[38;5;198m%.1f\033[0m\n", drain)
		return
	}

	if roll < 75 {
		fmt.Println("demon draw upon the power of \033[38;5;226mhell fire\033[0m!")
		target.attr().effects["burning"] = 3
		d.attackWith(target, d.strength*0.8)
		return
	}

	messages := []string{
		"demon charged at you!",
		"demon thrusts its trident at you!",
		"demon swung its trident and slashes you!",
		"demon lunged before punching you in the gut!",
	}

	fmt.Printf(messages[rand.IntN(len(messages))])
	d.attackWith(target, d.strength)
}

type shardling struct {
	attributes
}

func newShardling() entity {
	attr := attributes{
		name:     "shardling",
		hp:       scale(60, 6.4),
		hpcap:    scale(60, 6.4),
		defense:  scale(11, 0.4),
		strength: scale(8, 1.12),
		effects:  make(map[string]int),
	}

	return &shardling{attr}
}

func (s *shardling) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 20 {
		fmt.Printf("shardling rams itself onto you!")
		s.attackWith(target, s.strength)
	} else if roll < 40 {
		fmt.Printf("shardling strikes with its crystal limbs!")
		s.attackWith(target, s.strength*1.1)
	} else if roll < 80 {
		fmt.Printf("shardling launched volley of shards!")
		dmg := rand.Float32()*10 + s.strength
		s.attackWith(target, dmg)
	} else {
		fmt.Printf("shardling impales you with crystal spike!")
		s.attackWith(target, s.strength+20)
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
		name:     "genie",
		hp:       scale(100, 12),
		hpcap:    scale(100, 12),
		defense:  scale(2, 0.14),
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

		fmt.Print("genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("hp cap reduced by %.1f\n", curse)
	} else if roll < 10 {
		curse := 0.1 + rand.Float32()*1
		attr.strength = max(5, attr.strength-curse)

		fmt.Print("genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("strength reduced by %.1f\n", curse)
	} else if roll < 15 {
		curse := 0.1 + rand.Float32()*0.8
		attr.defense = max(1, attr.defense-curse)

		fmt.Print("genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("defense reduced by %.1f\n", curse)
	} else if roll < 18 && attr.name == "player" {
		p := target.(*Player)
		p.energycap = max(10, p.energycap-1)

		fmt.Print("genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Println("energy cap reduced by 1")
	} else if roll < 31 {
		fmt.Println("genie cast an \033[38;5;226millusion\033[0m!")
		fmt.Print("  self: ")
		target.attack(target)
	} else if roll < 38 {
		fmt.Println("genie cast \033[38;5;226mforce-field\033[0m! reduce incoming damage")
		g.effects["force-field"] = 5
	} else if roll < 60 {
		fmt.Printf("genie cast a \033[38;5;226msandstorm\033[0m!")
		dmg := g.strength*0.5 + scale(25, 0.3)
		g.attackWith(target, dmg)
	} else if roll < 85 {
		fmt.Printf("genie blast you with magical energy!")
		dmg := g.strength*0.8 + 15
		g.attackWith(target, dmg)
	} else {
		fmt.Printf("genie reach out for a punch!")
		g.attackWith(target, g.strength)
	}
}

type celestial struct {
	attributes
}

func newCelestial() entity {
	attr := attributes{
		name:     "the celestial being",
		hp:       scale(140, 20),
		hpcap:    scale(140, 20),
		defense:  0,
		strength: scale(9, 1),
		effects:  make(map[string]int),
	}

	return &celestial{attr}
}

func (c *celestial) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 12 {
		fmt.Println("the celestial being channels a \033[38;5;226mhealing aura\033[0m!")
		c.effects["heal aura"] = 5
	} else if roll < 37 {
		messages := []string{
			"the celestial being call upon the power of the \033[38;5;226mholy fire\033[0m!",
			"the celestial being draw upon the power of the \033[38;5;226msun\033[0m!",
			"the celestial being channels energy from the \033[38;5;226mstars\033[0m!",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		target.attr().effects["burning"] = 3
		c.attackWith(target, c.strength*0.8)
	} else {
		messages := []string{
			"the celestial being strikes you with holy blade",
			"the celestial being strikes you with a beam of light",
			"the celestial being swung its holy blade",
			"the celestial being knocked you with a burst of wind",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		c.attackWith(target, c.strength)
	}
}

type shapeshift struct {
	attributes
	mimic bool
}

func newShapeshift() entity {
	attr := attributes{
		name:     "shapeshift",
		hp:       80,
		hpcap:    80,
		defense:  1,
		strength: 1,
		effects:  make(map[string]int),
	}
	return &shapeshift{attr, false}
}

// future plan: morph into other enemies during battle
func (s *shapeshift) attack(target entity) {
	if !s.mimic {
		fmt.Print(success)
		s.hp = target.attr().hp * 1.5
		s.hpcap = target.attr().hpcap * 1.5
		s.strength = target.attr().strength * 1.5
		s.defense = target.attr().defense * 1.5
		s.mimic = true
		fmt.Println("shapeshift \033[38;5;226mmorph\033[0m itself!")
		return
	}

	s.attributes.attack(target)
}

type undead struct {
	attributes
}

func newUndead() entity {
	attr := attributes{
		name:     "undead",
		hp:       scale(78, 9),
		hpcap:    scale(78, 9),
		defense:  scale(2, 0.15),
		strength: scale(9, 1),
		effects:  make(map[string]int),
	}
	return &undead{attr}
}

func (u *undead) attack(target entity) {
	fmt.Print(success)
	roll := roll()
	strength := u.strength

	// penalty against deadman perk
	if v, ok := target.(*Player); ok && v.perk == 5 {
		strength -= strength * 0.33
	}

	if roll < 8 {
		fmt.Print("the undead vomits a stream of \033[38;5;226macidic bile\033[0m!")
		target.attr().effects["poisoned"] = 3
		u.attackWith(target, strength*0.4)
	} else if roll < 28 {
		fmt.Print("the undead calls in fellow undead from underground!")
		dmg := 5 + rand.Float32()*scale(5, 1)
		u.attackWith(target, dmg)
	} else if roll < 35 {
		fmt.Print("the undead bites down on the target's leg!")
		u.attackWith(target, strength*1.33)
	} else {
		messages := []string{
			"the undead rotting fist connects with its target!",
			"the undead reach out with a clumsy swing!",
			"the undead charges itself limply!",
			"a bony finger jabs into the target side!",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		u.attackWith(target, strength)
	}
}
