package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func runC1() {
	fmt.Println("Running Challenge 1")
	rotationList := readFile[[]rotation]("c1.txt", mapToRotation)
	positionList := dial(rotationList)
	password := countZeros(positionList)
	fmt.Println("Solution 1: ", password)

	passwordTwo := dialTwo(rotationList)
	fmt.Println("Solution 2: ", passwordTwo)
}

type rotation struct {
	direction rune
	count     int
}

func mapToRotation(content []byte) []rotation {
	NEWLINE := '\n'

	buf := [10]rune{}
	bufIndex := 0
	bufFlush := func(buf *[10]rune, bufIndex *int) {
		for i, _ := range buf {
			buf[i] = 0x00
		}
		*bufIndex = 0
	}

	rotationList := make([]rotation, 0, 100)
	appendToRotationList := func(buf *[10]rune, end int) {
		direction := buf[0]
		count, countErr := strconv.Atoi(string(buf[1:end]))
		if countErr != nil {
			panic(countErr)
		}

		rotationList = append(rotationList, rotation{direction, count})
	}

	for _, b := range bytes.Runes(content) {
		if b != NEWLINE {
			buf[bufIndex] = b
			bufIndex++
		} else {
			appendToRotationList(&buf, bufIndex-1)
			bufFlush(&buf, &bufIndex)
		}
	}

	appendToRotationList(&buf, bufIndex)
	bufFlush(&buf, &bufIndex)

	return rotationList
}

func dial(rotationList []rotation) []int {
	ADDITION := 'R'
	SUBTRACTION := 'L'

	position := 50
	positionList := make([]int, 0, 10)

	for _, rotation := range rotationList {
		if rotation.direction == ADDITION {
			position = positiveMod(position+rotation.count, 100)
		}
		if rotation.direction == SUBTRACTION {
			position = positiveMod(position-rotation.count, 100)
		}
		positionList = append(positionList, position)

	}

	return positionList
}

func countZeros(positionList []int) (acc int) {
	for _, position := range positionList {
		if position == 0 {
			acc++
		}
	}

	return acc
}

func positiveMod(a, b int) int {
	return (a%b + b) % b
}

func dialTwo(rotationList []rotation) (password int) {
	start := 50

	for _, rotation := range rotationList {
		password += rotate(&start, rotation.direction, rotation.count)
	}

	return password
}

func rotate(position *int, direction rune, count int) (zeros int) {
	var ops func(a, b int) int

	if direction == 'R' {
		ops = func(a, b int) int {
			return a + b
		}
	} else {
		ops = func(a, b int) int {
			return a - b
		}
	}

	for i := 0; i < count; i++ {
		*position = ops(*position, 1)
		if *position == -1 {
			*position = 99
		}
		if *position == 100 {
			*position = 0
		}
		if *position == 0 {
			zeros++
		}
	}

	return
}
