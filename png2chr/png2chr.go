package main

import "os"
import "fmt"
import "image"
import "image/png"
import "path"
import "errors"

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	pngFilePath := os.Args[1]

	var chrFilePath string

	if len(os.Args) >= 3 {
		chrFilePath = os.Args[2]
	} else {
		pngFileBase := path.Base(pngFilePath)
		chrFilePath = pngFileBase[:len(pngFileBase)-len(path.Ext(pngFileBase))] + ".chr"
	}

	if err := convertPngToChr(pngFilePath, chrFilePath); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %s file.png [file.chr]\n", os.Args[0])
	os.Exit(1)
}

func convertPngToChr(pngFilePath string, chrFilePath string) error {
	pngFile, err := os.Open(pngFilePath)
	if err != nil {
		return err
	}

	img, err := png.Decode(pngFile)
	if err != nil {
		return err
	}

	if _, ok := img.(image.PalettedImage); ok == false {
		return errors.New("Input image must be a paletted PNG.")
	}

	imageBounds := img.Bounds()

	width := imageBounds.Max.X
	height := imageBounds.Max.Y

	if (width%8) > 0 || (height%8) > 0 {
		return errors.New("Image's width and height must be divisible by 8.")
	}

	columns := width / 8
	rows := height / 8

	if rows*columns > 512 {
		return errors.New("There is too many tiles to fit in a 8KB CHR ROM.")
	}

	chrFile, err := os.Create(chrFilePath)
	if err != nil {
		return err
	}

	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			chr, err := convertTileToChr(column, row, img.(image.PalettedImage))
			if err != nil {
				return err
			}

			if _, err := chrFile.Write(chr); err != nil {
				return err
			}
		}
	}

	return nil
}

func convertTileToChr(column int, row int, img image.PalettedImage) ([]byte, error) {
	plane1 := make([]byte, 8)
	plane2 := make([]byte, 8)

	offsetX := column * 8
	offsetY := row * 8

	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			realX := x + offsetX
			realY := y + offsetY

			colorIndex := img.ColorIndexAt(realX, realY)

			if colorIndex > 4 {
				return nil, errors.New(fmt.Sprintf("Palette can have no greater than 4 colors. Found %d at %d,%d.", colorIndex, realX, realY))
			}

			bitPosition := uint(7 - x)

			plane1[y] |= ((colorIndex & 1) << bitPosition)
			plane2[y] |= (((colorIndex & 2) >> 1) << bitPosition)
		}
	}

	return append(plane1, plane2...), nil
}
