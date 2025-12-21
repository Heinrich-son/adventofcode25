package main

import (
	"os"
)

// Reads a file from ./assets in its entirety and applies a mapper function on its content.
// Returns the mapped content.
func readFile[T any](path string, mapper func([]byte) T) T {
	content, contentErr := os.ReadFile("assets/" + path)
	if contentErr != nil {
		panic(contentErr)
	}

	return mapper(content)
}

func writeTileDrawning(drawing []rune, maxY int, maxX int) {
	f, err := os.OpenFile("assets/out.txt", 2|64, 0666)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for l := range maxY {
		str := string(drawing[l*maxY : maxX])
		_, err := f.WriteString(str)
		if err != nil {
			panic(err)
		}
	}
}
