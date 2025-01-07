package main

import (
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/Noblefel/vivi"
)

var (
	out     io.Writer = os.Stdout
	player  *Player
	stage   int
	success = "\033[38;5;83m✔\033[0m "
	fail    = "\033[38;5;196m✘\033[0m "
)

func main() {
	clear()
	fmt.Fprintln(out, " \033[1m[ Welcome ]\033[0m select your player perk")
	fmt.Fprintln(out, "------------------------------------------")

	perk := vivi.Choices(
		"[1] 🛡️  Resiliency : increased survivability",
		"[2] ⚔️  Havoc      : extra damage, but low starting gold & max hp",
		"[3] 🐻 Berserk    : more powerful the lower your hp is",
	)

	player = NewPlayer(perk)
	menuMain()
}

func menuMain() {
	for {
		clear()
		fmt.Fprintln(out, " \033[1m[ Main Menu ]\033[0m")
		fmt.Fprintln(out, " --------")
		fmt.Fprintf(out, " Health : ")
		fmt.Fprintf(out, "%s %.1f\n", barhp(36, player.hp, player.hpcap), player.hp)
		fmt.Fprintf(out, " Energy : ")
		energybar := bars(36, 14, 20)
		fmt.Fprintln(out, "\033[38;5;226m"+energybar[0]+"\033[0m"+energybar[1], 14)
		fmt.Fprintf(out, " Gold   : %d \n", player.gold)
		fmt.Fprintf(out, " Stage  : %d \n", stage+1)
		fmt.Fprintln(out, " --------")

		fmt.Fprintf(out, "\033[s")
		choice := vivi.Choices(
			"[1] 🗺️  Battle",
			"[2] 🏕️  Deep forest",
			"[3] 📋 View Attributes",
			"[4] 🛏️  Rest ($5)",
			"[5] 💪 Train ($10)",
		)

		switch choice {
		case 0:
			enemy := randomEnemy()
			menuBattle(enemy, false)
		case 1:
			exploreDeepForest()
		case 2:
			menuAttributes()
		case 3:
			if player.gold < 5 {
				fmt.Fprint(out, fail)
				fmt.Fprintln(out, "You don't have enough money to rest")
				vivi.Choices("continue")
				continue
			}

			fmt.Fprintf(out, "\033[u\033[0J")
			fmt.Fprintf(out, "  resting... ")
			player.rest()
			vivi.Choices("continue")
		case 4:
			if player.gold < 10 {
				fmt.Fprint(out, fail)
				fmt.Fprintln(out, "You don't have enough money to train")
				vivi.Choices("continue")
				continue
			}

			fmt.Fprintf(out, "\033[u\033[0J")
			fmt.Fprintf(out, "  training... ")
			timer(2500)
			player.train()
			vivi.Choices("continue")
		}
	}
}

func menuAttributes() {
	clear()
	fmt.Fprintln(out, "\033[1m[ Player attributes ]\033[0m")
	fmt.Fprintln(out, "----------")

	hpcap := bars(40, player.hpcap, 300)
	strength := bars(40, player.strength, 80)
	defense := bars(40, player.defense, 30)

	fmt.Fprintf(out, "HP Cap   :")
	fmt.Fprintf(out, "%s %.1f\n", "\033[38;5;83m"+hpcap[0]+"\033[0m"+hpcap[1], player.hpcap)

	fmt.Fprintf(out, "Strength :")
	fmt.Fprintf(out, "%s %.1f\n", "\033[38;5;83m"+strength[0]+"\033[0m"+strength[1], player.strength)

	fmt.Fprintf(out, "Defense  :")
	fmt.Fprintf(out, "%s %.1f\n", "\033[38;5;83m"+defense[0]+"\033[0m"+defense[1], player.defense)

	fmt.Fprintf(out, "\n")
	vivi.Choices("Go back")
}

func menuBattle(enemy entity, exploring bool) {
	for {
		clear()

		if enemy.attr().hp <= 0 && exploring {
			fmt.Fprintf(out, success+"you've won against %s\n", enemy.attr().name)
			vivi.Choices("continue your journey")
			return
		}

		if player.hp <= 0 && exploring {
			fmt.Fprintf(out, fail+"you've lost against %s\n", enemy.attr().name)
			vivi.Choices("continue your journey")
			return
		}

		if enemy.attr().hp <= 0 {
			fmt.Fprintln(out, "You have \033[38;5;83mwon\033[0m the battle")
			gold := scale(10, 2) + rand.Float32()*scale(40, 5)
			player.gold += int(gold)
			fmt.Fprintf(out, success+"got %d gold\n\n", int(gold))
			vivi.Choices("return")
			stage++
			return
		}

		if player.hp <= 0 {
			fmt.Fprintln(out, "You have \033[38;5;196mlost\033[0m the battle")
			vivi.Choices("return")
			return
		}

		fmt.Fprintln(out, "\033[1mYou\033[0m")
		fmt.Fprintf(out, "Health : ")
		fmt.Fprintf(out, "%s %.1f\n", barhp(40, player.hp, player.hpcap), player.hp)
		fmt.Fprintf(out, "Energy : ")
		energybar := bars(40, 14, 20)
		fmt.Fprintln(out, "\033[38;5;226m"+energybar[0]+"\033[0m"+energybar[1], 14)
		fmt.Fprintln(out, "--------")
		fmt.Fprintf(out, "\033[1m%s\033[0m\n", enemy.attr().name)
		fmt.Fprintf(out, "Health : ")
		fmt.Fprintf(out, "%s %.1f\n", barhp(40, enemy.attr().hp, enemy.attr().hpcap), enemy.attr().hp)
		fmt.Fprintln(out, "--------")

		fmt.Fprintf(out, "\033[s")
		choice := vivi.Choices(
			"[1] ⚔️  Attack",
			"[2] 🔥 Skills",
			"[3] 🧰 Items",
			"[4] 🏃 Flee",
			"[5] ⌛ Skip",
		)
		fmt.Fprintf(out, "\033[u\033[0J")

		switch choice {
		case 0:
			player.attack(enemy)
		case 1, 2:
			fmt.Fprintln(out, "\033[38;5;196mNot implemented\033[0m")
			vivi.Choices("continue")
			continue
		case 3:
			fmt.Fprintf(out, "attempting to escape... ")
			timer(1700)

			if player.flee(enemy) {
				vivi.Choices("return")
				return
			}
		case 4:
			fmt.Fprintln(out, "You decided to do nothing")
		}

		fmt.Fprintf(out, "waiting for enemy... ")
		timer(2000)
		fmt.Fprintf(out, "\r\033[K")

		enemy.attack(player)
		vivi.Choices("next turn")
	}
}

func exploreDeepForest() {
	clear()
	fmt.Fprintln(out, "going into the \033[38;5;41mdeep forest\033[0m 🏕️")

	for i := 0; i < 10; i++ {
		fmt.Fprint(out, "exploring...")
		timer(1000 + rand.Float32()*2000)

		fmt.Fprint(out, "\r\033[2K")
		n := rand.IntN(80)

		if n < 14 {
			player.gold += 2
			fmt.Fprint(out, success)
			fmt.Fprintln(out, "You found 2 gold pieces")
		} else if n < 16 {
			player.gold += 5 + rand.IntN(10)
			fmt.Fprint(out, success)
			fmt.Fprintln(out, "You found a \033[38;5;226mpouch\033[0m of gold")
		} else if n < 17 {
			player.gold += 15 + rand.IntN(45)
			fmt.Fprint(out, success)
			fmt.Fprintln(out, "Jackpot! you found a \033[38;5;226mstash\033[0m of gold!")
		} else if n < 24 {
			fmt.Fprint(out, success)
			heal := 1 + rand.Float32()*5
			player.hp = min(player.hp+heal, player.hpcap)
			fmt.Fprintf(out, "You eat some berries, ")
			fmt.Fprintf(out, "recover \033[38;5;83m%.1f\033[0m hp\n", heal)
		} else if n < 28 {
			fmt.Fprint(out, fail)
			fmt.Fprintf(out, "You eat some poisounus berries,")
			player.damage(player.hpcap * 0.07)
		} else if n < 31 {
			fmt.Fprint(out, success)
			heal := 4 + rand.Float32()*5
			player.hp = min(player.hp+heal, player.hpcap)
			fmt.Fprintf(out, "You rest by a campfire, ")
			fmt.Fprintf(out, "recover \033[38;5;83m%.1f\033[0m hp\n", heal)
		} else if n < 35 {
			fmt.Fprint(out, success)
			fmt.Fprintln(out, "You climbed a tree, +\033[38;5;83m0.1\033[0m strength")
			player.strength += 0.1
		} else if n < 38 {
			fmt.Fprint(out, success)
			fmt.Fprintln(out, "You chop some woods, +\033[38;5;83m0.12\033[0m strength")
			player.strength += 0.12
		} else if n < 41 {
			fmt.Fprint(out, fail)
			fmt.Fprint(out, "You fell off a cliff, +\033[38;5;83m0.12\033[0m defense but took")
			player.defense += 0.12
			player.damage(14)
		} else if n < 44 {
			fmt.Fprint(out, success)
			fmt.Fprintln(out, "You endured the long trail, +\033[38;5;83m0.1\033[0m defense")
			player.defense += 0.1
		} else if n < 49 {
			fmt.Fprintln(out, "  You found a potion... \033[s")

			if vivi.Choices("drink it", "pass") == 1 {
				fmt.Fprintln(out, "\033[u\033[0Jyou ignore it")
				continue
			}

			val := player.hpcap * 0.05
			fmt.Fprint(out, "\033[u\033[0J")

			if rand.IntN(10) < 5 {
				fmt.Fprint(out, success)
				player.hp = min(player.hp+val, player.hpcap)
				player.hpcap++
				fmt.Fprintf(out, "It was magical, +\033[38;5;83m1\033[0m hp cap")
				fmt.Fprintf(out, " and recover \033[38;5;83m%.01f\033[0m hp\n", val)
			} else {
				fmt.Fprint(out, fail)
				player.hp = max(player.hp-val, 0)
				player.hpcap -= 2
				fmt.Fprintf(out, "It was cursed, -\033[38;5;198m2\033[0m hp cap and took")
				player.damage(val)
			}
		} else if n < 55 {
			fmt.Fprintln(out, "  You found a hot spring... \033[s")

			if vivi.Choices("swim", "pass") == 1 {
				fmt.Fprintln(out, "\033[u\033[0Jyou ignore it")
				continue
			}

			fmt.Fprint(out, "\033[u\033[0J")
			n := rand.IntN(10)

			if n < 1 {
				fmt.Fprint(out, success)
				fmt.Fprintln(out, "It felt refreshing, +\033[38;5;83m0.2\033[0m to strength and defense")
				player.defense += 0.2
				player.strength += 0.2
			} else if n < 4 {
				fmt.Fprint(out, fail)
				fmt.Fprintf(out, "It boiled you, took")
				player.damage(18)
			} else {
				fmt.Fprint(out, fail)
				fmt.Fprintln(out, "It was okay")
			}
		} else if n < 59 {
			fmt.Fprint(out, fail)
			fmt.Fprintln(out, "You were ambushed by wolves")
			vivi.Choices("oh oh")

			var wolves entity = &attributes{
				name:     "Wolves",
				hp:       20,
				hpcap:    20,
				strength: 6,
				defense:  3,
			}

			menuBattle(wolves, true)
			clear()
		} else if n < 63 {
			fmt.Fprint(out, fail)
			fmt.Fprintln(out, "You crashed into a beehive")
			vivi.Choices("oh oh")

			var bees entity = &attributes{
				name:     "Bee swarm",
				hp:       10,
				hpcap:    10,
				strength: 5,
				defense:  1,
			}

			menuBattle(bees, true)
			clear()
		} else if n < 67 {
			fmt.Fprint(out, fail)
			fmt.Fprintln(out, "You met a wild boar")
			vivi.Choices("oh oh")

			var boar entity = &attributes{
				name:     "Wild boar",
				hp:       30,
				hpcap:    30,
				strength: 6,
				defense:  4,
			}

			menuBattle(boar, true)
			clear()
		} else if n < 77 {
			fmt.Fprint(out, fail)
			dmg := player.hp * 0.15
			player.hp = max(player.hp-dmg, 0)
			fmt.Fprintf(out, "You were affected by some dark magic, ")
			fmt.Fprintf(out, "\033[38;5;198m%.1f\033[0m damage\n", dmg)
		} else {
			fmt.Fprint(out, fail)
			fmt.Fprintln(out, "theres nothing")
		}
	}

	vivi.Choices("You're done here")
}

func clear() {
	fmt.Fprintf(out, "\033[H")
	fmt.Fprintf(out, "\033[J")
}

func barhp(length int, val, cap float32) string {
	percentage := val / cap * 100
	bars := bars(length, val, cap)

	if percentage > 60 {
		bars[0] = "\033[38;5;83m" + bars[0] + "\033[0m"
	} else if percentage > 30 {
		bars[0] = "\033[38;5;226m" + bars[0] + "\033[0m"
	} else {
		bars[0] = "\033[38;5;196m" + bars[0] + "\033[0m"
	}

	return bars[0] + bars[1]
}

func bars(length int, val, cap float32) [2]string {
	val = max(val, 0)
	cap = max(cap, 0)

	if val > cap {
		val = cap
	}

	percentage := val / cap * 100
	colored := int(percentage) * length / 100
	bar1 := strings.Repeat("━", colored)
	bar2 := strings.Repeat("━", max(length-colored, 0))
	bar2 = "\033[38;5;240m" + bar2 + "\033[0m"

	return [2]string{bar1, bar2}
}

func timer(ms float32) {
	for ms > 0 {
		fmt.Fprintf(out, " %0.1fs", ms/1000)
		time.Sleep(90 * time.Millisecond)
		fmt.Fprintf(out, "\033[5D")
		fmt.Fprintf(out, "\033[K")
		ms -= 100
	}
}

func scale(base, growth float32) float32 {
	return base + growth*float32(stage)
}
