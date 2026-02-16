package models

type Board struct {
	Size      int
	Grid      [][]string
	Q         []int
	Iter      int
	Solutions int
}

func NewBoard(size int, gridData [][]string) *Board {
	return &Board{
		Size: size,
		Grid: gridData,
		Q:    make([]int, size), // jumlah Q = n (n x n grid)
	}
}
