package main

import (
	"io"
	"strings"
	"testing"
)

func init() { out = io.Discard }

func TestNewPlayer(t *testing.T) {
	player := NewPlayer(-1)

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

	t.Run("with resiliency perk", func(t *testing.T) {
		player := NewPlayer(0)

		if player.hpcap != 105 {
			t.Errorf("hpcap should be 105, got %.1f", player.hpcap)
		}

		if player.defense != 3.5 {
			t.Errorf("defense should be 3.5, got %.1f", player.defense)
		}
	})

	t.Run("with havoc perk", func(t *testing.T) {
		player := NewPlayer(1)

		if player.hpcap != 75 {
			t.Errorf("hpcap should be 75, got %.1f", player.hpcap)
		}

		if player.gold != 0 {
			t.Errorf("gold should be 0, got %d", player.gold)
		}
	})
}

func TestPlayer_Attack(t *testing.T) {
	//enemy attribute access
	attr := &attributes{
		hp:    200,
		hpcap: 200,
	}
	enemy := entity(attr)

	player := NewPlayer(0)
	player.attack(enemy)

	if attr.hp != 200-player.strength {
		t.Errorf("enemy hp should be reduced to %.1f", 200-player.strength)
	}

	t.Run("with havoc perk", func(t *testing.T) {
		attr.hp = 200
		player.perk = 1
		player.strength = 100
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
		player.strength = 100
		player.attack(enemy)
		dmg := attr.hpcap - attr.hp

		if dmg != 110 {
			t.Errorf("damage should be around 110 (10%% increase), got: %.1f", dmg)
		}
	})
}

func TestPlayer_Damage(t *testing.T) {
	player := Player{perk: -1}
	player.hp = 10
	player.defense = 5
	player.damage(10)

	if player.hp != 5 {
		t.Errorf("hp should be 5, got %.1f", player.hp)
	}

	t.Run("with resiliency perk", func(t *testing.T) {
		player.perk = 0
		player.defense = 1
		player.hp = 1000
		player.damage(100)
		want := 1000 - 90 + player.defense

		if player.hp != want {
			t.Errorf("hp should be %.1f (10%% dmg reduction), got: %.1f", want, player.hp)
		}
	})

	t.Run("with berserk perk and 20% hp", func(t *testing.T) {
		player.perk = 2
		player.defense = 1
		player.hp = 200
		player.hpcap = 1000
		player.damage(100)

		if player.hp != 120+player.defense {
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

	player.effects["cd_"+skills[0].name] = 1

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

	t.Run("with ingenious perk", func(t *testing.T) {
		i := find("charge")
		basecd := skills[0].cd
		player.perk = 3
		player.skill(i, e)
		cd := player.effects["cd_"+skills[0].name]

		if cd != basecd-1 {
			t.Errorf("cooldown should be reduced by one, want %d, got %d", basecd-1, cd)
		}
	})

	t.Run("with berserk perk and < 15% hp", func(t *testing.T) {
		i := find("charge")
		basecd := skills[0].cd
		player.perk = 2
		player.hp = 1
		player.skill(i, e)
		cd := player.effects["cd_"+skills[0].name]

		if cd != basecd-1 {
			t.Errorf("cooldown should be reduced by one, want %d, got %d", basecd-1, cd)
		}
	})
	t.Run("charge", func(t *testing.T) {
		i := find("charge")
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 13 {
			t.Errorf("damage should be 13 (130%% strength), got %.1f", dmg)
		}

		if player.strength != 10 {
			t.Errorf("strength should stay the same (10), got %.1f", player.strength)
		}
	})

	t.Run("heal", func(t *testing.T) {
		i := find("heal")
		player.hp = 0
		player.skill(i, e)

		if player.hp != 18 {
			t.Errorf("should heal by 10 + 8%% hpcap, got %.1f", player.hp)
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

	t.Run("vision", func(t *testing.T) {
		var sb strings.Builder
		out = &sb

		i := find("vision")
		attr.hpcap = 99
		attr.strength = 5.5
		attr.defense = 1.3
		player.skill(i, e)

		got := replacer.Replace(sb.String())
		dump(t, got)

		if !strings.Contains(got, "hp cap   : 99.0") {
			t.Error("incorrect hp cap")
		}

		if !strings.Contains(got, "strength : 5.5") {
			t.Error("incorrect strength")
		}

		if !strings.Contains(got, "defense  : 1.3") {
			t.Error("incorrect defense")
		}
	})

	t.Run("drain", func(t *testing.T) {
		i := find("drain")
		attr.hp = 100
		attr.hpcap = 200
		attr.defense = 5
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 15 {
			t.Errorf("damage should be 15 (20%% hp  - defense), got %.1f", dmg)
		}
	})

	t.Run("absorb", func(t *testing.T) {
		i := find("absorb")
		attr.hp = 100
		attr.hpcap = 200
		attr.defense = 99999
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 15 {
			t.Errorf("damage should be 15 (7.5%% hp cap & ignore defense), got %.1f", dmg)
		}
	})

	t.Run("trick", func(t *testing.T) {
		i := find("trick")
		attr.hp = 100
		attr.strength = 10
		attr.defense = 0
		player.strength = 99999 // make it clear its not us who attack
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 10 {
			t.Errorf("damage should be 10, got %.1f", dmg)
		}
	})

	t.Run("poison", func(t *testing.T) {
		i := find("poison")
		attr.hp = 100
		player.perk = -1
		player.strength = 10
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 6 {
			t.Errorf("damage should be 6 (60%% strength), got %.1f", dmg)
		}

		if e.attr().effects["poisoned"] != 3 {
			t.Error("should give poison effect for 3 turns, got", e.attr().effects["poisoned"])
		}
	})
	t.Run("stun", func(t *testing.T) {
		i := find("stun")
		attr.hp = 100
		player.perk = -1
		player.strength = 10
		player.skill(i, e)

		dmg := 100 - attr.hp
		if dmg != 3 {
			t.Errorf("damage should be 3 (30%% strength), got %.1f", dmg)
		}

		if e.attr().effects["stunned"] != 2 {
			t.Error("should give stunned effect for 2 turns, got", e.attr().effects["stunned"])
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

	if player.gold != 5 {
		t.Errorf("should cost 5 gold, got %d", 10-player.gold)
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

		if player.gold != 90 {
			t.Errorf("should cost 10 gold, got %d", 100-player.gold)
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
