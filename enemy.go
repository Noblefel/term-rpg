package main

import (
	"fmt"
	"math/rand/v2"
)

type Enemy interface {
	Attack()
	Attr() Attributes
	TakeDamage(float32)
}

// for enemy
type Attributes struct {
	name     string
	hp       float32
	hpcap    float32
	strength float32
	defense  float32
}

func (attr *Attributes) TakeDamage(dmg float32) {
	dmg = max(dmg-attr.defense, 1)
	attr.hp -= dmg
	fmt.Fprintf(out, " \033[38;5;198m%.1f\033[0m damage\n", dmg)
}

func (attr Attributes) Attr() Attributes { return attr }

func SpawnRandom() Enemy {
	var enemies = []func() Enemy{
		NewKnight,
		NewWizard,
		NewChangeling,
		NewVampire,
		NewDemon,
	}

	spawn := enemies[rand.IntN(len(enemies))]
	return spawn()
}

type Knight struct {
	Attributes
}

func NewKnight() Enemy {
	attr := Attributes{
		name:     "Knight 🛡️",
		hp:       scale(90, 1.5),
		hpcap:    scale(90, 1.5),
		defense:  scale(6, 0.3),
		strength: scale(9, 0.5),
	}
	return &Knight{attr}
}

func (k *Knight) Attack() {
	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")

	if rand.IntN(100) < 15 {
		def := rand.Float32() * scale(5, 0.3)
		k.defense += def
		fmt.Fprint(out, "The knight reinforced his armor! ")
		fmt.Fprintf(out, "\033[38;5;83m+%.1f\033[0m defense\n", def)
		return
	}

	messages := []string{
		"The knight charged at you!",
		"The knight cuts you with his sword!",
		"The knight slashes you with his sword!",
		"The knight smashed you with his shield!",
	}
	fmt.Fprint(out, messages[rand.IntN(len(messages))])
	player.TakeDamage(k.strength)
}

type Wizard struct {
	Attributes
}

func NewWizard() Enemy {
	attr := Attributes{
		name:     "Wizard 🧙",
		hp:       scale(55, 0.75),
		hpcap:    scale(55, 0.75),
		defense:  scale(2, 0.15),
		strength: scale(4, 0.35),
	}
	return &Wizard{attr}
}

func (w *Wizard) Attack() {
	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")

	rng := rand.IntN(100)

	if rng < 20 {
		fmt.Fprint(out, "Wizard cannot cast, hit you with his staff instead! ")
		player.TakeDamage(w.strength)
	} else if rng < 40 {
		fmt.Fprint(out, "Wizard cast \033[38;5;226mfireball\033[0m!")
		player.TakeDamage(rand.Float32() * scale(30, 1.75))
	} else if rng < 60 {
		fmt.Fprint(out, "Wizard cast \033[38;5;226mlightning\033[0m!")
		player.TakeDamage(10 + rand.Float32()*scale(20, 1))
	} else if rng < 80 {
		fmt.Fprint(out, "Wizard summons \033[38;5;226mmeteor\033[0m!")
		player.TakeDamage(35)
	} else {
		heal := w.hpcap * 0.20
		w.hp = min(w.hp+heal, w.hpcap)
		fmt.Fprint(out, "Wizard cast healing! ")
		fmt.Fprintf(out, "recover \033[38;5;83m%.1f\033[0m hp\n", heal)
	}
}

type Changeling struct {
	Attributes
	mimic bool
}

func NewChangeling() Enemy {
	attr := Attributes{
		name:     "Changeling 🎭",
		hp:       80,
		hpcap:    80,
		defense:  1,
		strength: 1,
	}
	return &Changeling{attr, false}
}

func (c *Changeling) Attack() {
	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")

	if !c.mimic {
		c.hp = player.hp
		c.hpcap = player.hpcap
		c.strength = player.strength
		c.defense = player.defense
		c.mimic = true
		fmt.Fprintln(out, "Changeling \033[38;5;226mmorph\033[0m itself to look like you!")
		return
	}

	fmt.Fprint(out, "Changeling attacked! ")
	player.TakeDamage(c.strength)
}

type Vampire struct {
	Attributes
}

func NewVampire() Enemy {
	attr := Attributes{
		name:     "Vampire 🧛",
		hp:       scale(80, 1.2),
		hpcap:    scale(80, 1.2),
		defense:  scale(4, 0.15),
		strength: scale(10, 0.9),
	}
	return &Vampire{attr}
}

func (v *Vampire) Attack() {
	rng := rand.IntN(100)

	if rng < 15 {
		fmt.Fprint(out, "\033[38;5;196m✘\033[0m ")
		fmt.Fprint(out, "Vampire get exposed to sunlight! getting")
		v.TakeDamage(scale(10, 0.9))
		return
	}

	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")

	if rng < 70 {
		temp := player.hp
		fmt.Fprint(out, "Vampire bites down hard! healing and dealing")
		player.TakeDamage(v.strength)
		v.hp = min(v.hp+temp-player.hp, v.hpcap)
	} else if rng < 85 {
		fmt.Fprint(out, "Vampire slashed you with claws!")
		player.TakeDamage(v.strength)
	} else {
		fmt.Fprint(out, "Vampire summoned then swarm you with bats!")
		player.TakeDamage(v.strength + rand.Float32()*v.strength)
	}
}

type Demon struct {
	Attributes
}

func NewDemon() Enemy {
	attr := Attributes{
		name:     "Demon 👹",
		hp:       scale(112, 1),
		hpcap:    scale(112, 1),
		defense:  scale(4, 0.2),
		strength: scale(10, 0.3),
	}
	return &Demon{attr}
}

func (d *Demon) Attack() {
	fmt.Fprint(out, "\033[38;5;83m✔\033[0m ")

	if rand.IntN(100) < 60 {
		messages := []string{
			"Demon drains your life force!",
			"Demon absors your life force!",
			"Demon reach out and take good portion of your hp!",
			"Demon conjure dark energy that weakens you whole!",
		}

		drain := d.strength + player.hp*0.05
		player.hp = max(player.hp-drain, 0) //to ignore defense
		fmt.Fprint(out, messages[rand.IntN(len(messages))])
		fmt.Fprintf(out, " \033[38;5;198m%.1f\033[0m damage\n", drain)
		return
	}

	messages := []string{
		"Demon charged at you!",
		"Demon thrusts its trident at you!",
		"Demon swung its trident and slashes you!",
		"Demon lunged before punching you in the gut!",
	}

	fmt.Fprint(out, messages[rand.IntN(len(messages))])
	player.TakeDamage(d.strength)
}
