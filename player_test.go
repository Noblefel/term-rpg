package main

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer()

	if p.hpcap != 200 {
		t.Errorf("default hpcap should be 200, got %.1f", p.hpcap)
	}

	if p.defense != 15 {
		t.Errorf("default defense should be 15, got %.1f", p.defense)
	}

	if p.strength != 50 {
		t.Errorf("default strength should be 50, got %.1f", p.strength)
	}

	if p.agility != 10 {
		t.Errorf("default agility should be 10, got %.1f", p.agility)
	}

	if p.energycap != 20 {
		t.Errorf("default energycap should be 20, got %d", p.energycap)
	}
}

func TestPlayer_AttackWith(t *testing.T) {
	p := newTestPlayer()
	p.attack(p)

	if p.hp != p.hpcap-p.strength {
		t.Errorf("enemy hp should be reduced to %.1f", p.hpcap-p.strength)
	}

	t.Run("with fencer perk", func(t *testing.T) {
		p := newTestPlayer()
		p.perk = 12
		p.attack(p)

		dmg := p.strength * 0.55 * 2
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (55%% twice), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("with havoc perk", func(t *testing.T) {
		p := newTestPlayer()
		p.perk = 2
		p.attackWith(p, p.strength)

		dmg := p.strength * 1.2
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (20%% increase), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("with berserk perk and 30% hp", func(t *testing.T) {
		p1 := newTestPlayer()
		p2 := newTestPlayer()
		p1.hp = p1.hpcap * 0.3
		p1.perk = 3
		p1.attackWith(p2, p1.strength)

		dmg := p1.strength * 1.15
		if 100-p2.hp != dmg {
			t.Errorf("damage should be around %.1f (15%% increase), got %.1f", dmg, 100-p2.hp)
		}
	})

	t.Run("if low energy", func(t *testing.T) {
		p := newTestPlayer()
		p.energy = 0
		p.attackWith(p, p.strength)

		dmg := p.strength * 0.8
		if 100-p.hp != dmg {
			t.Errorf("damage should be around %.1f (20%% decrease), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("insanity perk", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		p.perk = 8
		p.attackWith(p, p.strength)
		mindmg := p.strength - p.strength*0.3
		maxdmg := p.strength + p.strength*0.3
		dmg := 100 - p.hp

		if dmg < mindmg || dmg > maxdmg {
			t.Errorf(
				"damage should be between %.1f to %.1f (-30%% to 30%% MULTIPLIER VAL), got %.1f",
				mindmg,
				maxdmg,
				dmg,
			)
		}

		rolltest = 50
		p.hp = 100
		p.attackWith(p, p.strength)
		mindmg = p.strength - 10
		maxdmg = p.strength + 10
		dmg = 100 - p.hp

		if dmg < mindmg || dmg > maxdmg {
			t.Errorf(
				"damage should be between %.1f to %.1f (-10 to +10 FLAT VAL), got %.1f",
				mindmg,
				maxdmg,
				dmg,
			)
		}
	})

	t.Run("frigid perk", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		p.perk = 10
		p.attackWith(p, 1)

		// if p.effects["frozen"] != 2 {
		// 	t.Errorf("should inflict frozen effect for 2 turns, got %d", p.effects["frozen"])
		// }
	})
}

func TestPlayer_Damage(t *testing.T) {
	p := newTestPlayer()
	p.defense = 5
	p.damage(10)

	if p.hp != 95 {
		t.Errorf("hp should be 95, got %.1f", p.hp)
	}

	t.Run("with resilient perk", func(t *testing.T) {
		p := newTestPlayer()
		want := p.hp - 9
		p.perk = 1
		p.damage(10)

		if p.hp != want {
			t.Errorf("hp should be %.1f (10%% dmg reduction), got: %.1f", want, p.hp)
		}
	})

	t.Run("with berserk perk and 20% hp", func(t *testing.T) {
		p := newTestPlayer()
		p.perk = 3
		p.hp = p.hpcap * 0.2
		p.damage(10)

		if p.hp != 20-8.0 {
			t.Errorf("hp should be %.1f (20%% dmg reduction), got: %.1f", 20-8.0, p.hp)
		}
	})
}

func TestPlayer_Skill(t *testing.T) {
	find := func(name string) int {
		for i, s := range skills {
			if s.name == name {
				return i
			}
		}
		panic("no skill found")
	}

	t.Run("constraints", func(t *testing.T) {
		p := newTestPlayer()
		p.energy = 0

		if p.skill(0, p) {
			t.Error("should return false if player has no energy")
		}

		p.energy = 999
		p.effects["cd"+skills[0].name] = 1

		if p.skill(0, p) {
			t.Error("should return false if skill is in cooldown")
		}

		clear(p.effects)
		p.effects["disoriented"] = 1

		if p.skill(0, p) {
			t.Error("should return false if player is disoriented")
		}
	})

	t.Run("confused effect", func(t *testing.T) {
		i := find("charge")
		p := newTestPlayer()
		p.effects["confused"] = 1
		p.skill(i, p)

		en := skills[i].cost + 1
		if en != p.energycap-p.energy {
			t.Errorf("energy cost should increase by 1")
		}
	})

	t.Run("with ingenious perk", func(t *testing.T) {
		i := find("charge")
		p := newTestPlayer()
		p.perk = 4
		p.skill(i, p)
		basecd := skills[0].cd
		cd := p.effects["cd"+skills[0].name]

		if cd != basecd-2 {
			t.Errorf("cooldown should be reduced by 2, want %d, got %d", basecd-2, cd)
		}
	})

	t.Run("with berserk perk and < 25% hp", func(t *testing.T) {
		i := find("charge")
		p := newTestPlayer()
		p.perk = 3
		p.hp = p.hpcap * 0.25
		p.skill(i, p)
		basecd := skills[0].cd
		cd := p.effects["cd"+skills[0].name]

		if cd != basecd-1 {
			t.Errorf("cooldown should be reduced by one, want %d, got %d", basecd-1, cd)
		}
	})

	t.Run("with insanity perk", func(t *testing.T) {
		i := find("charge")
		skills[i].cd = 99999 // to make sure
		p := newTestPlayer()
		p.perk = 8
		p.skill(i, p)
		cd := p.effects["cd"+skills[0].name]

		if cd > 8 {
			t.Errorf("cooldown should be randomized from 0 to 8, got %d", cd)
		}
	})

	t.Run("with celestial staff weapon", func(t *testing.T) {
		p := newTestPlayer()

		for i, w := range weapons {
			if w.name == "celestial staff" {
				p.weapon = i
			}
		}

		i := find("charge")
		p.skill(i, p)
		basecd := skills[0].cd
		cd := p.effects["cd"+skills[0].name]

		if cd != basecd-1 {
			t.Errorf("cooldown should be reduced by one, want %d, got %d", basecd-1, cd)
		}
	})

	t.Run("charge", func(t *testing.T) {
		i := find("charge")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := p.strength * 1.3
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (130%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("frenzy", func(t *testing.T) {
		i := find("frenzy")
		p1 := newTestPlayer()
		p2 := newTestPlayer()
		sacrifice := p1.hpcap*0.05 + p1.hp*0.15
		p1.skill(i, p2)

		if 100-p1.hp != sacrifice {
			t.Errorf("should sacrifice %.1f hp (5%% hp cap + 15%% hp), got %.1f", sacrifice, 100-p1.hp)
		}

		dmg := p1.strength * 2.5
		if 100-p2.hp != dmg {
			t.Errorf("damage should be %.1f (250%% strength), got %.1f", dmg, 100-p2.hp)
		}
	})

	t.Run("great blow", func(t *testing.T) {
		i := find("great blow")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := p.strength * 2.1
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (210%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["stunned"] != 2 {
			t.Errorf("should get stunned effect for 2 turns, got %d", p.effects["stunned"])
		}
	})

	t.Run("poison", func(t *testing.T) {
		i := find("poison")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := p.strength * 0.85
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (85%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["poisoned"] != 3 {
			t.Error("should inflict poison for 3 turns, got", p.effects["poisoned"])
		}
	})
	t.Run("stun", func(t *testing.T) {
		i := find("stun")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := p.strength * 0.6
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (60%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["stunned"] != 2 {
			t.Error("should inflict stun for 2 turns, got", p.effects["stunned"])
		}
	})

	t.Run("icy blast", func(t *testing.T) {
		rolltest = 1
		i := find("icy blast")
		p := newTestPlayer()
		t.Log(p.attributes)
		p.skill(i, p)

		dmg := p.strength * 0.6
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (60%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["frozen"] != 2 {
			t.Error("should inflict frozen for 2 turns, got", p.effects["frozen"])
		}
	})

	t.Run("swift strike", func(t *testing.T) {
		i := find("swift strike")
		p := newTestPlayer()

		if p.skill(i, p) {
			t.Error("should return false because it doesnt consume turn")
		}

		dmg := p.strength * 0.85
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (85%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("knives throw", func(t *testing.T) {
		i := find("knives throw")
		p := newTestPlayer()

		if p.skill(i, p) {
			t.Error("should return false because it doesnt consume turn")
		}

		if 100-p.hp != 40 {
			t.Errorf("damage should be 40 (fixed), got %.1f", 100-p.hp)
		}
	})

	t.Run("fireball", func(t *testing.T) {
		i := find("fireball")
		p := newTestPlayer()
		p.skill(i, p)

		if 100-p.hp != 80 {
			t.Errorf("damage should be 80 (fixed), got %.1f", 100-p.hp)
		}

		if p.effects["burning"] != 2 {
			t.Errorf("should inflict burning for 2 turns, got %d", p.effects["burning"])
		}
	})

	t.Run("meteor strike", func(t *testing.T) {
		i := find("meteor strike")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := 100 - p.hp
		if dmg < 50 || dmg > 220 {
			t.Errorf("damage should be between 50-220, got %.1f", dmg)
		}
	})

	t.Run("strengthen", func(t *testing.T) {
		i := find("strengthen")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := p.strength
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["strengthen"] != 4 {
			t.Error("should get strengthen for 4 turns, got", p.effects["strengthen"])
		}
	})

	t.Run("focus attack", func(t *testing.T) {
		i := find("focus attack")
		p := newTestPlayer()
		p.skill(i, p)

		dmg := p.strength
		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["focus"] != 4 {
			t.Error("should get focus for 4 turns, got", p.effects["focus"])
		}
	})

	t.Run("devour", func(t *testing.T) {
		i := find("devour")
		p1 := newTestPlayer()
		p2 := newTestPlayer()
		p1.hp = 0
		p1.skill(i, p2)

		dmg := p1.strength * 1.5
		if 100-p2.hp != dmg {
			t.Errorf("damage should be %.1f (150%% strength), got %.1f", dmg, 100-p2.hp)
		}

		heal := p1.hpcap * 0.05
		if p1.hp != heal {
			t.Errorf("should heal by %.1f (5%% hpcap), got %.1f", heal, p1.hp)
		}
	})

	t.Run("vitality", func(t *testing.T) {
		i := find("vitality")
		p := newTestPlayer()
		p.skill(i, p)

		if p.effects["vitality"] != 5 {
			t.Error("should get vitality for 5 turns, got", p.effects["vitality"])
		}
	})

	t.Run("barrier", func(t *testing.T) {
		i := find("barrier")
		p := newTestPlayer()
		p.skill(i, p)

		if p.effects["barrier"] != 2 {
			t.Error("should get barrier for 2 turns, got", p.effects["barrier"])
		}
	})

	t.Run("force-field", func(t *testing.T) {
		i := find("force-field")
		p := newTestPlayer()
		p.skill(i, p)

		if p.effects["force-field"] != 5 {
			t.Error("should get force-field for 5 turns, got", p.effects["force-field"])
		}
	})

	t.Run("heal spell", func(t *testing.T) {
		i := find("heal spell")
		p := newTestPlayer()
		p.hp = 0
		p.skill(i, p)

		heal := p.hpcap * 0.15
		if p.hp != heal {
			t.Errorf("should heal by %.1f (15%% hpcap), got %.1f", heal, p.hp)
		}
	})

	t.Run("heal aura", func(t *testing.T) {
		i := find("heal aura")
		p := newTestPlayer()
		p.skill(i, p)

		if p.effects["heal aura"] != 4 {
			t.Error("should get heal aura for 4 turns, got", p.effects["heal aura"])
		}
	})

	t.Run("heal potion", func(t *testing.T) {
		i := find("heal potion")
		p := newTestPlayer()
		p.hp = 0
		p.skill(i, p)

		if p.hp != 40 {
			t.Errorf("should heal by 40 fixed, got %.1f", p.hp)
		}
	})

	t.Run("drain", func(t *testing.T) {
		i := find("drain")
		p := newTestPlayer()
		p.defense = 10
		dmg := p.hp*0.22 - p.defense
		p.skill(i, p)

		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (22%% hp), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("absorb", func(t *testing.T) {
		i := find("absorb")
		p := newTestPlayer()
		p.defense = 99999 // to make sure
		p.effects["immunity"] = 1
		dmg := p.hp * 0.1
		p.skill(i, p)

		if 100-p.hp != dmg {
			t.Errorf("damage should be %.1f (10%% hp cap ignore defense & effects), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("trick", func(t *testing.T) {
		i := find("trick")
		p1 := newTestPlayer()
		p2 := newTestPlayer()
		p1.strength = 99999 // to make it clear
		p1.skill(i, p2)

		dmg := p2.strength
		if 100-p2.hp != dmg {
			t.Errorf("damage should be %.1f (not us who attacked), got %.1f", dmg, 100-p2.hp)
		}
	})

	t.Run("vision", func(t *testing.T) {
		i := find("vision")
		p := newTestPlayer()

		if p.skill(i, p) {
			t.Error("should return false because it doesnt consume turn")
		}
	})
}

func TestPlayer_Rest(t *testing.T) {
	p := Player{gold: 10}
	p.hpcap = 100
	p.defense = 5
	p.energycap = 100
	p.rest()

	min := p.hpcap*0.1 + 20

	if p.hp < min {
		t.Errorf("should atleast heal %.1f hp (10%% hpcap + 20), got %.1f", min, p.hp)
	}

	if p.energy != 5 {
		t.Errorf("should restore 5 energy, got %d", p.energy)
	}
}

func TestPlayer_Train(t *testing.T) {
	t.Run("fail roll", func(t *testing.T) {
		rolltest = 1
		p := &Player{}
		p.train()

		v := p.hpcap
		v += p.strength
		v += p.defense
		v += p.agility
		v += float64(p.energycap)

		if v != 0 {
			t.Errorf("attributes should not be increased, %+v", p.attributes)
		}
	})

	t.Run("hp cap roll", func(t *testing.T) {
		rolltest = 51
		p := &Player{}
		p.train()

		if p.hpcap < 2.5 {
			t.Errorf("should atleast be increased by 2.5")
		}
	})

	t.Run("strength roll", func(t *testing.T) {
		rolltest = 62
		p := &Player{}
		p.train()

		if p.strength < 0.2 {
			t.Errorf("should atleast be increased by 0.2")
		}
	})

	t.Run("defense roll", func(t *testing.T) {
		rolltest = 73
		p := &Player{}
		p.train()

		if p.defense < 0.2 {
			t.Errorf("should atleast be increased by 0.2")
		}
	})

	t.Run("agility roll", func(t *testing.T) {
		rolltest = 84
		p := &Player{}
		p.train()

		if p.agility < 0.1 {
			t.Errorf("should atleast be increased by 0.1")
		}
	})

	t.Run("energy cap roll", func(t *testing.T) {
		rolltest = 95
		p := &Player{}
		p.train()

		if p.energycap != 1 {
			t.Errorf("should be increased by 1")
		}
	})
}

func TestPlayer_Flee(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		e := newTestPlayer()
		p := newTestPlayer()
		p.agility = 9999
		p.flee(e)

		if _, ok := p.effects["fled"]; !ok {
			t.Error("should get 'fled' effect")
		}
	})

	t.Run("too slow and get caught", func(t *testing.T) {
		rolltest = 25
		e := newTestPlayer()
		e.agility = 99999
		p := newTestPlayer()
		p.flee(e)

		if p.hp == 100 {
			t.Error("should be attacked by the enemy")
		}
	})

	t.Run("too slow and get caught but enemy is frozen", func(t *testing.T) {
		rolltest = 25
		e := newTestPlayer()
		e.effects["frozen"] = 1
		p := newTestPlayer()
		p.hp = 100
		p.flee(e)

		if p.hp != 100 {
			t.Error("should not be attacked by the enemy")
		}
	})

	t.Run("too slow and get caught but enemy is stunned", func(t *testing.T) {
		rolltest = 25
		e := newTestPlayer()
		e.effects["stunned"] = 1
		p := newTestPlayer()
		p.hp = 100
		p.flee(e)

		if p.hp != 100 {
			t.Error("should not be attacked by the enemy")
		}
	})

	t.Run("slipped in the mud", func(t *testing.T) {
		rolltest = 68
		p := newTestPlayer()
		p.hp = 100
		p.flee(p)

		if !equal(p.hp, 82) {
			t.Errorf("should take 18 damage, got %.1f", 100-p.hp)
		}
	})

	t.Run("fell into a ditch", func(t *testing.T) {
		rolltest = 76
		p := newTestPlayer()
		p.hp = 100
		p.flee(p)

		if p.hp != 64 {
			t.Errorf("should take 36 damage, got %.1f", 100-p.hp)
		}
	})

	t.Run("walked into a trap", func(t *testing.T) {
		rolltest = 84
		p := newTestPlayer()
		p.hp = 100
		p.flee(p)

		if p.hp != 95 {
			t.Errorf("should take 5%% hpcap damage, got %.1f", 100-p.hp)
		}
	})
}

func TestPlayer_useWeapon(t *testing.T) {
	find := func(name string) int {
		for i, w := range weapons {
			if w.name == name {
				return i
			}
		}
		panic("no weapon found")
	}

	t.Run("needle", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("needle")
		p.defense = 10

		dmg := 20 + p.defense*0.1
		got := p.useWeapon(20, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (ignore 25%% defense), got %.1f", dmg, got)
		}
	})

	t.Run("club", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("club")

		dmg := 100 + 100*0.05
		got := p.useWeapon(100, p)

		if got != dmg {
			t.Errorf("damage should be %.1f (5%% bonus), got %.1f", dmg, got)
		}
	})

	t.Run("flaming sword", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		p.weapon = find("flaming sword")
		p.useWeapon(10, p)

		if p.effects["burning"] != 2 {
			t.Errorf("should inflict burning for 2 turns, got %d", p.effects["burning"])
		}
	})

	t.Run("rapier", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("rapier")
		p.defense = 10

		dmg := 20 + p.defense*0.25
		got := p.useWeapon(20, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (ignore 25%% defense), got %.1f", dmg, got)
		}
	})

	t.Run("warhammer", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("warhammer")

		dmg := 100 + 100*0.1
		got := p.useWeapon(100, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (10%% bonus), got %.1f", dmg, got)
		}
	})

	t.Run("demonic blade", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("demonic blade")
		p.hp = 1000

		dmg := 10 + p.hp*0.05
		got := p.useWeapon(10, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (5%% target hp), got %.1f", dmg, got)
		}
	})

	t.Run("daunting mace", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("daunting mace")

		dmg := p.strength + p.hpcap*0.07
		got := p.useWeapon(p.strength, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (+7%% self hp cap as damage), got %.1f", dmg, got)
		}
	})

	t.Run("crimson blade", func(t *testing.T) {
		p := newTestPlayer()
		p.hp = 0
		p.weapon = find("crimson blade")
		p.useWeapon(10, p)

		if p.hp != 5 {
			t.Errorf("should recover hp by 5 (fixed), got %.1f", p.hp)
		}
	})

	t.Run("dragonscale blade", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		p.weapon = find("dragonscale blade")
		p.useWeapon(10, p)

		if p.effects["burning severe"] != 2 {
			t.Errorf("should inflict burning severe for 2 turns, got %d", p.effects["burning severe"])
		}

		rolltest = 10
		p.useWeapon(10, p)

		if p.effects["burning"] != 2 {
			t.Errorf("should inflict burning for 2 turns, got %d", p.effects["burning"])
		}
	})

	t.Run("astral rapier", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("astral rapier")
		p.defense = 10

		dmg := 20 + p.defense*0.4
		got := p.useWeapon(20, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (ignore 40%% defense), got %.1f", dmg, got)
		}
	})

	t.Run("lance", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("lance")

		dmg := p.strength * 1.5
		got := p.useWeapon(p.strength, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (50%% bonus), got %.1f", dmg, got)
		}

		got = p.useWeapon(p.strength, p)
		if got != p.strength {
			t.Errorf("damage should not be modified again, got %.1f", got)
		}
	})

	t.Run("obsidian warhammer", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("obsidian warhammer")

		dmg := 100 + 100*0.2
		got := p.useWeapon(100, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (20%% bonus), got %.1f", dmg, got)
		}
	})

	t.Run("voidforged rapier", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("voidforged rapier")
		p.defense = 10

		dmg := 20 + p.defense
		got := p.useWeapon(20, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (ignore defense), got %.1f", dmg, got)
		}
	})

	t.Run("soulreaper", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("soulreaper")
		p.hp = 1000

		dmg := 10 + p.hp*0.15
		got := p.useWeapon(10, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (15%% target hp), got %.1f", dmg, got)
		}
	})

	t.Run("celestial staff", func(t *testing.T) {
		p := newTestPlayer()
		p.hp = 0
		p.hpcap = 1000
		p.weapon = find("celestial staff")
		p.useWeapon(10, p)

		if p.hp != 20 {
			t.Errorf("should recover hp by 2%% hp cap, got %.1f", p.hp)
		}
	})

	t.Run("vanguard lance", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("vanguard lance")

		dmg := p.strength * 2
		got := p.useWeapon(p.strength, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (100%% bonus), got %.1f", dmg, got)
		}

		got = p.useWeapon(p.strength, p)
		if got != p.strength {
			t.Errorf("damage should not be modified again, got %.1f", got)
		}
	})

	t.Run("earthbreaker", func(t *testing.T) {
		p := newTestPlayer()
		p.weapon = find("earthbreaker")

		dmg := 100 + 100*0.25
		got := p.useWeapon(100, p)
		if got != dmg {
			t.Errorf("damage should be %.1f (25%% bonus), got %.1f", dmg, got)
		}
	})
}

func TestPlayer_UseArmor(t *testing.T) {
	find := func(name string) int {
		for i, w := range armory {
			if w.name == name {
				return i
			}
		}
		panic("no armor found")
	}

	t.Run("crystal armor", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("crystal armor")
		p.useArmor(1)

		if p.effects["reflect low"] == 0 {
			t.Error("should get reflect low effect")
		}
	})

	t.Run("enchanted plate", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("enchanted plate")

		dmg := 10 - 10*0.08
		got := p.useArmor(10)
		if got != dmg {
			t.Errorf("damage should be %.1f (8%% reduction), got %.1f", dmg, got)
		}
	})

	t.Run("deepsea mantle", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("deepsea mantle")

		dmg := p.hpcap * 0.2
		got := p.useArmor(99999) // to make sure
		if got != dmg {
			t.Errorf("damage should be %.1f (cannot exceed 20%% hp cap), got %.1f", dmg, got)
		}
	})

	t.Run("mythril plate", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("mythril plate")

		dmg := 10 - 10*0.16
		got := p.useArmor(10)
		if got != dmg {
			t.Errorf("damage should be %.1f (16%% reduction), got %.1f", dmg, got)
		}
	})

	t.Run("chainmail", func(t *testing.T) {
		p := newTestPlayer()
		p.defense = 10
		p.armor = find("chainmail")

		dmg := 10 - 10*0.15
		got := p.useArmor(10)
		if got != dmg {
			t.Errorf("damage should be %.1f (+15%% defense value), got %.1f", dmg, got)
		}
	})

	t.Run("void mantle", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("void mantle")
		p.useArmor(1)

		if p.effects["immunity"] == 0 {
			t.Errorf("should get immunity effect")
		}
	})

	t.Run("energy shield", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("energy shield")

		dmg := 10 - 10*0.35
		got := p.useArmor(10)
		if got != dmg {
			t.Errorf("damage should be %.1f (35%% reduction), got %.1f", dmg, got)
		}

		if p.effects["energy shield"] != 4 {
			t.Errorf("should get cooldown of 4 turn, got %d", p.effects["energy shield"])
		}
	})

	t.Run("amethyst armor", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("amethyst armor")
		p.useArmor(1)

		if p.effects["reflect high"] == 0 {
			t.Error("should get reflect high effect")
		}
	})

	t.Run("conqueror's armor", func(t *testing.T) {
		p := newTestPlayer()
		p.armor = find("conqueror's armor")

		dmg := 10 - 10*0.08
		got := p.useArmor(10)
		if got != dmg {
			t.Errorf("damage should be %.1f (8%% reduction), got %.1f", dmg, got)
		}
	})
}

func TestPlayer_SetPerk(t *testing.T) {
	t.Run("resilient", func(t *testing.T) {
		p := &Player{}
		p.hp = 100
		p.hpcap = 100
		p.defense = 1
		p.setPerk(1)

		if p.hpcap != 120 {
			t.Errorf("hpcap should be 120, got %.1f", p.hpcap)
		}

		if p.defense != 6 {
			t.Errorf("defense should be 6, got %.1f", p.defense)
		}

		p.setPerk(0)

		if p.hpcap != 100 || p.hp != 100 {
			t.Errorf("hpcap & hp should be back to 100, got %.1f %.1f", p.hpcap, p.hp)
		}

		if p.defense != 1 {
			t.Errorf("defense should be back to 1, got %.1f", p.defense)
		}
	})

	t.Run("havoc", func(t *testing.T) {
		p := &Player{}
		p.hp = 100
		p.hpcap = 100
		p.energycap = 20
		p.setPerk(2)

		if p.hpcap != 50 {
			t.Errorf("hpcap should be 50, got %.1f", p.hpcap)
		}

		if p.energycap != 16 {
			t.Errorf("energy cap should be 16, got %d", p.energycap)
		}

		p.setPerk(0)

		if p.hpcap != 100 {
			t.Errorf("hpcap should be back to 100, got %.1f", p.hpcap)
		}

		if p.energycap != 20 {
			t.Errorf("energy cap should be back to 20, got %d", p.energycap)
		}
	})

	t.Run("ingenious", func(t *testing.T) {
		p := &Player{}
		p.energycap = 20
		p.setPerk(4)

		if p.energycap != 22 {
			t.Errorf("energy cap should be 22, got %d", p.energycap)
		}

		p.setPerk(0)

		if p.energycap != 20 {
			t.Errorf("energy cap should be back to 20, got %d", p.energycap)
		}
	})

	t.Run("survivor", func(t *testing.T) {
		var p Player
		p.agility = 10
		p.setPerk(7)

		if p.agility != 15 {
			t.Errorf("agility should be 15, got %.1f", p.agility)
		}

		p.setPerk(0)

		if p.agility != 10 {
			t.Errorf("agility should be back to 10, got %.1f", p.agility)
		}
	})

	t.Run("shock", func(t *testing.T) {
		var p Player
		p.strength = 10
		p.setPerk(9)

		if p.strength != 15 {
			t.Errorf("strength should be 15, got %.1f", p.strength)
		}

		p.setPerk(0)

		if p.strength != 10 {
			t.Errorf("strength should be back to 10, got %.1f", p.strength)
		}
	})

	t.Run("frigid", func(t *testing.T) {
		var p Player
		p.agility = 10
		p.setPerk(10)

		if p.agility != 5 {
			t.Errorf("agility should be 5, got %.1f", p.agility)
		}

		p.setPerk(0)

		if p.agility != 10 {
			t.Errorf("agility should be back to 10, got %.1f", p.agility)
		}
	})
}

// test player with 100 hp. 10 strength. 20 energy
func newTestPlayer() *Player {
	var p Player
	p.hp = 100
	p.hpcap = 100
	p.strength = 10
	p.energy = 20
	p.energycap = 20
	p.effects = make(map[string]int)
	return &p
}
