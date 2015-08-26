package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	Wpm  = 20.0        //words per minutes
	Tone = 900         //frequency in Hertz
	Sps  = 44100       //samples per seconds
	Path = "morse.wav" //filepath
)

//directly compare byte for optimization
var m = map[string]string{
	":":  "---...",
	"?":  "..--..",
	"'":  ".----.",
	"/":  "-..-.",
	"(":  "-.--.-",
	")":  "-.--.-",
	"\"": ".-..-.",
	"@":  ".--.-.",
	"=":  "-...-",
	"-":  "-....-",
	",":  "--..--",
	".":  ".-.-.-",
	" ":  "/",
	"A":  ".-",
	"B":  "-...",
	"C":  "-.-.",
	"D":  "-..",
	"E":  ".",
	"F":  "..-.",
	"G":  "--.",
	"H":  "....",
	"I":  "..",
	"J":  ".---",
	"K":  "-.-",
	"L":  ".-..",
	"M":  "--",
	"N":  "-.",
	"O":  "---",
	"P":  ".--.",
	"Q":  "--.-",
	"R":  ".-.",
	"S":  "...",
	"T":  "-",
	"U":  "..-",
	"V":  "...-",
	"W":  ".--",
	"X":  "-..-",
	"Y":  "-.--",
	"Z":  "--..",
	"0":  "-----",
	"1":  ".----",
	"2":  "..---",
	"3":  "...--",
	"4":  "....-",
	"5":  ".....",
	"6":  "-....",
	"7":  "--...",
	"8":  "---..",
	"9":  "----.",
}

func main() {

	filename := os.Args[1]

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	var message string

	for scanner.Scan() {
		message += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}

	morseCodeChannel := make(chan string)

	go translateRuneToMorse(message, morseCodeChannel)

	for i := range morseCodeChannel {
		fmt.Print(i)
	}
}

func translateRuneToMorse(input string, output chan string) {
	for _, r := range input {
		char := strings.ToUpper(string(r))
		if val, ok := m[char]; ok {
			output <- val + " "
		}
	}
	close(output)
}
