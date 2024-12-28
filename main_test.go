package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Noblefel/vivi"
)

func TestMain(m *testing.M) {
	vivi.Out = io.Discard
	os.Exit(m.Run())
}

// func TestMenuMain(t *testing.T) {
// 	var sb strings.Builder
// 	out = &sb
// 	player = &Player{}
// 	go menuMain()
// 	time.Sleep(50 * time.Millisecond)

// 	t.Run("display info correctly", func(t *testing.T) {
// 		defer sb.Reset()

// 		hpbar := barhp(36, player.hp, player.hpcap)
// 		if !strings.Contains(sb.String(), hpbar) {
// 			t.Error("incorrect hpbar")
// 		}

// 		if !strings.Contains(sb.String(), "Gold   : 0") {
// 			t.Error("incorrect gold")
// 		}

// 		if !strings.Contains(sb.String(), "Stage  : 1") {
// 			t.Error("incorrect stage")
// 		}
// 	})

// 	t.Run("choosing rest with 0 gold", func(t *testing.T) {
// 		defer sb.Reset()
// 		defer keyboard.SimulateKeyPress(keys.Enter)

// 		keyboard.SimulateKeyPress(keys.Down)
// 		keyboard.SimulateKeyPress(keys.Down)
// 		keyboard.SimulateKeyPress(keys.Enter)
// 		time.Sleep(50 * time.Millisecond)

// 		if !strings.Contains(sb.String(), "You don't have enough money to rest") {
// 			t.Error("should fail when attempting to rest")
// 		}
// 	})

// 	t.Run("choosing rest with 0 gold", func(t *testing.T) {
// 		defer sb.Reset()
// 		defer keyboard.SimulateKeyPress(keys.Enter)

// 		keyboard.SimulateKeyPress(keys.Down)
// 		keyboard.SimulateKeyPress(keys.Down)
// 		keyboard.SimulateKeyPress(keys.Down)
// 		keyboard.SimulateKeyPress(keys.Enter)
// 		time.Sleep(50 * time.Millisecond)

// 		if !strings.Contains(sb.String(), "You don't have enough money to train") {
// 			t.Error("should fail when attempting to train")
// 		}
// 	})
// }

// func TestMenuAttributes(t *testing.T) {
// 	var sb strings.Builder
// 	out = &sb
// 	player = &Player{hpcap: 999, defense: 120, strength: 570}
// 	go menuAttributes()
// 	time.Sleep(50 * time.Millisecond)

// 	if !strings.Contains(sb.String(), "999") {
// 		t.Error("should display number of 999 (hp cap)")
// 	}

// 	if !strings.Contains(sb.String(), "120") {
// 		t.Error("should display number of 120 (defense)")
// 	}

// 	if !strings.Contains(sb.String(), "570") {
// 		t.Error("should display number of 570 (strength)")
// 	}
// }

// func TestMenuBattles(t *testing.T) {
// 	var sb strings.Builder
// 	out = &sb
// 	player = &Player{hpcap: 10, hp: 10}

// 	t.Run("entering battle with 0 hp", func(t *testing.T) {
// 		defer sb.Reset()
// 		defer keyboard.SimulateKeyPress(keys.Enter)
// 		player.hp = 0
// 		go menuBattle(nil)
// 		time.Sleep(50 * time.Millisecond)

// 		if !strings.Contains(sb.String(), "lost") {
// 			t.Error("should tell that player had lost")
// 		}

// 		os.WriteFile("testlog.txt", []byte(sb.String()), os.ModePerm)
// 	})

// 	t.Run("when enemy has 0 hp left", func(t *testing.T) {
// 		defer sb.Reset()
// 		defer keyboard.SimulateKeyPress(keys.Enter)
// 		player.hp = 10
// 		enemy := Attributes{hp: 0}
// 		go menuBattle(&enemy)
// 		time.Sleep(50 * time.Millisecond)

// 		if !strings.Contains(sb.String(), "won") {
// 			t.Error("should tell that player had won")
// 		}
// 	})
// }

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
