package main

import "fmt"
import "os"
import "path"
import "math"
import "image"
import "image/png"
import "image/color"

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	chrFilePath := os.Args[1]
	chrFileBase := path.Base(chrFilePath)

	var pngFilePath string

	if len(os.Args) >= 3 {
		pngFilePath = os.Args[2]

		if path.Ext(pngFilePath) != ".png" {
			printUsage()
		}
	} else {
		pngFilePath = chrFileBase[:len(chrFileBase)-len(path.Ext(chrFileBase))] + ".png"
	}

	if err := convertChrToPng(chrFilePath, pngFilePath); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %s file.chr [file.png]\n", os.Args[0])
	os.Exit(1)
}

func convertChrToPng(chrFilePath string, pngFilePath string) error {
	fileInfo, err := os.Stat(chrFilePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if fileInfo.Size()%16 != 0 {
		fmt.Printf("Invalid CHR file?\n")
		os.Exit(1)
	}

	amountTiles := int(fileInfo.Size() / 16)

	chrFile, err := os.Open(chrFilePath)
	if err != nil {
		return err
	}
	defer chrFile.Close()

	pngFile, err := os.Create(pngFilePath)
	if err != nil {
		return err
	}
	defer pngFile.Close()

	var height, width int

	if amountTiles < 16 {
		height = 8
		width = 8 * amountTiles
	} else {
		height = int(math.Ceil(float64(amountTiles)/16.0) * 8)
		width = 8 * 16
	}

	palette := color.Palette{
		color.Transparent,
		color.RGBA{0x99, 0xB0, 0xFF, 0xFF},
		color.RGBA{0xFF, 0x3A, 0x50, 0xFF},
		color.RGBA{0x99, 0x28, 0x50, 0xFF},
	}

	chrImage := image.NewPaletted(image.Rect(0, 0, width, height), palette)

	var tileBytes = make([]byte, 16)

	for tile := 0; tile < amountTiles; tile++ {
		_, err := chrFile.Read(tileBytes)

		if err != nil {
			return err
		}

		tileOffsetX := (tile % 16) * 8
		tileOffsetY := ((tile - (tile % 16)) / 16) * 8

		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				bitPosition := uint(7 - x)
				bit1 := (tileBytes[y] & (1 << bitPosition)) >> bitPosition
				bit2 := (tileBytes[y+8] & (1 << bitPosition)) >> bitPosition

				colorIndex := bit1 | (bit2 << 1)

				chrImage.SetColorIndex(x+tileOffsetX, y+tileOffsetY, colorIndex)
			}
		}
	}

	png.Encode(pngFile, chrImage)

	return nil
}
