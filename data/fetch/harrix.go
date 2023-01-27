package main

import (
	"bufio"
	"io"
	"log"
)

const harrixURL = "https://github.com/Harrix/Russian-Nouns/raw/main/dist/russian_nouns.txt"

// Fetch noun list from https://github.com/Harrix/Russian-Nouns (MIT)
type harrixIterator struct {
	reader  io.ReadCloser
	scanner *bufio.Scanner
	value   string
}

func (h *harrixIterator) init() {
	reader, err := open(harrixURL)
	if err != nil {
		log.Fatal(err)
	}
	h.reader = reader
	h.scanner = bufio.NewScanner(h.reader)
}

func (h *harrixIterator) Next() bool {
	if h.reader == nil {
		h.init()
	}
	if h.scanner.Scan() {
		h.value = h.scanner.Text()
		return true
	}
	if err := h.scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false
}

func (h *harrixIterator) Value() string {
	return h.value
}

func (h *harrixIterator) Close() error {
	err := h.reader.Close()
	h = &harrixIterator{}
	return err
}
