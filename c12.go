package main

import (
	"fmt"
	"strconv"
	"strings"
)

func runC12() {
	presentAndRegion := readFile("c12.txt", parseToPresents)
	regionList := presentAndRegion.regionList

	filterByTotalArea(regionList)

	acc := 0
	for _, region := range regionList {
		if region.doesFit == YES {
			acc++
		}
	}
	fmt.Println("Solution 1:", acc)
}

func parseToPresents(content []byte) struct {
	presentMap map[int]present_S
	regionList []*region_S
} {
	presentMap := make(map[int]present_S)
	regionList := make([]*region_S, 0)

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.Contains(line, "x") {
			regionList = append(regionList, parseRegion(line))
			continue
		}

		if strings.Contains(line, ":") {
			index, present := parsePresent(lines[i : i+4])
			presentMap[index] = present
			continue
		}
	}

	return struct {
		presentMap map[int]present_S
		regionList []*region_S
	}{
		presentMap,
		regionList,
	}
}

func parseRegion(line string) *region_S {
	region := region_S{}

	lineSplit := strings.Split(line, " ")
	gridSizeSlice := strings.Split(lineSplit[0][0:len(lineSplit[0])-1], "x")

	region.width = atoi(gridSizeSlice[0])
	region.length = atoi(gridSizeSlice[1])
	region.presentIndexes = make([]int, 0)

	indexesSlice := lineSplit[1:]
	for _, indexString := range indexesSlice {
		region.presentIndexes = append(region.presentIndexes, atoi(indexString))
	}

	return &region
}

func parsePresent(lines []string) (int, present_S) {
	index := atoi(lines[0][0:1])

	size := 0
	shape := make([][]rune, 0, 3)
	for i := 1; i < len(lines); i++ {
		size += strings.Count(lines[i], "#")
		shape = append(shape, []rune(lines[i]))
	}

	return index, present_S{size, shape}
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		panic("Invalid string.")
	}

	return i
}

type present_S struct {
	size  int
	shape [][]rune
}

type region_S struct {
	width          int
	length         int
	doesFit        fits
	presentIndexes []int
}

type fits int8

const (
	UNKNOWN fits = 0
	YES     fits = 1
	NO      fits = 2
)

func (r region_S) area() int {
	return r.width * r.length
}

func (r region_S) requiredArea() (required int) {
	for _, amount := range r.presentIndexes {
		required = required + amount*9
	}

	return required
}

func filterByTotalArea(regionList []*region_S) {
	for _, region := range regionList {
		if region.area() < region.requiredArea() {
			region.doesFit = NO
		} else {
			region.doesFit = YES
		}
	}
}
