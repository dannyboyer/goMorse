package main

import (
	"bufio"
	bin "encoding/binary"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

const (
	Wpm  = 20.0      //words per minutes
	Tone = 900       //frequency in Hertz
	Sps  = 44100     //samples per seconds
	Eps  = Wpm / 1.2 //elements per second (frequency of morse coding)
	Bit  = 1.2 / Wpm //seconds per element (period of morse coding)

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

	morseCode := translateRuneToMorse(message)
	fmt.Println(morseCode)

	freqSlice := translateMorseToFreq(morseCode)
	fmt.Println(freqSlice)

	data := translateFreqToData(freqSlice)

	writeWave("morse.wav", 8000, data)
}

func translateRuneToMorse(input string) string {
	var morseCode string
	for _, r := range input {
		char := strings.ToUpper(string(r))
		if val, ok := m[char]; ok {
			morseCode += val + " "
		}
	}
	return morseCode
}

func translateMorseToFreq(input string) []int {
	var freqSlice []int
	for _, c := range input {
		if c == '.' {
			freqSlice = append(freqSlice, 1, 0)
		} else if c == '-' {
			freqSlice = append(freqSlice, 1, 1, 1, 0)
		} else if c == ' ' {
			freqSlice = append(freqSlice, 0, 0) //end of letter
		}
	}
	freqSlice = append(freqSlice, 0, 0, 0, 0) //end of word
	return freqSlice
}

func translateFreqToData(input []int) []int16 {
	var data []int16

	ampl := 0.85 * 32767
	pi := 3.1415926535
	w := 2.0 * pi * Tone
	var i, n int32

	n = int32(Bit * Sps)
	for _, freq := range input {
		for i = 0; i < n; i += 1 {
			t := float64(i) / Sps
			data = append(data, int16(freq)*int16(ampl*math.Sin(w*t)))
		}
	}

	return data
}

func writeWave(fn string, sample_rate int, data []int16) {
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Error opening file for writing: %v\n", err)
	}

	var header Wav

	header.ChunkID = [4]byte{'R', 'I', 'F', 'F'}
	header.Format = [4]byte{'W', 'A', 'V', 'E'}
	header.Subchunk1ID = [4]byte{'f', 'm', 't', ' '}
	header.Subchunk1Size = 16
	header.AudioFormat = 1 // 1 == PCM
	header.NumChannels = 2 // 2 == Stereo
	header.SampleRate = uint32(sample_rate)
	header.BitsPerSample = 16 // 16bit integer samples
	header.ByteRate = uint32(int(header.SampleRate) * int(header.NumChannels) * int(header.BitsPerSample/8))
	header.BlockAlign = uint16(int(header.NumChannels) * int(header.BitsPerSample/8))
	header.Subchunk2ID = [4]byte{'d', 'a', 't', 'a'}

	header.Subchunk2Size = uint32(len(data) * int(header.BitsPerSample/8))
	header.ChunkSize = 4 + (8 + header.Subchunk1Size) + (8 + header.Subchunk2Size)

	bin.Write(file, bin.BigEndian, &header.ChunkID)
	bin.Write(file, bin.LittleEndian, &header.ChunkSize)
	bin.Write(file, bin.BigEndian, &header.Format)

	bin.Write(file, bin.BigEndian, &header.Subchunk1ID)
	bin.Write(file, bin.LittleEndian, &header.Subchunk1Size)
	bin.Write(file, bin.LittleEndian, &header.AudioFormat)
	bin.Write(file, bin.LittleEndian, &header.NumChannels)
	bin.Write(file, bin.LittleEndian, &header.SampleRate)
	bin.Write(file, bin.LittleEndian, &header.ByteRate)
	bin.Write(file, bin.LittleEndian, &header.BlockAlign)
	bin.Write(file, bin.LittleEndian, &header.BitsPerSample)

	bin.Write(file, bin.BigEndian, &header.Subchunk2ID)
	bin.Write(file, bin.LittleEndian, &header.Subchunk2Size)

	bin.Write(file, bin.LittleEndian, data)

	file.Close()

	return
}
