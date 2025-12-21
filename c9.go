package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func runC9() {
	coordinateList := readFile("c9.txt", mapToTileCoordinate)
	areas := calculateAllTileAreas(coordinateList)
	fmt.Println(areas[0])
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

	for _, line := range drawingBoard {
		fmt.Println(string(line))
	}
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

// I don't know any proper algorithm to make it "fast", so drawning it is
// Alright let's leave drawing, I literally don't have enough RAM

func execC9Two(coordinateList []coordinate) {
	yMap, xMap := createXYMaps(coordinateList)
	horizontalRanges := createHorizontalTileRanges(yMap)
	verticalRanges := createVerticalTileRanges(xMap)

	//draw(coordinateList, append(horizontalRanges, verticalRanges...))

	areas := calculateAllValidTileAreas(coordinateList, horizontalRanges, verticalRanges)
	fmt.Println(areas[0:5])
}

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

// calculate horizontal intervalls
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

// calculate vertical intervalls
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

// calculate coordinates of two other corners
func calculateAllValidTileAreas(coordinateList []coordinate, horizontalRanges []coordinate, verticalRanges []coordinate) []int {
	areas := make([]int, 0, len(coordinateList))
	coordinateBorders := slices.Concat(coordinateList, horizontalRanges, verticalRanges)

	fmt.Println(len(coordinateList))
	for i := range coordinateList {
		if i%50 == 0 {
			fmt.Println("at ", i)
		}
		for j := i + 1; j < len(coordinateList); j++ {

			x1, y1 := coordinateList[i].x, coordinateList[i].y
			x2, y2 := coordinateList[j].x, coordinateList[j].y

			mCoordinate := coordinate{x1, y2}
			nCoordinate := coordinate{x2, y1}

			if verifyOtherCorners(mCoordinate, nCoordinate, coordinateBorders...) {
				areas = append(areas, calculateTileArea(x1, y1, x2, y2))
			}
		}
	}

	slices.SortFunc(areas, func(a, b int) int {
		return b - a
	})

	return areas
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

func verifyOtherCorners(m coordinate, n coordinate, coordinateBorders ...coordinate) bool {
	if m.x == n.x || m.y == n.y {
		return true
	}

	mValid := false
	nValid := false

	mValidationMap := validCoordinateMap{}
	nValidationMap := validCoordinateMap{}

	for _, coordinate := range coordinateBorders {
		updateBorders(&mValidationMap, m, coordinate)
		updateBorders(&nValidationMap, n, coordinate)

		if m.equal(coordinate) {
			mValid = true
		}
		if n.equal(coordinate) {
			nValid = true
		}

		if mValid && nValid {
			break
		}
	}

	if (mValid || mValidationMap.isWithinBorder()) && (nValid || nValidationMap.isWithinBorder()) {
		return true
	}

	return false
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
