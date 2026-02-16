package main

import (
	"log"

	"github.com/aufafaza/tucil1-stima.git/src/utils"
)

func main() {
	data, _ := utils.ReadFile("../test/test1.txt")

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			log.Printf("%s ", data[i][j])
		}
		log.Println()
	}
}
