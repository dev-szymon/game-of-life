package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gdamore/tcell"
)

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

type Cell struct {
	isLive         bool
	liveNeighbours int
}
type Row = map[int]*Cell
type Matrix = map[int]Row
type Game struct {
	rows Matrix
}

func NewGame() *Game {
	game := &Game{rows: make(Matrix)}
	return game
}

func (g *Game) Matrix() Matrix {
	return g.rows
}

func (g *Game) initialiseSeed(seed [][]int) {
	for _, cell := range seed {
		row := cell[0]
		col := cell[1]
		if g.rows[row] == nil {
			g.rows[row] = make(Row)
		}
		g.rows[row][col] = &Cell{isLive: true, liveNeighbours: 0}
	}
}

func (g *Game) selectLiveCellsForNextIteration() {
	for row, cols := range g.rows {
		for col, cell := range cols {
			if cell.isLive && (cell.liveNeighbours == 2 || cell.liveNeighbours == 3) {
				cell.liveNeighbours = 0
			} else if !cell.isLive && cell.liveNeighbours == 3 {
				cell.isLive = true
				cell.liveNeighbours = 0
			} else {
				delete(cols, col)
			}
		}
		if len(cols) == 0 {
			delete(cols, row)
		}
	}
}

func (g *Game) increaseCellNeighbours(rows Matrix, row, col int) {
	isLive := g.rows[row] != nil && g.rows[row][col] != nil
	if rows[row][col] == nil {
		rows[row][col] = &Cell{isLive: isLive, liveNeighbours: 1}
	} else {
		rows[row][col].liveNeighbours++
	}
}

func (g *Game) calculateNextIterationNeighbours(width, height int) {
	nextRows := make(Matrix)
	for row, cols := range g.rows {
		for col := range cols {
			rowUp := row - 1
			rowDown := row + 1
			colLeft := col - 1
			colRight := col + 1
			if nextRows[row] == nil {
				nextRows[row] = make(Row)
			}
			if nextRows[rowUp] == nil {
				nextRows[rowUp] = make(Row)
			}
			if nextRows[rowDown] == nil {
				nextRows[rowDown] = make(Row)
			}

			if nextRows[row][col] == nil {
				nextRows[row][col] = &Cell{isLive: true, liveNeighbours: 0}
			}
			if rowUp >= 0 {
				g.increaseCellNeighbours(nextRows, rowUp, col)
				if colLeft >= 0 {
					g.increaseCellNeighbours(nextRows, rowUp, colLeft)
				}
				if colRight <= width {
					g.increaseCellNeighbours(nextRows, rowUp, colRight)
				}
			}
			if rowDown <= height {
				g.increaseCellNeighbours(nextRows, rowDown, col)
				if colLeft >= 0 {
					g.increaseCellNeighbours(nextRows, rowDown, colLeft)
				}
				if colRight <= width {
					g.increaseCellNeighbours(nextRows, rowDown, colRight)
				}
			}
			if colLeft >= 0 {
				g.increaseCellNeighbours(nextRows, row, colLeft)
			}
			if colRight <= width {
				g.increaseCellNeighbours(nextRows, row, colRight)
			}
		}
	}
	g.rows = nextRows
}

var aliveCellStyle = tcell.StyleDefault.Background(tcell.ColorTomato).Foreground(tcell.ColorTomato)

const statsRows = 5

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

		for row, cols := range game.Matrix() {
			for col, cell := range cols {
				if cell.isLive {
					s.SetContent(col, row, ' ', nil, aliveCellStyle)
				}
			}
		}
		game.calculateNextIterationNeighbours(width, height-statsRows)
		game.selectLiveCellsForNextIteration()
		PrintMemUsage(s)
		s.Show()
	}
}

func PrintMemUsage(screen tcell.Screen) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc := fmt.Sprintf("Alloc = %v MiB", m.Alloc/1024/1024)
	totalAlloc := fmt.Sprintf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	Mib := fmt.Sprintf("\tSys = %v MiB", m.Sys/1024/1024)
	numgc := fmt.Sprintf("\tNumGC = %v\n", m.NumGC)
	width, height := screen.Size()
	drawText(screen, 0, height-1, width, height, tcell.StyleDefault, alloc)
	drawText(screen, 0, height-2, width, height-1, tcell.StyleDefault, totalAlloc)
	drawText(screen, 0, height-3, width, height-2, tcell.StyleDefault, Mib)
	drawText(screen, 0, height-4, width, height-3, tcell.StyleDefault, numgc)
}

func drawText(screen tcell.Screen, startX, startY, endX, endY int, style tcell.Style, text string) {
	row := startY
	col := startX
	for _, r := range text {
		screen.SetContent(col, row, r, nil, style)
		col++
		if col >= endX {
			row++
			col = startX
		}
		if row > endY {
			break
		}
	}
}
