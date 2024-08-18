package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

type Choice struct {
	w       io.Writer
	options []string
	current int
}

func NewChoice(w io.Writer, options ...string) *Choice {
	return &Choice{
		w:       w,
		options: options,
	}
}

func (c *Choice) Select() int {
	defer c.reset()
	c.printOptions()

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.Enter, keys.Space:
			return true, nil
		case keys.Esc, keys.CtrlC:
			os.Exit(1)
		case keys.Down:
			if c.current >= len(c.options)-1 {
				c.current = 0
			} else {
				c.current++
			}
		case keys.Up:
			if c.current <= 0 {
				c.current = len(c.options) - 1
			} else {
				c.current--
			}
		}

		c.clearOptions()
		c.printOptions()
		return false, nil
	})

	return c.current
}

func (c *Choice) reset() { c.current = 0 }

func (c *Choice) clearOptions() {
	fmt.Fprintf(c.w, "\033[%dA\r", len(c.options)+2) // plus the info
}

func (c *Choice) printOptions() {
	var buf strings.Builder

	for i, o := range c.options {
		if i == c.current {
			buf.WriteString("\033[1m> ")
			buf.WriteString(o)
			buf.WriteString("\033[0m\n\r")
		} else {
			buf.WriteString("  ")
			buf.WriteString(o)
			buf.WriteString("\n\r")
		}
	}

	buf.WriteString("\npress ENTER or SPACE to confirm\n\r")
	fmt.Fprint(c.w, buf.String())
}
