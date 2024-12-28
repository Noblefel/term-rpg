package main

import (
	"io"
	"testing"
)

func init() { out = io.Discard }

func TestPlayer_Attack(t *testing.T) {
	enemy := new(Knight)
	enemy.hp = 2000
	enemy.hpcap = 2000

	before := enemy.hp
	player := NewPlayer(0)
	player.Attack(enemy)

	if enemy.hp == before {
		t.Error("enemy hp should be reduced")
	}

	t.Run("with havoc perk", func(t *testing.T) {
		enemy.hp = 2000
		player.perk = 1
		player.strength = 1000
		player.Attack(enemy)
		dmg := enemy.hpcap - enemy.hp

		if dmg < 1200 || dmg > 1207 {
			t.Errorf("damage should be around 1200 (20%% increase), got: %.1f", dmg)
		}
	})

	t.Run("with berserk perk and 30% hp", func(t *testing.T) {
		enemy.hp = 2000
		player.perk = 2
		player.hp = 30
		player.hpcap = 100
		player.strength = 1000
		player.Attack(enemy)
		dmg := enemy.hpcap - enemy.hp

		if dmg < 1100 || dmg > 1107 {
			t.Errorf("damage should be around 1300 (10%% increase), got: %.1f", dmg)
		}
	})
}

func TestPlayer_TakeDamage(t *testing.T) {
	player := Player{hp: 10, defense: 5, perk: -1}
	player.TakeDamage(10)

	if player.hp != 5 {
		t.Errorf("hp should be 5, got %.1f", player.hp)
	}

	t.Run("with resiliency perk", func(t *testing.T) {
		player.perk = 0
		player.defense = 1
		player.hp = 1000
		player.TakeDamage(100)
		want := 1000 - 90 + player.defense

		if player.hp != want {
			t.Errorf("hp should be %.1f (10%% dmg reduction), got: %.1f", want, player.hp)
		}
	})

	t.Run("with berserk perk and 20% hp", func(t *testing.T) {
		player.perk = 2
		player.defense = 1
		player.hp = 200
		player.hpcap = 1000
		player.TakeDamage(100)

		if player.hp != 120+player.defense {
			t.Errorf("hp should be %.1f (20%% dmg reduction), got: %.1f", 125+player.defense, player.hp)
		}
	})
}

func TestPlayer_Rest(t *testing.T) {
	player := Player{hp: 0, hpcap: 100, gold: 10}
	player.Rest()

	if player.hp < 15 {
		t.Errorf("should atleast heal 15 hp, got %.1f", player.hp)
	}

	if player.gold != 5 {
		t.Errorf("should atleast cost 5 gold, got %d", 10-player.gold)
	}
}
