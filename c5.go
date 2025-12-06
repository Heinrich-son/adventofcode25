package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// ~35min parsing
// ~15min part one
// >60min part two

func runC5() {
	_ = func(list []idRange) {
		for _, ele := range list {
			display := strings.Join([]string{"[", strconv.Itoa(ele.opening), "-", strconv.Itoa(ele.closing), "]"}, "")
			fmt.Println(display)
		}
	}
	_ = func(list []int) {
		fmt.Println(list)
	}

	tuple := readFile("c5.txt", mapRangeAndIdList)
	//printIdRangeList(tuple.idRangeList)
	//printIdList(tuple.idList)

	execC5(tuple.idRangeList, tuple.idList)
	execC5Two(tuple.idRangeList)
}

func mapRangeAndIdList(content []byte) struct {
	idRangeList []idRange
	idList      []int
} {
	NEWLINE := "\n"
	EMPTY := ""
	DASH := "-"

	idRangeList := make([]idRange, 0, 100)
	idList := make([]int, 0, 200)

	for s := range strings.SplitSeq(string(content), NEWLINE) {
		if strings.Contains(s, DASH) {
			idRangeList = append(idRangeList, mapRange(s))
		} else if s == EMPTY {
			continue
		} else {
			idList = append(idList, mapId(s))
		}
	}

	return struct {
		idRangeList []idRange
		idList      []int
	}{
		idRangeList: idRangeList,
		idList:      idList,
	}
}

func mapRange(rangeString string) idRange {
	DASH := "-"
	rangeArray := strings.Split(rangeString, DASH)
	opening := mapId(rangeArray[0])
	closing := mapId(rangeArray[1])

	return idRange{opening, closing}
}

func mapId(idString string) int {
	id, err := strconv.Atoi(idString)
	if err != nil {
		panic(err)
	}
	return id
}

type idRange struct {
	opening int
	closing int
}

func execC5(idRangeList []idRange, idList []int) {
	fresh := make([]int, 0, 100)
	for _, id := range idList {
		for _, iRange := range idRangeList {
			if containsId(iRange, id) {
				fresh = append(fresh, id)
				break
			}
		}
	}

	fmt.Println("Solution 1:", len(fresh))
}

func containsId(iRange idRange, id int) bool {
	if id >= iRange.opening && id <= iRange.closing {
		return true
	} else {
		return false
	}
}

func execC5Two(idRangeList []idRange) {
	groupedIdRangeList := make([][]idRange, 0, 100)

	sortedIdRangeList := slices.SortedFunc(slices.Values(idRangeList), func(a, b idRange) int {
		return a.opening - b.opening
	})

	groupOpeningIndex := 0
	groupClosingValue := sortedIdRangeList[groupOpeningIndex].closing
	currentGroup := []idRange{sortedIdRangeList[groupOpeningIndex]}
	for i := 1; i < len(sortedIdRangeList); i++ {
		if groupClosingValue >= sortedIdRangeList[i].opening {
			if groupClosingValue < sortedIdRangeList[i].closing {
				groupClosingValue = sortedIdRangeList[i].closing
			}
			currentGroup = append(currentGroup, sortedIdRangeList[i])
		} else {
			groupedIdRangeList = append(groupedIdRangeList, currentGroup)

			groupOpeningIndex = i
			groupClosingValue = sortedIdRangeList[i].closing
			currentGroup = []idRange{sortedIdRangeList[i]}
		}

		if i == len(sortedIdRangeList)-1 {
			groupedIdRangeList = append(groupedIdRangeList, currentGroup)

		}
	}

	mergedIdRangeList := make([]idRange, 0, 50)
	for _, group := range groupedIdRangeList {
		merged := mergeIdRanges(group)
		mergedIdRangeList = append(mergedIdRangeList, merged)
	}

	idsAmount := sumIdRanges(mergedIdRangeList)
	fmt.Println("Solution 2:", idsAmount)
}

func mergeIdRanges(idRangeList []idRange) idRange {
	opening := idRangeList[0].opening
	closing := slices.MaxFunc(idRangeList, func(a, b idRange) int {
		return a.closing - b.closing
	}).closing

	return idRange{opening, closing}
}

func sumIdRanges(idRangeList []idRange) int {
	acc := 0

	for _, iRange := range idRangeList {
		acc += (iRange.closing - iRange.opening) + 1
	}

	return acc
}
