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

	answerFile := pflag.String("answer-file", "", "Custom word list: file to load answers from")
	allowedFile := pflag.String("allowed-file", "", "Custom word list: file to load allowed guesses from (optional)")
	startDateString := pflag.String("start-date", "", "Custom word list: start date of the custom word list (optional if --random is provided)")
	pflag.Parse()

	list := getWordList(*answerFile, *allowedFile, *startDateString, *randomDate)
	date := today()

	if len(*dateString) > 0 && *randomDate {
		fmt.Fprintf(os.Stderr, "Cannot supply both --date and --random\n")
		os.Exit(1)
	} else if len(*dateString) > 0 {
		date = parseDate(*dateString)
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

func getWordList(answerFile string, allowedFile string, startDateString string, random bool) *wordlist.WordList {
	if answerFile == "" {
		if allowedFile != "" {
			fmt.Fprintf(os.Stderr, "--answer-file must be provided with --allowed-file\n")
			os.Exit(1)
		}
		if startDateString != "" {
			fmt.Fprintf(os.Stderr, "--answer-file must be provided with --start-date\n")
			os.Exit(1)
		}

		return wordlist.Original
	}

	startDate := today()
	if startDateString != "" {
		startDate = parseDate(startDateString)
	} else if startDateString == "" && !random {
		fmt.Fprintf(os.Stderr, "Either --start-date or --random must be provided with --answer-file\n")
		os.Exit(1)
	}

	list, err := wordlist.FromFile(answerFile, allowedFile, startDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	return list
}

func parseDate(s string) time.Time {
	parsed, err := time.Parse("2006-01-02", s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	return parsed
}

func today() time.Time {
	return time.Now().Truncate(24 * time.Hour)
}
