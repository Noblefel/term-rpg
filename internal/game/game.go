package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/Noblefel/term-rpg/internal/battle"
	"github.com/Noblefel/term-rpg/internal/display"
	"github.com/Noblefel/term-rpg/internal/entity"
)

type Game struct {
	scanner *bufio.Scanner
	dis     *display.Display
	p       *entity.Player
}

func New(scanner *bufio.Scanner, dis *display.Display) *Game {
	return &Game{
		scanner: scanner,
		dis:     dis,
	}
}

func (g *Game) Start() {
	if g.p == nil {
		perk := g.selectPerks()
		g.p = entity.NewPlayer(perk)
	}

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
	fmt.Printf("%s\n\n", perk)

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
	g.dis.Center(g.dis.Green, "Select your perk 🔥")
	g.dis.Printf(g.dis.Green, "══════════════════════════════════════════════════\n")

	fmt.Println("┌─────────────────────────────────────────────")
	fmt.Println("■ > 1. 💰 Greed (Gain 15% more loot)")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 2. 🛡️  Resiliency (+1 Def point and 10% dmg reduction)")
	fmt.Println("|─────────────────────────────────────────────")
	fmt.Println("■ > 3. ⚔️   Havoc (+25% Attack, but -15 HP cap)")
	fmt.Println("└─────────────────────────────────────────────")

	for {
		g.dis.Printf(g.dis.White, "Input: ")
		g.scanner.Scan()

		t := g.scanner.Text()
		switch t {
		case "1", "2", "3":
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

	if g.p.Money > 5 {
		g.p.Money -= 5
		n := 5 + (g.p.HpCap * 0.1) + rand.Float32()*8
		g.p.RecoverHP(n)

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
	if rand.Intn(100) < 20 {
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

func (g *Game) battle() {
	display.Clear()
	b := battle.New(entity.SpawnRandom())
	g.dis.Printf(g.dis.Red, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "You encountered %s", b.EnemyAttr.Name)
	g.dis.Printf(g.dis.Red, "══════════════════════════════════════════════════\n")
	g.dis.Center(nil, "■ > Press enter to proceed\n")
	g.scanner.Scan()

	for {
		display.Clear()
		g.dis.Center(g.dis.White, "Battle 🔥")
		g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
		g.dis.Center(nil, b.Log)
		fmt.Print("\n")

		g.dis.Bar(g.p.Hp, g.p.HpCap)
		g.dis.Printf(g.dis.White, "❤️  %.1f ", g.p.Hp)
		fmt.Printf("(You) \n\n")
		g.dis.Bar(b.EnemyAttr.Hp, b.EnemyAttr.HpCap)
		g.dis.Printf(g.dis.White, "❤️  %.1f ", b.EnemyAttr.Hp)
		fmt.Printf("(%s ) \n\n", b.EnemyAttr.Name)

		if b.IsEnemyTurn {
			g.enemyTurn(b)
		} else {
			g.playerTurn(b)
		}

		b.IsEnemyTurn = !b.IsEnemyTurn
		b.Turn++

		if b.Status != battle.NEXT {
			display.Clear()
			g.dis.Center(nil, b.Log)
			g.battleConclusion(b)
			return
		}
	}
}

func (g *Game) playerTurn(b *battle.Battle) {
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
			att = b.Enemy.TakeDamage(att)

			if b.EnemyAttr.Hp <= 0 {
				b.Status = battle.WIN
				b.Log = fmt.Sprintf("You've slained them with %.1f damage ⚔️  🩸", att)
			} else {
				b.Status = battle.NEXT
				b.Log = fmt.Sprintf("You attacked, dealing %.1f damage ⚔️", att)
			}
		case "2":
			b.Status = battle.FLED
			b.Log = "You decided to fight another day 🏃"
		default:
			g.dis.Printf(g.dis.Red, "Invalid Input\n")
			continue
		}
		return
	}
}

func (g *Game) enemyTurn(b *battle.Battle) {
	g.dis.Center(nil, "■ > Enemy's turn 🔶. Press enter to proceed")
	g.scanner.Scan()

	att := b.Enemy.Attack()
	att = g.p.TakeDamage(att)

	if g.p.Hp <= 0 {
		b.Status = battle.LOSE
		b.Log = fmt.Sprintf("You've been slained with %.1f damage ⚔️  🩸", att)
	} else {
		b.Status = battle.NEXT
		b.Log = fmt.Sprintf("%s  attacked, dealing %.1f damage", b.EnemyAttr.Name, att)
	}
}

func (g *Game) battleConclusion(b *battle.Battle) {
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "Battle is Over")
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════\n")

	switch b.Status {
	case battle.WIN:
		loot := b.Enemy.DropLoot()
		loot = g.p.AddMoney(loot)

		g.dis.Printf(g.dis.Green, "You have won the battle 🏆\n")
		g.dis.Printf(g.dis.Green, "Loot: ")
		fmt.Printf("%.1f 💰\n", loot)
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", b.EnemyAttr.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", b.Turn)

		g.dis.Center(nil, "■ > Press enter to go back to menu")
		g.scanner.Scan()
	case battle.LOSE:
		g.dis.Printf(g.dis.Red, "You have died ☠️\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", b.EnemyAttr.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", b.Turn)
		g.dis.Center(nil, "■ > Thanks for playing!")
		os.Exit(1)
	case battle.FLED:
		g.dis.Printf(g.dis.Green, "You have fled from the battle 🏃\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", b.EnemyAttr.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", b.Turn)
		g.dis.Center(nil, "■ > Press enter to go back to menu")

		g.scanner.Scan()
	}
}
