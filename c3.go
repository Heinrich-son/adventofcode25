package main

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func runC3() {
	batteries := readFile("c3.txt", mapToBatteries)
	joltage := execC3(batteries)
	fmt.Println("Solution 1: ", joltage)
	joltageTwo := execC3Two(batteries)
	fmt.Println("Solution 2: ", joltageTwo)
}

func mapToBatteries(content []byte) [][]rune {
	NEWLINE := '\r'

	contentLines := make([]string, 0, 100)
	for part := range strings.SplitSeq(string(content), string(NEWLINE)) {
		contentLines = append(contentLines, strings.ReplaceAll(part, string(NEWLINE), ""))
	}

	batteries := make([][]rune, 0, 100)

	for _, content := range contentLines {
		batteries = append(batteries, []rune(content))
	}

	return batteries
}

func execC3(batterieList [][]rune) (acc int) {
	joltageList := make([]int, 0, len(batterieList))
	getMsr := func(subBatterie []rune) rune {
		biggest := slices.MaxFunc(subBatterie, func(a, b rune) int {
			return cmp.Compare(getInt(a), getInt(b))
		})

		return biggest
	}

	for _, batterie := range batterieList {
		subBatterieOne := batterie[0 : len(batterie)-1]
		msrOne := getMsr(subBatterieOne)
		indexOne := slices.IndexFunc(subBatterieOne, func(n rune) bool {
			return n == msrOne
		})

		subBatterieTwo := batterie[indexOne+1:]
		msrTwo := getMsr(subBatterieTwo)

		joltage, err := strconv.Atoi(string([]rune{msrOne, msrTwo}))
		if err != nil {
			panic(err)
		}

		joltageList = append(joltageList, joltage)
	}

	for _, joltage := range joltageList {
		acc += joltage
	}

	return
}

func getInt(r rune) int {
	return int(r - '0')
}

func execC3Two(batterieList [][]rune) (acc int) {
	joltageList := make([]int, 0, 100)

	for _, batterie := range batterieList {
		msrList := calculateMsrList(batterie, 0)
		joltage, err := strconv.Atoi(string(msrList))
		if err != nil {
			panic(err)
		}
		joltageList = append(joltageList, joltage)
	}

	for _, joltage := range joltageList {
		acc += joltage
	}

	return
}

func calculateMsrList(batterieList []rune, depth int) []rune {
	getMsr := func(batterieList []rune) rune {
		biggest := slices.MaxFunc(batterieList, func(a, b rune) int {
			return cmp.Compare(getInt(a), getInt(b))
		})

		return biggest
	}

	getMsrIndex := func(msr rune, batterieList []rune) int {
		return slices.IndexFunc(batterieList, func(n rune) bool {
			return n == msr
		})
	}

	subBatterieList := batterieList[0 : len(batterieList)-(11-depth)]
	msr := getMsr(subBatterieList)
	msrIndex := getMsrIndex(msr, batterieList)

	nextBatterieList := batterieList[msrIndex+1:]

	msrList := []rune{msr}
	if depth == 11 {
		return msrList
	} else {
		return append(msrList, calculateMsrList(nextBatterieList, depth+1)...)
	}
}
