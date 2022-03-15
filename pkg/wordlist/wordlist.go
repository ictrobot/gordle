package wordlist

import (
	"github.com/ictrobot/gordle/internal/util"
	"math/rand"
	"strings"
	"time"
)

type WordList struct {
	Answers    []string
	Allowed    []string
	StartDate  time.Time
	WordLength int
}

func (w *WordList) GetAnswer(date time.Time) (int, string) {
	day := int(date.Sub(w.StartDate).Hours()) / 24
	if day < 0 {
		day = 0
	}

	day %= len(w.Answers)
	return day, strings.ToUpper(w.Answers[day])
}

func (w *WordList) IsAllowed(s string) bool {
	s = strings.ToLower(s)
	return util.SliceContains(w.Answers, s) || util.SliceContains(w.Allowed, s)
}

func (w *WordList) GetRandomDate() time.Time {
	rand.Seed(time.Now().UnixNano())
	return w.StartDate.AddDate(0, 0, rand.Intn(len(w.Answers)))
}

// TODO from file
