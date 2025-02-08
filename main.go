package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/Noblefel/vivi"
)

var (
	player  *Player
	stage   int
	success = "\033[38;5;83m✔\033[0m "
	fail    = "\033[38;5;196m✘\033[0m "
)

func main() {
	clearScreen()
	fmt.Println(" \033[1m[ Welcome ]\033[0m select your player perk")
	fmt.Println("------------------------------------------")

	perk := vivi.Choices(
		"[1] 🛡️  Resiliency : increased survivability",
		"[2] ⚔️  Havoc      : +strength damage, but low starting gold & max hp",
		"[3] 🐻 Berserk    : more powerful the lower your hp is",
		"[4] 🐇 Ingenious  : skill cooldown reduced by 1",
		"[5] 🍹 Poisoner   : give poison effect at the start of battle",
	)

	player = NewPlayer(perk)
	menuPoints()
	menuMain()
}

func menuPoints() {
	var (
		points    = 10
		temphpcap = player.hpcap
		tempdef   = player.defense
	)

	for {
		clearScreen()
		fmt.Printf("Modify your starting attributes - ")
		fmt.Printf("\033[38;5;83m%d\033[0m points left\n", points)
		fmt.Printf("HP cap     : %.1f\n", player.hpcap)
		fmt.Printf("Strength   : %.1f\n", player.strength)
		fmt.Printf("Defense    : %.1f\n", player.defense)
		fmt.Printf("Energy cap : %.d\n", player.energycap)
		fmt.Println("------------------------------------------")

		choice := vivi.Choices(
			"increase HP cap by 3",
			"increase strength by 0.25",
			"increase defense by 0.25",
			"increase energy cap by 1 (2 points)",
			"Reset",
			"Done",
		)

		if choice < 3 && points == 0 || choice == 3 && points < 2 {
			fmt.Println("\033[38;5;196mnot enough points\033[0m")
			vivi.Choices("continue")
			continue
		}

		switch choice {
		case 0:
			player.hpcap += 3
			points--
		case 1:
			player.strength += 0.25
			points--
		case 2:
			player.defense += 0.25
			points--
		case 3:
			player.energycap++
			points -= 2
		case 4:
			points = 10
			player.hpcap = temphpcap
			player.defense = tempdef
			player.strength = 20
			player.energycap = 20
		case 5:
			player.hp = player.hpcap
			player.energy = player.energycap
			return
		}
	}
}

func menuMain() {
	for {
		clearScreen()
		fmt.Println(" \033[1m[ Main Menu ]\033[0m")
		fmt.Println(" --------")
		fmt.Printf(" Health : %s %.1f\n", player.attr().hpbar(), player.hp)
		fmt.Printf(" Energy : %s %d\n", player.energybar(), player.energy)
		fmt.Printf(" Gold   : %d \n", player.gold)
		fmt.Printf(" Stage  : %d \n", stage+1)
		fmt.Println(" --------")

		choice := vivi.Choices(
			"[1] 🗺️  Battle",
			"[2] 🏕️  Deep forest",
			"[3] 📋 View attributes",
			"[4] 📋 Equip skills",
			"[5] 🏘️  Visit town",
		)

		switch choice {
		case 0:
			enemy := randomEnemy()

			if player.perk == 4 {
				enemy.attr().effects["poisoned"] = 3
			}

			menuBattle(enemy, false)
		case 1:
			exploreDeepForest()
		case 2:
			menuAttributes()
		case 3:
			menuSkills()
		case 4:
			menuTown()
		}
	}
}

func menuAttributes() {
	clearScreen()
	fmt.Println("\033[1m[ Player attributes ]\033[0m")
	fmt.Println("----------")

	var (
		hpcap     = bars(40, player.hpcap, 300)
		strength  = bars(40, player.strength, 80)
		defense   = bars(40, player.defense, 30)
		energycap = bars(40, float32(player.energycap), 40)
	)

	fmt.Printf("HP cap    :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+hpcap[0]+"\033[0m"+hpcap[1], player.hpcap)

	fmt.Printf("Strength  :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+strength[0]+"\033[0m"+strength[1], player.strength)

	fmt.Printf("Defense   :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+defense[0]+"\033[0m"+defense[1], player.defense)

	fmt.Printf("Energy cap:")
	fmt.Printf("%s %d\n", "\033[38;5;83m"+energycap[0]+"\033[0m"+energycap[1], player.energycap)

	fmt.Printf("\n")
	vivi.Choices("Go back")
}

func menuSkills() {
	var allskills strings.Builder
	tw := tabwriter.NewWriter(&allskills, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(tw, "no\t skill\t description\t energy\t cooldown")
	fmt.Fprintln(tw, "--\t --\t --\t --\t --")

	for i, s := range skills {
		fmt.Fprintf(tw, "%d\t %s\t %s\t %d\t %d\n",
			i+1,
			s.name,
			s.desc,
			s.cost,
			s.cd,
		)
	}
	tw.Flush()

	for {
		clearScreen()
		fmt.Println("\033[1m[ Equipped skills ]\033[0m")
		fmt.Println("-----")
		var choices []string

		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.Debug)
		fmt.Fprintln(tw, "skill\t description\t energy\t cooldown")
		for _, i := range player.skills {
			choices = append(choices, skills[i].name)

			fmt.Fprintf(tw, "%s\t %s\t %d\t %d\n",
				skills[i].name,
				skills[i].desc,
				skills[i].cost,
				skills[i].cd,
			)
		}
		tw.Flush()
		fmt.Println("\nselect the skill you want to change")

		choices = append(choices, "go back")
		i := vivi.Choices(choices...)

		if i == len(choices)-1 {
			return
		}

		clearScreen()
		fmt.Println(allskills.String())

		for {
			fmt.Printf("switch \033[38;5;226m%s\033[0m with... (pick a number) > ", skills[i].name)
			// using fmt.Scan didnt work somehow
			newindex, _ := strconv.Atoi(vivi.Password("#"))

			if newindex > 0 && newindex <= len(skills) {
				player.skills[i] = int(newindex) - 1
				break
			}

			fmt.Println("\033[38;5;196minvalid skill\033[0m")
		}
	}
}

func menuTown() {
	for {
		clearScreen()
		fmt.Println(" \033[1m[ Town square ]\033[0m")
		fmt.Println(" --------")
		fmt.Printf(" Health : %s %.1f\n", player.attr().hpbar(), player.hp)
		fmt.Printf(" Energy : %s %d\n", player.energybar(), player.energy)
		fmt.Printf(" Gold   : %d \n", player.gold)
		fmt.Println(" --------")
		fmt.Printf("\033[s")

		choices := vivi.Choices(
			"[1] 🛏️  Guest house ($5)",
			"[2] 💪 Training grounds ($10)",
			"[3] go back",
		)

		switch choices {
		case 0:
			fmt.Printf("\033[u\033[0J")
			if player.gold < 5 {
				fmt.Print(fail)
				fmt.Println("You don't have enough money to rest")
				vivi.Choices("continue")
				continue
			}

			fmt.Printf("  resting... ")
			player.rest()
			vivi.Choices("continue")
		case 1:
			fmt.Printf("\033[u\033[0J")
			if player.gold < 10 {
				fmt.Print(fail)
				fmt.Println("You don't have enough money to train")
				vivi.Choices("continue")
				continue
			}

			fmt.Printf("  training... ")
			timer(2500)
			player.train()
			vivi.Choices("continue")
		case 2:
			return
		}
	}
}

func menuBattle(enemy entity, exploring bool) {
	defer clear(player.effects)
	defer clear(enemy.attr().effects)

	for {
		clearScreen()
		fmt.Println("\033[1mYou\033[0m")
		fmt.Printf("Health : %s %.1f\n", player.attr().hpbar(), player.hp)
		fmt.Printf("Energy : %s %d\n", player.energybar(), player.energy)
		fmt.Println("--------")
		fmt.Printf("\033[1m%s\033[0m\n", enemy.attr().name)
		fmt.Printf("Health : %s %.1f\n", enemy.attr().hpbar(), enemy.attr().hp)
		fmt.Println("--------")

		if player.effects["poisoned"] > 0 {
			dmg := player.hp * 0.07
			fmt.Print("  You took ", enemy.attr().name)
			fmt.Printf("\033[38;5;198m%.1f\033[0m poison damage\n", dmg)
			player.hp = max(player.hp-dmg, 0)
		}

		if player.effects["stunned"] > 0 {
			fmt.Println(fail + "You are stunned")
		} else {
			menuPlayerActions(enemy)
		}

		if player.effects["fled"] > 0 {
			vivi.Choices("return")
			return
		}

		if enemy.attr().effects["poisoned"] > 0 {
			dmg := enemy.attr().hp * 0.07
			fmt.Printf("  %s took ", enemy.attr().name)
			fmt.Printf("\033[38;5;198m%.1f\033[0m poison damage\n", dmg)
			enemy.setHP(enemy.attr().hp - dmg)
		}

		if enemy.attr().hp <= 0 && exploring {
			fmt.Printf(success+"you've won against %s\n", enemy.attr().name)
			vivi.Choices("continue your journey")
			return
		}

		if enemy.attr().hp <= 0 {
			fmt.Println("You have \033[38;5;83mwon\033[0m the battle")
			gold := scale(10, 2) + rand.Float32()*scale(20, 5)
			player.gold += int(gold)
			fmt.Printf("got %.0f gold\n", gold)
			vivi.Choices("return")
			stage++
			return
		}

		if enemy.attr().effects["stunned"] > 0 {
			fmt.Printf(fail+"%s is stunned\n", enemy.attr().name)
		} else {
			fmt.Printf("waiting for enemy... ")
			timer(2000)
			fmt.Printf("\r\033[K")
			enemy.attack(player)
		}

		if player.hp <= 0 && exploring {
			fmt.Printf(fail+"you've lost against %s\n", enemy.attr().name)
			vivi.Choices("continue your journey")
			return
		}

		if player.hp <= 0 {
			fmt.Println("You have \033[38;5;196mlost\033[0m the battle")
			vivi.Choices("return")
			return
		}

		enemy.attr().decrementEffect()
		player.decrementEffect()
		player.energy = min(player.energy+1, player.energycap)
		vivi.Choices("next turn")
	}
}

func menuPlayerActions(enemy entity) {
	fmt.Printf("\033[s")

	for {
		choice := vivi.Choices(
			"[1] ⚔️  Attack",
			"[2] 🔥 Skills",
			"[3] 🧰 Items",
			"[4] 🏃 Flee",
			"[5] ⌛ Skip",
		)
		fmt.Printf("\033[u\033[0J")

		switch choice {
		case 0:
			player.attack(enemy)
			return
		case 1:
			var choices []string

			for _, i := range player.skills {
				choice := fmt.Sprintf(
					"%s\033[0m: %d energy",
					skills[i].name,
					skills[i].cost,
				)

				if skills[i].cost > player.energy || player.effects["cd_"+skills[i].name] > 0 {
					choice = "\033[38;5;196m" + choice
				} else {
					choice = "\033[38;5;226m" + choice
				}

				choices = append(choices, choice)
			}

			choices = append(choices, "cancel")
			i := vivi.Choices(choices...)
			fmt.Printf("\033[u\033[0J")

			if i < len(choices)-1 && player.skill(player.skills[i], enemy) {
				return
			}
		case 2:
			fmt.Println("\033[38;5;196mNot implemented\033[0m")
		case 3:
			fmt.Printf("attempting to escape... ")
			timer(1700)
			player.flee(enemy)
			return
		case 4:
			fmt.Println("  You decided to do nothing")
			return
		}
	}
}

func exploreDeepForest() {
	clearScreen()
	fmt.Println("going into the \033[38;5;41mdeep forest\033[0m 🏕️")

	for i := 0; i < 10; i++ {
		fmt.Print("exploring...")
		timer(1000 + rand.Float32()*2000)

		fmt.Print("\r\033[K")
		n := rand.IntN(80)

		if n < 14 {
			player.gold += 2
			fmt.Print(success)
			fmt.Println("You found 2 gold pieces")
		} else if n < 16 {
			player.gold += 5 + rand.IntN(10)
			fmt.Print(success)
			fmt.Println("You found a \033[38;5;226mpouch\033[0m of gold")
		} else if n < 17 {
			player.gold += 15 + rand.IntN(45)
			fmt.Print(success)
			fmt.Println("Jackpot! you found a \033[38;5;226mstash\033[0m of gold!")
		} else if n < 24 {
			fmt.Print(success)
			heal := 1 + rand.Float32()*5
			player.hp = min(player.hp+heal, player.hpcap)
			fmt.Printf("You eat some berries, ")
			fmt.Printf("recover \033[38;5;83m%.1f\033[0m hp\n", heal)
		} else if n < 28 {
			fmt.Print(fail)
			fmt.Printf("You eat some poisounus berries,")
			player.damage(player.hpcap * 0.07)
		} else if n < 31 {
			fmt.Print(success)
			heal := 4 + rand.Float32()*5
			player.hp = min(player.hp+heal, player.hpcap)
			fmt.Printf("You rest by a campfire, ")
			fmt.Printf("recover \033[38;5;83m%.1f\033[0m hp\n", heal)
		} else if n < 35 {
			fmt.Print(success)
			fmt.Println("You climbed a tree, +\033[38;5;83m0.1\033[0m strength")
			player.strength += 0.1
		} else if n < 38 {
			fmt.Print(success)
			fmt.Println("You chop some woods, +\033[38;5;83m0.12\033[0m strength")
			player.strength += 0.12
		} else if n < 41 {
			fmt.Print(fail)
			fmt.Print("You fell off a cliff, +\033[38;5;83m0.12\033[0m defense but took")
			player.defense += 0.12
			player.damage(14)
		} else if n < 44 {
			fmt.Print(success)
			fmt.Println("You endured the long trail, +\033[38;5;83m0.1\033[0m defense")
			player.defense += 0.1
		} else if n < 49 {
			fmt.Println("  You found a potion... \033[s")

			if vivi.Choices("drink it", "pass") == 1 {
				fmt.Println("\033[u\033[0Jyou ignore it")
				continue
			}

			val := player.hpcap * 0.05
			fmt.Print("\033[u\033[0J")

			if rand.IntN(10) < 5 {
				fmt.Print(success)
				player.hp = min(player.hp+val, player.hpcap)
				player.hpcap++
				fmt.Printf("It was magical, +\033[38;5;83m1\033[0m hp cap")
				fmt.Printf(" and recover \033[38;5;83m%.01f\033[0m hp\n", val)
			} else {
				fmt.Print(fail)
				player.hp = max(player.hp-val, 0)
				player.hpcap -= 2
				fmt.Printf("It was cursed, -\033[38;5;198m2\033[0m hp cap and took")
				player.damage(val)
			}
		} else if n < 55 {
			fmt.Println("  You found a hot spring... \033[s")

			if vivi.Choices("swim", "pass") == 1 {
				fmt.Println("\033[u\033[0Jyou ignore it")
				continue
			}

			fmt.Print("\033[u\033[0J")
			n := rand.IntN(10)

			if n < 1 {
				fmt.Print(success)
				fmt.Println("It felt refreshing, +\033[38;5;83m0.2\033[0m to strength and defense")
				player.defense += 0.2
				player.strength += 0.2
			} else if n < 4 {
				fmt.Print(fail)
				fmt.Printf("It boiled you, took")
				player.damage(18)
			} else {
				fmt.Print(fail)
				fmt.Println("It was okay")
			}
		} else if n < 59 {
			fmt.Print(fail)
			fmt.Println("You were ambushed by wolves")
			vivi.Choices("oh oh")

			var wolves entity = &attributes{
				name:     "Wolves",
				hp:       20,
				hpcap:    20,
				strength: 6,
				defense:  3,
				effects:  make(map[string]int),
			}

			menuBattle(wolves, true)
			clearScreen()
		} else if n < 63 {
			fmt.Print(fail)
			fmt.Println("You crashed into a beehive")
			vivi.Choices("oh oh")

			var bees entity = &attributes{
				name:     "Bee swarm",
				hp:       10,
				hpcap:    10,
				strength: 5,
				defense:  1,
				effects:  make(map[string]int),
			}

			menuBattle(bees, true)
			clearScreen()
		} else if n < 67 {
			fmt.Print(fail)
			fmt.Println("You met a wild boar")
			vivi.Choices("oh oh")

			var boar entity = &attributes{
				name:     "Wild boar",
				hp:       30,
				hpcap:    30,
				strength: 6,
				defense:  4,
				effects:  make(map[string]int),
			}

			menuBattle(boar, true)
			clearScreen()
		} else if n < 70 {
			fmt.Print(fail)
			dmg := player.hp * 0.15
			player.hp = max(player.hp-dmg, 0)
			fmt.Printf("You were affected by some dark magic, ")
			fmt.Printf("\033[38;5;198m%.1f\033[0m damage\n", dmg)
		} else {
			fmt.Print(fail)
			fmt.Println("theres nothing")
		}
	}

	vivi.Choices("You're done here")
}
