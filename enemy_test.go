package main

import (
	"math"
	"testing"
)

func TestKnight(t *testing.T) {
	var knight knight
	knight.hpcap = 100
	knight.strength = 10
	knight.effects = make(map[string]int)

	t.Run("reinforce armor", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		tempdef := knight.defense
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
		p := newTestPlayer()
		knight.attack(p)

		dmg := knight.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}

		if knight.effects["strengthen"] != 4 {
			t.Errorf("should get strengthen effect for 4 turns, got %d", knight.effects["strengthen"])
		}

		clear(knight.effects)
	})

	t.Run("attack roll", func(t *testing.T) {
		rolltest = 30
		p := newTestPlayer()
		knight.attack(p)

		dmg := knight.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestWizard(t *testing.T) {
	var wizard wizard
	wizard.hpcap = 100
	wizard.strength = 10
	wizard.effects = make(map[string]int)

	t.Run("cast healing", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		wizard.attack(p)

		heal := wizard.hpcap * 0.2

		if wizard.hp != heal {
			t.Errorf("should heal by 20 (20%% hpcap), got %.1f", wizard.hp)
		}
	})

	t.Run("immunity", func(t *testing.T) {
		rolltest = 15
		p := newTestPlayer()
		wizard.attack(p)

		if wizard.effects["immunity"] != 2 {
			t.Errorf("should get immunity effect for 2 turns, got %d", wizard.effects["immunity"])
		}
	})

	t.Run("barrier", func(t *testing.T) {
		rolltest = 20
		p := newTestPlayer()
		wizard.attack(p)

		if wizard.effects["barrier"] != 4 {
			t.Errorf("should get barrier effect for 4 turns, got %d", wizard.effects["barrier"])
		}
	})

	t.Run("confuse", func(t *testing.T) {
		rolltest = 30
		p := newTestPlayer()
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
		p := newTestPlayer()
		wizard.attack(p)

		dmg := wizard.strength * 3
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (300%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("fireball", func(t *testing.T) {
		rolltest = 50
		p := newTestPlayer()
		wizard.attack(p)

		dmg := 100 - p.hp
		// is scaled but assume its stage 0
		if dmg < 0 || dmg > 100 {
			t.Errorf("damage should be between 0-100, got %.1f", dmg)
		}
	})

	t.Run("lightning", func(t *testing.T) {
		rolltest = 70
		p := newTestPlayer()
		wizard.attack(p)

		dmg := 100 - p.hp
		// is scaled but assume its stage 0
		if dmg < 40 || dmg > 115 {
			t.Errorf("damage should be between 40-115, got %.1f", dmg)
		}
	})

	t.Run("summon meteor", func(t *testing.T) {
		rolltest = 85
		p := newTestPlayer()
		p.hp = 200
		p.hpcap = 200
		wizard.attack(p)

		if 200-p.hp != 120 {
			t.Errorf("damage should be 120, got %.1f", 100-p.hp)
		}
	})

	t.Run("staff attack", func(t *testing.T) {
		rolltest = 90
		p := newTestPlayer()
		wizard.attack(p)

		dmg := wizard.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f, got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestChangeling(t *testing.T) {
	p := newTestPlayer()
	attr := &changeling{}
	attr.attack(p)

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
	var vampire vampire
	vampire.hp = 100
	vampire.hpcap = 100
	vampire.strength = 10

	t.Run("exposed to sunlight", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		vampire.defense = 999
		vampire.attack(p)
		vampire.defense = 0

		if 100-vampire.hp != 10 {
			t.Errorf("damage should be 10 (10%% of hp and ignore defense), got %.1f", 100-vampire.hp)
		}
	})

	t.Run("bites lifesteal", func(t *testing.T) {
		rolltest = 8
		p := newTestPlayer()
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
		p := newTestPlayer()
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
		p := newTestPlayer()
		vampire.attack(p)

		dmg := vampire.attr().strength * 1.2
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (120%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("claws", func(t *testing.T) {
		rolltest = 85
		p := newTestPlayer()
		vampire.attack(p)

		dmg := vampire.attr().strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

}

func TestDemon(t *testing.T) {
	var demon demon
	demon.hpcap = 100
	demon.strength = 10

	t.Run("soul absorption", func(t *testing.T) {
		p := newTestPlayer()
		p.defense = 99999 // to make sure
		rolltest = 1
		dmg := demon.strength + p.hp*0.03
		demon.attack(p)

		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (take 3%% hp and ignore defense), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("hell fire burning", func(t *testing.T) {
		rolltest = 60
		p := newTestPlayer()
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
		p := newTestPlayer()
		demon.attack(p)

		dmg := demon.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestShardling(t *testing.T) {
	t.Run("spawner", func(t *testing.T) {
		if newShardling().attr().effects["reflect"] == 0 {
			t.Error("should start with reflect")
		}
	})

	shardling := &shardling{}
	shardling.strength = 10
	shardling.hp = 100

	t.Run("ram", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		shardling.attack(p)

		dmg := shardling.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("crystal limbs", func(t *testing.T) {
		rolltest = 40
		p := newTestPlayer()
		shardling.attack(p)

		dmg := shardling.strength * 1.1
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (110%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("volley of shards", func(t *testing.T) {
		rolltest = 50
		p := newTestPlayer()
		shardling.attack(p)

		mindmg := shardling.strength
		maxdmg := mindmg + 20
		dmg := 100 - p.hp
		if dmg < mindmg || dmg > maxdmg {
			t.Errorf("damage should be between %.1f-%.1f, got %.1f", mindmg, maxdmg, dmg)
		}
	})

	t.Run("spike", func(t *testing.T) {
		rolltest = 80
		p := newTestPlayer()
		shardling.attack(p)

		dmg := shardling.strength * 1.25
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (125%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestGenie(t *testing.T) {
	var genie genie
	genie.strength = 10
	genie.effects = make(map[string]int)

	t.Run("hp curse", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		genie.attack(p)

		if p.hpcap >= 100 {
			t.Error("hp cap should be reduced")
		}
	})

	t.Run("strength curse", func(t *testing.T) {
		rolltest = 5
		p := newTestPlayer()
		genie.attack(p)

		if p.strength >= 50 {
			t.Error("strength should be reduced")
		}
	})

	t.Run("defense curse", func(t *testing.T) {
		rolltest = 10
		p := newTestPlayer()
		genie.attack(p)

		if p.defense >= 15 {
			t.Error("defense should be reduced")
		}
	})

	t.Run("energy cap curse", func(t *testing.T) {
		rolltest = 15
		p := newTestPlayer()
		p.name = "player"
		genie.attack(p)

		if p.energycap >= 20 {
			t.Error("energy cap should be reduced")
		}
	})

	t.Run("illusion", func(t *testing.T) {
		p := newTestPlayer()
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

	t.Run("force-field", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 31
		genie.attack(p)

		if genie.effects["force-field"] != 5 {
			t.Errorf("should get force-field effect for 5 turns, got %d", genie.effects["force-field"])
		}
	})

	t.Run("sandstorm", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 38
		genie.attack(p)

		dmg := genie.strength*0.5 + 40
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (50%% strength + 40), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("blast", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 60
		genie.attack(p)

		dmg := genie.strength * 1.13
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (113%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("punch", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 85
		genie.attack(p)

		dmg := genie.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestCelestial(t *testing.T) {
	var celestial celestial
	celestial.strength = 10
	celestial.effects = make(map[string]int)

	t.Run("healing aura", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		celestial.attack(p)

		if celestial.effects["heal aura"] != 5 {
			t.Errorf("should get heal aura effect for 5 turns, got %d", celestial.effects["heal aura"])
		}
	})

	t.Run("blinding light", func(t *testing.T) {
		rolltest = 10
		p := newTestPlayer()
		celestial.attack(p)

		dmg := celestial.strength * 0.7
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f 780%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["disoriented"] != 3 {
			t.Errorf("should give disoriented effect for 3 turns, got %d", p.effects["disoriented"])
		}
	})

	t.Run("holy fire burning", func(t *testing.T) {
		rolltest = 30
		p := newTestPlayer()
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
		rolltest = 40
		p := newTestPlayer()
		celestial.attack(p)

		dmg := celestial.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

}

func TestShapeshift(t *testing.T) {
	p := newTestPlayer()
	attr := &shapeshift{}
	attr.attack(p)

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
	var undead undead
	undead.strength = 10

	t.Run("vomit poison", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		undead.attack(p)

		dmg := undead.strength * 0.4
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (40%% strength), got %.1f", dmg, 100-p.hp)
		}

		if p.effects["poisoned"] != 3 {
			t.Errorf("should give poisoned effect for 3 turns, got %d", p.effects["poisoned"])
		}
	})

	t.Run("fellow undeads", func(t *testing.T) {
		rolltest = 8
		p := newTestPlayer()
		undead.attack(p)

		mindmg := undead.strength
		maxdmg := mindmg + 30
		dmg := 100 - p.hp
		if dmg < mindmg || dmg > maxdmg {
			t.Errorf("damage should be between %.1f-%.1f, got %.1f", mindmg, maxdmg, dmg)
		}
	})

	t.Run("bite", func(t *testing.T) {
		rolltest = 28
		p := newTestPlayer()
		undead.attack(p)

		dmg := undead.strength * 1.33
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (133%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		rolltest = 35
		p := newTestPlayer()
		undead.attack(p)

		dmg := undead.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestScorpion(t *testing.T) {
	var scorpion scorpion
	scorpion.strength = 20

	t.Run("basic", func(t *testing.T) {
		p := newTestPlayer()
		p.defense = 10
		scorpion.attack(p)

		dmg := scorpion.strength - p.defense + p.defense*0.3
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
		rolltest = 1
		p := newTestPlayer()
		scorpion.attack(p)

		if p.effects["poisoned severe"] != 3 {
			t.Errorf("should give poisoned severe effect for 3 turns, got %d", p.effects["poisoned severe"])
		}
	})
}

func TestGoblin(t *testing.T) {
	var goblin goblin
	goblin.strength = 10

	t.Run("powder", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		goblin.attack(p)

		if p.effects["confused"] != 3 {
			t.Errorf("should get confused effect for 3 turns, got %d", p.effects["confused"])
		}

		dmg := goblin.strength * 0.45
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (45%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("leap", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 12
		goblin.attack(p)

		dmg := goblin.strength * 1.25
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (125%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("rapid strike", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 24
		goblin.attack(p)

		dmg := goblin.strength * 1.1
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (110%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 36
		goblin.attack(p)

		dmg := goblin.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestInfernal(t *testing.T) {
	t.Run("spawner", func(t *testing.T) {
		if newInfernal().attr().effects["burning immunity"] == 0 {
			t.Error("should start with burning immunity effect")
		}
	})

	var infernal infernal
	infernal.strength = 10

	t.Run("burning severe", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		infernal.attack(p)

		if p.effects["burning severe"] != 2 {
			t.Errorf("should inflict severe burning for 2 turns, got %d", p.effects["burning severe"])
		}

		dmg := infernal.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("enhanced energy", func(t *testing.T) {
		str := infernal.strength
		rolltest = 30
		infernal.attack(&infernal)

		if infernal.strength == str {
			t.Error("should buff strength")
		}
	})

	t.Run("magma shield", func(t *testing.T) {
		def := infernal.defense
		rolltest = 40
		infernal.attack(&infernal)

		if infernal.defense == def {
			t.Error("should buff defense ")
		}
	})

	t.Run("basic", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 50
		infernal.attack(p)

		if p.effects["burning"] != 2 {
			t.Errorf("should inflict burning for 2 turns, got %d", p.effects["burning"])
		}

		dmg := infernal.strength
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestVineMonster(t *testing.T) {
	t.Run("spawner", func(t *testing.T) {
		if newVineMonster().attr().effects["reflect low"] == 0 {
			t.Error("should start with reflect low effect")
		}
	})

	var vine vineMonster
	vine.strength = 10

	t.Run("ensnare", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		vine.attack(p)

		if p.effects["stunned"] != 2 {
			t.Errorf("should inflict stun for 2 turns, got %d", p.effects["stunned"])
		}

		dmg := vine.strength * 0.6
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (60%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("burst of thorns", func(t *testing.T) {
		rolltest = 30
		p := newTestPlayer()
		vine.attack(p)

		dmg := vine.strength * 0.6
		if 100-p.hp < dmg {
			t.Errorf("damage should be atleast %.1f (60%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 45
		vine.attack(p)

		dmg := vine.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestArcticWarrior(t *testing.T) {
	t.Run("spawner", func(t *testing.T) {
		if newArcticWarrior().attr().effects["frozen immunity"] == 0 {
			t.Error("should start with frozen immunity effect")
		}
	})

	var arctic arcticWarrior
	arctic.strength = 10

	t.Run("frozen attack", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		arctic.attack(p)

		if p.effects["frozen"] != 2 {
			t.Errorf("should inflict frozen for 2 turns, got %d", p.effects["frozen"])
		}
	})

	t.Run("snowstorm", func(t *testing.T) {
		rolltest = 20
		p := newTestPlayer()
		arctic.attack(p)

		dmg := arctic.strength * 0.2
		if 100-p.hp < dmg {
			t.Errorf("damage should be atleast %.1f (20%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("avalanche", func(t *testing.T) {
		rolltest = 35
		p := newTestPlayer()
		arctic.attack(p)

		dmg := arctic.strength
		if 100-p.hp < dmg {
			t.Errorf("damage should be atleast %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 43
		arctic.attack(p)

		dmg := arctic.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestJungleWarrior(t *testing.T) {
	t.Run("spawner", func(t *testing.T) {
		if newJungleWarrior().attr().effects["poison immunity"] == 0 {
			t.Error("should start with poison immunity effect")
		}
	})

	var jungle jungleWarrior
	jungle.strength = 10
	jungle.effects = make(map[string]int)

	t.Run("venom spear", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		jungle.attack(p)

		if p.effects["shiver"] != 3 {
			t.Errorf("should inflict shiver for 3 turns, got %d", p.effects["shiver"])
		}

		dmg := jungle.strength * 1
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("ancestral spirits", func(t *testing.T) {
		rolltest = 14
		p := newTestPlayer()
		jungle.attack(p)

		if jungle.effects["vitality"] != 4 {
			t.Errorf("should get vitality for 4 turns, got %d", jungle.effects["vitality"])
		}
		clear(jungle.effects)

		dmg := jungle.strength * 0.8
		dmg += dmg * 0.05 // vitality
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (80%% strength + 5%%), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("ancestral roar", func(t *testing.T) {
		rolltest = 28
		p := newTestPlayer()
		jungle.attack(p)

		if jungle.effects["vitality"] != 2 {
			t.Errorf("should get vitality for 2 turns, got %d", jungle.effects["vitality"])
		}
		clear(jungle.effects)

		if p.effects["shiver"] != 2 {
			t.Errorf("should inflict shiver for 2 turns, got %d", p.effects["shiver"])
		}

		dmg := jungle.strength * 1.3
		dmg += dmg * 0.05 // vitality
		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (130%% strength + 5%%), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("leap", func(t *testing.T) {
		rolltest = 35
		p := newTestPlayer()
		jungle.attack(p)

		dmg := jungle.strength * 1.2
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (120%% strength), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		p := newTestPlayer()
		rolltest = 45
		jungle.attack(p)

		dmg := jungle.strength
		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength), got %.1f", dmg, 100-p.hp)
		}
	})
}

func TestLeechMonster(t *testing.T) {
	var leech leechMonster
	leech.strength = 10

	t.Run("bleeding", func(t *testing.T) {
		rolltest = 1
		p := newTestPlayer()
		p.defense = 1
		dmg := leech.strength + p.hp*0.1 - p.defense*0.75
		dmg -= p.defense
		leech.attack(p)

		if p.effects["bleeding"] != 10 {
			t.Errorf("should inflict bleeding with 10 severity, got %d", p.effects["bleeding"])
		}

		if dmg != 100-p.hp {
			t.Errorf("damage should be %.1f (100%% strength, 10%% target hp, x1.75 target defense value), got %.1f", dmg, 100-p.hp)
		}
	})

	t.Run("basic", func(t *testing.T) {
		rolltest = 25
		p := newTestPlayer()
		p.defense = 1

		dmg := leech.strength + p.hp*0.2 - p.defense*1.3
		dmg -= p.defense
		leech.attack(p)

		if !equal(dmg, 100-p.hp) {
			t.Errorf("damage should be %.1f (100%% strength, 20%% target hp, x2.3 target defense value), got %.1f", dmg, 100-p.hp)
		}
	})
}

// quick fix floating issue
func equal(a, b float64) bool {
	return math.Round(a*100) == math.Round(b*100)
}
