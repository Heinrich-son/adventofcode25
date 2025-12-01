package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(os.Args)
	runChallenge(args[1])
}

func runChallenge(input string) {
	switch input {
	case "1":
		runC1()
	default:
		fmt.Println("No corresponding Challenge found. Exiting")
	}
}
