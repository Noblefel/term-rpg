package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

var rolltest = -1

func clearScreen() {
	fmt.Printf("\033[H")
	fmt.Printf("\033[J")
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
		fmt.Printf(" %0.1fs", ms/1000)
		time.Sleep(90 * time.Millisecond)
		fmt.Printf("\033[5D")
		fmt.Printf("\033[K")
		ms -= 100
	}
}

func scale(base, growth float32) float32 {
	return base + growth*float32(stage)
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
	Skills [5]int  `json:"skills"`
	Hp     float32 `json:"hp"`
	Hpc    float32 `json:"hpc"`
	Str    float32 `json:"str"`
	Def    float32 `json:"def"`
	En     int     `json:"en"`
	Enc    int     `json:"enc"`
}

func save() error {
	data := savedata{
		Stage:  stage,
		Perk:   player.perk,
		Gold:   player.gold,
		Skills: player.skills,
		Hp:     player.hp,
		Hpc:    player.hpcap,
		Str:    player.strength,
		Def:    player.defense,
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
	load.skills = data.Skills
	load.hp = data.Hp
	load.hpcap = data.Hpc
	load.strength = data.Str
	load.defense = data.Def
	load.energy = data.En
	load.energycap = data.Enc
	load.effects = make(map[string]int)

	stage = data.Stage
	player = &load
	return nil
}
