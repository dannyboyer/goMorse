package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/cryptix/wav"
)

const (
	bits = 32
	rate = 8000
)

//directly use byte for optimization
var m = map[string]string{
	"-": "-....-",
	",": "--..--",
	".": ".-.-.-",
	" ": "/",
	"A": ".-",
	"B": "-...",
	"C": "-.-.",
	"D": "-..",
	"E": ".",
	"F": "..-.",
	"G": "--.",
	"H": "....",
	"I": "..",
	"J": ".---",
	"K": "-.-",
	"L": ".-..",
	"M": "--",
	"N": "-.",
	"O": "---",
	"P": ".--.",
	"Q": "--.-",
	"R": ".-.",
	"S": "...",
	"T": "-",
	"U": "..-",
	"V": "...-",
	"W": ".--",
	"X": "-..-",
	"Y": "-.--",
	"Z": "--..",
	"0": "-----",
	"1": ".----",
	"2": "..---",
	"3": "...--",
	"4": "....-",
	"5": ".....",
	"6": "-....",
	"7": "--...",
	"8": "---..",
	"9": "----.",
}

func main() {

	filename := os.Args[1]

	inputFile, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}

	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	var morseCode string

	for scanner.Scan() {
		morseCode += translateToMorse(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
	}

	fmt.Println(morseCode)
}

func translateToMorse(line string) string {

	var morseCode string
	for _, r := range line {
		char := strings.ToUpper(string(r))
		if val, ok := m[char]; ok {
			morseCode += val + " "
		}
	}

	return morseCode
}

func translateToFreq(morseCode string) []int32 {

	var freqSlice []int32

	for _, code := range morseCode {
		if code == ' ' {
			freqSlice = append(freqSlice, 0)
		} else {
			freqSlice = append(freqSlice, 1)
		}
	}

	return freqSlice
}

func writeWavFile(morseCode []int32) {

	wavOut, err := os.Create("output.wav")
	checkErr(err)
	defer wavOut.Close()

	meta := wav.File{
		Channels:        1,
		SampleRate:      rate,
		SignificantBits: bits,
	}

	writer, err := meta.NewWriter(wavOut)
	checkErr(err)
	defer writer.Close()

	var freq float64 = 0.5

	for n := 0; n < len(morseCode); n += 1 {

		y := int32(0.8 * math.Pow(2, bits-1) * math.Sin(float64(morseCode[n])*freq*float64(n)))

		for n := 0; n < 4000; n += 1 {
			err = writer.WriteInt32(y)
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
