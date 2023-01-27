package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode"
	"os"

	"github.com/sio/wordle"
)

func main() {
	sources := map[string]fetchIterator{
		"harrix.txt": &harrixIterator{},
	}
	for filename, iter := range sources {
		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		defer iter.Close()

		var count uint
		for iter.Next() {
			word := strings.TrimSpace(iter.Value())
			if len(word) != wordle.WordSize {
				continue
			}
			for _, char := range word {
				if !unicode.IsLetter(char) || !unicode.IsLower(char) {
					continue
				}
			}
			_, err = fmt.Fprintln(out, word)
			if err != nil {
				log.Fatalf("writing to %s failed: %v", filename, err)
			}
			count++
		}
		log.Printf("wrote %d words to %s", count, filename)
	}
}

func open(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return resp.Body, err
	}
	if resp.StatusCode != 200 {
		return resp.Body, fmt.Errorf("HTTP %s: %s", resp.Status, url)
	}
	return resp.Body, nil
}
