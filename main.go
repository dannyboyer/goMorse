package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

//directly use byte for optimization
var m = map[string]string{
	"A": ".-",
	"B": "-...",
}

func main() {

	filename := os.Args[1]

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for scanner.Scan() {
		fmt.Println(translateToMorse(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}
}

func translateToMorse(line string) string {

	var morseCode string
	for _, r := range line {
		char := strings.ToUpper(string(r))
		if val, ok := m[char]; ok {
			morseCode += val
		}
	}

	return morseCode
}
