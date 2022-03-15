package wordlist

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"time"
)

func FromFile(answerFile string, allowedFile string, startDate time.Time) (*WordList, error) {
	if answerFile == "" {
		return nil, errors.New("loading wordlist: no answer file provided")
	}
	answers, err := loadWords(answerFile)
	if err != nil {
		return nil, err
	}
	if len(answers) == 0 {
		return nil, errors.New("loading wordlist: no valid words found")
	}

	allowed := make([]string, 0)
	if allowedFile != "" {
		allowed, err = loadWords(allowedFile)
		if err != nil {
			return nil, err
		}
	}

	wordLength := len(answers[0])
	if !checkEveryLengthEquals(answers, wordLength) || !checkEveryLengthEquals(allowed, wordLength) {
		return nil, errors.New("loading wordlist: words of different length found")
	}

	return &WordList{
		Answers:    answers,
		Allowed:    allowed,
		StartDate:  startDate,
		WordLength: wordLength,
	}, nil
}

func loadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	words := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ToLower(strings.TrimSpace(scanner.Text()))

		onlyLetters := true
		for _, c := range []rune(line) {
			if c < 'a' || c > 'z' {
				onlyLetters = false
				break
			}
		}

		if onlyLetters {
			words = append(words, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

func checkEveryLengthEquals(strings []string, length int) bool {
	for _, s := range strings {
		if len(s) != length {
			return false
		}
	}
	return true
}
