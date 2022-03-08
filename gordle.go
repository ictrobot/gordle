package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/ictrobot/gordle/pkg/game"
	"github.com/ictrobot/gordle/pkg/wordlist"
	"github.com/spf13/pflag"
	"os"
	"time"
)

type Gordle struct {
	Date       time.Time
	List       *wordlist.WordList
	NumGuesses int

	game   *game.Game
	screen tcell.Screen
}

func main() {
	dateString := pflag.String("date", "", "Date")
	randomDate := pflag.Bool("random", false, "Random Date")
	numGuesses := pflag.Int("guesses", 6, "Number of guesses")
	pflag.Parse()

	list := &wordlist.Original
	date := time.Now().Truncate(24 * time.Hour)

	if len(*dateString) > 0 && *randomDate {
		fmt.Fprintf(os.Stderr, "Cannot supply both --date and --random\n")
		os.Exit(1)
	} else if len(*dateString) > 0 {
		parsed, err := time.Parse("2006-01-02", *dateString)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		} else {
			date = parsed
		}
	} else if *randomDate {
		date = list.GetRandomDate()
	}

	m := Gordle{
		List:       list,
		Date:       date,
		NumGuesses: *numGuesses,
	}
	m.Main()
}

func (m *Gordle) Main() {
	m.newGame()
	m.initScreen()

	m.game.Draw(m.screen)
	for {
		switch ev := m.screen.PollEvent().(type) {
		case *tcell.EventResize:
			// TODO check window size
			m.game.Draw(m.screen)
			m.screen.Sync()
		case *tcell.EventKey:
			m.handleKey(ev)
		}
	}
}

func (m *Gordle) handleKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC || (m.game.IsFinished() && ev.Rune() == 'q') {
		m.screen.Fini()
		os.Exit(0)
	} else if m.game.IsFinished() {
		// handle menu input

		if ev.Rune() == 'n' {
			m.Date = m.Date.AddDate(0, 0, 1)
			m.newGame()
		} else if ev.Rune() == 'p' {
			m.Date = m.Date.AddDate(0, 0, -1)
			m.newGame()
		} else if ev.Rune() == 'r' {
			m.Date = m.List.GetRandomDate()
			m.newGame()
		}
	} else {
		// handle game input

		if ev.Rune() >= 'a' && ev.Rune() <= 'z' {
			m.game.Keypress(ev.Rune() + ('A' - 'a'))
		} else if ev.Rune() >= 'A' && ev.Rune() <= 'Z' {
			m.game.Keypress(ev.Rune())
		} else if ev.Key() == tcell.KeyEnter {
			m.game.Enter()
		} else if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
			m.game.Backspace()
		}
	}

	m.game.Draw(m.screen)
	m.screen.Show()
}

func (m *Gordle) newGame() {
	if m.List == nil {
		panic("Word list is not set")
	}
	if m.Date == (time.Time{}) {
		panic("Date is not set")
	}
	if m.NumGuesses < 1 {
		panic("Number of guesses is less than 1")
	}

	m.game = game.NewGame(m.List, m.Date, m.NumGuesses)
}

func (m *Gordle) initScreen() {
	s, e := tcell.NewScreen()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	m.screen = s
}
