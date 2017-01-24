package main

import "os"
import "io"
import "encoding/binary"
import "fmt"

type iNESHeader struct {
	Constant  [4]byte
	AmountPRG uint8
	AmountCHR uint8
	Flags6    byte
	_         [9]byte
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	inesFilePath := os.Args[1]

	inesFile, err := os.Open(inesFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer inesFile.Close()

	var header iNESHeader
	binary.Read(inesFile, binary.LittleEndian, &header)

	if header.Constant != [4]byte{0x4E, 0x45, 0x53, 0x1A} {
		fmt.Println("Not a valid iNES file.")
		os.Exit(1)
	}

	hasTrainer := (header.Flags6 & 4) > 0
	if hasTrainer {
		inesFile.Seek(512, 1)
	}

	inesFile.Seek(16384*int64(header.AmountPRG), 1)

	io.Copy(os.Stdout, inesFile)
	os.Stdout.Sync()

	os.Exit(0)
}

func printUsage() {
	fmt.Printf("Usage: %s file.nes\n", os.Args[0])
	os.Exit(1)
}
