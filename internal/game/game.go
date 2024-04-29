package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/Noblefel/term-rpg/internal/display"
	"github.com/Noblefel/term-rpg/internal/entity"
)

type Game struct {
	scanner *bufio.Scanner
	dis     *display.Display
	p       *entity.Player
	stage   int
}

func New(scanner *bufio.Scanner, dis *display.Display) *Game {
	return &Game{
		scanner: scanner,
		dis:     dis,
		stage:   1,
	}
}

func (g *Game) Start() {
	perk := g.selectPerks()
	g.p = entity.NewPlayer(perk)

	for {
		g.menu()
	}
}

func (g *Game) menu() {
	display.Clear()
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "Main Menu")
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════\n")

	g.dis.Bar(g.p.Hp, g.p.HpCap)
	g.dis.Printf(g.dis.White, "❤️  Hp: ")
	g.dis.Printf(g.dis.Red, "%.1f\n", g.p.Hp)
	g.dis.Printf(g.dis.White, "💰 Money: ")
	g.dis.Printf(g.dis.Green, " %.1f\n", g.p.Money)
	g.dis.Printf(g.dis.White, "✨ Perk: ")
	perk := strings.Split(entity.Perks[g.p.Perk], "(")[0]
	fmt.Printf("%s\n", perk)
	g.dis.Printf(g.dis.White, "🌡️  Stage: ")
	fmt.Printf("%d\n\n", g.stage)

	fmt.Println("┌─────────────────────────────────────────────")
	fmt.Println("■ > 1. 🗺️   Next Stage")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 2. 📋  View Attributes")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 3. 🛏️   Rest ($5)")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 4. 💪  Train ($10)")
	fmt.Println("└─────────────────────────────────────────────")

	for {
		g.dis.Printf(g.dis.White, "Input: ")
		g.scanner.Scan()

		switch g.scanner.Text() {
		case "1":
			g.battle()
		case "2":
			g.attributes()
		case "3":
			g.rest()
		case "4":
			g.train()
		default:
			g.dis.Printf(g.dis.Red, "Invalid Option\n")
			continue
		}
		return
	}
}

func (g *Game) selectPerks() int {
	display.Clear()
	g.dis.Printf(g.dis.Green, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.Green, "Select your perk ✨")
	g.dis.Printf(g.dis.Green, "══════════════════════════════════════════════════\n")

	fmt.Println("┌─────────────────────────────────────────────")
	fmt.Printf("■ > 1. %s\n", entity.Perks[entity.GREED])
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Printf("■ > 2. %s\n", entity.Perks[entity.RESILIENCY])
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Printf("■ > 3. %s\n", entity.Perks[entity.HAVOC])
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Printf("■ > 4. %s\n", entity.Perks[entity.TEMPORAL])
	fmt.Println("└─────────────────────────────────────────────")

	for {
		g.dis.Printf(g.dis.White, "Input: ")
		g.scanner.Scan()

		t := g.scanner.Text()
		switch t {
		case "1", "2", "3", "4":
			perk, _ := strconv.Atoi(t)
			return perk
		}

		g.dis.Printf(g.dis.Red, "Invalid Perk\n")
	}
}

func (g *Game) attributes() {
	display.Clear()

	g.dis.Center(g.dis.Green, "Your attributes 📋")
	fmt.Println("──────────────────────────────────────────────")

	g.dis.Printf(g.dis.White, "Hp Cap: ")
	fmt.Printf("%.1f\n", g.p.HpCap)
	g.dis.Printf(g.dis.White, "Attack: ")
	fmt.Printf("%.1f\n", g.p.Att)
	g.dis.Printf(g.dis.White, "Defense: ")
	fmt.Printf("%.1f\n", g.p.Def)
	g.dis.Printf(g.dis.White, "Dmg Reduction: ")
	fmt.Printf("%.1f%%\n", g.p.DmgReduc*100)

	g.dis.Center(nil, "■ > Press enter to continue")
	g.scanner.Scan()
}

func (g *Game) rest() {
	display.Clear()

	if g.p.Money >= 5 {
		g.p.Money -= 5
		n := 5 + (g.p.HpCap * 0.1) + rand.Float32()*8
		g.p.Heal(n)

		g.dis.Center(g.dis.Green, "You decided to rest 💤\n")
		fmt.Println("──────────────────────────────────────────────")
		g.dis.Center(g.dis.White, "Effect: Recovered %.1f hp", n)
	} else {
		g.dis.Center(g.dis.Red, "You don't have enough money to rest ⚠️\n")
		fmt.Println("──────────────────────────────────────────────")
	}

	g.dis.Center(nil, "■ > Press enter to continue")
	g.scanner.Scan()
}

func (g *Game) train() {
	display.Clear()

	if g.p.Money < 10 {
		g.dis.Center(g.dis.Red, "You don't have enough money to train ⚠️\n")
		fmt.Println("──────────────────────────────────────────────")
		g.dis.Center(nil, "■ > Press enter to go back")
		g.scanner.Scan()
		return
	}

	g.p.Money -= 10
	if rand.Intn(100) < 30 {
		g.dis.Center(g.dis.Green, "Hard work pays off 💪\n")
		fmt.Println("──────────────────────────────────────────────")
		g.dis.Center(g.dis.White, g.p.Train(rand.Intn(4)))
	} else {
		g.dis.Center(g.dis.Red, "Training did not yield any result")
		fmt.Println("──────────────────────────────────────────────")
	}

	g.dis.Center(nil, "■ > Press enter to continue")
	g.scanner.Scan()
}
