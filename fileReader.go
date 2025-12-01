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
