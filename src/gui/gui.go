package gui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aufafaza/tucil1-stima.git/src/models"
	"github.com/aufafaza/tucil1-stima.git/src/solver"
	"github.com/aufafaza/tucil1-stima.git/src/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 800
)

type Game struct {
	Board              *models.Board
	Found              bool
	IterationsPerFrame int
	TileSize           float32
	InputName          string
	StartTime          time.Time
	EndTime            time.Duration
}
type Color struct {
	R, G, B uint8
}

var palette = map[rune]color.RGBA{
	'A': {255, 0, 0, 255},
	'B': {0, 255, 0, 255},
	'C': {0, 0, 255, 255},
	'D': {255, 255, 0, 255},
	'E': {255, 0, 255, 255},
	'F': {0, 255, 255, 255},
	'G': {255, 165, 0, 255},
	'H': {128, 0, 128, 255},
	'I': {0, 128, 0, 255},
	'J': {165, 42, 42, 255},
	'K': {255, 192, 203, 255},
	'L': {128, 128, 128, 255},
	'M': {0, 0, 128, 255},
	'N': {128, 128, 0, 255},
	'O': {0, 128, 128, 255},
	'P': {192, 192, 192, 255},
	'.': {230, 230, 230, 255}, // Added dot for empty cells
}
var DefaultColor = color.RGBA{200, 200, 200, 255}

//init board

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	for row := 0; row < g.Board.Size; row++ {
		for col := 0; col < g.Board.Size; col++ {
			x := float32(col) * g.TileSize
			y := float32(row) * g.TileSize

			// Ensure char is a rune for the palette lookup
			char := rune(g.Board.Grid[row][col][0])

			cellColor, poly := palette[char]
			if !poly {
				cellColor = DefaultColor
			}

			vector.FillRect(
				screen,
				x, y,
				g.TileSize-1, g.TileSize-1,
				cellColor,
				false,
			)
			if g.Board.Q[row] == col {
				queenColor := color.Black

				if g.Found {
					queenColor = color.White
				}
				padding := g.TileSize * 0.2
				vector.DrawFilledCircle(
					screen,
					x+(g.TileSize/2),
					y+(g.TileSize/2),
					(g.TileSize/2)-padding,
					queenColor,
					true,
				)
			}
		}
	}
	var elapsed time.Duration
	if g.Found && g.EndTime != 0 {
		elapsed = g.EndTime
	} else {
		elapsed = time.Since(g.StartTime)
	}

	stats := fmt.Sprintf(
		"File: %s\nIterations: %d\nRuntime: %v\nStatus: %s",
		filepath.Base(g.InputName),
		g.Board.Iter,
		elapsed.Round(time.Millisecond),
		func() string {
			if g.Found {
				return "FINISHED"
			}
			return "SEARCHING..."
		}(),
	)

	ebitenutil.DebugPrint(screen, stats)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// This defines the coordinate system inside your Draw function
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	if g.Found {
		return nil
	}
	for i := 0; i < g.IterationsPerFrame; i++ {
		g.Board.Iter++

		// per 100k iterations, add logging
		if g.Board.Iter%100000 == 0 {
			PrintBoardCLI(g.Board)
		}
		if solver.CheckValid(g.Board) {
			fmt.Println("Solution found")
			PrintBoardCLI(g.Board)
			g.Found = true
			g.EndTime = time.Since(g.StartTime)
			g.Board.SolCount++
			solCopy := make([]int, g.Board.Size)
			copy(solCopy, g.Board.Q)
			g.Board.Solutions = [][]int{solCopy}

			baseName := filepath.Base(g.InputName)
			nameOnly := strings.TrimSuffix(baseName, filepath.Ext(baseName))
			outputName := "solution_" + nameOnly + ".txt"
			imageName := "solution_" + nameOnly + ".png"
			outputDir := "output"
			os.MkdirAll(outputDir, 0755)

			finalPath := filepath.Join(outputDir, outputName)
			utils.WriteFile(finalPath, g.Board)
			g.SaveImage(filepath.Join("output", imageName))
			return nil

		}

		if !solver.NextState(g.Board) {
			g.Found = true
			return nil
		}
	}
	return nil
}

func StartGame(grid [][]string, originalPath string) {
	size := len(grid)

	if size == 0 {
		log.Fatal("board is empty")
	}

	for _, row := range grid {
		if len(row) != size {
			log.Fatalf("unsolvable, not an nxn board. Row: %v, Expected Row: %v", len(row), size)
		}
	}
	colorMap := make(map[string]bool)
	for _, row := range grid {
		for _, cell := range row {
			if cell != "" {
				colorMap[cell] = true
			}
		}
	}
	if len(colorMap) < size {
		log.Fatalf("unsolvable, only %v colors found, needed %v colors", len(colorMap), size, size)
	}
	if size > 1 && size < 4 {
		log.Fatal("unsolvable, no solutions for size %v", size)
	}
	board := &models.Board{
		Grid: grid,
		Size: size,
		Q:    make([]int, size), // Init queens at column 0
	}

	game := &Game{
		Board:              board,
		Found:              false,
		IterationsPerFrame: 1000000,
		TileSize:           float32(screenWidth) / float32(size),
		InputName:          originalPath,
		StartTime:          time.Now(),
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("N-Queens")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) SaveImage(path string) {
	outImg := ebiten.NewImage(screenWidth, screenHeight)

	outImg.Fill(color.White)

	for row := 0; row < g.Board.Size; row++ {
		for col := 0; col < g.Board.Size; col++ {
			x := float32(col) * g.TileSize
			y := float32(row) * g.TileSize

			char := rune(g.Board.Grid[row][col][0])
			cellColor, ok := palette[char]
			if !ok {
				cellColor = DefaultColor
			}

			vector.FillRect(
				outImg,
				x, y,
				g.TileSize-1, g.TileSize-1,
				cellColor,
				false,
			)

			if g.Board.Q[row] == col {
				queenColor := color.Black
				if g.Found {
					queenColor = color.White
				}
				padding := g.TileSize * 0.2

				vector.DrawFilledCircle(
					outImg,
					x+(g.TileSize/2),
					y+(g.TileSize/2),
					(g.TileSize/2)-padding,
					queenColor,
					true,
				)
			}
		}
	}

	rect := outImg.Bounds()
	rgba := image.NewRGBA(rect)

	outImg.ReadPixels(rgba.Pix)

	f, err := os.Create(path)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return
	}
	defer f.Close()

	if err := png.Encode(f, rgba); err != nil {
		log.Printf("Failed to encode PNG: %v", err)
	}
}

func PrintBoardCLI(b *models.Board) {
	for r := 0; r < b.Size; r++ {
		for c := 0; c < b.Size; c++ {
			if b.Q[r] == c {
				fmt.Print("# ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println("-------------------")
}
