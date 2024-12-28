package main

import "testing"

func TestEnemy_TakeDamage(t *testing.T) {
	enemy := NewKnight()
	before := enemy.Attr().hp
	enemy.TakeDamage(10)

	if enemy.Attr().hp == before {
		t.Error("hp should be reduced")
	}
}
