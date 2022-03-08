package game

import (
	"github.com/gdamore/tcell/v2"
)

type keyboardRow struct {
	letters []rune
	offset  int
}

var keyboardLayout = [...]keyboardRow{
	{
		letters: []rune{'Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'},
		offset:  0,
	},
	{
		letters: []rune{'A', 'S', 'D', 'F', 'G', 'H', 'J', 'K', 'L'},
		offset:  1,
	},
	{
		letters: []rune{'Z', 'X', 'C', 'V', 'B', 'N', 'M'},
		offset:  3,
	},
}

func drawKeyboard(s tcell.Screen, x int, y int, results map[rune]result) (int, int) {
	maxX := 0
	maxY := 0

	for i, row := range keyboardLayout {
		for j, letter := range row.letters {
			keyX := x + row.offset + 4*j
			keyY := y + 4*i
			drawKey(s, keyX, keyY, letter, results[letter])

			if keyX > maxX {
				maxX = keyX
			}
			if keyY > maxY {
				maxY = keyY
			}
		}
	}

	return maxX + 3, maxY + 3
}
