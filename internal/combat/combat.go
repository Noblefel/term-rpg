package combat

import (
	"math/rand"
)

type Combatant interface {
	// Attack sums the att stat with random value. For testing would sum with 10 instead
	Attack() float32
	// TakeDamage reduces dmg with the def stat and decrements the hp
	TakeDamage(dmg float32) float32

	Data() Base
}

type Base struct {
	Name      string
	Hp        float32
	Att       float32
	Def       float32
	HpCap     float32
	DmgReduc  float32
	isTesting bool
}

func (b Base) Data() Base { return b }

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
