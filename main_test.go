package main

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/Noblefel/vivi"
)

var replacer = strings.NewReplacer(
	"\033[38;5;83m", "",
	"\033[38;5;196m", "",
	"\033[38;5;198m", "",
	"\033[38;5;226m", "",
	"\033[38;5;240m", "",
	"\033[s", "",
	"\033[u", "",
	"\033[K", "",
	"\033[5D", "",
	"\033[H", "",
	"\033[J", "",
	"\033[0J", "",
	"\033[1m", "",
	"\033[0m", "",
	"\r", "",
)

func dump(t *testing.T, s string) {
	t.Helper()
	tname := strings.Replace(t.Name(), "/", "_", -1)
	os.WriteFile("./testdump/"+tname, []byte(s), os.ModePerm)
}

func TestMain(m *testing.M) {
	vivi.Out = io.Discard
	os.Exit(m.Run())
}

func TestMenuPoints(t *testing.T) {
	var sb strings.Builder
	out = &sb
	player = new(Player)
	player.hpcap = 100
	player.strength = 9
	player.defense = 5
	player.energycap = 10

	go menuPoints()
	time.Sleep(50 * time.Millisecond)

	t.Run("display correctly", func(t *testing.T) {
		defer sb.Reset()
		got := replacer.Replace(sb.String())
		dump(t, got)

		if !strings.Contains(got, "10 points left") {
			t.Error("incorrect points left")
		}

		if !strings.Contains(got, "HP cap     : 100.0") {
			t.Error("incorrect hp cap")
		}

		if !strings.Contains(got, "Strength   : 9.0") {
			t.Error("incorrect strength")
		}

		if !strings.Contains(got, "Defense    : 5.0") {
			t.Error("incorrect defense")
		}

		if !strings.Contains(got, "Energy cap : 10") {
			t.Error("incorrect energy cap")
		}
	})

	t.Run("choosing increase hp", func(t *testing.T) {
		defer sb.Reset()

		player.gold = 0
		keyboard.SimulateKeyPress(keys.Enter)
		time.Sleep(40 * time.Millisecond)

		got := replacer.Replace(sb.String())
		dump(t, got)

		if !strings.Contains(got, "9 points left") {
			t.Error("incorrect points left")
		}

		if !strings.Contains(got, "HP cap     : 103.0") {
			t.Error("hp cap should increase by 3")
		}
	})

	keyboard.SimulateKeyPress(keys.Up)
	keyboard.SimulateKeyPress(keys.Enter) //exit
}

func TestMenuMain(t *testing.T) {
	var sb strings.Builder
	out = &sb

	player = new(Player)
	player.hp = 105
	player.energy = 13
	player.gold = 50

	go menuMain()
	time.Sleep(50 * time.Millisecond)

	t.Run("display info correctly", func(t *testing.T) {
		defer sb.Reset()
		got := replacer.Replace(sb.String())
		dump(t, got)

		want := "Health : ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 105.0"
		if !strings.Contains(got, want) {
			t.Error("incorrect hp bar")
		}

		want = "Energy : ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 13"
		if !strings.Contains(got, want) {
			t.Error("incorrect energy bar")
		}

		if !strings.Contains(got, "Gold   : 50") {
			t.Error("incorrect gold")
		}

		if !strings.Contains(got, "Stage  : 1") {
			t.Error("incorrect stage")
		}
	})

	t.Run("choosing rest with 0 gold", func(t *testing.T) {
		defer sb.Reset()

		player.gold = 0
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Down)
		keyboard.SimulateKeyPress(keys.Enter)
		time.Sleep(40 * time.Millisecond)

		got := replacer.Replace(sb.String())
		keyboard.SimulateKeyPress(keys.Enter) //go back to menu
		dump(t, got)

		if !strings.Contains(got, "You don't have enough money to rest") {
			t.Error("should fail when attempting to rest")
		}
	})

	t.Run("choosing train with 0 gold", func(t *testing.T) {
		keyboard.SimulateKeyPress(keys.Up)
		keyboard.SimulateKeyPress(keys.Up)
		keyboard.SimulateKeyPress(keys.Enter)
		time.Sleep(40 * time.Millisecond)

		got := replacer.Replace(sb.String())
		keyboard.SimulateKeyPress(keys.Enter) //go back to menu
		dump(t, got)

		if !strings.Contains(got, "You don't have enough money to train") {
			t.Error("should fail when attempting to train")
		}
	})

	keyboard.SimulateKeyPress(keys.Up)
	keyboard.SimulateKeyPress(keys.Enter) //exit
}

func TestMenuAttributes(t *testing.T) {
	var sb strings.Builder
	out = &sb

	player = new(Player)
	player.hpcap = 120
	player.energycap = 5
	player.strength = 14
	player.defense = 7.5

	go menuAttributes()
	time.Sleep(50 * time.Millisecond)

	t.Run("display info correctly", func(t *testing.T) {
		defer sb.Reset()
		got := replacer.Replace(sb.String())
		dump(t, got)

		want := "HP cap    :━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 120.0"
		if !strings.Contains(got, want) {
			t.Error("incorrect hp cap")
		}

		want = "Strength  :━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 14.0"
		if !strings.Contains(got, want) {
			t.Error("incorrect strength")
		}

		want = "Defense   :━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 7.5"
		if !strings.Contains(got, want) {
			t.Error("incorrect defense")
		}

		want = "Energy cap:━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 5"
		if !strings.Contains(got, want) {
			t.Error("incorrect energy cap")
		}
	})
}

func TestBars(t *testing.T) {
	t.Run("10 length with max value", func(t *testing.T) {
		bars := bars(10, 10, 10)

		got := strings.Count(bars[0], "━")
		if got != 10 {
			t.Errorf("1st bar should be 10, got: %d", got)
		}

		got = strings.Count(bars[1], "━")
		if got != 0 {
			t.Errorf("2nd bar should be 0, got: %d", got)
		}
	})

	t.Run("10 length with empty value", func(t *testing.T) {
		bars := bars(10, 0, 10)

		got := strings.Count(bars[0], "━")
		if got != 0 {
			t.Errorf("1st bar should be 0, got: %d", got)
		}

		got = strings.Count(bars[1], "━")
		if got != 10 {
			t.Errorf("2nd bar should be 10, got: %d", got)
		}
	})

	t.Run("edge case if value > cap", func(t *testing.T) {
		bars := bars(10, 99999, 10)

		got := strings.Count(bars[0], "━")
		if got != 10 {
			t.Errorf("1st bar should be capped at 10, got: %d", got)
		}

		got = strings.Count(bars[1], "━")
		if got != 0 {
			t.Errorf("2nd bar should be 0, got: %d", got)
		}
	})
}
