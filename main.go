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
	out    io.Writer = os.Stdout
	player *Player
	stage  int
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
		fmt.Fprint(out, " Health : ")
		fmt.Fprintf(out, "%s %.1f\n", barhp(36, player.hp, player.hpcap), player.hp)
		fmt.Fprint(out, " Energy : ")
		energybar := bars(36, 14, 20)
		fmt.Fprintln(out, "\033[38;5;226m"+energybar[0]+"\033[0m"+energybar[1], 14)
		fmt.Fprintf(out, " Gold   : %d \n", player.gold)
		fmt.Fprintf(out, " Stage  : %d \n", stage+1)
		fmt.Fprintln(out, " --------")

		choice := vivi.Choices(
			"[1] 🗺️  Next Stage",
			"[2] 📋 View Attributes",
			"[3] 🛏️  Rest ($5)",
			"[4] 💪 Train ($10)",
		)

		switch choice {
		case 0:
			menuBattle(nil)
		case 1:
			menuAttributes()
		case 2:
			if player.gold < 5 {
				fmt.Fprint(out, "\033[38;5;196m✘\033[0m ")
				fmt.Fprintln(out, "You don't have enough money to rest")
				vivi.Choices("continue")
				continue
			}

			fmt.Fprint(out, "  resting... ")
			player.Rest()
			vivi.Choices("continue")
		case 3:
			if player.gold < 10 {
				fmt.Fprint(out, "\033[38;5;196m✘\033[0m ")
				fmt.Fprintln(out, "You don't have enough money to train")
				vivi.Choices("continue")
				continue
			}

			fmt.Fprint(out, "  training... ")
			timer(2500)
			player.Train()
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

	fmt.Fprint(out, "HP Cap   :")
	fmt.Fprintf(out, "%s %.1f\n", "\033[38;5;83m"+hpcap[0]+"\033[0m"+hpcap[1], player.hpcap)

	fmt.Fprint(out, "Strength :")
	fmt.Fprintf(out, "%s %.1f\n", "\033[38;5;83m"+strength[0]+"\033[0m"+strength[1], player.strength)

	fmt.Fprint(out, "Defense  :")
	fmt.Fprintf(out, "%s %.1f\n", "\033[38;5;83m"+defense[0]+"\033[0m"+defense[1], player.defense)

	fmt.Fprint(out, "\n")
	vivi.Choices("Go back")
}

func menuBattle(enemy Enemy) {
	if enemy == nil {
		enemy = SpawnRandom()
	}

	for {
		clear()

		if enemy.Attr().hp <= 0 {
			fmt.Fprintln(out, "You have \033[38;5;83mwon\033[0m the battle")
			gold := scale(10, 5) + rand.Float32()*scale(40, 10)
			player.gold += int(gold)
			fmt.Fprintf(out, "\033[38;5;83m✔\033[0m got %d gold\n\n", int(gold))
			vivi.Choices("return to menu")
			stage++
			return
		}

		if player.hp <= 0 {
			fmt.Fprintln(out, "You have \033[38;5;196mlost\033[0m the battle")
			vivi.Choices("return to menu")
			return
		}

		fmt.Fprintf(out, "                    \033[4;38;5;219mBattle - stage %d\033[0m\n", stage+1)
		fmt.Fprintln(out, "\033[1mYou\033[0m")
		fmt.Fprint(out, "Health : ")
		fmt.Fprintf(out, "%s %.1f\n", barhp(40, player.hp, player.hpcap), player.hp)
		fmt.Fprint(out, "Energy : ")
		energybar := bars(40, 14, 20)
		fmt.Fprintln(out, "\033[38;5;226m"+energybar[0]+"\033[0m"+energybar[1], 14)
		fmt.Fprintln(out, "--------")
		fmt.Fprintf(out, "\033[1m%s\033[0m\n", enemy.Attr().name)
		fmt.Fprint(out, "Health : ")
		fmt.Fprintf(out, "%s %.1f\n", barhp(40, enemy.Attr().hp, enemy.Attr().hpcap), enemy.Attr().hp)
		fmt.Fprintln(out, "--------")

		fmt.Fprint(out, "\033[s")
		choice := vivi.Choices(
			"[1] ⚔️  Attack",
			"[2] 🔥 Skills",
			"[3] 🧰 Items",
			"[4] 🏃 Flee",
			"[5] ⌛ Skip",
		)
		fmt.Fprint(out, "\033[u\033[0J")

		switch choice {
		case 0:
			player.Attack(enemy)
		case 1, 2:
			fmt.Fprintln(out, "\033[38;5;196mNot implemented\033[0m")
			vivi.Choices("continue")
			continue
		case 3:
			fmt.Fprint(out, "attempting to escape... ")
			timer(1700)

			if player.Flee(enemy) {
				vivi.Choices("return to main menu")
				return
			}
		case 4:
			fmt.Fprintln(out, "You decided to do nothing")
		}

		fmt.Fprint(out, "waiting for enemy... ")
		timer(2000)
		fmt.Fprint(out, "\r\033[K")

		enemy.Attack()
		vivi.Choices("next turn")
	}
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
