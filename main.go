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
	fmt.Println("------------------------------------------")
	fmt.Println("load save file?")

	if vivi.Choices("yes", "no") == 0 {
		if err := load(); err != nil {
			fmt.Println(fail, "error loading savefile")
			fmt.Println(fail, err)
			vivi.Choices("start a new game")
			// player = nil
		}
	}

	if player == nil {
		perk := selectPerks()
		player = NewPlayer()
		player.setPerk(perk)
		menuPoints()
	}

	menuMain()
}

func selectPerks() int {
	clearScreen()
	fmt.Println("please select your player perk")
	fmt.Println("------------------------------------------")

	return vivi.Choices(
		"[1] 🛡️  Resilient  : increase overall defense",
		"[2] ⚔️  Havoc      : +20% damage, but low starting gold & energy cap",
		"[3] 🐻 Berserk    : more powerful the lower your hp is",
		"[4] 🐇 Ingenious  : +2 energy cap, skill cooldown reduced by 2",
		"[5] 🍹 Poisoner   : inflict poisoning effect at the start of battle",
		"[6] ⚰️  Deadman    : inflict weaken effect at the start of battle",
		"[7] 🏃 Survivor   : almost always succeed when fleeing",
		"[8] 🎃 Insanity   : it can either go really well or really bad",
	)
}

func menuPoints() {
	var (
		points    = 10
		temphpcap = player.hpcap
		tempstr   = player.strength
		tempdef   = player.defense
		tempagi   = player.agility
		tempencap = player.energycap
	)

	for {
		clearScreen()
		fmt.Printf("Modify your starting attributes - ")
		fmt.Printf("\033[38;5;83m%d\033[0m points left\n", points)
		fmt.Printf("HP cap     : %.1f\n", player.hpcap)
		fmt.Printf("Strength   : %.1f\n", player.strength)
		fmt.Printf("Defense    : %.1f\n", player.defense)
		fmt.Printf("Agility    : %.1f\n", player.agility)
		fmt.Printf("Energy cap : %.d\n", player.energycap)
		fmt.Println("------------------------------------------")
		fmt.Printf("\033[s")

		choice := vivi.Choices(
			"increase HP cap by 7.5",
			"increase strength by 1",
			"increase defense by 1",
			"increase agility by 1",
			"increase energy cap by 1 (3 points)",
			"Reset",
			"Done",
		)

		if choice < 4 && points == 0 || choice == 4 && points < 3 {
			fmt.Print("\033[u\033[0J")
			fmt.Println("\033[38;5;196mnot enough points\033[0m")
			vivi.Choices("continue")
			continue
		}

		switch choice {
		case 0:
			player.hpcap += 7.5
			points--
		case 1:
			player.strength++
			points--
		case 2:
			player.defense++
			points--
		case 3:
			player.agility++
			points--
		case 4:
			player.energycap++
			points -= 3
		case 5:
			points = 10
			player.hpcap = temphpcap
			player.strength = tempstr
			player.defense = tempdef
			player.agility = tempagi
			player.energycap = tempencap
		case 6:
			player.hp = player.hpcap
			player.energy = player.energycap
			return
		}
	}
}

func menuMain() {
	for {
		clearScreen()
		fmt.Printf(" Health : %s %.1f\n", player.attr().hpbar(), player.hp)
		fmt.Printf(" Energy : %s %d\n", player.energybar(), player.energy)
		fmt.Printf(" Perk   : %s \n", player.getperk())
		fmt.Printf(" Gold   : %d \n", player.gold)
		fmt.Printf(" Stage  : %d \n", stage+1)
		fmt.Println(" --------")

		choice := vivi.Choices(
			"[1] 🗺️  battle",
			"[2] 🏕️  deep forest",
			"[3] 📋 view attributes",
			"[4] 📋 equip skills",
			"[5] 🏘️  visit town",
			"[6] save game",
		)

		switch choice {
		case 0:
			menuBattle(spawn(), false)
		case 1:
			exploreDeepForest()
		case 2:
			menuAttributes()
		case 3:
			menuSkills()
		case 4:
			menuTown()
		case 5:
			if err := save(); err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println("progress saved")
			}

			vivi.Choices("ok")
		}
	}
}

func menuAttributes() {
	clearScreen()
	fmt.Println("\033[1mPlayer attributes\033[0m")
	fmt.Println("----------")

	var (
		hpcap     = bars(40, player.hpcap, 750)
		strength  = bars(40, player.strength, 400)
		defense   = bars(40, player.defense, 150)
		agility   = bars(40, player.agility, 100)
		energycap = bars(40, float32(player.energycap), 40)
	)

	fmt.Printf("HP cap     :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+hpcap[0]+"\033[0m"+hpcap[1], player.hpcap)

	fmt.Printf("Strength   :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+strength[0]+"\033[0m"+strength[1], player.strength)

	fmt.Printf("Defense    :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+defense[0]+"\033[0m"+defense[1], player.defense)

	fmt.Printf("Agility    :")
	fmt.Printf("%s %.1f\n", "\033[38;5;83m"+agility[0]+"\033[0m"+agility[1], player.agility)

	fmt.Printf("Energy cap :")
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
		fmt.Println(" \033[1mTown square\033[0m")
		fmt.Println(" --------")
		fmt.Printf(" Health : %s %.1f\n", player.attr().hpbar(), player.hp)
		fmt.Printf(" Energy : %s %d\n", player.energybar(), player.energy)
		fmt.Printf(" Perk   : %s \n", player.getperk())
		fmt.Printf(" Gold   : %d \n", player.gold)
		fmt.Println(" -------- ")
		fmt.Print("\033[s")

		choices := vivi.Choices(
			"[1] guest house ($5)",
			"[2] training grounds ($10)",
			"[3] switch perk ($34)",
			"[4] go back",
		)

		fmt.Print("\033[u\033[0J")
		switch choices {
		case 0:
			if player.gold < 5 {
				fmt.Println(fail + "not enough gold to rest")
				vivi.Choices("continue")
				continue
			}

			fmt.Print("  resting... ")
			player.gold -= 5
			player.rest()
			vivi.Choices("continue")
		case 1:
			if player.gold < 10 {
				fmt.Println(fail + "not enough gold to train")
				vivi.Choices("continue")
				continue
			}

			fmt.Printf("  training... ")
			timer(2500)
			player.gold -= 10
			player.train()
			vivi.Choices("continue")
		case 2:
			if player.gold < 34 {
				fmt.Println(fail + "not enough gold to switch perk")
				vivi.Choices("continue")
				continue
			}

			newperk := selectPerks()

			if newperk == player.perk {
				fmt.Println(fail + "you already have that perk")
				vivi.Choices("go back")
				continue
			}

			player.gold -= 34
			player.setPerk(newperk)
			fmt.Println(success + "perk changed")
			vivi.Choices("continue")
		case 3:
			return
		}
	}
}

func menuBattle(enemy entity, exploring bool) {
	defer clear(player.effects)
	defer clear(enemy.attr().effects)

	if player.perk == 4 {
		enemy.attr().effects["poisoned"] = 6
	}

	if player.perk == 5 {
		enemy.attr().effects["weakened"] = 3
	}

	if _, ok := enemy.(*undead); ok {
		player.effects["weakened"] = 3
	}

	for {
		clearScreen()
		fmt.Println("\033[1mYou\033[0m")
		fmt.Printf("health : %s %.1f\n", player.attr().hpbar(), player.hp)
		fmt.Printf("energy : %s %d\n", player.energybar(), player.energy)
		fmt.Println("--------")
		fmt.Printf("\033[1m%s\033[0m\n", enemy.attr().name)
		fmt.Printf("health : %s %.1f\n", enemy.attr().hpbar(), enemy.attr().hp)
		fmt.Println("--------")

		if player.perk == 7 {
			roll := roll()

			if roll < 1 {
				fmt.Print("  \033[38;5;226minsanity\033[0m: you somehow annihilate the enemy!")
				enemy.setHP(0)
			} else if roll < 10 {
				player.effects["stunned"] = 1
				fmt.Println("  \033[38;5;226minsanity\033[0m: your mind is in disarray")
			} else if roll < 17 {
				n := -0.5 + rand.Float32()*1
				player.hpcap = max(50, player.hpcap+n)
				fmt.Printf("  \033[38;5;226minsanity\033[0m: you get %.1f hp cap\n", n)
			} else if roll < 24 {
				n := -0.25 + rand.Float32()*0.5
				player.strength = max(5, player.strength+n)
				fmt.Printf("  \033[38;5;226minsanity\033[0m: you get %.1f strength\n", n)
			} else if roll < 31 {
				n := -0.25 + rand.Float32()*0.5
				player.defense = max(1, player.defense+n)
				fmt.Printf("  \033[38;5;226minsanity\033[0m: you get %.1f defense\n", n)
			}
		}

		player.applyEffects()

		if player.effects["stunned"] > 0 {
			fmt.Println(fail + "you are stunned")
		} else {
			menuPlayerActions(enemy)
		}

		if player.effects["fled"] > 0 {
			vivi.Choices("return")
			return
		}

		enemy.applyEffects()

		if enemy.attr().hp <= 0 && exploring {
			fmt.Printf(success+"you've won against %s\n", enemy.attr().name)
			vivi.Choices("continue your journey")
			return
		}

		if enemy.attr().hp <= 0 {
			fmt.Println("you have \033[38;5;83mwon\033[0m the battle")
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
			fmt.Println("you have \033[38;5;196mlost\033[0m the battle")
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
			"[1] attack",
			"[2] skills",
			"[3] items",
			"[4] flee",
			"[5] skip",
		)
		fmt.Printf("\033[u\033[0J")

		switch choice {
		case 0:
			player.attack(enemy)
			return
		case 1:
			if _, ok := enemy.(*celestial); ok {
				fmt.Println("  \033[38;5;196mcannot use skill in the presence of the divine\033[0m")
				continue
			}

			var choices []string

			for _, i := range player.skills {
				choice := fmt.Sprintf(
					"%s\033[0m: %d energy",
					skills[i].name,
					skills[i].cost,
				)

				cost := skills[i].cost

				if player.effects["confused"] > 0 {
					cost++
				}

				if cost > player.energy || player.effects["cd"+skills[i].name] > 0 {
					choice = "\033[38;5;196m" + choice
				} else {
					choice = "\033[38;5;226m" + choice
				}

				choices = append(choices, choice)
			}

			if player.effects["confused"] > 0 {
				fmt.Println("warning, you are \033[38;5;226mconfused\033[0m! energy cost increased by 1")
			}

			if player.perk == 2 && player.hp/player.hpcap <= 0.2 {
				fmt.Println("berserk perk bonus! cooldown decreased by 1")
			}

			if player.perk == 7 {
				fmt.Println("insanity perk! cooldown will be randomized")
			}

			choices = append(choices, "cancel")
			i := vivi.Choices(choices...)
			fmt.Printf("\033[u\033[0J")

			if i < len(choices)-1 && player.skill(player.skills[i], enemy) {
				return
			}
		case 2:
			fmt.Println("\033[38;5;196mnot implemented\033[0m")
		case 3:
			fmt.Printf("attempting to escape... ")
			timer(1700)
			player.flee(enemy)
			return
		case 4:
			fmt.Println("  you decided to do nothing")
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
			heal := 5 + rand.Float32()*5
			player.hp = min(player.hp+heal, player.hpcap)
			fmt.Printf("You eat some berries, ")
			fmt.Printf("recover \033[38;5;83m%.1f\033[0m hp\n", heal)
		} else if n < 28 {
			fmt.Print(fail)
			fmt.Printf("You eat some poisounus berries,")
			player.damage(player.hpcap * 0.04)
		} else if n < 31 {
			fmt.Print(success)
			heal := 10 + rand.Float32()*10
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
			player.damage(30)
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
				player.hpcap += 4
				fmt.Printf("It was magical, +\033[38;5;83m4\033[0m hp cap")
				fmt.Printf(" and recover \033[38;5;83m%.01f\033[0m hp\n", val)
			} else {
				fmt.Print(fail)
				player.hp = max(player.hp-val, 0)
				player.hpcap -= 4
				fmt.Printf("It was cursed, -\033[38;5;198m4\033[0m hp cap and took")
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
				player.damage(50)
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
				hp:       60,
				hpcap:    60,
				strength: 50,
				defense:  8,
				agility:  8,
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
				hp:       37,
				hpcap:    37,
				strength: 37,
				defense:  1,
				agility:  15,
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
				hp:       70,
				hpcap:    70,
				strength: 40,
				defense:  20,
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
