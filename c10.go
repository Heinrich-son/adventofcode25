package main

import (
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"
)

/*
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
indicator - [.##.] where . is off and # is on
button set - (3) (1,3) (2)... available buttons to change the indicator
button - (1,3) contains the indexes of affected indicators
joltage requirements - {3,5,4,7} part 2

Problem: Each indicator light is initially off. Determine the fewest possible
total presses to reach the wanted indicator state
*/

func runC10() {
	indicatorList := readFile("c10.txt", mapToIndicatorList)
	/*for _, indicator := range indicatorList {
		indicator.displayInitial()
	}*/

	execC10(indicatorList)
	execC10Two(indicatorList)
}

func mapToIndicatorList(content []byte) (indicatorList []indicator_S) {
	NEWLINE := "\n"
	SPACE := " "
	COMMA := ","

	parseButton := func(buttonStringSlice []string) [][]int {
		buttonSetList := make([][]int, 0, len(buttonStringSlice))

		for _, buttonStringSet := range buttonStringSlice {
			buttonStringSetTrimmed := strings.TrimFunc(buttonStringSet, func(r rune) bool {
				if r == '(' || r == ')' {
					return true
				} else {
					return false
				}
			})

			buttonStringSetValues := strings.Split(buttonStringSetTrimmed, COMMA)
			buttonSet := make([]int, 0, len(buttonStringSetValues))
			for _, buttonStringValue := range buttonStringSetValues {
				button, err := strconv.Atoi(buttonStringValue)
				if err != nil {
					panic(err)
				}
				buttonSet = append(buttonSet, button)
			}

			buttonSetList = append(buttonSetList, buttonSet)
		}

		return buttonSetList
	}

	parseJoltage := func(joltageString string) []int {
		joltageStringTrimmed := strings.TrimFunc(joltageString, func(r rune) bool {
			if r == '{' || r == '}' {
				return true
			} else {
				return false
			}
		})

		joltageStringList := strings.Split(joltageStringTrimmed, COMMA)
		joltageSet := make([]int, 0, len(joltageStringList))
		for _, joltageStringValue := range joltageStringList {
			joltage, err := strconv.Atoi(joltageStringValue)
			if err != nil {
				panic(err)
			}
			joltageSet = append(joltageSet, joltage)
		}

		return joltageSet
	}

	lines := strings.SplitSeq(string(content), NEWLINE)
	for machine := range lines {
		sections := strings.Split(machine, SPACE)
		indicatorRunes := []rune(sections[0])
		buttonSlice := sections[1 : len(sections)-1]
		joltageString := sections[len(sections)-1]

		buttonSetList := parseButton(buttonSlice)
		joltageSet := parseJoltage(joltageString)

		indicator := indicator_S{joltageSet: joltageSet, buttonSetList: buttonSetList}.fromRunes(indicatorRunes)
		indicatorList = append(indicatorList, *indicator)
	}

	return
}

type indicator_S struct {
	closingValue int
	currentValue int

	joltagePresses int

	closingRunes  []rune
	joltageSet    []int
	buttonSetList [][]int
}

func (i indicator_S) fromRunes(r []rune) *indicator_S {
	rTrimmed := r[1 : len(r)-1]

	i.closingRunes = make([]rune, len(rTrimmed))
	copy(i.closingRunes, rTrimmed)

	acc := 0
	for x := range len(i.closingRunes) {
		if i.closingRunes[x] == '#' {
			power := int(math.Pow(2, float64(x)))
			acc += power
		}
	}

	i.closingValue = acc
	return &i
}

func (i *indicator_S) displayInitial() {
	fmt.Println(i.toStringClosing(), i.closingValue, i.buttonSetList, i.joltageSet)
}

func (i *indicator_S) toString() string {
	runes := []rune(strconv.FormatInt(int64(i.currentValue), 2))
	for x := range runes {
		if runes[x] == '0' {
			runes[x] = '.'
		} else {
			runes[x] = '#'
		}
	}

	return "[" + string(runes) + "]"
}

func (i *indicator_S) toStringClosing() string {
	return "[" + string(i.closingRunes) + "]"
}

func getJoltageParity(joltageSet []int) int {
	acc := 0
	for i, joltage := range joltageSet {
		if parity := joltage % 2; parity == 1 {
			power := int(math.Pow(2, float64(i)))
			acc += power
		}
	}

	return acc
}

func execC10(indicatorList []indicator_S) {
	var minimumPresses []int
	for _, indicator := range indicatorList {
		presses := calculateCombinations(indicator)
		minimumPresses = append(minimumPresses, presses)
	}

	//fmt.Println(minimumPresses)
	var acc int
	for _, presses := range minimumPresses {
		acc += presses
	}

	fmt.Println("Solution 1:", acc)
}

func calculateCombinations(indicator indicator_S) int {
	unfilteredIndexCombinations := startGeneratingIndexCombinations(len(indicator.buttonSetList))
	indexCombinations := sortIndexCombinations(unfilteredIndexCombinations)
	buttonPresses := 0

	for _, indexCombination := range indexCombinations {
		initialValue := 0
		for _, index := range indexCombination {
			buttonSet := indicator.buttonSetList[index]
			for _, button := range buttonSet {
				temp := initialValue
				initialValue = temp ^ int(math.Pow(2, float64(button)))
			}
		}
		if initialValue == indicator.closingValue {
			buttonPresses = len(indexCombination)
			break
		}
	}

	return buttonPresses
}

func startGeneratingIndexCombinations(buttonLen int) [][]int {
	set := make([]int, 0, buttonLen)

	for i := range buttonLen {
		set = append(set, i)
	}

	var indexCombinations [][]int
	generateCombinations(set, 0, []int{}, &indexCombinations)

	return indexCombinations
}

func generateCombinations(set []int, index int, current []int, combinations *[][]int) {
	if index == len(set) {
		*combinations = append(*combinations, append([]int(nil), current...))
		return
	}

	generateCombinations(set, index+1, current, combinations)
	generateCombinations(set, index+1, append(current, set[index]), combinations)
}

func sortIndexCombinations(indexCombinations [][]int) [][]int {
	var filteredIndexCombinations [][]int

	for _, combination := range indexCombinations {
		if len(combination) >= 1 {
			filteredIndexCombinations = append(filteredIndexCombinations, combination)
		}
	}

	slices.SortFunc(filteredIndexCombinations, func(a, b []int) int {
		cmp := len(a) - len(b)
		if cmp == 0 {
			for i := range len(a) {
				cmpI := a[i] - b[i]
				if cmpI != 0 {
					return cmpI
				}
			}
		}
		return cmp
	})

	return filteredIndexCombinations
}

/*
Used u/tenthmascot bifurcate solution

See https://www.reddit.com/r/adventofcode/comments/1pk87hl/2025_day_10_part_2_bifurcate_your_way_to_victory/
*/
func execC10Two(indicatorList []indicator_S) {
	acc := 0
	for _, indicator := range indicatorList {
		calculateJoltagePresses(&indicator)
		//fmt.Println(indicator.joltageSet, indicator.joltagePresses)
		acc += indicator.joltagePresses

	}

	fmt.Println("Solution 2:", acc)
}

type joltages_S struct {
	presses  int
	joltages []int
}

func calculateJoltagePresses(indicator *indicator_S) {
	parityMap := findAllParityCombinations(indicator)

	var solveCombinations func([]int) int
	solveCombinations = func(joltageSet []int) int {
		if zeroedJoltage(joltageSet...) {
			return 0
		}

		min := math.MaxInt32

		parityJoltages := parityMap[getJoltageParity(joltageSet)]
		for _, parityJoltage := range parityJoltages {
			if isLeq(parityJoltage.joltages, joltageSet) {
				newJoltageSet := subtractAndHalfJoltage(joltageSet, parityJoltage.joltages)
				min = slices.Min([]int{min, parityJoltage.presses + 2*solveCombinations(newJoltageSet)})
			}
		}

		return min
	}

	indicator.joltagePresses = solveCombinations(indicator.joltageSet)
}

// returns a map with key parity and value all joltage combinations with # presses to reach said parity
func findAllParityCombinations(indicator *indicator_S) map[int][]joltages_S {
	combinationList := sortIndexCombinations(startGeneratingIndexCombinations(len(indicator.buttonSetList)))
	parityMap := make(map[int][]joltages_S, 10)

	for _, combination := range combinationList {
		joltages := make([]int, len(indicator.joltageSet))

		for _, index := range combination {
			// select one button from the buttonSetList
			button := indicator.buttonSetList[index]
			// iterate over all button values within the set
			for _, buttonValue := range button {
				joltages[buttonValue] += 1
			}
		}

		parity := getJoltageParity(joltages)
		appendToParityMap(parityMap, parity, joltages_S{presses: len(combination), joltages: joltages})
	}

	appendToParityMap(parityMap, 0, joltages_S{presses: 0, joltages: slices.Repeat([]int{0}, len(indicator.joltageSet))})

	return parityMap
}

func appendToParityMap(parityMap map[int][]joltages_S, parity int, parityJoltage joltages_S) {
	if _, ok := parityMap[parity]; ok {
		for i, entry := range parityMap[parity] {
			if slices.Equal(entry.joltages, parityJoltage.joltages) {
				if entry.presses > parityJoltage.presses {
					parityMap[parity][i] = parityJoltage
				}
				return
			}
		}
		parityMap[parity] = append(parityMap[parity], parityJoltage)
	} else {
		parityMap[parity] = append([]joltages_S{}, parityJoltage)
	}
}

func zeroedJoltage(joltage ...int) bool {
	for _, j := range joltage {
		if j != 0 {
			return false
		}
		if j < 0 {
			panic("Error invalid joltage found")
		}
	}

	return true
}

func isLeq(a []int, b []int) bool {
	for i := range len(a) {
		if !(a[i] <= b[i]) {
			return false
		}
	}

	return true
}

func subtractAndHalfJoltage(a, b []int) []int {
	result := make([]int, len(a))
	for i := range len(a) {
		result[i] = (a[i] - b[i]) / 2
	}

	return result
}
