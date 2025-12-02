package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func runC2() {
	fmt.Println("Running Challenge 2")

	intervalList := readFile[[]interval]("c2.txt", mapToIntervals)
	invalidSum := execC2(intervalList, testPattern)
	fmt.Println("Solution 1: ", invalidSum)

	invalidSumTwo := execC2(intervalList, testPatternTwo)
	fmt.Println("Solution 2: ", invalidSumTwo)
}

type interval struct {
	start int
	end   int
}

func mapToIntervals(content []byte) []interval {
	COMMA := ','

	buf := [50]rune{}
	bufIndex := 0
	bufFlush := func(buf *[50]rune, bufIndex *int) {
		for i, _ := range buf {
			buf[i] = 0x00
		}
		*bufIndex = 0
	}

	intervalList := make([]interval, 0, 100)

	for i, b := range bytes.Runes(content) {
		if b != COMMA {
			buf[bufIndex] = b
			bufIndex++
		}

		if b == COMMA || i == len(content)-1 {
			intervalList = append(intervalList, createInterval(&buf))
			bufFlush(&buf, &bufIndex)
		}
	}

	return intervalList
}

func createInterval(buf *[50]rune) interval {
	DASH := '-'
	NIL := rune(0x00)

	findSymbol := func(buf *[50]rune, symbol rune) int {
		for i, r := range buf {
			if r == symbol {
				return i
			}
		}
		return -1
	}

	dashIndex := findSymbol(buf, DASH)
	endIndex := findSymbol(buf, NIL)
	startString := string(buf[0:dashIndex])
	endString := string(buf[dashIndex+1 : endIndex])

	start, err := strconv.Atoi(startString)
	end, err := strconv.Atoi(endString)
	if err != nil {
		panic(err)
	}

	return interval{start, end}
}

func execC2(intervalList []interval, test func(int) bool) (sum int) {
	invalidList := make([]int, 0, 100)

	for _, interval := range intervalList {
		// fmt.Println("Calculating interval: ", interval)
		for j := interval.start; j <= interval.end; j++ {
			if test(j) {
				invalidList = append(invalidList, j)
			}
		}
	}

	for _, value := range invalidList {
		sum += value
	}

	return
}

func testPattern(id int) bool {
	idRune := []rune(strconv.Itoa(id))
	if len(idRune)%2 == 1 {
		return false
	}

	middle := len(idRune) / 2
	first := string(idRune[0:middle])
	second := string(idRune[middle:len(idRune)])

	return first == second
}

func testPatternTwo(id int) bool {
	idString := strconv.Itoa(id)
	idRune := []rune(idString)

	for size := 1; size <= len(idString)/2; size++ {
		pattern := string(idRune[0:size])
		isPattern := strings.ReplaceAll(idString, pattern, "")
		if isPattern == "" {
			return true
		}
	}

	return false
}
