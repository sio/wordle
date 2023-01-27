package wordle

import "embed"

const WordSize = 5

//go:embed data/harrix.txt
//go:embed data/opencorpora.txt
var RussianWordLists embed.FS
