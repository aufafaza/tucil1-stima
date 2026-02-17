package main

import (
	"fmt"
	"github.com/aufafaza/tucil1-stima.git/src/models"
	"github.com/aufafaza/tucil1-stima.git/src/solver"
	"github.com/aufafaza/tucil1-stima.git/src/utils"
	"log"
)

func main() {
	grid, err := utils.ReadFile("../test/test1.txt")
	if err != nil {
		log.Fatal(err)
	}
	N := len(grid)
	if N == 0 {
		log.Fatal("empty board")
	}
	log.Printf("Board size: %dx%d\n", N, N)
	board := &models.Board{
		Size: N,
		Grid: grid,
		Q:    make([]int, N),
		Iter: 0,
	}
	found := solver.Solver(board)

	if found {
		log.Printf("Solution found in %d interations\n", board.Iter)
		printBoard(board)
	} else {
		log.Println("solution not found")
	}
}

func printBoard(b *models.Board) {
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			if b.Q[r] == c {
				fmt.Print("# ") // Queen
			} else {
				// Print the color (or '.' for empty to see clearly)
				fmt.Printf("%s ", b.Grid[r][c])
			}
		}
		fmt.Println()
	}
}
