package game

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/Noblefel/term-rpg/internal/display"
	"github.com/Noblefel/term-rpg/internal/entity"
)

const (
	NEXT int = iota + 1
	WIN
	LOSE
	FLED
)

type battle struct {
	turn        int
	isEnemyTurn bool
	log         string
	status      int
	enemy       entity.Enemy
	enemyAttr   entity.EnemyBase
	player      *entity.Player
}

func newBattle(p *entity.Player) *battle {
	e := entity.SpawnRandom()
	return &battle{
		log:       "-- No Recent Log --",
		status:    NEXT,
		enemy:     e,
		enemyAttr: e.Attr(),
		player:    p,
	}
}

func (b *battle) Next() bool {
	if b.player.Hp > 0 && b.enemyAttr.Hp <= 0 {
		b.status = WIN
		return false
	} else if b.player.Hp <= 0 {
		b.status = LOSE
		return false
	} else if b.player.HasFled {
		b.status = FLED
		return false
	}

	return true
}

func (b *battle) endTurn() {
	b.enemyAttr = b.enemy.Attr()
	b.isEnemyTurn = !b.isEnemyTurn
	b.turn++
}

func (g *Game) battle() {
	display.Clear()
	b := newBattle(g.p)
	g.dis.Printf(g.dis.Red, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "You encountered %s", b.enemyAttr.Name)
	g.dis.Printf(g.dis.Red, "══════════════════════════════════════════════════\n")
	g.dis.Center(nil, "■ > Press enter to proceed\n")
	g.scanner.Scan()

	for b.Next() {
		display.Clear()
		g.dis.Center(g.dis.White, "Battle 🔥")
		g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
		g.dis.Center(nil, b.log)
		fmt.Print("\n")

		g.dis.Bar(g.p.Hp, g.p.HpCap)
		g.dis.Printf(g.dis.White, "❤️  %.1f ", g.p.Hp)
		fmt.Println("(You)")
		g.dis.Printf(g.dis.White, "📈 effects: ")
		if g.p.GuardTurns > 0 {
			fmt.Print("🛡️")
		}
		if g.p.FuryTurns > 0 {
			fmt.Printf("🔥")
		}
		fmt.Printf("\n\n")

		g.dis.Bar(b.enemyAttr.Hp, b.enemyAttr.HpCap)
		g.dis.Printf(g.dis.White, "❤️  %.1f ", b.enemyAttr.Hp)
		fmt.Printf("(%s) \n", b.enemyAttr.Name)
		g.dis.Printf(g.dis.White, "📈 effects: ")
		if b.enemyAttr.GuardTurns > 0 {
			fmt.Print("🛡️")
		}
		if b.enemyAttr.FuryTurns > 0 {
			fmt.Printf("🔥")
		}
		fmt.Printf("\n\n")

		g.battleTurn(b)
	}

	display.Clear()
	g.dis.Center(nil, b.log)
	g.battleConclusion(b)
}

func (g *Game) battleTurn(b *battle) {
	defer b.endTurn()

	if b.isEnemyTurn {
		g.dis.Center(nil, "■ > Enemy's turn 🔶. Press enter to proceed")
		g.scanner.Scan()

		b.log = entity.EnemyTakeAction(b.enemy, g.p, rand.Intn(100))
	} else {
		fmt.Println("┌─────────────────────────────────────────────")
		fmt.Println("| ■ > 1. ⚔️   Attack	 ■ > 3. 🔥  Fury")
		fmt.Println("|")
		fmt.Println("| ■ > 2. 🛡️  Guard	 ■ > 4. 🏃  Flee")
		fmt.Println("└─────────────────────────────────────────────")

		for {
			g.dis.Printf(g.dis.White, "Input: ")
			g.scanner.Scan()

			n, err := strconv.Atoi(g.scanner.Text())
			if err != nil {
				g.dis.Printf(g.dis.Red, "Input should be a number\n")
				continue
			}

			s, ok := g.p.TakeAction(b.enemy, n)
			if !ok {
				g.dis.Printf(g.dis.Red, s)
				continue
			}

			b.log = s
			break
		}
	}
}

func (g *Game) battleConclusion(b *battle) {
	g.p.FuryTurns = 0
	g.p.GuardTurns = 0
	g.p.HasFled = false

	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════")
	g.dis.Center(g.dis.White, "Battle is Over")
	g.dis.Printf(g.dis.White, "══════════════════════════════════════════════════\n")

	switch b.status {
	case WIN:
		loot := float32(5*g.stage) + rand.Float32()*20
		loot = g.p.AddMoney(loot)

		g.dis.Printf(g.dis.Green, "You have won the battle 🏆\n")
		g.dis.Printf(g.dis.White, "Loot: ")
		fmt.Printf("%.1f 💰\n", loot)
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", b.enemyAttr.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", b.turn)

		g.dis.Center(nil, "■ > Press enter to go back to menu")
		g.scanner.Scan()
		g.stage++
	case LOSE:
		g.dis.Printf(g.dis.Red, "You have died ☠️\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", b.enemyAttr.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", b.turn)
		g.dis.Center(nil, "■ > Thanks for playing!")
		os.Exit(1)
	case FLED:
		g.dis.Printf(g.dis.Green, "You have fled from the battle 🏃\n")
		g.dis.Printf(g.dis.White, "Enemy: ")
		fmt.Printf("%s\n", b.enemyAttr.Name)
		g.dis.Printf(g.dis.White, "Total turns: ")
		fmt.Printf("%d\n", b.turn)
		g.dis.Center(nil, "■ > Press enter to go back to menu")

		g.scanner.Scan()
	}
}
