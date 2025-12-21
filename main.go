package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Println(os.Args)

	if len(args) > 1 {
		runChallenge(args[1])
	} else {
		runChallenge("-1")
	}
}

func runChallenge(input string) {
	switch input {
	case "1":
		runC1()
	case "2":
		runC2()
	case "3":
		runC3()
	case "4":
		runC4()
	case "5":
		runC5()
	case "6":
		runC6()
	case "7":
		runC7()
	case "8":
		runC8()
	case "9":
		runC9()
	default:
		runC9()
	}
}
