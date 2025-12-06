package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ~30min parsing
// ~25min part one
// ~35min part two

func runC6() {
	mathProblemList := readFile("c6.txt", mapToMathProblem)
	solution := resolveAllMathProblems(mathProblemList)
	fmt.Println("Solution 1:", solution)

	mathProblemListTwo := readFile("c6.txt", mapToMathProblemTwo)
	solutionTwo := resolveAllMathProblems(mathProblemListTwo)
	fmt.Println("Solution 2:", solutionTwo)
}

func mapToMathProblem(content []byte) [][]string {
	NEWLINE := "\n"
	SPACE := " "

	lines := strings.Split(string(content), NEWLINE)

	linesByWords := make([][]string, 0, len(lines))
	for _, line := range lines {
		linesUncleaned := strings.Split(line, SPACE)
		linesCleaned := slices.DeleteFunc(linesUncleaned, func(s string) bool {
			return s == ""
		})

		linesByWords = append(linesByWords, linesCleaned)
	}

	mathProblemList := make([][]string, 0, len(linesByWords[0]))
	for i := range len(linesByWords[0]) {
		mathProblem := make([]string, 0, len(linesByWords))
		for _, line := range linesByWords {
			mathProblem = append(mathProblem, line[i])
		}
		mathProblemList = append(mathProblemList, mathProblem)
	}

	return mathProblemList
}

func resolveAllMathProblems(mathProblemList [][]string) int {
	solutionList := make([]int, 0, len(mathProblemList[0]))

	for _, mathProblem := range mathProblemList {
		solution := resolveMathProblem(mathProblem)
		solutionList = append(solutionList, solution)
	}

	finalSolution := applyOperation(solutionList, "+")
	return finalSolution
}

func resolveMathProblem(mathProblem []string) int {
	operation := mathProblem[len(mathProblem)-1]
	numbersString := mathProblem[:len(mathProblem)-1]
	numbersInt := toInt(numbersString...)

	solution := applyOperation(numbersInt, operation)
	return solution
}

func toInt(sliceString ...string) []int {
	sliceInt := make([]int, 0, len(sliceString))
	for _, s := range sliceString {
		s := strings.TrimSpace(s)
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		sliceInt = append(sliceInt, i)
	}

	return sliceInt
}

func applyOperation(numbers []int, operator string) int {
	sum := func(numbers []int) (acc int) {
		for _, number := range numbers {
			acc += number
		}
		return acc
	}
	mul := func(numbers []int) int {
		acc := 1
		for _, number := range numbers {
			acc *= number
		}
		return acc
	}

	if operator == "+" {
		return sum(numbers)
	} else {
		return mul(numbers)
	}
}

// split by newline and convert lines to []runes
// read first 3 lines simultanously from right to left => len(x)-1 to 0
// from top to bottom convert to string, strip spaces and append to "mathProblem"
// if line 4 is present => operation is there and append to "mathProblem" and then "mathProblemList"
// if all lines are spaces prepare for new "mathProblem"
func mapToMathProblemTwo(content []byte) [][]string {
	NEWLINE := '\n'
	SUM := '+'
	MUL := '*'

	linesString := strings.Split(string(content), string(NEWLINE))
	linesRunes := make([][]rune, 0, len(linesString))
	for _, line := range linesString {
		linesRunes = append(linesRunes, []rune(line))
	}

	mathProblemList := make([][]string, 0, len(linesRunes[0]))
	mathProblem := make([]string, 0, len(linesRunes)-1)
	for i := len(linesRunes[0]) - 1; i >= 0; i-- {
		var oparation rune
		numberRune := make([]rune, 0, len(linesRunes)-1)
		for lineIndex, line := range linesRunes {
			if lineIndex < len(linesRunes)-1 {
				numberRune = append(numberRune, line[i])
			} else {
				oparation = line[i]
			}
		}

		numberString := strings.TrimSpace(string(numberRune))
		if numberString != "" {
			mathProblem = append(mathProblem, numberString)
			if oparation == SUM || oparation == MUL {
				mathProblem = append(mathProblem, string(oparation))
				mathProblemList = append(mathProblemList, mathProblem)
			}
		} else {
			mathProblem = make([]string, 0, len(linesRunes)-1)

		}
	}

	return mathProblemList
}
