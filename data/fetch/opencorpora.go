package main

import (
	"bufio"
	"compress/bzip2"
	"io"
	"log"
	"regexp"
	"strings"
)

// Fetch nouns from http://www.opencorpora.org/dict.php (CC BY-SA 3.0)
type opencorporaIterator struct {
	webReader io.ReadCloser
	reader    io.Reader
	scanner   *bufio.Scanner
	value     string
}

const opencorporaURL = "http://www.opencorpora.org/files/export/dict/dict.opcorpora.txt.bz2"

func (i *opencorporaIterator) init() {
	webReader, err := open(opencorporaURL)
	if err != nil {
		log.Fatal(err)
	}
	i.webReader = webReader
	i.reader = bzip2.NewReader(i.webReader)
	i.scanner = bufio.NewScanner(i.reader)
}

// Regex to find nouns in Opencorpora dictionary dump
var opencorporaRegex = regexp.MustCompile(`^(\S+)\s+NOUN,.*sing,nomn$`)

// Noun categories to drop <http://www.opencorpora.org/dict.php?act=gram>
var opencorporaDrop = regexp.MustCompile(`Abbr|Dist|Erro|Geox|Infr|Init|Name|Orgn|Patr|Slng|Surn`)

func (i *opencorporaIterator) Next() bool {
	if i.reader == nil {
		i.init()
	}
	for i.scanner.Scan() {
		line := i.scanner.Text()
		match := opencorporaRegex.FindStringSubmatch(line)
		if len(match) < 2 {
			continue
		}
		if opencorporaDrop.MatchString(line) {
			continue
		}
		i.value = strings.ToLower(match[1])
		return true
	}
	if err := i.scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false
}

func (i *opencorporaIterator) Value() string {
	return i.value
}

func (i *opencorporaIterator) Close() error {
	err := i.webReader.Close()
	i = &opencorporaIterator{}
	return err
}
