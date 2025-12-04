package main

import (
	"fmt"
	"strings"
)

func runC4() {
	grid := readFile("c4.txt", mapToGrid)
	gridTwo := grid
	_, forkable := execC4(grid)
	fmt.Println("Solution 1: ", forkable)

	gridMarkedTwo, forkableTwo := execC4Two(gridTwo)
	printGrid(gridMarkedTwo)
	fmt.Println("Solution 2: ", forkableTwo)
}

func printGrid(grid [][]rune) {
	for _, g := range grid {
		fmt.Println(string(g))
	}
}

func mapToGrid(content []byte) [][]rune {
	NEWLINE := '\n'

	contentLines := make([]string, 0, 100)
	for part := range strings.SplitSeq(string(content), string(NEWLINE)) {
		contentLines = append(contentLines, strings.ReplaceAll(part, string(NEWLINE), ""))
	}

	grid := make([][]rune, 0, 100)

	for _, lines := range contentLines {
		grid = append(grid, []rune(lines))
	}

	return grid
}

type calc struct {
	col int
	row int
}

func execC4(grid [][]rune) ([][]rune, int) {
	newGrid := copyGrid(grid)

	forkable := 0
	for i, col := range newGrid {
		for j, row := range col {
			if row == '@' {
				if isForkable(calc{i, j}, newGrid) {
					newGrid[i][j] = 'x'
					forkable++
				}
			}
		}
	}

	return newGrid, forkable
}

func isForkable(current calc, grid [][]rune) bool {
	CALCULATIONS := []calc{{-1, -1}, {-1, 0}, {0, -1}, {-1, 1}, {1, -1}, {1, 0}, {0, 1}, {1, 1}}

	rolls := 0

	for _, calc := range CALCULATIONS {
		neighbor := look(current, calc, grid)
		if neighbor == '@' || neighbor == 'x' {
			rolls++
		}
	}

	if rolls < 4 {
		return true
	} else {
		return false
	}
}

func look(position calc, neighbor calc, grid [][]rune) rune {
	col := position.col + neighbor.col
	row := position.row + neighbor.row

	outsideCol := func(col int) bool {
		if col < 0 || col >= len(grid) {
			return true
		}
		return false
	}

	outsideRow := func(row int) bool {
		if row < 0 || row >= len(grid[0]) {
			return true
		}
		return false
	}

	if outsideCol(col) || outsideRow(row) {
		return 0x00
	}

	return grid[col][row]
}

func execC4Two(prev [][]rune) ([][]rune, int) {
	forkable := 0

	for {
		next, nextForkable := execC4(prev)
		next = removeRolls(next)
		forkable += nextForkable

		if isEqualGrids(prev, next) {
			return next, forkable
		} else {
			prev = next
		}

	}
}

func removeRolls(grid [][]rune) [][]rune {
	for i, col := range grid {
		for j, row := range col {
			if row == 'x' {
				grid[i][j] = '.'
			}
		}
	}

	return grid
}

func isEqualGrids(prev [][]rune, current [][]rune) bool {
	for i, line := range current {
		for j := range line {
			if prev[i][j] != current[i][j] {
				return false
			}
		}
	}

	return true
}

func copyGrid(original [][]rune) [][]rune {
	newGrid := make([][]rune, 0, len(original))
	for _, line := range original {
		newLine := make([]rune, len(line))
		copy(newLine, line)
		newGrid = append(newGrid, newLine)
	}

	return newGrid
}
