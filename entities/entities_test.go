package entities

import "testing"

func TestBaseAttack(t *testing.T) {
	var player Player
	player.Strength = 10

	var enemy Player
	enemy.HP = 100

	player.Attack(&enemy)

	if enemy.HP == 100 {
		t.Fatal("enemy hp is not affected")
	}
}

func TestBaseTakeDamage(t *testing.T) {
	var player Player
	player.HP = 100
	player.Defense = 10
	player.TakeDamage(20)

	if player.HP != 90 {
		t.Errorf("want player hp to be 90, got %.0f", player.HP)
	}
}
