package entity

import (
	"fmt"
	"math/rand"
)

type Entity interface {
	// Attack sums the attack with random value. In tests, replaced random value with 10.
	Attack(e Entity) (float32, string)
	// TakeDamage reduces dmg with the defense and decrements the hp
	TakeDamage(dmg float32) float32
	// DropLoot multiplies random value with drop rate. In tests, replaced random value with 10.
	DropLoot() float32
	// RecoverHP adds hp to the entity, cannot exceeds hp cap
	RecoverHP(n float32)

	Attr() *Base
}

type Base struct {
	Name     string
	Hp       float32
	Att      float32
	Def      float32
	HpCap    float32
	DmgReduc float32
	DropRate float32

	IsDefending bool
	isTesting   bool
}

func (b *Base) Attack(e Entity) (float32, string) {
	dmg := b.Att

	if b.isTesting {
		dmg += 10
	} else {
		dmg += rand.Float32() * 10
	}

	dmg = e.TakeDamage(dmg)

	return dmg, fmt.Sprintf("%s attacked (%.1f dmg)", b.Name, dmg)
}

func (b *Base) TakeDamage(dmg float32) float32 {
	dmg -= b.Def + (dmg * b.DmgReduc)

	if b.IsDefending {
		dmg -= dmg * 0.2
	}

	if dmg < 0 {
		return 0
	}

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

func (b *Base) RecoverHP(n float32) {
	hp := b.Hp + n

	if hp > b.HpCap {
		hp = b.HpCap
	}

	b.Hp = hp
}

func (b *Base) Attr() *Base { return b }
