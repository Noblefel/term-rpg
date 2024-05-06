package entity

import (
	"math/rand"
)

type base struct {
	Hp         float32
	Att        float32
	Def        float32
	HpCap      float32
	DmgReduc   float32
	GuardTurns int
	FuryTurns  int
	isTesting  bool
}

func (b *base) attack() float32 {
	dmg := b.Att

	if b.isTesting {
		dmg += 10
	} else {
		dmg += rand.Float32() * 10
	}

	return max(dmg, 0)
}

func (b *base) takeDamage(dmg float32) float32 {
	dmg -= b.Def + (dmg * b.DmgReduc)

	if b.GuardTurns > 0 {
		dmg -= dmg * 0.2
	}

	dmg = max(dmg, 0)
	b.Hp -= min(dmg, b.Hp)

	return dmg
}

func (b *base) heal(n float32) {
	b.Hp = min(b.HpCap, b.Hp+n)
}

func (b *base) guard(extraTurn int) {
	b.GuardTurns = 3 + extraTurn
}

func (b *base) fury(extraTurn int) float32 {
	b.FuryTurns = max(3, b.FuryTurns) + extraTurn
	sacrifice := 1 + (b.Hp * 0.1) + (rand.Float32() * 4)
	b.Hp -= sacrifice
	return sacrifice
}

func (b *base) decrementEffects() {
	b.GuardTurns--
	b.FuryTurns--
}
