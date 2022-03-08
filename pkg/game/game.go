package game

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/ictrobot/gordle/pkg/wordlist"
	"time"
)

type Game struct {
	WordList     *wordlist.WordList
	guesses      [][]rune
	guessResults [][]result
	keyResults   map[rune]result
	answer       []rune
	date         time.Time
	day          int
	won          bool
	gameOver     bool
	currentGuess int
}

func NewGame(wordlist *wordlist.WordList, date time.Time, numGuesses int) *Game {
	g := new(Game)
	g.WordList = wordlist
	g.guesses = make([][]rune, numGuesses)
	g.guessResults = make([][]result, numGuesses)
	g.keyResults = make(map[rune]result)

	day, answer := wordlist.GetAnswer(date)
	g.answer = []rune(answer)
	g.date = date.Truncate(24 * time.Hour)
	g.day = day

	return g
}

func (g *Game) IsFinished() bool {
	return g.won || g.gameOver
}

func (g *Game) Keypress(x rune) {
	if x < 'A' || x > 'Z' {
		return
	}
	if g.won {
		return
	}
	if g.currentGuess >= len(g.guesses) {
		return
	}
	if len(g.guesses[g.currentGuess]) >= g.WordList.WordLength {
		return
	}

	g.guesses[g.currentGuess] = append(g.guesses[g.currentGuess], x)
}

func (g *Game) Enter() {
	if g.won {
		return
	}
	if g.currentGuess >= len(g.guesses) {
		return
	}
	if len(g.guesses[g.currentGuess]) != g.WordList.WordLength {
		return
	}
	if !g.WordList.IsAllowed(string(g.guesses[g.currentGuess])) {
		return
	}

	won := true
	for i, r := range g.guesses[g.currentGuess] {
		if r != g.answer[i] {
			won = false
			break
		}
	}

	g.calculateGuessResults()
	g.updateKeyResults()
	g.currentGuess++

	if won {
		g.won = true
	} else if g.currentGuess >= len(g.guesses) {
		g.gameOver = true
	}
}

func (g *Game) Backspace() {
	if g.currentGuess < len(g.guesses) && len(g.guesses[g.currentGuess]) > 0 {
		g.guesses[g.currentGuess] = g.guesses[g.currentGuess][:len(g.guesses[g.currentGuess])-1]
	}
}

func (g *Game) Draw(s tcell.Screen) {
	s.Clear()

	rhsX, y := drawKeyboard(s, 0, 4, g.keyResults)

	drawText(s, rhsX/2-3, 1, "gordle", tcell.StyleDefault.Bold(true))

	g.drawGuesses(s, rhsX+3, 0)

	dayStr := fmt.Sprintf("Day #%d: %s", g.day, g.date.Format("2006/01/02"))
	drawText(s, 0, y+1, dayStr, tcell.StyleDefault)

	if g.won {
		drawText(s, 0, y+3, "You won!", tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
	}
	if g.gameOver {
		gameOverStr := fmt.Sprintf("Game over - Word was '%s'", string(g.answer))
		drawText(s, 0, y+3, gameOverStr, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true))
	}

	if g.won || g.gameOver {
		drawText(s, 0, y+5, "Press n for the next day", tcell.StyleDefault)
		drawText(s, 0, y+6, "Press p for the previous day", tcell.StyleDefault)
		drawText(s, 0, y+7, "Press r for a random day", tcell.StyleDefault)
	}
}

func (g *Game) drawGuesses(s tcell.Screen, x int, y int) {
	for i := 0; i < len(g.guesses); i++ {
		for j := 0; j < g.WordList.WordLength; j++ {
			letter := ' '
			if j < len(g.guesses[i]) {
				letter = g.guesses[i][j]
			}

			result := unknown
			if j < len(g.guessResults[i]) {
				result = g.guessResults[i][j]
			}

			drawKey(s, x+j*4, y+i*4, letter, result)
		}
	}
}

func (g *Game) calculateGuessResults() {
	g.guessResults[g.currentGuess] = make([]result, g.WordList.WordLength)
	// remaining correct letters, used to ensure the correct number of yellows
	correctLetters := make(map[rune]int)

	for i, correctLetter := range g.answer {
		if g.guesses[g.currentGuess][i] == correctLetter {
			g.guessResults[g.currentGuess][i] = correct
		} else {
			g.guessResults[g.currentGuess][i] = notIncluded
			correctLetters[correctLetter]++
		}
	}

	for i, guessedLetter := range g.guesses[g.currentGuess] {
		if g.guessResults[g.currentGuess][i] == notIncluded && correctLetters[guessedLetter] > 0 {
			g.guessResults[g.currentGuess][i] = wrongPos
			correctLetters[guessedLetter]--
		}
	}
}

func (g *Game) updateKeyResults() {
	for i, guessedLetter := range g.guesses[g.currentGuess] {
		if g.answer[i] == guessedLetter {
			g.keyResults[guessedLetter] = correct
		} else if g.inAnswer(guessedLetter) {
			// don't overwrite correct
			if g.keyResults[guessedLetter] < wrongPos {
				g.keyResults[guessedLetter] = wrongPos
			}
		} else {
			g.keyResults[guessedLetter] = notIncluded
		}
	}
}

func (g *Game) inAnswer(r rune) bool {
	for _, l := range g.answer {
		if l == r {
			return true
		}
	}
	return false
}
