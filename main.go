package main

import (
	"log"
	"time"

	"github.com/gdamore/tcell"
)

type Cell struct {
	isLive         bool
	liveNeighbours int
}
type Row = map[int]*Cell
type Matrix = map[int]Row
type Game struct {
	state Matrix
}

func increaseNeighbours(matrix Matrix, row, col int) {
	rowToChange, ok := matrix[row]
	if !ok {
		rowToChange = make(Row)
	}
	colToChange, ok := rowToChange[col]
	if !ok {
		colToChange = &Cell{isLive: false}
	}
	colToChange.liveNeighbours++

	// Then we modify the copy
	rowToChange[col] = colToChange

	// Then we reassign map entry
	matrix[row] = rowToChange
}

func setIsLive(matrix Matrix, row, col int) {
	if matrix[row] == nil {
		matrix[row] = make(Row)
	}
	cell := matrix[row][col]
	if cell == nil {
		matrix[row][col] = &Cell{isLive: true, liveNeighbours: 0}
	} else {
		cell.isLive = true
	}
}

func selectLiveCellsForNextIteration(matrix Matrix) Matrix {
	nextMatrix := make(Matrix)
	for row, cols := range matrix {
		for col, cell := range cols {
			if cell.isLive && (cell.liveNeighbours == 2 || cell.liveNeighbours == 3) {
				if nextMatrix[row] == nil {
					nextMatrix[row] = make(Row)
				}
				nextMatrix[row][col] = &Cell{isLive: true}
			}
			if !cell.isLive && cell.liveNeighbours == 3 {
				if nextMatrix[row] == nil {
					nextMatrix[row] = make(Row)
				}
				nextMatrix[row][col] = &Cell{isLive: true}
			}
		}
	}

	return nextMatrix
}

func NewGame() *Game {
	game := &Game{state: make(Matrix)}
	return game
}

func NewScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := s.Init(); err != nil {
		return nil, err
	}
	s.SetStyle(tcell.StyleDefault)
	return s, err

}

func (g *Game) initialiseSeed(seed [][]int) {
	for _, cell := range seed {
		row := cell[0]
		col := cell[1]
		if g.state[row] == nil {
			g.state[row] = make(Row)
		}
		g.state[row][col] = &Cell{isLive: true, liveNeighbours: 0}
	}
}

func calculateNextIterationNeighbours(prevState Matrix, width, height int) Matrix {
	nextMatrix := make(Matrix)
	for row, cols := range prevState {
		for col := range cols {
			if nextMatrix[row] == nil {
				nextMatrix[row] = make(Row)
			}
			rowUp := row - 1
			rowDown := row + 1
			colLeft := col - 1
			colRight := col + 1

			setIsLive(nextMatrix, row, col)
			if rowUp >= 0 {
				if nextMatrix[rowUp] == nil {
					nextMatrix[rowUp] = make(Row)
				}

				increaseNeighbours(nextMatrix, rowUp, col)
				if colLeft >= 0 {
					increaseNeighbours(nextMatrix, rowUp, colLeft)
				}
				if colRight <= width {
					increaseNeighbours(nextMatrix, rowUp, colRight)
				}
			}
			if rowDown <= height {
				if nextMatrix[rowDown] == nil {
					nextMatrix[rowDown] = make(Row)
				}

				increaseNeighbours(nextMatrix, rowDown, col)
				if colLeft >= 0 {
					increaseNeighbours(nextMatrix, rowDown, colLeft)
				}
				if colRight <= width {
					increaseNeighbours(nextMatrix, rowDown, colRight)
				}
			}
			if colLeft >= 0 {
				increaseNeighbours(nextMatrix, row, colLeft)
			}
			if colRight <= width {
				increaseNeighbours(nextMatrix, row, colRight)
			}
		}
	}

	return nextMatrix
}

var aliveCellStyle = tcell.StyleDefault.Background(tcell.ColorTomato).Foreground(tcell.ColorTomato)

func main() {
	seed := [][]int{{1, 1}, {2, 1}, {2, 5}, {2, 7}, {8, 7}, {8, 9}, {8, 8}, {3, 2}, {3, 6}, {3, 7}, {8, 2}, {8, 3}, {8, 4}, {9, 5}, {9, 4}, {9, 2}}
	game := NewGame()
	game.initialiseSeed(seed)

	s, err := NewScreen()
	if err != nil {
		log.Fatalf("error initialising screen: %+v", err)
	}
	width, height := s.Size()

	ticker := time.NewTicker(300 * time.Millisecond)
	for range ticker.C {
		s.Clear()

		for row, cols := range game.state {
			for col, cell := range cols {
				if cell.isLive {
					s.SetContent(col, row, ' ', nil, aliveCellStyle)
				}
			}
		}
		nextIteration := calculateNextIterationNeighbours(game.state, width, height)
		game.state = selectLiveCellsForNextIteration(nextIteration)
		s.Show()
	}
}
