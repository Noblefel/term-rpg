package combat

import (
	"math/rand"
)

type Combatant interface {
	// Attack sums the attack with random value. In tests, replaced random value with 10.
	Attack() float32
	// TakeDamage reduces dmg with the defense and decrements the hp
	TakeDamage(dmg float32) float32
	// DropLoot multiplies random value with drop rate. In tests, replaced random value with 10.
	// As of now, it will only drop money.
	DropLoot() float32

	Attr() Base
}

type Base struct {
	Name      string
	Hp        float32
	Att       float32
	Def       float32
	HpCap     float32
	DmgReduc  float32
	DropRate  float32
	isTesting bool
}

func (b Base) Attr() Base { return b }

func (b *Base) Attack() float32 {
	dmg := b.Att

	if b.isTesting {
		dmg += 10
	} else {
		dmg += rand.Float32() * 10
	}

	if dmg <= 0 {
		return 0
	}

	return dmg
}

func (b *Base) TakeDamage(dmg float32) float32 {
	dmg -= b.Def

	if dmg <= 0 {
		return 0
	}

	dmg -= dmg * b.DmgReduc

	b.Hp -= dmg
	return dmg
}

func (b *Base) DropLoot() float32 {
	var loot float32

	if b.isTesting {
		loot = 10
	} else {
		loot += 5 + rand.Float32()*30
	}

	loot *= b.DropRate

	if loot <= 0 {
		return 0
	}

	return loot
}
