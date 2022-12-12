package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	_ "embed"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/fatih/color"
)

func isBackspace(key keys.Key) bool {
	return key.Code == keys.CtrlH || key.Code == keys.Backspace
}

func isDebug(c int, r byte, k keys.Key) {
	cursor.StartOfLineUp(0)
	cursor.ClearLine()
	fmt.Printf("debug: %d %s %v", c, string(r), k)
}

//go:embed dict.txt
var dict string

func generateWords() string {
	// lazy seed
	rand.Seed(time.Now().Unix())

	lines := strings.Split(dict, "\n")
	left := rand.Intn(len(lines) - 10)
	words := lines[left : left+10]
	for i := range words {
		r := []byte(words[i])
		// switch first character to UPPER CASE
		r[0] -= 32
		words[i] = string(r)
	}
	return strings.Join(words, " ")
}

func main() {
	correctTextColor := color.New(color.FgHiBlue)
	wrongTextColor := color.New(color.FgHiYellow, color.BgHiRed)
	idleTextColor := color.New(color.FgHiYellow)

	var c int
	// isDebug(c, 0, keys.Key{})
	// fmt.Println()
	words := generateWords()
	idleTextColor.Print(words)
	cursor.StartOfLine()

	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		if key.Code == keys.CtrlC {
			return true, nil
		}
		if isBackspace(key) {
			c = int(math.Max(0, float64(c-1)))
			cursor.HorizontalAbsolute(c)
			return false, nil
		}
		if c == len(words)-1 {
			cursor.ClearLine()
			cursor.StartOfLine()
			words = generateWords()
			idleTextColor.Print(words)
			cursor.StartOfLine()
			c = 0
			return false, nil
		}
		c++

		// check correct characters
		w := key.String()
		if key.Code == keys.Space {
			w = " "
		}
		r := string(words[c-1])
		if w != r {
			wrongTextColor.Print(r)
		} else {
			correctTextColor.Print(r)
		}
		cursor.HorizontalAbsolute(c)

		// isDebug(c, words[c-1], key)
		// cursor.StartOfLineDown(0)
		cursor.HorizontalAbsolute(c)

		return false, nil
	})
}
