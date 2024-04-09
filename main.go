package main

import (
	"bufio"
	"os"

	"github.com/Noblefel/term-rpg/internal/display"
	"github.com/Noblefel/term-rpg/internal/game"
)

func main() {
	d := display.New(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)

	app := game.New(scanner, d)
	app.Start()
}
