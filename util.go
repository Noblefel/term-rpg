package main

import (
	"fmt"
	"math/rand/v2"
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
