package main

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// damn... what a mess

func runC9() {
	coordinateList := readFile("c9.txt", mapToTileCoordinate)
	areas := calculateAllTileAreas(coordinateList)
	fmt.Println("Solution 1: ", areas[0])
	execC9Two(coordinateList)
}

func mapToTileCoordinate(content []byte) []coordinate {
	NEWLINE := "\n"
	COMMA := ","

	lines := strings.Split(string(content), NEWLINE)
	coordinateList := make([]coordinate, 0, len(lines)*2)
	for _, line := range lines {
		numbers := strings.Split(line, COMMA)
		if len(numbers) != 2 {
			panic("Unexpected amount of numbers, expected 2")
		}
		x, _ := strconv.Atoi(numbers[0])
		y, _ := strconv.Atoi(numbers[1])
		coordinateList = append(coordinateList, coordinate{x, y})
	}

	return coordinateList
}

type coordinate struct {
	x int
	y int
}

type coordinatePair struct {
	a coordinate
	b coordinate
}

func (c coordinate) equal(o coordinate) bool {
	return c.x == o.x && c.y == o.y
}

func calculateAllTileAreas(coordinateList []coordinate) []int {
	areas := make([]int, 0, len(coordinateList)*len(coordinateList))
	for i := range coordinateList {
		for j := i + 1; j < len(coordinateList); j++ {
			x1, y1 := coordinateList[i].x, coordinateList[i].y
			x2, y2 := coordinateList[j].x, coordinateList[j].y
			areas = append(areas, calculateTileArea(x1, y1, x2, y2))
		}
	}

	slices.SortFunc(areas, func(a, b int) int {
		return b - a
	})

	return areas
}

func calculateTileArea(x1, y1, x2, y2 int) int {
	var xSmall, xBig, ySmall, yBig int
	if x1 < x2 {
		xSmall = x1
		xBig = x2
	} else {
		xSmall = x2
		xBig = x1
	}

	if y1 < y2 {
		ySmall = y1
		yBig = y2
	} else {
		ySmall = y2
		yBig = y1
	}

	area := ((xBig - xSmall) + 1) * ((yBig - ySmall) + 1)
	return area
}

// only use for small board
func draw(coordinateList []coordinate, rangeList []coordinate) {
	maxY := maxYTile(coordinateList) + 1
	maxX := maxXTile(coordinateList) + 1
	drawingBoard := make([][]rune, maxY)
	for n := range drawingBoard {
		drawingBoard[n] = make([]rune, maxX)
	}

	for i := range maxY {
		for j := range maxX {
			drawingBoard[i][j] = '.'
		}
	}

	for _, coor := range coordinateList {
		drawingBoard[coor.y][coor.x] = '#'
	}

	for _, r := range rangeList {
		drawingBoard[r.y][r.x] = 'X'
	}

	/*for _, line := range drawingBoard {
		fmt.Println(string(line))
	}*/

	writeTileDrawning(drawingBoard)
}

func maxXTile(coordinateList []coordinate) int {
	maxX := slices.MaxFunc(coordinateList, func(a, b coordinate) int {
		return a.x - b.x
	})

	return maxX.x
}

func maxYTile(coordinateList []coordinate) int {
	maxX := slices.MaxFunc(coordinateList, func(a, b coordinate) int {
		return a.y - b.y
	})

	return maxX.y
}

// so this became an utter mess

func execC9Two(coordinateList []coordinate) {
	yCompressed, xCompressed := compressCoordinates(coordinateList)
	compressedList := mapCoordinates(coordinateList, yCompressed, xCompressed)

	yMap, xMap := createXYMaps(compressedList)
	horizontalRanges := createHorizontalTileRanges(yMap)
	verticalRanges := createVerticalTileRanges(xMap)

	draw(compressedList, append(horizontalRanges, verticalRanges...))

	pairs := calculateAllValidTileAreasV2(compressedList, horizontalRanges, verticalRanges)
	areas := calculateFromAllPairs(pairs, yCompressed, xCompressed)
	fmt.Println(areas[0])
}

// creates two maps with y and x each as key with all coordinates falling under each values (as slice)
// all entries are sorted in ascending order
func createXYMaps(coordinateList []coordinate) (map[int][]coordinate, map[int][]coordinate) {
	yMap := make(map[int][]coordinate, 100)
	xMap := make(map[int][]coordinate, 100)

	for _, coor := range coordinateList {
		_, yOk := yMap[coor.y]
		_, xOk := xMap[coor.x]

		if !yOk {
			yMap[coor.y] = make([]coordinate, 0, 2)
		}
		if !xOk {
			xMap[coor.x] = make([]coordinate, 0, 2)
		}

		yMap[coor.y] = append(yMap[coor.y], coor)
		xMap[coor.x] = append(xMap[coor.x], coor)
	}

	for yKey := range yMap {
		slices.SortFunc(yMap[yKey], func(a, b coordinate) int {
			return a.x - b.x
		})
	}

	for xKey := range xMap {
		slices.SortFunc(xMap[xKey], func(a, b coordinate) int {
			return a.y - b.y
		})
	}

	return yMap, xMap
}

// calculates all horizontal edges as a list of coordinates excluding the corners
func createHorizontalTileRanges(coordinateMap map[int][]coordinate) []coordinate {
	yRanges := make([]coordinate, 0, 500)

	for key := range coordinateMap {
		entry := coordinateMap[key]
		start := entry[0]
		for i := 1; i < len(entry); i++ {
			end := entry[i]
			for n := start.x + 1; n < end.x; n++ {
				yRanges = append(yRanges, coordinate{n, key})
			}
			start = end
		}
	}

	return yRanges
}

// calculate all vertical edges as a list of coordinates excluding the corners
func createVerticalTileRanges(coordinateMap map[int][]coordinate) []coordinate {
	xRanges := make([]coordinate, 0, 500)

	for key := range coordinateMap {
		entry := coordinateMap[key]
		start := entry[0]
		for i := 1; i < len(entry); i++ {
			end := entry[i]
			for n := start.y + 1; n < end.y; n++ {
				xRanges = append(xRanges, coordinate{key, n})
			}
			start = end
		}
	}

	return xRanges
}

type validCoordinateMap struct {
	hasTopBorder    bool
	hasBottomBorder bool
	hasLeftBorder   bool
	hasRightBorder  bool
}

func (v validCoordinateMap) isWithinBorder() bool {
	return v.hasTopBorder && v.hasBottomBorder && v.hasLeftBorder && v.hasRightBorder
}

func updateBorders(validationMap *validCoordinateMap, point coordinate, border coordinate) {
	// left border exists (border - point)
	if point.y == border.y && border.x < point.x {
		validationMap.hasLeftBorder = true
	}

	// right border exists (point - border)
	if point.y == border.y && point.x < border.x {
		validationMap.hasRightBorder = true
	}

	// top border exists (border | point)
	if point.x == border.x && border.y < point.y {
		validationMap.hasTopBorder = true
	}

	// bottom border exists (point | border)
	if point.x == border.x && point.y < border.y {
		validationMap.hasBottomBorder = true
	}
}

func compressCoordinates(coordinateList []coordinate) (map[int]int, map[int]int) {
	yMap := make(map[int]int, len(coordinateList))
	xMap := make(map[int]int, len(coordinateList))

	slices.SortFunc(coordinateList, func(a, b coordinate) int {
		return a.y - b.y
	})

	var i, j int

	for _, c := range coordinateList {
		if _, ok := yMap[c.y]; !ok {
			yMap[c.y] = i
			i++
		}
	}

	slices.SortFunc(coordinateList, func(a, b coordinate) int {
		return a.x - b.x
	})

	for _, c := range coordinateList {
		if _, ok := xMap[c.x]; !ok {
			xMap[c.x] = j
			j++
		}
	}

	fmt.Printf("Compressed to y:%d, x:%d\n", i, j)
	return yMap, xMap
}

func mapCoordinates(coordinateList []coordinate, yMap, xMap map[int]int) []coordinate {
	compressedList := make([]coordinate, 0, len(coordinateList))

	for _, c := range coordinateList {
		compressed := coordinate{xMap[c.x], yMap[c.y]}
		compressedList = append(compressedList, compressed)
	}

	return compressedList
}

func unmapCoordinate(coor coordinate, yMap, xMap map[int]int) coordinate {
	var y, x int

	for keyY, valY := range yMap {
		if coor.y == valY {
			y = keyY
			break
		}
	}

	for keyX, valX := range xMap {
		if coor.x == valX {
			x = keyX
		}
	}

	return coordinate{x, y}
}

func createSquareList(coor1 coordinate, coor2 coordinate) []coordinate {
	var xSmall, xBig, ySmall, yBig int
	if coor1.x < coor2.x {
		xSmall = coor1.x
		xBig = coor2.x
	} else {
		xSmall = coor2.x
		xBig = coor1.x
	}

	if coor1.y < coor2.y {
		ySmall = coor1.y
		yBig = coor2.y
	} else {
		ySmall = coor2.y
		yBig = coor1.y
	}

	squareMap := make(map[coordinate]bool, (xBig-xSmall)*2+(yBig-ySmall)*2)
	squareMap[coor1] = true
	squareMap[coor2] = true

	for i := ySmall; i <= yBig; i++ {
		squareMap[coordinate{xSmall, i}] = true
		squareMap[coordinate{xBig, i}] = true
	}

	for j := xSmall; j <= xBig; j++ {
		squareMap[coordinate{j, ySmall}] = true
		squareMap[coordinate{j, yBig}] = true
	}

	squareIter := maps.Keys(squareMap)

	squareList := slices.SortedFunc(squareIter, func(a, b coordinate) int {
		i := a.x - b.x
		if i == 0 {
			return a.y - b.y
		} else {
			return i
		}
	})

	return squareList
}

func calculateAllValidTileAreasV2(coordinateList []coordinate, horizontalRanges []coordinate, verticalRanges []coordinate) []coordinatePair {
	pairs := make([]coordinatePair, 0, len(coordinateList)/2)
	coordinateBorders := slices.Concat(coordinateList, horizontalRanges, verticalRanges)

	for i := range coordinateList {
		if i%50 == 0 {
			fmt.Println("Loop at", i, "from", len(coordinateList))
		}
		for j := i + 1; j < len(coordinateList); j++ {
			squareList := createSquareList(coordinateList[i], coordinateList[j])

			allEdgesValid := true
			for _, s := range squareList {
				sIsValid := verifyOtherCornersV2(s, coordinateBorders...)
				if !sIsValid {
					allEdgesValid = false
					break
				}
			}

			if allEdgesValid {
				pairs = append(pairs, coordinatePair{coordinateList[i], coordinateList[j]})
			}
		}
	}

	return pairs
}

func verifyOtherCornersV2(m coordinate, coordinateBorders ...coordinate) bool {
	mValid := false

	mValidationMap := validCoordinateMap{}

	for _, coordinate := range coordinateBorders {
		updateBorders(&mValidationMap, m, coordinate)

		if m.equal(coordinate) {
			mValid = true
		}

		if mValid {
			break
		}
	}

	if mValid || mValidationMap.isWithinBorder() {
		return true
	}

	return false
}

func calculateFromAllPairs(coodinatePairList []coordinatePair, yMap, xMap map[int]int) []int {
	areas := make([]int, 0, len(coodinatePairList))

	for _, pair := range coodinatePairList {
		rA := unmapCoordinate(pair.a, yMap, xMap)
		rB := unmapCoordinate(pair.b, yMap, xMap)

		area := calculateTileArea(rA.x, rA.y, rB.x, rB.y)
		areas = append(areas, area)
	}

	slices.SortFunc(areas, func(a, b int) int {
		return b - a
	})

	return areas
}
