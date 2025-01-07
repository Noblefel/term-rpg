package main

import "testing"

func TestEnemy_TakeDamage(t *testing.T) {
	enemy := newKnight()
	before := enemy.attr().hp
	enemy.damage(10)

	if enemy.attr().hp == before {
		t.Error("hp should be reduced")
	}
}
