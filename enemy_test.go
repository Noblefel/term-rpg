package main

import "testing"

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

		dmg := demon.attr().strength + 50
		if dmg != 1000-p.hp {
			t.Errorf("damage should be %.1f (take 5%% hpcap and ignore defense), got %.1f", dmg, 1000-p.hp)
		}
	})
}
