package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/Noblefel/term-rpg/internal/combat"
	"github.com/Noblefel/term-rpg/internal/display"
)

const (
	WIN int = iota
	LOSE
	FLED
	NEXT
)

type Game struct {
	scanner *bufio.Scanner
	dis     *display.Display
	p       *combat.Player
	e       combat.Combatant
}

func New(scanner *bufio.Scanner, dis *display.Display) *Game {
	return &Game{
		scanner: scanner,
		dis:     dis,
	}
}

func (g *Game) Menu() {
	if g.p == nil {
		perk := g.selectPerks()
		g.p = combat.NewPlayer(perk)
	}

	display.Clear()
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "Main Menu")
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════\n")

	g.dis.Bar(g.p.Hp, g.p.HpCap)
	g.dis.Printf(g.dis.White, "\n❤️  Hp: ")
	g.dis.Printf(g.dis.Red, "%.1f\n", g.p.Hp)
	g.dis.Printf(g.dis.White, "💰 Money: ")
	g.dis.Printf(g.dis.Green, " %.1f\n", g.p.Money)
	g.dis.Printf(g.dis.White, "✨ Perk: ")
	perk := strings.Split(combat.Perks[g.p.Perk], "(")[0]
	fmt.Printf("%s\n", perk)

	fmt.Println("┌─────────────────────────────────────────────")
	fmt.Println("■ > 1. 🗺️   Next Stage")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 2. 🛏️   Rest ($5)")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 3. 🛑  Quit")
	fmt.Println("└─────────────────────────────────────────────")

	for {
		g.dis.Printf(g.dis.White, "Input: ")
		g.scanner.Scan()

		switch g.scanner.Text() {
		case "1":
			g.battle()
			g.Menu()
		case "2":
			g.rest()
			g.Menu()
		case "3":
			os.Exit(1)
		default:
			g.dis.Printf(g.dis.Red, "Invalid Option\n")
			continue
		}
	}
}

func (g *Game) selectPerks() int {
	display.Clear()
	g.dis.Printf(g.dis.Green, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.Green, "Select your perk 🔥")
	g.dis.Printf(g.dis.Green, "══════════════════════════════════════════════════\n")

	fmt.Println("┌─────────────────────────────────────────────")
	fmt.Println("■ > 1.", combat.Perks[combat.GREED])
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 2.", combat.Perks[combat.RESILIENCY])
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 3.", combat.Perks[combat.HAVOC])
	fmt.Println("└─────────────────────────────────────────────")

	var perk int
	for {
		g.dis.Printf(g.dis.White, "Input: ")
		g.scanner.Scan()

		n, err := strconv.Atoi(g.scanner.Text())
		if err != nil {
			g.dis.Printf(g.dis.Red, "Should be a number\n")
			continue
		}

		_, ok := combat.Perks[n]
		if !ok {
			g.dis.Printf(g.dis.Red, "Invalid Perk\n")
			continue
		}

		perk = n
		break
	}

	return perk
}

func (g *Game) rest() {
	display.Clear()

	if g.p.Money > 5 {
		g.p.Money -= 5
		n := 3 + (g.p.HpCap * 0.05) + rand.Float32()*5
		g.p.RecoverHP(n)

		g.dis.Center(g.dis.Green, "You decided to rest 💤\n")
		fmt.Println("──────────────────────────────────────────────")
		g.dis.Center(g.dis.White, "Effect: Recovered %.1f hp", n)
	} else {
		g.dis.Center(g.dis.Red, "You don't have enough money to rest 💤\n")
		fmt.Println("──────────────────────────────────────────────")
	}

	g.dis.Center(nil, "■ > Press enter to continue")
	g.scanner.Scan()
}

func (g *Game) battle() {
	var turn int
	var isEnemyTurn bool
	var res int
	log := "-- No recent log --"
	display.Clear()

	if g.e == nil {
		g.e = combat.SpawnRandom()
		g.dis.Printf(g.dis.Red, "══════════════════════════════════════════════════")
		g.dis.Center(g.dis.White, "You encountered %s", g.e.Data().Name)
		g.dis.Printf(g.dis.Red, "══════════════════════════════════════════════════\n")
		g.dis.Center(nil, "■ > Press enter to proceed\n")
		g.scanner.Scan()
		display.Clear()
	}

	for {
		display.Clear()
		g.dis.Center(g.dis.White, "Battle 🔥")
		g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
		g.dis.Center(nil, log)
		fmt.Print("\n")

		e := g.e.Data()
		g.dis.Bar(g.p.Hp, g.p.HpCap)
		g.dis.Printf(g.dis.White, "❤️  %.1f ", g.p.Hp)
		fmt.Printf("(You) \n\n")
		g.dis.Bar(e.Hp, e.HpCap)
		g.dis.Printf(g.dis.White, "❤️  %.1f ", e.Hp)
		fmt.Printf("(%s ) \n\n", e.Name)

		if isEnemyTurn {
			res, log = g.enemyTurn(e)
		} else {
			res, log = g.playerTurn(e)
		}

		isEnemyTurn = !isEnemyTurn
		turn++

		if res != NEXT {
			g.battleConclusion(res, turn, e)
			return
		}
	}
}

func (g *Game) playerTurn(e combat.Base) (int, string) {
	fmt.Println("┌─────────────────────────────────────────────")
	fmt.Println("■ > 1. ⚔️   Attack")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 2. 🏃  Flee ")
	fmt.Println("└─────────────────────────────────────────────")

	for {
		g.dis.Printf(g.dis.White, "Input: ")
		g.scanner.Scan()

		switch g.scanner.Text() {
		case "1":
			att := g.p.Attack()
			att = g.e.TakeDamage(att)

			if e.Hp-att <= 0 {
				return WIN, fmt.Sprintf("You delivered the last blow with %.1f damage ⚔️  🩸", att)
			}

			return NEXT, fmt.Sprintf("You attacked, dealing %.1f damage ⚔️", att)
		case "2":
			return FLED, "You decided to fight another day 🏃"
		}

		g.dis.Printf(g.dis.Red, "Invalid Input\n")
	}
}

func (g *Game) enemyTurn(e combat.Base) (int, string) {
	g.dis.Center(nil, "■ > Enemy's turn 🔶. Press enter to proceed")
	g.scanner.Scan()

	att := g.e.Attack()
	att = g.p.TakeDamage(att)

	if g.p.Hp <= 0 {
		return LOSE, fmt.Sprintf("You've been slained with %.1f damage ⚔️  🩸", att)
	}

	return NEXT, fmt.Sprintf("The %s  attacked, dealing %.1f damage", e.Name, att)
}

func (g *Game) battleConclusion(res, turn int, e combat.Base) {
	display.Clear()
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "Battle is Over")
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════\n")

	switch res {
	case WIN:
		g.dis.Printf(g.dis.Green, "You have won the battle 🏆\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", e.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", turn)
		g.dis.Center(nil, "■ > Press enter to go back to menu")

		g.e = nil
		g.scanner.Scan()
	case LOSE:
		g.dis.Printf(g.dis.Red, "You have died ☠️\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", e.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", turn)
		g.dis.Center(nil, "■ > Thanks for playing!")

		g.scanner.Scan()
		os.Exit(1)
	case FLED:
		g.dis.Printf(g.dis.Green, "You have fled from the battle 🏃\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", e.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", turn)
		g.dis.Center(nil, "■ > Press enter to go back to menu")

		g.e = nil
		g.scanner.Scan()
	}
}
