package main

import (
	"testing"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer()

	if player.hpcap != 100 {
		t.Errorf("default hpcap should be 100, got %.1f", player.hpcap)
	}

	if player.defense != 1 {
		t.Errorf("default defense should be 1, got %.1f", player.defense)
	}

	if player.strength != 20 {
		t.Errorf("default strength should be 20, got %.1f", player.strength)
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
	enemy := entity(attr)

	player := NewPlayer()
	player.attack(enemy)

	if attr.hp != 200-player.strength {
		t.Errorf("enemy hp should be reduced to %.1f", 200-player.strength)
	}

	player.strength = 100
	t.Run("with havoc perk", func(t *testing.T) {
		attr.hp = 200
		player.perk = 1
		player.attack(enemy)
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
		player.attack(enemy)
		dmg := attr.hpcap - attr.hp

		if dmg != 110 {
			t.Errorf("damage should be around 110 (10%% increase), got: %.1f", dmg)
		}
	})

	t.Run("if low energy", func(t *testing.T) {
		attr.hp = 100
		attr.hpcap = 100
		player.perk = -1
		player.energy = 0
		player.attack(enemy)
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
		player.damage(100)
		var want float32 = 1000 - 90

		if player.hp != want {
			t.Errorf("hp should be %.1f (10%% dmg reduction), got: %.1f", want, player.hp)
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

	t.Run("with barrier effect", func(t *testing.T) {
		player.perk = -1
		player.hp = 1000
		player.effects["barrier"] = 1
		player.damage(100)
		var want float32 = 1000 - 60

		if player.hp != want {
			t.Errorf("hp should be %.1f (40%% dmg reduction), got: %.1f", want, player.hp)
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

	t.Run("with berserk perk and < 15% hp", func(t *testing.T) {
		i := find("charge")
		basecd := skills[0].cd
		player.perk = 2
		player.hp = 1
		player.skill(i, e)
		cd := player.effects["cd"+skills[0].name]
		clear(player.effects)

		if cd != basecd-1 {
			t.Errorf("cooldown should be reduced by one, want %d, got %d", basecd-1, cd)
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
		player.skill(i, e)

		// 20% 100 = 20
		// 5% 200 = 10
		if player.hp != 70 {
			t.Errorf("should sacrifice 20%% hp (20) + 5%% hpcap (4), got %.1f", player.hp)
		}

		clear(player.effects)
		player.perk = 1
		player.hp = 100
		player.skill(i, e)

		if player.hp != 70 {
			t.Error("should not be affected by perk modifier", player.hp)
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
		if dmg != 15 {
			t.Errorf("damage should be 15 (fixed), got %.1f", dmg)
		}
	})

	t.Run("fireball", func(t *testing.T) {
		i := find("fireball")
		attr.hp = 100
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 24 {
			t.Errorf("damage should be 24 (fixed), got %.1f", dmg)
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

		heal := player.hpcap*0.12 + 5
		if player.hp != heal {
			t.Errorf("should heal by %.1f (5 + 12%% hpcap), got %.1f", heal, player.hp)
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

		if player.hp != 34 {
			t.Errorf("should heal by 34 fixed, got %.1f", player.hp)
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

	if player.hp < 15 {
		t.Errorf("should atleast heal 15 hp, got %.1f", player.hp)
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
		v += float32(player.energycap)

		if v != 0 {
			t.Errorf("attributes should not be increased, %+v", player.attributes)
		}
	})

	t.Run("hp cap roll", func(t *testing.T) {
		rolltest = 60
		player.train()

		if player.hpcap < 1 {
			t.Errorf("should atleast be increased by 1")
		}
	})

	t.Run("strength roll", func(t *testing.T) {
		rolltest = 71
		player.train()

		if player.strength < 0.1 {
			t.Errorf("should atleast be increased by 0.1")
		}
	})

	t.Run("defense roll", func(t *testing.T) {
		rolltest = 82
		player.train()

		if player.defense < 0.1 {
			t.Errorf("should atleast be increased by 0.1")
		}
	})

	t.Run("energy cap roll", func(t *testing.T) {
		rolltest = 93
		player.train()

		if player.energycap != 1 {
			t.Errorf("should be increased by 1")
		}
	})
}

func TestPlayer_Flee(t *testing.T) {
	player := &Player{perk: -1}
	player.hp = 100
	player.hpcap = 100
	player.effects = make(map[string]int)

	var e entity = &attributes{
		strength: 20,
		effects:  make(map[string]int),
	}

	t.Run("success roll", func(t *testing.T) {
		rolltest = 1
		player.flee(e)

		if _, ok := player.effects["fled"]; !ok {
			t.Error("should get 'fled' effect")
		}
	})

	t.Run("too slow and get caught", func(t *testing.T) {
		rolltest = 60
		player.flee(e)

		if player.hp != 80 {
			t.Error("should be attacked by the enemy")
		}
	})

	t.Run("too slow and get caught but enemy is stunned", func(t *testing.T) {
		rolltest = 60
		player.hp = 100
		e.attr().effects["stunned"] = 1
		player.flee(e)

		if player.hp != 100 {
			t.Error("should not be attacked by the enemy")
		}
	})

	t.Run("slipped in the mud", func(t *testing.T) {
		rolltest = 68
		player.hp = 100
		player.flee(e)

		if player.hp != 98 {
			t.Errorf("should take 2 damage, got %.1f", 100-player.hp)
		}
	})

	t.Run("fell into a ditch", func(t *testing.T) {
		rolltest = 76
		player.hp = 100
		player.flee(e)

		if player.hp != 94 {
			t.Errorf("should take 6 damage, got %.1f", 100-player.hp)
		}
	})

	t.Run("walked into a trap", func(t *testing.T) {
		rolltest = 84
		player.hp = 100
		player.flee(e)

		if player.hp != 95 {
			t.Errorf("should take 5%% hpcap damage, got %.1f", 100-player.hp)
		}
	})
}

func TestPlayer_SetPerk(t *testing.T) {
	t.Run("resilient", func(t *testing.T) {
		player := NewPlayer()
		player.setPerk(0)

		if player.hpcap != 105 {
			t.Errorf("hpcap should be 105, got %.1f", player.hpcap)
		}

		if player.defense != 3.5 {
			t.Errorf("defense should be 3.5, got %.1f", player.defense)
		}

		player.setPerk(-1)

		if player.hpcap != 100 || player.hp != 100 {
			t.Errorf("hpcap & hp should be back to 100, got %.1f %.1f", player.hpcap, player.hp)
		}

		if player.defense != 1 {
			t.Errorf("defense should be back to 1, got %.1f", player.defense)
		}
	})

	t.Run("havoc", func(t *testing.T) {
		player := NewPlayer()
		player.setPerk(1)

		if player.hpcap != 75 {
			t.Errorf("hpcap should be 75, got %.1f", player.hpcap)
		}

		if player.energycap != 16 {
			t.Errorf("energy cap should be 16, got %d", player.energycap)
		}

		player.setPerk(-1)

		if player.hpcap != 100 || player.hp != 100 {
			t.Errorf("hpcap & hp should be back to 100, got %.1f %.1f", player.hpcap, player.hp)
		}

		if player.energy != 20 || player.energycap != 20 {
			t.Errorf("energy & cap should be back to 20, got %d %d", player.energycap, player.energy)
		}
	})

	t.Run("ingenious", func(t *testing.T) {
		player := NewPlayer()
		player.setPerk(3)

		if player.energycap != 22 {
			t.Errorf("energy cap should be 22, got %d", player.energycap)
		}

		player.setPerk(-1)

		if player.energycap != 20 {
			t.Errorf("energy cap should be back to 20, got %d", player.energycap)
		}
	})
}
