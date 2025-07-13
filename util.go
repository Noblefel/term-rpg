package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

var rolltest = -1 // for unit tests

func clearScreen() {
	fmt.Printf("\033[H")
	fmt.Printf("\033[J")
}

func bars(length int, val, cap float64) [2]string {
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

func timer(ms float64) {
	for ms > 0 {
		fmt.Printf(" %0.1fs", ms/1000)
		time.Sleep(90 * time.Millisecond)
		fmt.Printf("\033[5D")
		fmt.Printf("\033[K")
		ms -= 100
	}
}

func scale(base, growth float64) float64 {
	return base + growth*float64(stage)
}

func roll() int {
	if rolltest >= 0 {
		return rolltest
	}
	return rand.IntN(100)
}

type savedata struct {
	Stage  int     `json:"stage"`
	Perk   int     `json:"perk"`
	Gold   int     `json:"gold"`
	Weapon int     `json:"weapon"`
	Armor  int     `json:"armor"`
	Skills [5]int  `json:"skills"`
	Hp     float64 `json:"hp"`
	Hpc    float64 `json:"hpc"`
	Str    float64 `json:"str"`
	Def    float64 `json:"def"`
	Agi    float64 `json:"agi"`
	En     int     `json:"en"`
	Enc    int     `json:"enc"`
}

func save() error {
	data := savedata{
		Stage:  stage,
		Perk:   player.perk,
		Gold:   player.gold,
		Weapon: player.weapon,
		Armor:  player.armor,
		Skills: player.skills,
		Hp:     player.hp,
		Hpc:    player.hpcap,
		Str:    player.strength,
		Def:    player.defense,
		Agi:    player.agility,
		En:     player.energy,
		Enc:    player.energycap,
	}

	f, err := os.Create("savegame.json")
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func load() error {
	f, err := os.Open("savegame.json")
	if err != nil {
		return err
	}
	defer f.Close()

	var data savedata

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		return err
	}

	var load Player
	load.name = "player"
	load.perk = data.Perk
	load.gold = data.Gold
	load.weapon = data.Weapon
	load.armor = data.Armor
	load.skills = data.Skills
	load.hp = data.Hp
	load.hpcap = data.Hpc
	load.strength = data.Str
	load.defense = data.Def
	load.agility = data.Agi
	load.energy = data.En
	load.energycap = data.Enc
	load.effects = make(map[string]int)

	stage = data.Stage
	player = &load
	return nil
}
