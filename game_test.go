package main

import (
	"io"
	"strings"
	"testing"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/Noblefel/term-rpg/entities"
)

func TestBar(t *testing.T) {
	tests := []struct {
		name  string
		want  string
		input []float32
	}{
		{
			name:  "full green bar",
			want:  "\033[38;2;80;230;95m██████████\033[0m\033[38;2;60;60;60m\033[0m\n",
			input: []float32{10, 10, 10},
		},
		{
			name:  "half the bar should be orange",
			want:  "\033[38;2;245;185;0m█████\033[0m\033[38;2;60;60;60m█████\033[0m\n",
			input: []float32{10, 5, 10},
		},
		{
			name:  "2/10 bar should be red",
			want:  "\033[38;2;247;17;17m██\033[0m\033[38;2;60;60;60m████████\033[0m\n",
			input: []float32{10, 2, 10},
		},
		{
			name:  "if empty or somehow negative, should be full grey",
			want:  "\033[38;2;247;17;17m\033[0m\033[38;2;60;60;60m██████████\033[0m\n",
			input: []float32{10, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sb strings.Builder
			game := Game{writer: &sb}
			game.printBar(int(tt.input[0]), tt.input[1], tt.input[2])

			if tt.want != sb.String() {
				t.Errorf("\nwant: %q\ngot: %q", tt.want, sb.String())
			}
		})
	}
}

func TestSelectionPerks(t *testing.T) {
	var sb strings.Builder
	game := Game{writer: &sb}

	go func() {
		time.Sleep(time.Millisecond)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	game.selectPerks()

	if game.player == nil {
		t.Error("player is nil")
	}
}

func TestRest(t *testing.T) {
	player := entities.Player{}
	player.HPCap = 100
	player.Gold = 5
	game := Game{writer: io.Discard, player: &player}

	go func() {
		time.Sleep(time.Millisecond)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
	}()

	game.menu()

	if player.HP == 0 {
		t.Errorf("player hp is still zero")
	}
}

func TestTrain(t *testing.T) {
	player := entities.Player{}
	player.Gold = 1000
	game := Game{writer: io.Discard, player: &player}

	for player.Gold > 10 {
		go func() {
			time.Sleep(time.Millisecond)
			keyboard.SimulateKeyPress(keys.Down)
			keyboard.SimulateKeyPress(keys.Down)
			keyboard.SimulateKeyPress(keys.Down)
			keyboard.SimulateKeyPress(keys.Enter)
		}()
		game.menu()
	}

	if player.Defense == 0 || player.HPCap == 0 || player.Strength == 0 {
		t.Errorf("player attributes not affected")
	}
}

func TestBattle(t *testing.T) {
	player := entities.Player{}
	player.HP = 100
	player.HPCap = 100
	game := Game{writer: io.Discard, player: &player}

	go func() {
		game.battle()
	}()

	keyboard.SimulateKeyPress(keys.Enter)
	keyboard.SimulateKeyPress(keys.Enter)
	keyboard.SimulateKeyPress(keys.Enter)
	keyboard.SimulateKeyPress(keys.Enter)
	keyboard.SimulateKeyPress(keys.Enter)

	if player.HP == 100 {
		t.Errorf("player hp is not affected")
	}
}
