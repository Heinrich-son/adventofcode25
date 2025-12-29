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
button set - (3) where each index will toggle the corresponding indicator
joltage requirements - {3} unknown

Problem: Each indicator light is initially off. Determine the fewest possible
total presses to reach the wanted indicator state
*/

func runC10() {
	indicatorList := readFile("c10_test.txt", mapToIndicatorList)
	/*for _, indicator := range indicatorList {
		indicator.displayInitial()
	}*/

	//execC10(indicatorList)
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

	lines := strings.Split(string(content), NEWLINE)
	for _, machine := range lines {
		sections := strings.Split(machine, SPACE)
		indicatorRunes := []rune(sections[0])
		buttonSlice := sections[1 : len(sections)-1]
		joltageString := sections[len(sections)-1]

		buttonSetList := parseButton(buttonSlice)
		joltageSet := parseJoltage(joltageString)

		indicator := indicator_S{joltageSet: joltageSet, remainingJoltage: joltageSet, buttonSetList: buttonSetList, joltageCache: make(map[int][]int, 10)}.fromRunes(indicatorRunes)
		indicatorList = append(indicatorList, *indicator)
	}

	return
}

type indicator_S struct {
	closingValue int
	currentValue int

	joltagePresses   int
	remainingJoltage []int
	joltageCache     map[int][]int

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
	fmt.Println(i.toStringClosing(), i.closingValue, i.buttonSetList, i.joltageSet, i.getJoltageParity())
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

func (i *indicator_S) getJoltageParity() int {
	return getJoltageParity(i.remainingJoltage)
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

func execC10Two(indicatorList []indicator_S) {
	acc := 0
	for i, indicator := range indicatorList {
		if i == 2 {
			calculateJoltagePresses(&indicator)
			fmt.Println(indicator.joltagePresses)
			acc += indicator.joltagePresses
		}
	}

	fmt.Println("Solution 2:", acc)
}

func generateCombinationsRepeat(n, k int) [][]int {
	var result [][]int
	var comb []int

	var dfs func(start int)
	dfs = func(start int) {
		if len(comb) > 0 {
			result = append(result, append([]int{}, comb...))
		}

		if len(comb) == k {
			return
		}

		for i := start; i < n; i++ {
			comb = append(comb, i)
			dfs(i)
			comb = comb[:len(comb)-1]
		}
	}

	dfs(0)
	return result
}

func calculateJoltagePresses(indicator *indicator_S) {
	if zeroedJoltage(indicator.remainingJoltage...) {
		return
	}

	calculatedValue, combination := findCombination(indicator)
	fmt.Print(indicator.remainingJoltage)

	if calculatedValue != math.MaxInt && combination != nil {
		indicator.joltageCache[calculatedValue] = combination
		for _, index := range combination {
			// select the buttonSet via index
			buttonSet := indicator.buttonSetList[index]
			fmt.Print(" (", buttonSet, ")")
			// iterate over all button values within the set
			for _, buttonValue := range buttonSet {
				indicator.remainingJoltage[buttonValue] = indicator.remainingJoltage[buttonValue] - 1
			}
		}
		fmt.Print(" => ", indicator.remainingJoltage)
		halfJoltage(indicator.remainingJoltage)
		fmt.Print(" => ", indicator.remainingJoltage, "\n")
		calculateJoltagePresses(indicator)
		indicator.joltagePresses = len(combination) + 2*indicator.joltagePresses
	} else {
		indicator.joltagePresses += bruteForceJoltagePresses(indicator)
	}
}

func findCombination(indicator *indicator_S) (int, []int) {
	if cachedCombination, ok := indicator.joltageCache[indicator.getJoltageParity()]; ok {
		return indicator.getJoltageParity(), cachedCombination
	}

	combinationList := sortIndexCombinations(generateCombinationsRepeat(len(indicator.buttonSetList), slices.Max(indicator.joltageSet)))

	// iterate over all index combinations
	for _, combination := range combinationList {
		calculatedValue := 0

		for _, index := range combination {
			// select the buttonSet via index
			buttonSet := indicator.buttonSetList[index]
			// iterate over all button values within the set
			for _, buttonValue := range buttonSet {
				temp := calculatedValue
				calculatedValue = temp ^ int(math.Pow(2, float64(buttonValue)))
			}
		}

		if calculatedValue == indicator.getJoltageParity() {
			return calculatedValue, combination
		}

	}

	return math.MaxInt, nil
}

func halfJoltage(joltageSet []int) []int {
	for i := range joltageSet {
		joltageSet[i] = joltageSet[i] / 2
	}

	return joltageSet
}

func bruteForceJoltagePresses(indicator *indicator_S) int {
	maxPresses := slices.Max(indicator.remainingJoltage) * 2
	combinationList := sortIndexCombinations(generateCombinationsRepeat(len(indicator.buttonSetList), maxPresses))

	for _, combination := range combinationList {
		remainingJoltage := make([]int, len(indicator.remainingJoltage))
		copy(remainingJoltage, indicator.remainingJoltage)

		calculatedValue := 0
		for _, index := range combination {
			buttonSet := indicator.buttonSetList[index]
			for _, buttonValue := range buttonSet {
				temp := calculatedValue
				calculatedValue = temp ^ int(math.Pow(2, float64(buttonValue)))

				remainingJoltage[buttonValue] = remainingJoltage[buttonValue] - 1
			}
		}

		if zeroedJoltage(remainingJoltage...) {
			indicator.joltageCache[calculatedValue] = combination
			return len(combination)
		}
	}

	// should not happen
	return math.MaxInt
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
