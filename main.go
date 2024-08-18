package main

import (
	"os"
)

func main() {
	game := Game{
		writer: os.Stdout,
		stage:  1,
	}

	game.printf("\033[H\033[J")
	game.selectPerks()

	for {
		game.menu()
	}
}
