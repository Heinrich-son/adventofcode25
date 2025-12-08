package main

import (
	"fmt"
	"strings"
)

// ~90min parsing
// ~15min part one
// >60min part two

func runC7() {
	diagramRune := readFile("c7.txt", mapToDiagramRune)
	diagram := mapToDiagram(diagramRune)
	//printTachyonDiagram(diagram)
	connectTachyonDiagram(diagram)
	//printTachyonDiagramWithConnection(diagram)
	startY, startX := findStart(diagram)
	diagram[startY][startX].shoot()
	//printTachyonDiagram(diagram)
	countSplitting(diagram)
	walkBackInTime(diagram)
}

func mapToDiagramRune(content []byte) [][]rune {
	NEWLINE := "\n"

	contentString := strings.Split(string(content), NEWLINE)
	contentRune := make([][]rune, 0, len(contentString))
	for _, line := range contentString {
		contentRune = append(contentRune, []rune(line))
	}

	return contentRune
}

func mapToDiagram(content [][]rune) [][]tachyonField {
	tachyonDiagram := make([][]tachyonField, 0, len(content))

	for _, line := range content {
		tachyonLine := make([]tachyonField, 0, len(content[0]))
		for _, s := range line {
			start := false
			splitter := false

			if s == 'S' {
				start = true
			}
			if s == '^' {
				splitter = true
			}

			tachyon := tachyonField{start: start, splitter: splitter}
			tachyonLine = append(tachyonLine, tachyon)
		}
		tachyonDiagram = append(tachyonDiagram, tachyonLine)
	}

	return tachyonDiagram
}

func printTachyonDiagram(diagram [][]tachyonField) {
	displayDiagram := make([]string, 0, len(diagram))
	for _, line := range diagram {
		displayLine := ""
		for _, tachyon := range line {
			displayLine = displayLine + tachyon.display()
		}
		displayDiagram = append(displayDiagram, displayLine)
	}

	displayString := strings.Join(displayDiagram, "\n")
	fmt.Println(displayString)
}

func printTachyonDiagramWithConnection(diagram [][]tachyonField) {
	displayDiagram := make([]string, 0, len(diagram))
	for _, line := range diagram {
		displayTopLine := ""
		displayBottomLine := ""
		for _, tachyon := range line {
			cube := tachyon.displayWithConnction()
			displayTopLine += cube[0]
			displayBottomLine += cube[1]
		}
		displayDiagram = append(displayDiagram, displayTopLine, displayBottomLine)
	}

	displayString := strings.Join(displayDiagram, "\n")
	fmt.Println(displayString)
}

func connectTachyonDiagram(diagram [][]tachyonField) [][]tachyonField {
	for y := range diagram {
		line := diagram[y]
		for x := range line {
			neighbors := retrieveNeighbors(&diagram, y, x)
			diagram[y][x].top = neighbors.top
			diagram[y][x].right = neighbors.right
			diagram[y][x].bottom = neighbors.bottom
			diagram[y][x].left = neighbors.left
		}
	}

	return diagram
}

func retrieveNeighbors(diagramPointer *[][]tachyonField, ySource int, xSource int) (neighbors struct {
	top    *tachyonField
	left   *tachyonField
	bottom *tachyonField
	right  *tachyonField
}) {
	diagram := *diagramPointer
	TRBL := []struct{ y, x int }{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	isInBound := func(y, x int) bool {
		maxY := len(diagram)
		maxX := len(diagram[0])

		if y < 0 || y >= maxY {
			return false
		}
		if x < 0 || x >= maxX {
			return false
		}
		return true
	}

	for i, direction := range TRBL {
		y := ySource + direction.y
		x := xSource + direction.x
		if isInBound(y, x) {
			switch i {
			case 0:
				neighbors.top = &diagram[y][x]
			case 1:
				neighbors.right = &diagram[y][x]
			case 2:
				neighbors.bottom = &diagram[y][x]
			case 3:
				neighbors.left = &diagram[y][x]
			default:
				panic("No corresponding direction found!")
			}
		}
	}

	return
}

type tachyonField struct {
	start    bool
	splitter bool
	beam     bool

	isSplitting bool

	leftTimeline  bool
	rightTimeline bool
	timelines     int

	top    *tachyonField
	left   *tachyonField
	bottom *tachyonField
	right  *tachyonField
}

func (t tachyonField) display() string {
	if t.start {
		return "S"
	}
	if t.splitter {
		return "^"
	}
	if t.beam {
		return "|"
	}
	return "."
}

func (t *tachyonField) shoot() {
	if t.splitter {
		t.isSplitting = true
		if t.left != nil && !t.left.beam {
			t.left.shoot()
		}
		if t.right != nil && !t.right.beam {
			t.right.shoot()
		}
	} else {
		t.beam = true
		if t.bottom != nil {
			t.bottom.shoot()
		}
	}
}

func (t *tachyonField) revert(timeline int) {
	t.timelines = timeline
	if t.left != nil && t.left.splitter && t.left.isSplitting {
		t.left.revertSplitter(t.timelines, 'l')
	}
	if t.right != nil && t.right.splitter && t.right.isSplitting {
		t.right.revertSplitter(t.timelines, 'r')
	}
	if t.top != nil && t.top.beam {
		t.top.revert(t.timelines)
	}
}

func (t *tachyonField) revertSplitter(timeline int, source rune) {
	t.timelines += timeline
	if source == 'l' {
		t.leftTimeline = true
	}
	if source == 'r' {
		t.rightTimeline = true
	}
	if t.leftTimeline && t.rightTimeline {
		t.revert(t.timelines)
	}
}

func (t tachyonField) displayWithConnction() []string {
	middle := t.display()
	right, bottom := "", "  "
	if t.right != nil && t.right.left != nil {
		right = "-"
	}
	if t.bottom != nil && t.bottom.top != nil {
		bottom = "| "
	}

	return []string{strings.Join([]string{middle, right}, ""), bottom}
}

func findStart(diagram [][]tachyonField) (int, int) {
	line := diagram[0]

	for i := range line {
		if line[i].display() == "S" {
			return 0, i
		}
	}

	panic("No starting point found!")
}

func countSplitting(diagram [][]tachyonField) {
	acc := 0
	for y := range diagram {
		line := diagram[y]
		for x := range line {
			if diagram[y][x].isSplitting {
				acc++
			}
		}
	}

	fmt.Println("Solution 1:", acc)
}

func walkBackInTime(diagram [][]tachyonField) {
	end := diagram[len(diagram)-1]
	acc := 0
	for i := range end {
		if end[i].beam {
			end[i].revert(1)
			acc++
			//fmt.Printf("%d... ", acc)
		}
		//fmt.Println()
	}

	yStart, xStart := findStart(diagram)
	start := diagram[yStart][xStart]
	fmt.Println("Solution 2:", start.timelines)
}
