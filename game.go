package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"strings"
	"text/tabwriter"

	"github.com/Noblefel/term-rpg/entities"
)

type Game struct {
	writer io.Writer
	player *entities.Player
	stage  int
	info   string
}

// wrapper for Fprintf
func (game *Game) printf(s string, a ...any) {
	fmt.Fprintf(game.writer, s, a...)
}

// prints info if there is any and clears it
func (game *Game) printInfo() {
	if game.info == "" {
		return
	}

	game.printf(game.info)
	game.printf("\n")
	game.info = ""
}

// prints a colored bar: full = green, half = orange, low = red
func (game *Game) printBar(length int, n, cap float32) {
	n = max(n, 0)
	cap = max(cap, 0)

	if n > cap {
		n = cap
	}

	percentage := (n / cap) * 100
	bars := int(percentage) * length / 100
	color := "38;2;247;17;17"

	if percentage > 60 {
		color = "38;2;80;230;95"
	} else if percentage > 30 {
		color = "38;2;245;185;0"
	}

	game.printf("\033[%sm%s\033[0m\033[%sm%s\033[0m\n",
		color,
		strings.Repeat("█", bars),
		"38;2;60;60;60",
		strings.Repeat("█", max(length-bars, 0)),
	)
}

// MENU: perk selection
func (game *Game) selectPerks() {
	game.printf("┌───────────────────────────────────────────┐\n")
	game.printf("│              \033[1;33mSelect your perk\033[0m             │\n")
	game.printf("└───────────────────────────────────────────┘\n")

	choice := NewChoice(game.writer,
		"[1] 🛡️  Resiliency (+2 defense point)",
		"[2] ⚔️  Havoc (+20% damage, but -15 HP cap)",
		// "[3] ⌛ Temporal (+1 extra turn for bonus effects)",
	)

	perk := choice.Select()
	game.player = entities.NewPlayer(perk)
}

// MENU: main menu
func (game *Game) menu() {
	game.printf("\033[H\033[J")
	game.printf("  ┌───────────────────────────────────────────┐\n")
	game.printf("  │                 \033[1mMain Menu\033[0m                 │\n")
	game.printf("  └───────────────────────────────────────────┘\n")
	game.printBar(50, game.player.HP, game.player.HPCap)
	game.printf("--------\n")
	game.printf("❤️  \033[1;31m%0.1f\033[0m\n", game.player.HP)
	game.printf("💰 \033[33m$%0.f\033[0m\n", game.player.Gold)
	game.printf("🗺️  Stage: %d\n", game.stage)
	game.printf("--------\n")
	game.printInfo()

	choice := NewChoice(game.writer,
		"[1] 🗺️   Next Stage",
		"[2] 📋  View Attributes",
		"[3] 🛏️   Rest ($5)",
		"[4] 💪  Train ($10)",
	)

	switch choice.Select() {
	case 0:
		game.battle()
	case 1:
		game.attributes()
	case 2:
		if game.player.Gold < 5 {
			game.info = "⚠️  \033[31mYou don't have enough money\033[0m\n"
			return
		}

		n := 10 + rand.Float32()*15
		game.player.HP = min(game.player.HPCap, n+game.player.HP)
		game.player.Gold -= 5
		game.info = fmt.Sprintf("💬 Your rest helped you recover \033[1;32m%.1f hp\033[0m\n", n)
	case 3:
		if game.player.Gold < 10 {
			game.info = "⚠️  \033[31mYou don't have enough money\033[0m\n"
			return
		}

		game.player.Gold -= 10

		if rand.IntN(100) > 30 {
			game.info = "💬 Training did not yield any result\n"
			return
		}

		switch rand.IntN(3) {
		case 0:
			n := 1 + rand.Float32()*5
			game.player.HPCap += n
			game.info = fmt.Sprintf("💪 Hp cap increased by \033[1m%.1f\033[0m\n", n)
		case 1:
			n := 0.5 + rand.Float32()*2
			game.player.Strength += n
			game.info = fmt.Sprintf("💪 Strength increased by \033[1m%.1f\033[0m\n", n)
		case 2:
			n := 0.5 + rand.Float32()*2
			game.player.Defense += n
			game.info = fmt.Sprintf("💪 Defense increased by \033[1m%.1f\033[0m\n", n)
		}
	}
}

// MENU: player attributes list
func (game *Game) attributes() {
	game.printf("\033[H\033[J")
	game.printf("┌───────────────────────────────────────────┐\n")
	game.printf("│              \033[1mYour attributes\033[0m              │\n")
	game.printf("└───────────────────────────────────────────┘\n")

	tw := tabwriter.NewWriter(game.writer, 5, 0, 5, ' ', 0)
	fmt.Fprint(tw, "Max HP\tStrength\tDefense\n")
	fmt.Fprint(tw, "------\t------\t-------\n")
	fmt.Fprintf(tw, "%.1f\t%.1f\t%.1f\n\n",
		game.player.HPCap,
		game.player.Strength,
		game.player.Defense,
	)
	tw.Flush()

	game.printf("\n")
	NewChoice(game.writer, "Go back").Select()
}

// MENU: battle
func (game *Game) battle() {
	var playerLog string
	var enemyLog string

	enemy := entities.SpawnRandom()

	choice := NewChoice(game.writer,
		"[1] ⚔️  Attack",
		"[2] 🛡️  Guard",
		"[3] 🔥 Fury",
		"[4] 🏃 Flee",
	)

	for {
		enemyAttributes := enemy.GetAttributes()

		if game.player.HP <= 0 {
			game.player.Gold += 10 + rand.Float32()*10
			game.info = "You have lost\n"
			return
		} else if enemyAttributes.HP <= 0 {
			game.player.Gold += 10 + rand.Float32()*50
			game.info = "You have won\n"
			return
		}

		game.printf("\033[H\033[J")
		game.printBar(50, game.player.HP, game.player.HPCap)
		game.printf("----------------\n")
		game.printf("❤️  \033[1;31m%0.1f\033[0m | \033[1mYou\033[0m\n", game.player.HP)
		game.printf("-> %s\n", playerLog)
		playerLog = ""
		game.printf("\n")

		game.printBar(50, enemyAttributes.HP, enemyAttributes.HPCap)
		game.printf("----------------\n")
		game.printf("❤️  \033[1;31m%0.1f\033[0m | \033[1m%s\033[0m\n",
			enemyAttributes.HP,
			enemyAttributes.Name,
		)
		game.printf("-> %s\n", enemyLog)
		enemyLog = ""
		game.printInfo()
		game.printf("\n")

		switch choice.Select() {
		case 0:
			dmg := game.player.Attack(enemy)
			playerLog = fmt.Sprintf("You attacked, dealing %.1f damage", dmg)
		case 1, 2:
			game.info = "\033[31mNot implemented\033[0m"
			continue
		case 3:
			game.info = "You have escaped\n"
			return
		}

		dmg := enemy.Attack(game.player)
		enemyLog = fmt.Sprintf("Enemy attacked, dealing %.1f damage", dmg)
	}
}
