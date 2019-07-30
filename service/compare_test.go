package service

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"testing"
)

var (
	testImage1 = "test1.png"
	testImage2 = "test2.png"
	output     = "output.png"
)

type direction int

const (
	DIRECTION_UP direction = iota
	DIRECTION_DOWN
)

func TestComparer_Compare(t *testing.T) {

	createTestImage(testImage1, 150, 150, DIRECTION_UP)
	createTestImage(testImage2, 150, 150, DIRECTION_DOWN)

	comparer, err := NewComparer(testImage1, testImage2, output)
	if err != nil {
		t.Fatal(err)
	}
	// do the comparison
	if err = comparer.Compare(MATCH_MODE_BLACK); err != nil {
		t.Fatal(err)
	}
}

// creates image into defined filename, with dimensions wMax:hMax, and a grayscale gradient
// if dir == DIRECTION_UP, gradient is from top-left to right-bottom
// if dir == DIRECTION_DOWN, gradient is from right-bottom to top-left
func createTestImage(filename string, wMax, hMax int, dir direction) {
	resultImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{wMax, hMax}})

	// prepare test image
	for y := 0; y < hMax; y++ {
		for x := 0; x < wMax; x++ {
			var level uint8
			var levelFloat = float64(x+y) / float64(wMax+hMax)
			if dir == DIRECTION_UP {
				level = uint8(math.MaxUint8 * levelFloat)
			} else {
				level = uint8(math.MaxUint8 * (1 - levelFloat))
			}
			newColor := color.RGBA{
				R: level,
				G: level,
				B: level,
				A: math.MaxUint8,
			}
			resultImage.Set(x, y, newColor)
		}
	}

	// save test image
	newfile, err := os.Create(filename)
	if err != nil {
		log.Printf("failed creating %s: %s", filename, err)
		panic(err.Error())
	}
	defer newfile.Close()
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
		BufferPool:       nil,
	}
	err = encoder.Encode(newfile, resultImage)
	if err != nil {
		log.Fatal(err)
	}
}
