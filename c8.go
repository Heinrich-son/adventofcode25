package main

import (
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
)

func runC8() {
	coordinateList := readFile("c8.txt", mapToCoordinates)
	pairList := calculateAllDistances(coordinateList)
	//printPairList(pairList)
	collectToCircuit(pairList, 1000)
	collectToCircuitTwo(pairList, len(coordinateList))
}

func printPairList(pairList []pairs) {
	for _, pair := range pairList {
		fmt.Println(pair)
	}
}

func mapToCoordinates(content []byte) []coordinates {
	NEWLINE := "\n"
	COMMA := ","

	lines := strings.Split(string(content), NEWLINE)
	coordinateList := make([]coordinates, 0, len(lines))
	for _, line := range lines {
		coordinate := toCoordinate(strings.Split(line, COMMA))
		coordinateList = append(coordinateList, coordinate)
	}

	return coordinateList
}

func toCoordinate(line []string) coordinates {
	if len(line) != 3 {
		panic("Line has to have 3 elements only")
	}

	x, _ := strconv.Atoi(line[0])
	y, _ := strconv.Atoi(line[1])
	z, _ := strconv.Atoi(line[2])

	return coordinates{x, y, z}
}

type coordinates struct {
	x int
	y int
	z int
}

func (c coordinates) hash() int {
	s := strconv.Itoa(c.x) + strconv.Itoa(c.y) + strconv.Itoa(c.z)
	h, _ := strconv.Atoi(s)
	return h
}

type pairs struct {
	a        coordinates
	b        coordinates
	distance float64
}

func calculateAllDistances(coordinateList []coordinates) []pairs {
	pairList := make([]pairs, 0, len(coordinateList)/2)
	for i := range coordinateList {
		for j := i + 1; j < len(coordinateList); j++ {
			a := coordinateList[i]
			b := coordinateList[j]
			distance := calculateDistance(a, b)
			pair := pairs{a, b, distance}
			pairList = append(pairList, pair)
		}
	}

	return slices.SortedFunc(slices.Values(pairList), func(a, b pairs) int {
		distance := a.distance - b.distance
		return int(distance)
	})
}

func calculateDistance(a, b coordinates) float64 {
	d1 := math.Pow(float64(a.x-b.x), 2)
	d2 := math.Pow(float64(a.y-b.y), 2)
	d3 := math.Pow(float64(a.z-b.z), 2)

	d := math.Sqrt(d1 + d2 + d3)
	return d
}

func collectToCircuit(pairList []pairs, amountOfPairs int) {
	circuitMap := make(map[int][]coordinates, len(pairList))

	if len(pairList) < amountOfPairs {
		panic("Intended amount of pairs exceeds actual amount of pairs")
	}

	for i := range amountOfPairs {
		pair := pairList[i]
		aSlice, aInMap := circuitMap[pair.a.hash()]
		bSlice, bInMap := circuitMap[pair.b.hash()]

		if aInMap && bInMap {
			if !slices.Equal(aSlice, bSlice) {
				mergeEntries(circuitMap, pair.a, pair.b)
			}
		}
		if aInMap && !bInMap {
			updateEntry(circuitMap, aSlice, pair.b)
		}
		if !aInMap && bInMap {
			updateEntry(circuitMap, bSlice, pair.a)
		}
		if !aInMap && !bInMap {
			newSlice := append(make([]coordinates, 0, 10), pair.a, pair.b)
			circuitMap[pair.a.hash()] = newSlice
			circuitMap[pair.b.hash()] = newSlice

		}
	}

	circuitSlice := uniqueSet(circuitMap)

	largestThree := 1
	for i := range 3 {
		largestThree *= len(circuitSlice[i])
	}

	fmt.Println("Solution 1:", largestThree)
}

func mergeEntries(circuitMap map[int][]coordinates, a coordinates, b coordinates) {
	aSlice := circuitMap[a.hash()]
	bSlice := circuitMap[b.hash()]
	merged := append(aSlice, bSlice...)
	for _, c := range merged {
		circuitMap[c.hash()] = merged
	}
}

func updateEntry(circuitMap map[int][]coordinates, currentSlice []coordinates, c coordinates) {
	newSlice := append(currentSlice, c)
	for _, coordinate := range newSlice {
		circuitMap[coordinate.hash()] = newSlice
	}
}

func uniqueSet(circuitMap map[int][]coordinates) [][]coordinates {
	circuitSet := make(map[int][]coordinates)
	for _, coordinateList := range circuitMap {
		acc := 0
		for _, coordinate := range coordinateList {
			acc += coordinate.hash()
		}
		circuitSet[acc] = coordinateList
	}

	circuitSlice := slices.SortedFunc(maps.Values(circuitSet), func(a, b []coordinates) int {
		return len(b) - len(a)
	})

	return circuitSlice
}

func collectToCircuitTwo(pairList []pairs, amountOfBoxes int) {
	circuitMap := make(map[int][]coordinates, len(pairList))

	for _, pair := range pairList {
		aSlice, aInMap := circuitMap[pair.a.hash()]
		bSlice, bInMap := circuitMap[pair.b.hash()]

		if aInMap && bInMap {
			if !slices.Equal(aSlice, bSlice) {
				mergeEntries(circuitMap, pair.a, pair.b)
				if checkIsSingleCircuit(uniqueSet(circuitMap), amountOfBoxes, pair) {
					break
				}
			}
		}
		if aInMap && !bInMap {
			updateEntry(circuitMap, aSlice, pair.b)
			if checkIsSingleCircuit(uniqueSet(circuitMap), amountOfBoxes, pair) {
				break
			}
		}
		if !aInMap && bInMap {
			updateEntry(circuitMap, bSlice, pair.a)
			if checkIsSingleCircuit(uniqueSet(circuitMap), amountOfBoxes, pair) {
				break
			}
		}
		if !aInMap && !bInMap {
			newSlice := append(make([]coordinates, 0, 10), pair.a, pair.b)
			circuitMap[pair.a.hash()] = newSlice
			circuitMap[pair.b.hash()] = newSlice
			if checkIsSingleCircuit(uniqueSet(circuitMap), amountOfBoxes, pair) {
				break
			}
		}
	}

}

func checkIsSingleCircuit(coordinate [][]coordinates, amountOfBoxes int, pair pairs) bool {
	if len(coordinate) == 1 {
		if len(coordinate[0]) == amountOfBoxes {
			doubleX := pair.a.x * pair.b.x
			fmt.Println("Solution 2:", doubleX)
			return true
		}
	}

	return false
}
