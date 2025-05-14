package main

import (
	"math"
	"testing"
)

func TestKnight(t *testing.T) {
	p := &Player{}
	p.hpcap = 100
	knight := newKnight()

	t.Run("reinforce armor", func(t *testing.T) {
		tempdef := knight.attr().defense
		knight.setHP(0)

		rolltest = 1
		knight.attack(p)

		buff := knight.attr().defense - tempdef
		if buff == 0 {
			t.Error("should buff defense")
		}

		if knight.attr().hp < 5 {
			t.Error("should heal atleast 5 hp")
		}
	})

	t.Run("attack roll", func(t *testing.T) {
		rolltest = 15
		knight.attack(p)

		if p.hp == 100 {
			t.Error("player is not damaged")
		}
	})
}

func TestWizard(t *testing.T) {
	p := &Player{}
	p.hpcap = 100
	p.perk = -1
	p.hp = 100
	wizard := newWizard()

	t.Run("staff attack", func(t *testing.T) {
		rolltest = 1
		wizard.attack(p)

		dmg := wizard.attr().strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f, got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("enhanced attack", func(t *testing.T) {
		rolltest = 20
		p.hp = 100
		wizard.attack(p)

		dmg := wizard.attr().strength * 3
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (300%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("summon meteor", func(t *testing.T) {
		rolltest = 70
		p.hp = 100
		wizard.attack(p)

		if p.hp != 65 {
			t.Errorf("damage should be 35, got %.1f", 100-p.hp)
		}
	})

	t.Run("cast healing", func(t *testing.T) {
		rolltest = 80
		wizard.setHP(0)
		wizard.attack(p)

		heal := wizard.attr().hpcap * 0.2

		if wizard.attr().hp != heal {
			t.Errorf("should heal by 20%% hp cap (20), got %.1f", wizard.attr().hp)
		}
	})
}

func TestChangeling(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.defense = 8
	p.strength = 35

	changeling := newChangeling()
	changeling.attack(p)

	attr := changeling.attr()

	if attr.hp != p.hp || attr.defense != p.defense || attr.strength != p.strength {
		t.Errorf("should copy player's attribute\nwant: %+v\ngot: %+v", attr, p.attr())
	}
}

func TestVampire(t *testing.T) {
	p := &Player{}
	p.perk = -1
	p.hp = 100
	p.hpcap = 100

	vampire := newVampire()

	t.Run("exposed to sunlight", func(t *testing.T) {
		rolltest = 1
		temp := vampire.attr().hp
		vampire.attack(p)

		if vampire.attr().hp == temp {
			t.Error("should take damage")
		}
	})

	t.Run("bites", func(t *testing.T) {
		rolltest = 15
		strength := vampire.attr().strength
		vampire.setHP(0)
		vampire.attack(p)

		if vampire.attr().hp != strength {
			t.Errorf("should heal based on damage %.1f, got %.1f", strength, vampire.attr().hp)
		}
	})

	t.Run("claws", func(t *testing.T) {
		rolltest = 70
		p.hp = 100
		vampire.attack(p)

		dmg := vampire.attr().strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("bat swarm", func(t *testing.T) {
		rolltest = 85
		p.hp = 100
		vampire.attack(p)

		dmg := vampire.attr().strength * 1.2
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (120%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestDemon(t *testing.T) {
	p := &Player{}
	p.perk = -1
	p.hp = 1000
	p.hpcap = 1000

	demon := newDemon()

	t.Run("soul absorption", func(t *testing.T) {
		p.defense = 99999
		rolltest = 1
		demon.attack(p)

		dmg := demon.attr().strength + 70
		if dmg != 1000-p.hp {
			t.Errorf("damage should be %.1f (take 7%% hpcap and ignore defense), got %.1f", dmg, 1000-p.hp)
		}
	})

	t.Run("basic attacks", func(t *testing.T) {
		rolltest = 60
		p.hp = 100
		p.defense = 0
		demon.attack(p)

		dmg := demon.attr().strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestShardling(t *testing.T) {
	p := &Player{}
	p.perk = -1
	p.hp = 100
	p.hpcap = 100
	p.strength = 10
	p.energy = 10 //avoid exhaustion

	shardling := newShardling().(*shardling)
	shardling.defense = 0

	t.Run("damage reflection", func(t *testing.T) {
		player = p // needed for this
		p.attack(shardling)
		dmg := p.strength * 0.3

		if dmg != 100-p.hp {
			t.Errorf("reflected damage should be %.1f (30%% of the damage), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("ram", func(t *testing.T) {
		p.hp = 100
		rolltest = 1
		shardling.attack(p)

		dmg := shardling.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("crystal limbs", func(t *testing.T) {
		p.hp = 100
		rolltest = 20
		shardling.attack(p)

		dmg := shardling.strength * 1.1
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (110%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("spike", func(t *testing.T) {
		p.hp = 100
		rolltest = 80
		shardling.attack(p)

		dmg := 20 + shardling.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength + 20), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestGenie(t *testing.T) {
	p := &Player{}
	p.name = "player"
	p.perk = -1
	p.hp = 100
	p.hpcap = 100
	p.strength = 20
	p.defense = 5
	p.energycap = 20

	genie := newGenie()

	t.Run("hp curse", func(t *testing.T) {
		rolltest = 1
		genie.attack(p)

		if p.hpcap >= 100 {
			t.Error("hp cap should be reduced")
		}
	})

	t.Run("strength curse", func(t *testing.T) {
		rolltest = 5
		genie.attack(p)

		if p.strength >= 20 {
			t.Error("strength should be reduced")
		}
	})

	t.Run("defense curse", func(t *testing.T) {
		rolltest = 10
		genie.attack(p)

		if p.defense >= 5 {
			t.Error("defense should be reduced")
		}
	})

	t.Run("energy cap curse", func(t *testing.T) {
		rolltest = 15
		genie.attack(p)

		if p.energycap >= 20 {
			t.Error("energy cap should be reduced")
		}
	})

	t.Run("illusion", func(t *testing.T) {
		p.hp = 999
		p.strength = 9999
		rolltest = 18
		genie.attack(p)

		// player should be dead because genie make us self attack
		// strength set to 9999 to speed up
		if p.hp > 0 {
			t.Error("player should be dead")
		}
	})

	p.defense = 0

	t.Run("sandstorm", func(t *testing.T) {
		p.hp = 100
		rolltest = 31
		genie.attack(p)

		dmg := 25 + genie.attr().strength*0.8
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (80%% strength + 25), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("blast", func(t *testing.T) {
		p.hp = 100
		rolltest = 53
		genie.attack(p)

		dmg := 10 + genie.attr().strength*1.4
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (140%% strength + 10), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("punch", func(t *testing.T) {
		p.hp = 100
		rolltest = 80
		genie.attack(p)

		dmg := genie.attr().strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

// quick fix floating issue
func equal(n, n2 float32) bool {
	x := math.Round(float64(n) * 100)
	y := math.Round(float64(n2) * 100)
	return x == y
}
