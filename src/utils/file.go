package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func ReadFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("couldn't open file %s\n", path)
	}
	defer file.Close()

	var data [][]string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		fields := strings.Fields(line)
		data = append(data, fields)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil

}
