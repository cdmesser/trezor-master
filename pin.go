// Package trezor implements a Trezor PIN entry UI.
//
// Unlocking a Trezor requires entering a PIN. The Trezor device itself shows a scrambled
// PIN entry pad on the screen, and the client software (eg this library) is responsible
// for providing a blank PIN entry pad for the user to input their PIN with. The user treats
// that blank entry pad as if it was labeled with numbers in the same order as the scramble
// they see on their Trezor screen. But the client software (this library) does not know
// the scramble, so it reports to Trezor the PIN entered by the user as if the entry pad
// had been a normal 1-9 grid.
//
// This is easier to understand if you have gone through the process yourself. You can
// do so by using a physical Trezor and the official Trezor wallet software.
//
// This library implements a full-screen terminal-based UI for PIN entry.

package trezor

import (
	"errors"
	"strings"

	termbox "github.com/sml/termbox-go"
)

var ErrUserCancelledInput = errors.New("user cancelled PIN entry")

// GetPIN implements Trezor PIN entry with a full-screen terminal-based UI.
func GetPIN(prompt string) (string, error) {
	err := termbox.Init()
	if err != nil {
		return "", err
	}
	defer termbox.Close()

	var (
		cursorX = 1
		cursorY = 1
		pin     string
	)

	clamp := func(x, min, max int) int {
		if x < min {
			return min
		}
		if x > max {
			return max
		}
		return x
	}

	printStr := func(x, y int, s string) {
		i := 0 // Rune index. (Not using the index from `range` because it's a byte index)
		for _, r := range s {
			termbox.SetCell(x+i, y, r, termbox.ColorDefault, termbox.ColorDefault)
			i++
		}
	}

	var keypad = [][]string{
		[]string{"7", "8", "9"},
		[]string{"4", "5", "6"},
		[]string{"1", "2", "3"},
	}

	for {
		// Render.
		{
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

			printStr(0, 0, prompt)

			printStr(3, 2, "●●●")
			printStr(3, 3, "●●●")
			printStr(3, 4, "●●●")

			printStr(0, 6, "PIN: "+strings.Map(func(rune) rune { return '●' }, pin))

			printStr(0, 9, "[Arrow keys]: move cursor")
			printStr(0, 10, "[Space]:      press button under cursor")
			printStr(0, 11, "[Backspace]:  delete")
			printStr(0, 12, "[Enter]:      submit PIN")
			printStr(0, 13, "[q]:          exit without submitting PIN")

			termbox.SetCursor(cursorX+3, cursorY+2)
			termbox.Flush()
		}

		// Update state.
		{
			event := termbox.PollEvent()
			if event.Type != termbox.EventKey {
				continue
			}

			if event.Ch == 'q' || event.Key == termbox.KeyCtrlC {
				return "", ErrUserCancelledInput
			}

			switch event.Key {
			case termbox.KeyArrowUp:
				cursorY--
			case termbox.KeyArrowDown:
				cursorY++
			case termbox.KeyArrowLeft:
				cursorX--
			case termbox.KeyArrowRight:
				cursorX++
			case termbox.KeySpace:
				pin += keypad[cursorY][cursorX]
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if len(pin) > 0 {
					pin = pin[:len(pin)-1]
				}
			case termbox.KeyEnter:
				return pin, nil
			}

			cursorY = clamp(cursorY, 0, 2)
			cursorX = clamp(cursorX, 0, 2)
		}
	}
}
