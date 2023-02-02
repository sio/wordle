package wordle

import (
	"bufio"
	"fmt"
	"log"
)

type Word [WordSize]rune

func (w Word) String() string {
	return string(w[:])
}

func (w *Word) Parse(input string) error {
	var index int
	var char rune
	for _, char = range input {
		if index > WordSize {
			return fmt.Errorf("word too long: %s (%d characters instead of %d)", input, len([]rune(input)), WordSize)
		}
		(*w)[index] = char
		index++
	}
	if index != WordSize-1 {
		return fmt.Errorf("word too short: %s (%d characters instead of %d)", input, len([]rune(input)), WordSize)
	}
	return nil
}

func (w *Word) Contains(char rune) bool {
	for _, c := range *w {
		if c == char {
			return true
		}
	}
	return false
}

type void struct{}

func RussianWords() *[]Word {
	const dataDir = "data"
	fs := RussianWordLists

	entries, err := fs.ReadDir(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	words := make(map[Word]void)
	for _, entry := range entries {
		file, err := fs.Open(fmt.Sprintf("%s/%s", dataDir, entry.Name()))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var count int
		for scanner.Scan() {
			var word Word
			word.Parse(scanner.Text())
			words[word] = void{}
			count++
		}
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	output := make([]Word, len(words))
	var i int
	for word := range words {
		output[i] = word
		i++
	}
	return &output
}
