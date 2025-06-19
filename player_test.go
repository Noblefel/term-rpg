package main

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer()

	if player.hpcap != 250 {
		t.Errorf("default hpcap should be 250, got %.1f", player.hpcap)
	}

	if player.defense != 15 {
		t.Errorf("default defense should be 15, got %.1f", player.defense)
	}

	if player.strength != 50 {
		t.Errorf("default strength should be 50, got %.1f", player.strength)
	}

	if player.agility != 5 {
		t.Errorf("default agility should be 5, got %.1f", player.agility)
	}

	if player.energycap != 20 {
		t.Errorf("default energycap should be 20, got %d", player.energycap)
	}
}

func TestPlayer_Attack(t *testing.T) {
	//enemy attribute access
	attr := &attributes{
		hp:    200,
		hpcap: 200,
	}

	player := NewPlayer()
	player.attack(attr)

	if attr.hp != 200-player.strength {
		t.Errorf("enemy hp should be reduced to %.1f", 200-player.strength)
	}

	player.strength = 100
	t.Run("with havoc perk", func(t *testing.T) {
		attr.hp = 200
		player.perk = 1
		player.attack(attr)
		dmg := attr.hpcap - attr.hp

		if dmg != 120 {
			t.Errorf("damage should be 120 (20%% increase), got: %.1f", dmg)
		}
	})

	t.Run("with berserk perk and 30% hp", func(t *testing.T) {
		attr.hp = 200
		player.perk = 2
		player.hp = 30
		player.hpcap = 100
		player.attack(attr)
		dmg := attr.hpcap - attr.hp

		if dmg != 115 {
			t.Errorf("damage should be around 115 (15%% increase), got: %.1f", dmg)
		}
	})

	t.Run("if low energy", func(t *testing.T) {
		attr.hp = 100
		attr.hpcap = 100
		player.perk = -1
		player.energy = 0
		player.attack(attr)
		dmg := attr.hpcap - attr.hp

		if dmg != 90 {
			t.Errorf("damage should be 90 (10%% decrease), got: %.1f", dmg)
		}
	})

	t.Run("deadman perk against undead", func(t *testing.T) {
		undead := undead{}
		undead.hp = 100
		undead.hpcap = 100

		player.perk = 5
		player.energy = 10
		player.strength = 10
		player.attack(&undead)

		dmg := undead.hpcap - undead.hp
		if !equal(dmg, 13.3) {
			t.Errorf("damage should be 13.3 (33%% increase), got: %.1f", dmg)
		}
	})

	t.Run("insanity perk", func(t *testing.T) {
		player.perk = 7
		player.strength = 10
		attr.hp = 100
		attr.hpcap = 100

		rolltest = 1
		player.attack(attr)
		dmg := 100 - attr.hp

		if dmg < 7 || dmg > 13 {
			t.Errorf("damage should be between 7 to 13 (-30%% to 30%% MULTIPLIER VAL), got %.1f", dmg)
		}

		rolltest = 50
		attr.hp = 100
		player.strength = 20
		player.attack(attr)
		dmg = 100 - attr.hp

		if dmg < 10 || dmg > 30 {
			t.Errorf("damage should be between 10 to 30 (-10 to 10 FLAT VAL), got %.1f", dmg)
		}
	})
}

func TestPlayer_Damage(t *testing.T) {
	player := Player{perk: -1}
	player.hp = 10
	player.defense = 5
	player.damage(10)
	player.effects = make(map[string]int)

	if player.hp != 5 {
		t.Errorf("hp should be 5, got %.1f", player.hp)
	}

	player.defense = 0
	t.Run("with resilient perk", func(t *testing.T) {
		player.perk = 0
		player.hp = 1000
		player.hpcap = 1000
		player.damage(100)
		var want float32 = 1000 - 90

		if player.hp != want {
			t.Errorf("hp should be %.1f (10%% dmg reduction), got: %.1f", want, player.hp)
		}

		player.hp = 100
		player.hpcap = 100
		player.damage(999)
		want = 100 - 16

		if player.hp != want {
			t.Errorf("hp should be %.1f (cannot exceed 16%% hp cap), got: %.1f", want, player.hp)
		}
	})

	t.Run("with berserk perk and 20% hp", func(t *testing.T) {
		player.perk = 2
		player.hp = 200
		player.hpcap = 1000
		player.damage(100)

		if player.hp != 120 {
			t.Errorf("hp should be %.1f (20%% dmg reduction), got: %.1f", 120+player.defense, player.hp)
		}
	})
}

func TestPlayer_Skill(t *testing.T) {
	player := Player{}
	player.strength = 10
	player.hpcap = 100
	player.effects = make(map[string]int)

	//enemy attribute access
	attr := &attributes{
		hp:      100,
		hpcap:   100,
		effects: make(map[string]int),
	}
	var e entity = attr

	if player.skill(0, e) {
		t.Error("should return false if player has no energy")
	}

	player.effects["cd"+skills[0].name] = 1

	if player.skill(0, e) {
		t.Error("should return false if skill is in cooldown")
	}

	clear(player.effects)
	player.energy = 9999

	find := func(name string) int {
		for i, s := range skills {
			if s.name == name {
				return i
			}
		}
		panic("no skill found")
	}

	t.Run("confused effect", func(t *testing.T) {
		i := find("charge")
		player.effects["confused"] = 1
		player.energy = 10
		player.skill(i, e)

		en := skills[i].cost + 1
		if en != 10-player.energy {
			t.Errorf("energy cost should increase by 1")
		}

		clear(player.effects)
		player.energy = 9999
	})

	t.Run("with ingenious perk", func(t *testing.T) {
		i := find("charge")
		basecd := skills[0].cd
		player.perk = 3
		player.skill(i, e)
		cd := player.effects["cd"+skills[0].name]
		clear(player.effects)

		if cd != basecd-2 {
			t.Errorf("cooldown should be reduced by 2, want %d, got %d", basecd-2, cd)
		}
	})

	t.Run("with berserk perk and < 25% hp", func(t *testing.T) {
		i := find("charge")
		basecd := skills[0].cd
		player.perk = 2
		player.hp = 2.5
		player.hp = 10
		player.skill(i, e)
		cd := player.effects["cd"+skills[0].name]
		clear(player.effects)

		if cd != basecd-1 {
			t.Errorf("cooldown should be reduced by one, want %d, got %d", basecd-1, cd)
		}
	})

	t.Run("with insanity perk", func(t *testing.T) {
		i := find("charge")
		skills[i].cd = 99999 // to make sure
		player.perk = 7
		player.skill(i, e)
		cd := player.effects["cd"+skills[0].name]
		clear(player.effects)

		if cd > 7 {
			t.Errorf("cooldown should be randomized from 0 to 7, got %d", cd)
		}
	})

	t.Run("charge", func(t *testing.T) {
		attr.hp = 100
		i := find("charge")
		player.perk = -1
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 13 {
			t.Errorf("damage should be 13 (130%% strength), got %.1f", dmg)
		}

		if player.strength != 10 {
			t.Errorf("strength should stay the same (10), got %.1f", player.strength)
		}
	})

	t.Run("frenzy", func(t *testing.T) {
		i := find("frenzy")
		player.hp = 100
		player.hpcap = 200
		dmg := player.hpcap*0.05 + player.hp*0.15

		player.skill(i, e)

		if 100-player.hp != dmg {
			t.Errorf("should sacrifice 15%% hp (15) + 5%% hpcap (10), got %.1f", 100-player.hp)
		}
	})

	t.Run("great blow", func(t *testing.T) {
		i := find("great blow")
		attr.hp = 100
		player.perk = -1
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 21 {
			t.Errorf("damage should be 21 (210%% strength), got %.1f", dmg)
		}

		if player.effects["stunned"] != 2 {
			t.Errorf("should get stun effect for 2 turns, got %d", player.effects["stunned"])
		}
	})

	t.Run("poison", func(t *testing.T) {
		i := find("poison")
		attr.hp = 100
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 8.5 {
			t.Errorf("damage should be 8.5 (85%% strength), got %.1f", dmg)
		}

		if e.attr().effects["poisoned"] != 3 {
			t.Error("should give poison effect for 3 turns, got", e.attr().effects["poisoned"])
		}
	})
	t.Run("stun", func(t *testing.T) {
		i := find("stun")
		attr.hp = 100
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 6 {
			t.Errorf("damage should be 6 (60%% strength), got %.1f", dmg)
		}

		if e.attr().effects["stunned"] != 2 {
			t.Error("should give stunned effect for 2 turns, got", e.attr().effects["stunned"])
		}
	})

	t.Run("swift strike", func(t *testing.T) {
		i := find("swift strike")
		attr.hp = 100

		if player.skill(i, e) {
			t.Error("should return false because it doesnt consume turn")
		}

		dmg := 100 - attr.hp
		if dmg != 8.5 {
			t.Errorf("damage should be 8.5 (85%% strength), got %.1f", dmg)
		}
	})

	t.Run("knives throw", func(t *testing.T) {
		i := find("knives throw")
		attr.hp = 100

		if player.skill(i, e) {
			t.Error("should return false because it doesnt consume turn")
		}

		dmg := 100 - attr.hp
		if dmg != 40 {
			t.Errorf("damage should be 40 (fixed), got %.1f", dmg)
		}
	})

	t.Run("fireball", func(t *testing.T) {
		i := find("fireball")
		attr.hp = 100
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 80 {
			t.Errorf("damage should be 80 (fixed), got %.1f", dmg)
		}

		if e.attr().effects["burning"] != 2 {
			t.Errorf("should give burning effect for 2 turns, got %d", e.attr().effects["burning"])
		}
	})

	t.Run("strengthen", func(t *testing.T) {
		i := find("strengthen")
		attr.hp = 100
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 10 {
			t.Errorf("damage should be 10 (100%% strength), got %.1f", dmg)
		}

		if player.effects["strengthen"] != 4 {
			t.Error("should get strengthen effect for 4 turns, got", player.effects["strengthen"])
		}
	})

	t.Run("barrier", func(t *testing.T) {
		i := find("barrier")
		player.skill(i, e)

		if player.effects["barrier"] != 2 {
			t.Error("should get barrier effect for 2 turns, got", player.effects["barrier"])
		}
	})

	t.Run("force-field", func(t *testing.T) {
		i := find("force-field")
		player.skill(i, e)

		if player.effects["force-field"] != 5 {
			t.Error("should get force-field effect for 5 turns, got", player.effects["force-field"])
		}
	})

	t.Run("heal spell", func(t *testing.T) {
		i := find("heal spell")
		player.hp = 0
		player.skill(i, e)

		heal := player.hpcap * 0.15
		if player.hp != heal {
			t.Errorf("should heal by %.1f (15%% hpcap), got %.1f", heal, player.hp)
		}
	})

	t.Run("heal aura", func(t *testing.T) {
		i := find("heal aura")
		player.hp = 0
		player.skill(i, e)

		if player.effects["heal aura"] != 4 {
			t.Error("should get heal aura effect for 4 turns, got", player.effects["heal aura"])
		}
	})

	t.Run("heal potion", func(t *testing.T) {
		i := find("heal potion")
		player.hp = 0
		player.skill(i, e)

		if player.hp != 40 {
			t.Errorf("should heal by 40 fixed, got %.1f", player.hp)
		}
	})

	t.Run("drain", func(t *testing.T) {
		i := find("drain")
		attr.hp = 100
		attr.hpcap = 200
		attr.defense = 5
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 17 {
			t.Errorf("damage should be 17 (22%% hp - defense), got %.1f", dmg)
		}
	})

	t.Run("absorb", func(t *testing.T) {
		i := find("absorb")
		attr.hp = 100
		attr.hpcap = 200
		attr.defense = 99999
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 20 {
			t.Errorf("damage should be 20 (10%% hp cap & ignore defense), got %.1f", dmg)
		}
	})

	t.Run("trick", func(t *testing.T) {
		i := find("trick")
		attr.hp = 100
		attr.strength = 10
		attr.defense = 0
		player.strength = 99999 // to make it clear
		player.skill(i, e)

		dmg := attr.strength
		if dmg != 100-attr.hp {
			t.Errorf("damage should be 10 (because its not us who attacked), got %.1f", 100-attr.hp)
		}
	})

}

func TestPlayer_Rest(t *testing.T) {
	player := Player{gold: 10}
	player.hpcap = 100
	player.defense = 5
	player.energycap = 100
	player.rest()

	if player.hp < 20 {
		t.Errorf("should atleast heal 20 hp, got %.1f", player.hp)
	}

	if player.energy != 5 {
		t.Errorf("should restore 5 energy, got %d", player.energy)
	}
}

func TestPlayer_Train(t *testing.T) {
	player := Player{gold: 100}

	t.Run("fail roll", func(t *testing.T) {
		rolltest = 1
		player.train()

		v := player.hpcap
		v += player.strength
		v += player.defense
		v += player.agility
		v += float32(player.energycap)

		if v != 0 {
			t.Errorf("attributes should not be increased, %+v", player.attributes)
		}
	})

	t.Run("hp cap roll", func(t *testing.T) {
		rolltest = 51
		player.train()

		if player.hpcap < 2.5 {
			t.Errorf("should atleast be increased by 2.5")
		}
	})

	t.Run("strength roll", func(t *testing.T) {
		rolltest = 62
		player.train()

		if player.strength < 0.2 {
			t.Errorf("should atleast be increased by 0.2")
		}
	})

	t.Run("defense roll", func(t *testing.T) {
		rolltest = 73
		player.train()

		if player.defense < 0.2 {
			t.Errorf("should atleast be increased by 0.2")
		}
	})

	t.Run("agility roll", func(t *testing.T) {
		rolltest = 84
		player.train()

		if player.agility < 0.1 {
			t.Errorf("should atleast be increased by 0.1")
		}
	})

	t.Run("energy cap roll", func(t *testing.T) {
		rolltest = 95
		player.train()

		if player.energycap != 1 {
			t.Errorf("should be increased by 1")
		}
	})
}

func TestPlayer_Flee(t *testing.T) {
	p := &Player{perk: -1}
	p.hp = 100
	p.hpcap = 100
	p.effects = make(map[string]int)

	var e entity = &attributes{
		strength: 20,
		effects:  make(map[string]int),
	}

	t.Run("success", func(t *testing.T) {
		p.agility = 9999
		p.flee(e)

		if _, ok := p.effects["fled"]; !ok {
			t.Error("should get 'fled' effect")
		}
	})

	t.Run("too slow and get caught", func(t *testing.T) {
		rolltest = 20
		p.agility = 0
		p.flee(e)

		if p.hp == 100 {
			t.Error("should be attacked by the enemy")
		}
	})

	t.Run("too slow and get caught but enemy is stunned", func(t *testing.T) {
		rolltest = 20
		p.hp = 100
		e.attr().effects["stunned"] = 1
		p.flee(e)

		if p.hp != 100 {
			t.Error("should not be attacked by the enemy")
		}
	})

	t.Run("slipped in the mud", func(t *testing.T) {
		rolltest = 68
		p.hp = 100
		p.flee(e)

		if !equal(p.hp, 82) {
			t.Errorf("should take 18 damage, got %.1f", 100-p.hp)
		}
	})

	t.Run("fell into a ditch", func(t *testing.T) {
		rolltest = 76
		p.hp = 100
		p.flee(e)

		if p.hp != 64 {
			t.Errorf("should take 36 damage, got %.1f", 100-p.hp)
		}
	})

	t.Run("walked into a trap", func(t *testing.T) {
		rolltest = 84
		p.hp = 100
		p.flee(e)

		if p.hp != 95 {
			t.Errorf("should take 5%% hpcap damage, got %.1f", 100-p.hp)
		}
	})
}

func TestPlayer_SetPerk(t *testing.T) {
	t.Run("resilient", func(t *testing.T) {
		p := &Player{perk: -1}
		p.hp = 100
		p.hpcap = 100
		p.defense = 1
		p.setPerk(0)

		if p.hpcap != 120 {
			t.Errorf("hpcap should be 120, got %.1f", p.hpcap)
		}

		if p.defense != 6 {
			t.Errorf("defense should be 6, got %.1f", p.defense)
		}

		p.setPerk(-1)

		if p.hpcap != 100 || p.hp != 100 {
			t.Errorf("hpcap & hp should be back to 100, got %.1f %.1f", p.hpcap, p.hp)
		}

		if p.defense != 1 {
			t.Errorf("defense should be back to 1, got %.1f", p.defense)
		}
	})

	t.Run("havoc", func(t *testing.T) {
		p := &Player{perk: -1}
		p.hp = 100
		p.hpcap = 100
		p.energycap = 20
		p.setPerk(1)

		if p.hpcap != 50 {
			t.Errorf("hpcap should be 50, got %.1f", p.hpcap)
		}

		if p.energycap != 16 {
			t.Errorf("energy cap should be 16, got %d", p.energycap)
		}

		p.setPerk(-1)

		if p.hpcap != 100 || p.hp != 100 {
			t.Errorf("hpcap & hp should be back to 100, got %.1f %.1f", p.hpcap, p.hp)
		}

		if p.energycap != 20 {
			t.Errorf("energy cap should be back to 20, got %d", p.energycap)
		}
	})

	t.Run("ingenious", func(t *testing.T) {
		p := &Player{perk: -1}
		p.energycap = 20
		p.setPerk(3)

		if p.energycap != 22 {
			t.Errorf("energy cap should be 22, got %d", p.energycap)
		}

		p.setPerk(-1)

		if p.energycap != 20 {
			t.Errorf("energy cap should be back to 20, got %d", p.energycap)
		}
	})

	t.Run("survivor", func(t *testing.T) {
		var p Player
		p.agility = 10
		p.setPerk(6)

		if p.agility != 15 {
			t.Errorf("agility should be 15, got %.1f", p.agility)
		}

		p.setPerk(-1)

		if p.agility != 10 {
			t.Errorf("agility should be back to 10, got %.1f", p.agility)
		}
	})
}
