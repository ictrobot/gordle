package game

import "github.com/gdamore/tcell/v2"

var styles = map[result]tcell.Style{
	unknown:     tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack),
	notIncluded: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite),
	wrongPos:    tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack),
	correct:     tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorWhite),
}

func drawKey(s tcell.Screen, x int, y int, r rune, result result) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			char := ' '
			if i == 1 && j == 1 {
				char = r
			}

			s.SetContent(x+i, y+j, char, nil, styles[result])
		}
	}
}

func drawText(s tcell.Screen, x int, y int, str string, style tcell.Style) {
	for i, l := range str {
		s.SetContent(x+i, y, l, nil, style)
	}
}
