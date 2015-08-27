package main

type Wav struct {
	ChunkID   [4]byte // B
	ChunkSize uint32  // L
	Format    [4]byte // B

	Subchunk1ID   [4]byte // B
	Subchunk1Size uint32  // L
	AudioFormat   uint16  // L
	NumChannels   uint16  // L
	SampleRate    uint32  // L
	ByteRate      uint32  // L
	BlockAlign    uint16  // L
	BitsPerSample uint16  // L

	Subchunk2ID   [4]byte // B
	Subchunk2Size uint32  // L

	NumSamples uint32
}
