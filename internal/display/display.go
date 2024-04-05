package display

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

type Display struct {
	White, Red, Green, Black, Yellow ColorFunc
	w                                io.Writer
	isTesting                        bool
}

type ColorFunc func(w io.Writer, format string, a ...interface{})

func New(w io.Writer) *Display {
	return &Display{
		Black:  color.New(color.FgHiBlack, color.Bold).FprintfFunc(),
		White:  color.New(color.FgHiWhite, color.Bold).FprintfFunc(),
		Red:    color.New(color.FgHiRed, color.Bold).FprintfFunc(),
		Green:  color.New(color.FgHiGreen, color.Bold).FprintfFunc(),
		Yellow: color.New(color.FgHiYellow, color.Bold).FprintfFunc(),
		w:      w,
	}
}

// Printf to write using ColorFunc, otherwise just use fmt
func (d *Display) Printf(f ColorFunc, format string, a ...interface{}) {
	f(d.w, format, a...)
}

// Center to centerize the given string format. Length is 46
func (d *Display) Center(f ColorFunc, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)

	if f == nil {
		fmt.Fprintf(d.w, "\n%*s\n", len(s)+(46-len(s))/2, s)
	} else {
		f(d.w, "\n%*s\n", len(s)+(46-len(s))/2, s)
	}
}

// Bar prints the percentage bar of current hp relative to it's max capacity.
// Length is 50 blocks.
func (d *Display) Bar(currHp, maxHp float32) {
	if currHp <= 0 {
		currHp = 0 // avoid panic
	}

	n := currHp / maxHp * 50

	if n > 35 {
		d.Printf(d.Green, "%s", strings.Repeat("█", int(n)))
	} else if n > 15 {
		d.Printf(d.Yellow, "%s", strings.Repeat("█", int(n)))
	} else {
		d.Printf(d.Red, "%s", strings.Repeat("█", int(n)))
	}

	if !d.isTesting {
		d.Printf(d.Black, "%s\n", strings.Repeat("█", 50-int(n)))
	}
}

func Clear() {
	if runtime.GOOS != "windows" {
		return
	}

	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout

	if err := c.Run(); err != nil {
		log.Fatal(err)
	}
}
