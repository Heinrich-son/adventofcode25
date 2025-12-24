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

func writeTileDrawning(drawingBoard [][]rune) {
	f, err := os.OpenFile("assets/out.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for _, line := range drawingBoard {
		str := string(line)
		_, err := f.WriteString(str + "\n")
		if err != nil {
			panic(err)
		}
	}
}
