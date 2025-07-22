package main

import (
	"fmt"
	"math/rand/v2"
)

func spawn() entity {
	var enemies = []func() entity{
		spawnKnight,
		spawnWizard,
		spawnChangeling,
		spawnVampire,
		spawnDemon,
		spawnShardling,
		spawnGenie,
		spawnCelestial,
		spawnShapeshift,
		spawnUndead,
		spawnScorpion,
		spawnGoblin,
		spawnInfernal,
		spawnVineMonster,
		spawnArctic,
		spawnChieftain,
		spawnGiantLeech,
		spawnPirate,
		spawnGladiator,
		spawnSorcerer,
	}

	spawn := enemies[rand.IntN(len(enemies))]
	return spawn()
}

type knight struct {
	attributes
}

func spawnKnight() entity {
	attr := attributes{
		name:     "knight",
		hp:       scale(210, 25),
		hpcap:    scale(210, 25),
		defense:  scale(23, 2.4),
		strength: scale(40, 3.1),
		agility:  scale(7, 0.18),
		effects:  make(map[string]int),
	}
	return &knight{attr}
}

func (k *knight) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 15 {
		def := 0.5 + rand.Float64()*(k.defense/10)
		k.defense += def
		k.hp = min(k.hp+def+20, k.hpcap)

		fmt.Print("the knight reinforced his armor! ")
		fmt.Printf("\033[38;5;83m+%.1f\033[0m defense\n", def)
	} else if roll < 30 {
		fmt.Print("the knight \033[38;5;226mstrengthen\033[0m himself!")
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

func spawnWizard() entity {
	attr := attributes{
		name:     "wizard",
		hp:       scale(132, 11.5),
		hpcap:    scale(132, 11.5),
		defense:  scale(12, 0.75),
		strength: scale(25, 2.66),
		agility:  scale(3, 0.1),
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
		fmt.Print("wizard cast \033[38;5;226mconfusion\033[0m!")
		target.damage(w.strength * 0.6)
	} else if roll < 50 {
		fmt.Printf("wizard cast \033[38;5;226menhanced\033[0m attack!")
		w.attackWith(target, w.strength*3)
	} else if roll < 70 {
		fmt.Printf("wizard cast \033[38;5;226mfireball\033[0m!")
		dmg := rand.Float64() * scale(100, 4)
		w.attackWith(target, dmg)
	} else if roll < 85 {
		fmt.Printf("wizard cast \033[38;5;226mlightning\033[0m!")
		dmg := 40 + rand.Float64()*scale(75, 2)
		w.attackWith(target, dmg)
	} else if roll < 90 {
		fmt.Printf("wizard summons \033[38;5;226mmeteor\033[0m!")
		w.attackWith(target, 120)
	} else {
		fmt.Printf("wizard cannot cast, he use his staff instead!")
		w.attackWith(target, w.strength)
	}
}

type changeling struct {
	attributes
	mimic bool
}

func spawnChangeling() entity {
	attr := attributes{
		name:     "changeling",
		hp:       100,
		hpcap:    100,
		defense:  scale(100, 1), // prevent instant kill before morph
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

func spawnVampire() entity {
	attr := attributes{
		name:     "vampire",
		hp:       scale(180, 21.4),
		hpcap:    scale(180, 21.4),
		defense:  scale(14, 1.9),
		strength: scale(43, 3.2),
		agility:  scale(10, 0.29),
		effects:  make(map[string]int),
	}
	return &vampire{attr}
}

func (v *vampire) attack(target entity) {
	roll := roll()

	if roll < 8 {
		fmt.Print(fail)
		fmt.Printf("vampire get exposed to sunlight! getting")
		v.damage(v.defense + v.hp*0.1)
		return
	}

	fmt.Print(success)

	if roll < 60 {
		oldhp := target.attr().hp
		fmt.Printf("vampire bites down hard!")
		dmg := v.strength + oldhp*0.01
		v.attackWith(target, dmg)

		heal := (oldhp-target.attr().hp)/2 + 15
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

func spawnDemon() entity {
	attr := attributes{
		name:     "demon",
		hp:       scale(212, 31),
		hpcap:    scale(212, 31),
		defense:  scale(20, 2.5),
		strength: scale(45, 3.8),
		agility:  scale(8.5, 0.27),
		effects:  make(map[string]int),
	}
	return &demon{attr}
}

func (d *demon) attack(target entity) {
	fmt.Print(success)
	roll := roll()
	ignored := target.attr().defense / 2

	if roll < 60 {
		messages := []string{
			"demon \033[38;5;226mdrains\033[0m the target life force!",
			"demon \033[38;5;226mabsorbs\033[0m the target life force!",
			"demon conjure \033[38;5;226mlife draining\033[0m magic!",
		}

		fmt.Printf(messages[rand.IntN(len(messages))])
		dmg := d.strength + ignored
		dmg += target.attr().hp * 0.03
		d.attackWith(target, dmg)
	} else if roll < 75 {
		fmt.Print("demon draw upon the power of \033[38;5;226mhell fire\033[0m!")
		target.attr().effects["burning"] = 3
		dmg := d.strength*0.8 + ignored
		d.attackWith(target, dmg)
	} else {
		messages := []string{
			"demon charged at you!",
			"demon thrusts its trident at you!",
			"demon swung its trident and slashes you!",
			"demon lunged before punching you in the gut!",
		}

		fmt.Printf(messages[rand.IntN(len(messages))])
		d.attackWith(target, d.strength+ignored)
	}
}

type shardling struct {
	attributes
}

func spawnShardling() entity {
	attr := attributes{
		name:     "shardling",
		hp:       scale(100, 10),
		hpcap:    scale(100, 10),
		defense:  scale(30, 2.65),
		strength: scale(27, 2.42),
		agility:  scale(6, 0.18),
		effects:  make(map[string]int),
	}

	attr.effects["reflect"] = 99
	return &shardling{attr}
}

func (s *shardling) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 40 {
		fmt.Printf("shardling rams itself onto you!")
		s.attackWith(target, s.strength)
	} else if roll < 50 {
		fmt.Printf("shardling strikes with its crystal limbs!")
		s.attackWith(target, s.strength*1.1)
	} else if roll < 80 {
		fmt.Printf("shardling launched volley of shards!")
		dmg := rand.Float64()*20 + s.strength
		s.attackWith(target, dmg)
	} else {
		fmt.Printf("shardling impales you with crystal spike!")
		s.attackWith(target, s.strength*1.25)
	}
}

type genie struct {
	attributes
}

func spawnGenie() entity {
	attr := attributes{
		name:     "genie",
		hp:       scale(200, 20),
		hpcap:    scale(200, 20),
		defense:  scale(17, 2.22),
		strength: scale(36, 3.13),
		agility:  scale(8, 0.275),
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
		curse := 1 + rand.Float64()*10
		attr.hpcap = max(50, attr.hpcap-curse)

		fmt.Print("genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("hp cap reduced by %.1f\n", curse)
	} else if roll < 10 {
		curse := 0.25 + rand.Float64()*1.5
		attr.strength = max(5, attr.strength-curse)

		fmt.Print("genie placed a \033[38;5;226mcurse mark\033[0m! ")
		fmt.Printf("strength reduced by %.1f\n", curse)
	} else if roll < 15 {
		curse := 0.1 + rand.Float64()*1
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
		dmg := g.strength*0.5 + 40
		g.attackWith(target, dmg)
	} else if roll < 85 {
		fmt.Printf("genie blast you with magical energy!")
		dmg := g.strength * 1.13
		g.attackWith(target, dmg)
	} else {
		fmt.Printf("genie reach out for a punch!")
		g.attackWith(target, g.strength)
	}
}

type celestial struct {
	attributes
}

func spawnCelestial() entity {
	attr := attributes{
		name:     "celestial being",
		hp:       scale(300, 40),
		hpcap:    scale(300, 40),
		defense:  scale(7.5, 0.5),
		strength: scale(41, 3),
		agility:  scale(7, 0.28),
		effects:  make(map[string]int),
	}

	return &celestial{attr}
}

func (c *celestial) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 10 {
		fmt.Println("the celestial being channels a \033[38;5;226mhealing aura\033[0m!")
		c.effects["heal aura"] = 5
	} else if roll < 30 {
		fmt.Print("the celestial being cast a \033[38;5;226mblinding light\033[0m!")
		target.attr().effects["disoriented"] = 3
		c.attackWith(target, c.strength*0.7)
	} else if roll < 40 {
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

func spawnShapeshift() entity {
	attr := attributes{
		name:     "shapeshift",
		hp:       80,
		hpcap:    80,
		defense:  scale(100, 1), // prevent instant kill before morph
		strength: 1,
		effects:  make(map[string]int),
	}
	return &shapeshift{attr, false}
}

// future plan: morph into other enemies during battle
func (s *shapeshift) attack(target entity) {
	if !s.mimic {
		fmt.Print(success)
		s.hp = target.attr().hp * 1.25
		s.hpcap = target.attr().hpcap * 1.25
		s.strength = target.attr().strength * 1.25
		s.defense = target.attr().defense * 1.25
		s.mimic = true
		fmt.Println("shapeshift \033[38;5;226mmorph\033[0m itself!")
		return
	}

	s.attributes.attack(target)
}

type undead struct {
	attributes
}

func spawnUndead() entity {
	attr := attributes{
		name:     "undead",
		hp:       scale(181, 20.8),
		hpcap:    scale(181, 20.8),
		defense:  scale(16.7, 2.3),
		strength: scale(38, 3.125),
		agility:  scale(5, 0.25),
		effects:  make(map[string]int),
	}
	return &undead{attr}
}

func (u *undead) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 8 {
		fmt.Print("the undead vomits a stream of \033[38;5;226macidic bile\033[0m!")
		target.attr().effects["poisoned"] = 3
		u.attackWith(target, u.strength*0.4)
	} else if roll < 28 {
		fmt.Print("the undead calls in fellow undead from underground!")
		dmg := u.strength + rand.Float64()*scale(30, 4)
		u.attackWith(target, dmg)
	} else if roll < 35 {
		fmt.Print("the undead bites down on the target's leg!")
		u.attackWith(target, u.strength*1.33)
	} else {
		messages := []string{
			"the undead rotting fist connects with its target!",
			"the undead reach out with a clumsy swing!",
			"the undead charges itself limply!",
			"a bony finger jabs into the target side!",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		u.attackWith(target, u.strength)
	}
}

type scorpion struct {
	attributes
}

func spawnScorpion() entity {
	attr := attributes{
		name:     "scorpion",
		hp:       scale(150, 12),
		hpcap:    scale(150, 12),
		defense:  scale(15, 1.95),
		strength: scale(49, 3.75),
		agility:  scale(9.2, 0.29),
		effects:  make(map[string]int),
	}
	return &scorpion{attr}
}

func (s *scorpion) attack(target entity) {
	fmt.Print(success)
	roll := roll()
	// ignore 30% def
	strength := s.strength + target.attr().defense*0.3

	if roll < 24 {
		messages := []string{
			"the scorpion charges forward and use its \033[38;5;226mstinger\033[0m",
			"the scorpion stungs hard, injecting deadly \033[38;5;226mvenom\033[0m",
			"the scorpion reach out to inject \033[38;5;226mvenom\033[0m",
			"the scorpion emits \033[38;5;226mpoisonous\033[0m gas",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		target.attr().effects["poisoned severe"] = 3
		s.attackWith(target, strength)
		return
	}

	messages := []string{
		"the scorpion lashes out with its tail",
		"the scorpion impales you from underground",
		"the scorpion pinches you with its massive claws",
	}

	fmt.Print(messages[rand.IntN(len(messages))])
	s.attackWith(target, strength)
}

type goblin struct {
	attributes
}

func spawnGoblin() entity {
	attr := attributes{
		name:     "goblin",
		hp:       scale(130, 11),
		hpcap:    scale(130, 11),
		defense:  scale(10, 1),
		strength: scale(38, 3.1),
		agility:  scale(25, 0.4),
		effects:  make(map[string]int),
	}
	return &goblin{attr}
}

func (g *goblin) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 12 {
		fmt.Print("the goblin blows a puff of \033[38;5;226mmysterious powder\033[0m!")
		target.attr().effects["confused"] = 3
		g.attackWith(target, g.strength*0.45)
	} else if roll < 24 {
		fmt.Print("the goblin leaps into the air and bring its club down!")
		g.attackWith(target, g.strength*1.25)
	} else if roll < 36 {
		fmt.Print("the goblin strikes rapidly!")
		g.attackWith(target, g.strength*1.1)
	} else {
		messages := []string{
			"the goblin stabs with a crude spear",
			"the goblin hurls shards of broken glass",
			"the goblin launches itself wildly",
			"the goblin swung its crude club",
			"the goblin's rusty dagger reach out for a slash",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		g.attackWith(target, g.strength)
	}
}

type infernal struct {
	attributes
}

func spawnInfernal() entity {
	attr := attributes{
		name:     "infernal",
		hp:       scale(163, 19),
		hpcap:    scale(163, 19),
		defense:  scale(12, 0.75),
		strength: scale(35, 2.8),
		agility:  scale(4, 0.11),
		effects:  make(map[string]int),
	}
	attr.effects["burning immunity"] = 99
	return &infernal{attr}
}

func (i *infernal) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 30 {
		messages := []string{
			"the infernal drew upon the power of the \033[38;5;226mdeepest fire\033[0m!",
			"the infernal summons a rain of \033[38;5;226mdeadly flame\033[0m!",
			"the infernal unleash a wave of \033[38;5;226mblack flame\033[0m!",
			"the infernal turned the ground into \033[38;5;226mmolten lava\033[0m!",
		}

		target.attr().effects["burning severe"] = 2
		fmt.Print(messages[rand.IntN(len(messages))])
		i.attackWith(target, i.strength)
	} else if roll < 40 {
		str := 1 + rand.Float64()*(i.strength/10)
		i.strength += str
		fmt.Printf("the infernal \033[38;5;226menhanced\033[0m itself with molten energy, +%.1f strength!\n", str)
	} else if roll < 50 {
		def := 4 + rand.Float64()*(i.defense/10)
		i.defense += def
		fmt.Printf("the infernal uses \033[38;5;226mmagma shield\033[0m, +%.1f defense!\n", def)
	} else {
		messages := []string{
			"the infernal sends a blast of flames",
			"the infernal's burning fist reaches out!",
			"the infernal surround the area with searing heat!",
			"fiery claws of the inferno tears through the flesh!",
		}

		target.attr().effects["burning"] = 2
		fmt.Print(messages[rand.IntN(len(messages))])
		i.attackWith(target, i.strength)
	}
}

type vineMonster struct {
	attributes
}

func spawnVineMonster() entity {
	attr := attributes{
		name:     "vine monster",
		hp:       scale(200, 23),
		hpcap:    scale(200, 23),
		defense:  scale(25, 2.34),
		strength: scale(33, 2.74),
		agility:  scale(7, 0.25),
		effects:  make(map[string]int),
	}
	attr.effects["reflect low"] = 99
	return &vineMonster{attr}
}

func (v *vineMonster) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 30 {
		messages := []string{
			"The vine monster ensnares you!",
			"The vine monster unleash a swarm of roots!",
			"The vine monster roots wrap around your legs!",
			"The vine monster locks you in place with thick vines!",
		}

		target.attr().effects["stunned"] = 2
		fmt.Print(messages[rand.IntN(len(messages))])
		v.attackWith(target, v.strength*0.6)
		return
	}

	if roll < 45 {
		fmt.Print("The vine monster erupted a burst of \033[38;5;226mthorns\033[0m!")
		v.attackWith(target, v.strength*0.75+rand.Float64()*33)
		return
	}

	messages := []string{
		"The vine monster slams you to the ground!",
		"The vine monster slams you with massive tree hand!",
		"The vine monster impales you with sharp roots!",
		"The vine monster whips you with thick vines!",
		"The vine monster whips you with sharp roots!",
	}

	fmt.Print(messages[rand.IntN(len(messages))])
	v.attackWith(target, v.strength)
}

type arctic struct {
	attributes
}

func spawnArctic() entity {
	attr := attributes{
		name:     "arctic warrior",
		hp:       scale(185, 22),
		hpcap:    scale(185, 22),
		defense:  scale(20, 2.35),
		strength: scale(40, 3.05),
		agility:  scale(10, 0.27),
		effects:  make(map[string]int),
	}
	attr.effects["frozen immunity"] = 99
	return &arctic{attr}
}

func (a *arctic) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 20 {
		messages := []string{
			"arctic warrior summons a wave of \033[38;5;226mfrost\033[0m!",
			"arctic warrior conjure a beam of \033[38;5;226mice\033[0m!",
			"arctic warrior calls upon the power of the \033[38;5;226mcold\033[0m!",
		}

		target.attr().effects["frozen"] = 2
		fmt.Println(messages[rand.IntN(len(messages))])
	} else if roll < 35 {
		fmt.Print("arctic warrior cast \033[38;5;226msnowstorm\033[0m!")
		dmg := a.strength*0.2 + rand.Float64()*60
		a.attackWith(target, dmg)
	} else if roll < 43 {
		fmt.Print("arctic warrior cast \033[38;5;226mavalanche\033[0m!")
		dmg := a.strength + rand.Float64()*scale(20, 1)
		a.attackWith(target, dmg)
	} else {
		messages := []string{
			"arctic warrior unleash a sweeping arc of ice!",
			"arctic warrior slashes you with icy sword!",
			"arctic warrior reach out for a punch!",
			"arctic warrior hurls icy spikes!",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		a.attackWith(target, a.strength)
	}
}

type chieftain struct {
	attributes
}

func spawnChieftain() entity {
	attr := attributes{
		name:     "chieftain",
		hp:       scale(185, 22),
		hpcap:    scale(185, 22),
		defense:  scale(18, 2.2),
		strength: scale(42, 3.15),
		agility:  scale(11, 0.275),
		effects:  make(map[string]int),
	}
	attr.effects["poison immunity"] = 99
	return &chieftain{attr}
}

func (c *chieftain) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 14 {
		fmt.Print("chieftain strike \033[38;5;226mferociously\033[0m!")
		target.attr().effects["shiver"] = 3
		c.attackWith(target, c.strength)
	} else if roll < 28 {
		fmt.Print("chieftain calls upon the help of \033[38;5;226mancestral spirits\033[0m!")
		c.effects["vitality"] = 4
		c.attackWith(target, c.strength*0.8)
	} else if roll < 35 {
		fmt.Print("chieftain let out \033[38;5;226mancestral roar\033[0m as he charge!")
		target.attr().effects["shiver"] = 2
		c.effects["vitality"] = 2
		c.attackWith(target, c.strength*1.3)
	} else if roll < 45 {
		fmt.Print("chieftain leap from a tree and slam onto you!")
		c.attackWith(target, c.strength*1.2)
	} else {
		messages := []string{
			"chieftain landed a blow with his primal club!",
			"chieftain thrust his primal spear!",
			"chieftain swung his primal club!",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		c.attackWith(target, c.strength)
	}
}

type giantLeech struct {
	attributes
}

func spawnGiantLeech() entity {
	attr := attributes{
		name:     "giant leech",
		hp:       scale(158, 14),
		hpcap:    scale(158, 14),
		defense:  scale(14, 1.8),
		strength: scale(36, 2.82),
		agility:  scale(7, 0.25),
		effects:  make(map[string]int),
	}
	return &giantLeech{attr}
}

func (lm *giantLeech) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if roll < 25 {
		messages := []string{
			"giant leech \033[38;5;226mbites\033[0m down hard!",
			"giant leech \033[38;5;226mbites\033[0m tear through your flesh",
			"giant leech draining attack left you \033[38;5;226mwounded\033[0m!",
		}

		target.attr().effects["bleeding"] += 10
		drain := target.attr().hp*0.1 + lm.strength
		drain -= target.attr().defense * 0.75 // target defense value/effectiveness +75%
		fmt.Print(messages[rand.IntN(len(messages))])
		lm.attackWith(target, drain)
		return
	}

	messages := []string{
		"giant leech drains your life blood",
		"giant leech viciously drain large amount of blood",
		"giant leech latches on before absorbing your blood",
		"giant leech lacthes on and suck large amount of blood",
	}

	drain := target.attr().hp*0.2 + lm.strength
	drain -= target.attr().defense * 1.3 // target defense value/effectiveness +130%
	fmt.Print(messages[rand.IntN(len(messages))])
	lm.attackWith(target, drain)
}

type pirate struct {
	attributes
}

func spawnPirate() entity {
	attr := attributes{
		name:     "pirate",
		hp:       scale(160, 18),
		hpcap:    scale(160, 18),
		defense:  scale(14, 1.8),
		strength: scale(38, 3.05),
		agility:  scale(11, 0.277),
		effects:  make(map[string]int),
	}
	return &pirate{attr}
}

func (p *pirate) attack(target entity) {
	fmt.Print(success)
	roll := roll()

	if p, ok := target.(*Player); ok && roll < 30 {
		messages := []string{
			"in a swift move, the pirate yanks away your satchel!",
			"in a swift move, the pirate snatches your satchel!",
			"in a swift move, the pirate grabs your loot bag!",
			"pirate shoots a grappling hook onto your loot bag!",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		// 5 + 1-5% player's current gold
		steal := 5
		steal += p.gold * (1 + rand.IntN(5)) / 100
		fmt.Printf(" \033[38;5;226mstole\033[0m %d gold\n", steal)
		p.gold = max(p.gold-steal, 0)
	} else if roll < 60 {
		fmt.Print("the pirate whirls his scimitar in a savage arc!")
		p.effects["bleeding"] += 5
		p.attackWith(target, p.strength*1.1)
	} else {
		messages := []string{
			"the pirate lunges with a rusty cutlass",
			"the pirate slashes you with a scimitar",
			"the pirate shoots his flintlock at you",
			"the pirate heavy hook swings down",
		}

		fmt.Print(messages[rand.IntN(len(messages))])
		p.attackWith(target, p.strength)
	}
}

type gladiator struct {
	attributes
}

func spawnGladiator() entity {
	attr := attributes{
		name:     "gladiator",
		hp:       scale(200, 23.5),
		hpcap:    scale(200, 23.5),
		defense:  scale(19, 2.35),
		strength: scale(43, 3.2),
		agility:  scale(10, 0.26),
		effects:  make(map[string]int),
	}
	return &gladiator{attr}
}

func (g *gladiator) attack(target entity) {
	g._attack(target)
	g._attack(target)
}

func (g *gladiator) _attack(target entity) {
	fmt.Print(success)

	if roll() < 15 {
		fmt.Println("the gladiator raise its trident and \033[38;5;226mtaunted\033[0m you!")
		g.effects["strengthen"] = 2
		return
	}

	messages := []string{
		"the gladiator charges forward!",
		"the gladiator bashes with his shield!",
		"the gladiator swung its greatsword!",
		"the gladiator delivers a sweeping kick!",
		"the gladiator thrusts his trident!",
		"the gladiator brings down its mighty trident!",
		"the gladiator strikes from behind!",
		"the gladiator pivots and brings down its greatsword!",
		"the gladiator sidestepped and swung its trident!",
		"the gladiator lunges again for the second time!",
		"the gladiator stepped back, throwing its trident!",
	}

	fmt.Print(messages[rand.IntN(len(messages))])
	g.attackWith(target, g.strength*0.6)
}

type sorcerer struct {
	attributes
}

func spawnSorcerer() entity {
	attr := attributes{
		name:     "sorcerer",
		hp:       scale(140, 12),
		hpcap:    scale(140, 12),
		defense:  scale(12, 0.9),
		strength: scale(20, 2.5),
		agility:  scale(5, 0.13),
		effects:  make(map[string]int),
	}
	return &sorcerer{attr}
}

func (s *sorcerer) attack(target entity) {
	fmt.Print(success)
	// no variation attack because this enemy ignores everything anyway, why make it even broken

	messages := []string{
		"the sorcerer casts a ray of sickness",
		"the sorcerer taps directly into your soul!",
		"the sorcerer conjures the power of pure fire!",
		"the sorcerer hurls shards made of void energy!",
		"the sorcerer casts a barrage of astral attacks!",
		"the sorcerer crushes you with telekinetic force!",
		"the sorecrer sends a barrage of astral projectiles!",
		"the sorcerer sends a ripple of astral force down on you!",
	}

	fmt.Print(messages[rand.IntN(len(messages))])
	newhp := target.attr().hp - s.strength
	newhp = max(newhp, 0)
	target.setHP(newhp)
	fmt.Printf(" \033[38;5;198m%.1f\033[0m\n", s.strength)
}
