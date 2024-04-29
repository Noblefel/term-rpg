package entity

import "math/rand"

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

	if dmg < 0 {
		dmg = 0
	}

	return dmg
}

func (b *base) takeDamage(dmg float32) float32 {
	dmg -= b.Def + (dmg * b.DmgReduc)

	if b.GuardTurns > 0 {
		dmg -= dmg * 0.2
	}

	if dmg < 0 {
		return 0
	}

	if b.Hp-dmg < 0 {
		b.Hp = 0
	} else {
		b.Hp -= dmg
	}

	return dmg
}

func (b *base) heal(n float32) {
	hp := b.Hp + n

	if hp > b.HpCap {
		hp = b.HpCap
	}

	b.Hp = hp
}
