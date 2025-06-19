package main

import (
	"testing"
)

func TestKnight(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.hpcap = 100
	p.perk = -1

	var knight knight
	knight.hpcap = 100
	knight.strength = 10
	knight.effects = make(map[string]int)

	t.Run("reinforce armor", func(t *testing.T) {
		tempdef := knight.defense
		rolltest = 1
		knight.attack(p)

		buff := knight.defense - tempdef
		if buff == 0 {
			t.Error("should buff defense")
		}

		if knight.hp < 20 {
			t.Error("should heal atleast 20 hp")
		}
	})

	t.Run("strengthen", func(t *testing.T) {
		rolltest = 15
		knight.attack(p)

		dmg := knight.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}

		eff := knight.effects["strengthen"]
		clear(knight.effects)
		if eff != 4 {
			t.Errorf("should get strengthen effect for 4 turns, got %d", eff)
		}
	})

	t.Run("attack roll", func(t *testing.T) {
		rolltest = 30
		p.hp = 100
		knight.attack(p)

		dmg := knight.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestWizard(t *testing.T) {
	p := &Player{}
	p.hpcap = 100
	p.perk = -1
	p.hp = 100
	p.effects = make(map[string]int)

	var wizard wizard
	wizard.hpcap = 100
	wizard.strength = 10
	wizard.effects = make(map[string]int)

	t.Run("cast healing", func(t *testing.T) {
		rolltest = 1
		wizard.setHP(0)
		wizard.attack(p)

		heal := wizard.hpcap * 0.2

		if wizard.hp != heal {
			t.Errorf("should heal by 20 (20%% hpcap), got %.1f", wizard.hp)
		}
	})

	t.Run("immunity", func(t *testing.T) {
		rolltest = 15
		wizard.attack(p)

		if wizard.effects["immunity"] != 2 {
			t.Errorf("should get immunity effect for 2 turns, got %d", wizard.effects["immunity"])
		}
	})

	t.Run("barrier", func(t *testing.T) {
		rolltest = 20
		wizard.attack(p)

		if wizard.effects["barrier"] != 4 {
			t.Errorf("should get barrier effect for 4 turns, got %d", wizard.effects["barrier"])
		}
	})

	t.Run("confuse", func(t *testing.T) {
		rolltest = 30
		p.hp = 100
		wizard.attack(p)

		if p.effects["confused"] != 5 {
			t.Errorf("should get confused effect for 5 turns, got %d", p.effects["confused"])
		}

		dmg := wizard.strength * 0.6
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (60%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("enhanced attack", func(t *testing.T) {
		rolltest = 35
		p.hp = 100
		wizard.attack(p)

		dmg := wizard.strength * 3
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (300%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("summon meteor", func(t *testing.T) {
		rolltest = 85
		p.hp = 200
		wizard.attack(p)

		if 200-p.hp != 120 {
			t.Errorf("damage should be 120, got %.1f", 100-p.hp)
		}
	})

	t.Run("staff attack", func(t *testing.T) {
		rolltest = 90
		p.hp = 100
		wizard.attack(p)

		dmg := wizard.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f, got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestChangeling(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.hpcap = 100
	p.defense = 8
	p.strength = 35

	changeling := newChangeling()
	changeling.attack(p)

	attr := changeling.attr()

	if attr.hp != p.hp*0.75 {
		t.Errorf("should copy hp by 75%%, want %.1f, got: %.1f", p.hp*0.75, attr.hp)
	}

	if attr.hpcap != p.hpcap*0.75 {
		t.Errorf("should copy hpcap by 75%%, want %.1f, got: %.1f", p.hpcap*0.75, attr.hpcap)
	}

	if attr.strength != p.strength*0.75 {
		t.Errorf("should copy strength by 75%%, want %.1f, got: %.1f", p.strength*0.75, attr.strength)
	}

	if attr.defense != p.defense*0.75 {
		t.Errorf("should copy defense by 75%%, want %.1f, got: %.1f", p.defense*0.75, attr.defense)
	}

	if attr.agility != p.agility*0.75 {
		t.Errorf("should copy agility by 75%%, want %.1f, got: %.1f", p.agility*0.75, attr.agility)
	}
}

func TestVampire(t *testing.T) {
	p := &Player{}
	p.perk = -1
	p.hp = 100
	p.hpcap = 100
	p.effects = make(map[string]int)

	var vampire vampire
	vampire.hp = 100
	vampire.hpcap = 100
	vampire.strength = 10

	t.Run("exposed to sunlight", func(t *testing.T) {
		rolltest = 1
		vampire.defense = 999
		vampire.attack(p)
		vampire.defense = 0

		if 100-vampire.hp != 10 {
			t.Errorf("damage should be 10 (10%% of hp and ignore defense), got %.1f", 100-vampire.hp)
		}
	})

	t.Run("bites lifesteal", func(t *testing.T) {
		rolltest = 8
		dmg := vampire.strength + p.hp*0.01
		vampire.hp = 0
		vampire.attack(p)

		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength + 1%% target hp), got %.1f", dmg, 100-p.hp)
		}

		heal := dmg/2 + 15
		if vampire.hp != heal {
			t.Errorf("should heal by %.1f (half the damage + 15), got %.1f", heal, vampire.hp)
		}
	})

	t.Run("bites poison", func(t *testing.T) {
		rolltest = 60
		p.hp = 100
		dmg := vampire.strength + p.hp*0.01
		vampire.attack(p)

		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength + 1%% target hp), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["poisoned"] != 3 {
			t.Errorf("should give poison effect for 3 turns, got %d", p.effects["poisoned"])
		}
	})

	t.Run("bat swarm", func(t *testing.T) {
		rolltest = 70
		p.hp = 100
		vampire.attack(p)

		dmg := vampire.attr().strength * 1.2
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (120%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("claws", func(t *testing.T) {
		rolltest = 85
		p.hp = 100
		vampire.attack(p)

		dmg := vampire.attr().strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

}

func TestDemon(t *testing.T) {
	p := &Player{}
	p.perk = -1
	p.effects = make(map[string]int)

	var demon demon
	demon.hpcap = 100
	demon.strength = 10

	t.Run("soul absorption", func(t *testing.T) {
		p.hp = 1000
		p.hpcap = 1000
		p.defense = 99999
		rolltest = 1
		demon.attack(p)

		dmg := demon.strength + 40
		if dmg != 1000-p.hp {
			t.Errorf("damage should be %.1f (take 4%% hp and ignore defense), got %.1f", dmg, 1000-p.hp)
		}
	})

	t.Run("hell fire burning", func(t *testing.T) {
		rolltest = 60
		p.hp = 100
		p.defense = 0
		demon.attack(p)

		dmg := demon.strength * 0.8
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (80%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["burning"] != 3 {
			t.Errorf("should give burning effect for 3 turns, got %d", p.effects["burning"])
		}
	})

	t.Run("basic attacks", func(t *testing.T) {
		rolltest = 75
		p.hp = 100
		demon.attack(p)

		dmg := demon.strength
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

	shardling := &shardling{}
	shardling.strength = 10
	shardling.hp = 100

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
		rolltest = 40
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

		dmg := shardling.strength * 1.25
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (125%% strength), got %.1f", dmg, 100-p.hp)
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

	var genie genie
	genie.strength = 10
	genie.effects = make(map[string]int)

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

	t.Run("force-field", func(t *testing.T) {
		p.hp = 100
		rolltest = 31
		genie.attack(p)

		if genie.effects["force-field"] != 5 {
			t.Errorf("should get force-field effect for 5 turns, got %d", genie.effects["force-field"])
		}
	})

	t.Run("sandstorm", func(t *testing.T) {
		p.hp = 100
		rolltest = 38
		genie.attack(p)

		dmg := genie.strength*0.5 + 40
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (50%% strength + 40), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("blast", func(t *testing.T) {
		p.hp = 100
		rolltest = 60
		genie.attack(p)

		dmg := genie.strength * 1.13
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (113%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("punch", func(t *testing.T) {
		p.hp = 100
		rolltest = 85
		genie.attack(p)

		dmg := genie.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestCelestial(t *testing.T) {
	p := &Player{}
	p.name = "player"
	p.perk = -1
	p.hp = 100
	p.hpcap = 100
	p.defense = 0
	p.effects = make(map[string]int)

	var celestial celestial
	celestial.strength = 10
	celestial.effects = make(map[string]int)

	t.Run("healing aura", func(t *testing.T) {
		rolltest = 1
		celestial.attack(p)

		if celestial.effects["heal aura"] != 5 {
			t.Errorf("should get heal aura effect for 5 turns, got %d", celestial.effects["heal aura"])
		}
	})

	t.Run("holy fire burning", func(t *testing.T) {
		rolltest = 12
		p.hp = 100
		p.defense = 0
		celestial.attack(p)

		dmg := celestial.strength * 0.8
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (80%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["burning"] != 3 {
			t.Errorf("should give burning effect for 3 turns, got %d", p.effects["burning"])
		}
	})

	t.Run("basic attack", func(t *testing.T) {
		rolltest = 37
		p.hp = 100
		celestial.attack(p)

		dmg := celestial.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

}

func TestShapeshift(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.hpcap = 100
	p.defense = 8
	p.strength = 35

	shapeshift := newShapeshift()
	shapeshift.attack(p)

	attr := shapeshift.attr()

	if attr.hp != p.hp*1.25 {
		t.Errorf("should copy hp by 125%%, want %.1f, got: %.1f", p.hp*1.25, attr.hp)
	}

	if attr.hpcap != p.hpcap*1.25 {
		t.Errorf("should copy hpcap by 125%%, want %.1f, got: %.1f", p.hpcap*1.25, attr.hpcap)
	}

	if attr.strength != p.strength*1.25 {
		t.Errorf("should copy strength by 125%%, want %.1f, got: %.1f", p.strength*1.25, attr.strength)
	}

	if attr.defense != p.defense*1.25 {
		t.Errorf("should copy defense by 125%%, want %.1f, got: %.1f", p.defense*1.25, attr.defense)
	}

	if attr.agility != p.agility*1.25 {
		t.Errorf("should copy agility by 125%%, want %.1f, got: %.1f", p.agility*1.25, attr.agility)
	}
}

func TestUndead(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.perk -= 1
	p.hpcap = 100
	p.effects = make(map[string]int)

	var undead undead
	undead.strength = 10

	t.Run("vomit poison", func(t *testing.T) {
		rolltest = 1
		undead.attack(p)

		dmg := undead.strength * 0.4
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (40%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["poisoned"] != 3 {
			t.Errorf("should give poisoned effect for 3 turns, got %d", p.effects["poisoned"])
		}
	})

	t.Run("bite", func(t *testing.T) {
		rolltest = 28
		p.hp = 100
		undead.attack(p)

		dmg := undead.strength * 1.33
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (133%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		rolltest = 35
		p.hp = 100
		undead.attack(p)

		dmg := undead.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("against deadman perk", func(t *testing.T) {
		rolltest = 35
		p.perk = 5
		p.hp = 100
		undead.attack(p)

		dmg := undead.strength * 0.67
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (33%% reduction), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestScorpion(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.perk -= 1
	p.hpcap = 100
	p.defense = 10
	p.effects = make(map[string]int)

	var scorpion scorpion
	scorpion.strength = 20

	t.Run("basic", func(t *testing.T) {
		rolltest = 1
		scorpion.attack(p)

		dmg := 20 - 10 + p.defense*0.3
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (ignore 30%% defense), got %.1f", dmg, 100-p.hp)
		}

		p.hp = 100
		p.defense = 100
		scorpion.strength = 100
		scorpion.attack(p)

		if p.hp != 70 {
			t.Errorf("damage should be 30 (ignore 30%% defense), got %.1f", 100-p.hp)
		}
	})

	t.Run("venom", func(t *testing.T) {
		rolltest = 24
		scorpion.attack(p)

		if p.effects["poisoned severe"] != 3 {
			t.Errorf("should give poisoned severe effect for 3 turns, got %d", p.effects["poisoned severe"])
		}
	})
}

func TestGoblin(t *testing.T) {
	p := &Player{}
	p.hp = 100
	p.perk -= 1
	p.hpcap = 100
	p.effects = make(map[string]int)

	var goblin goblin
	goblin.strength = 10

	t.Run("powder", func(t *testing.T) {
		rolltest = 1
		goblin.attack(p)

		if p.effects["confused"] != 3 {
			t.Errorf("should get confused effect for 3 turns, got %d", p.effects["confused"])
		}

		dmg := goblin.strength * 0.45
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (45%% strentgh), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("leap", func(t *testing.T) {
		p.hp = 100
		rolltest = 12
		goblin.attack(p)

		dmg := goblin.strength * 1.25
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (125%% strentgh), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("rapid strike", func(t *testing.T) {
		p.hp = 100
		rolltest = 24
		goblin.attack(p)

		dmg := goblin.strength * 1.1
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (110%% strentgh), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		p.hp = 100
		rolltest = 36
		goblin.attack(p)

		dmg := goblin.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% defense), got %.1f", dmg, 100-p.hp)
		}
	})
}
