package entities

import "math/rand/v2"

type Entity interface {
	Attack(Entity) float32
	TakeDamage(float32) float32
	GetAttributes() Attributes
}

type Attributes struct {
	Name     string
	HP       float32
	HPCap    float32
	Strength float32
	Defense  float32
}

func (attr Attributes) Attack(e Entity) float32 {
	dmg := attr.Strength + rand.Float32()*5
	return e.TakeDamage(dmg)
}

func (attr *Attributes) TakeDamage(dmg float32) float32 {
	dmg = max(dmg-attr.Defense, 1)
	attr.HP -= dmg
	return dmg
}

func (attr Attributes) GetAttributes() Attributes {
	return attr
}

type Player struct {
	Attributes
	Gold    float32
	isHavoc bool
}

func NewPlayer(perk int) *Player {
	var player Player
	player.Name = "Player"
	player.HP = 100
	player.HPCap = 100
	player.Defense = 1
	player.Strength = 10
	player.Gold = 50

	if perk == 0 {
		player.Defense += 2
	}

	if perk == 1 {
		player.HP = 85
		player.HPCap = 85
		player.isHavoc = true
	}

	return &player
}

func (p *Player) Attack(e Entity) float32 {
	if p.isHavoc {
		dmg := p.Strength + rand.Float32()*5
		dmg += dmg * 0.2
		return e.TakeDamage(dmg)
	}

	return p.Attributes.Attack(e)
}

func SpawnRandom() Entity {
	switch rand.IntN(1) {
	case 0:
		return NewKnight()
	}

	panic("unhandled spawn")
}

type Knight struct {
	Attributes
}

func NewKnight() *Knight {
	var knight Knight
	knight.Name = "Knight"
	knight.HP = 100
	knight.HPCap = 100
	knight.Defense = 5
	knight.Strength = 5

	return &knight
}
