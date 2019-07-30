package service

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/gabriel-vasile/mimetype"
)

type Comparer struct {
	imageFile1 string
	imageFile2 string
	image1     *IMG
	image2     *IMG
	outputFile string
}

// list of supported mime type (can be used for future extension of functionality)
var supportedMimes = []string{"image/png"}

//var supportedMimes = []string{"image/png", "image/jpeg", "image/jpg"}

// Comparer constructor. Performs images validation and prepares them for comparison
func NewComparer(image1 string, image2 string, output string) (*Comparer, error) {
	c := Comparer{imageFile1: image1, imageFile2: image2, outputFile: output}
	err := c.validateImages()
	if err != nil {
		return nil, err
	}
	return &c, err
}

func (c *Comparer) Compare() error {
	log.Println("compare start")

	bounds := c.image1.image.Bounds()
	wMin, hMin := bounds.Min.X, bounds.Min.Y
	wMax, hMax := bounds.Max.X, bounds.Max.Y

	resultImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{wMax, hMax}})

	// run loop through all pixels in original images
	for y := hMin; y < hMax; y++ {
		for x := wMin; x < wMax; x++ {

			// Since the input image is grayscale we can assume that all RGB are equal
			// Perhaps need to add validation for non-greyscale images
			r1, _, _, _ := c.image1.image.At(x, y).RGBA()
			r2, _, _, _ := c.image2.image.At(x, y).RGBA()

			var newColor color.RGBA
			if r1 == r2 {
				// keep the color if it's the same
				newColor = color.RGBA{
					R: uint8(r1 >> 8),
					G: uint8(r1 >> 8),
					B: uint8(r1 >> 8),
					A: math.MaxUint8,
				}
			} else if r1 < r2 {
				//highlight pixels with green if second image pixel is brighter
				newColor = color.RGBA{
					R: 0,
					G: uint8((r2 - r1) >> 8),
					B: 0,
					A: math.MaxUint8,
				}
			} else {
				//highlight pixels with red if second image pixel is darker
				newColor = color.RGBA{
					R: uint8((r1 - r2) >> 8),
					G: 0,
					B: 0,
					A: math.MaxUint8,
				}
			}
			// set the defined color of a pixel
			resultImage.Set(x, y, newColor)

		}
	}

	// save the result image file
	err := c.saveResult(resultImage)
	if err != nil {
		return err
	} else {
		log.Printf("result saved successfully into %s", c.outputFile)
		return nil
	}
}

func (c *Comparer) saveResult(resultImage image.Image) error {
	// save the result image file
	newFile, err := os.Create(c.outputFile)
	if err != nil {
		log.Printf("failed creating %s: %s", c.outputFile, err)
		panic(err.Error())
	}
	defer newFile.Close()
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
		BufferPool:       nil,
	}
	err = encoder.Encode(newFile, resultImage)
	if err != nil {
		return errors.Wrap(err, "can't save the result file")
	}
	return nil
}

func (c *Comparer) validateImages() error {

	// validate image 1
	mime1, err := getImageFormat(c.imageFile1)
	if err != nil {
		return err
	}

	if !isMimeSupported(mime1) {
		return errors.New("image 1 mime not supported")
	}

	// validate image 2
	mime2, err := getImageFormat(c.imageFile2)
	if err != nil {
		return err
	}

	if !isMimeSupported(mime2) {
		return errors.New("image 2 mime not supported")
	}

	c.image1, err = NewIMG(mime1, c.imageFile1)
	if err != nil {
		return errors.Wrap(err, "parse image 1:")
	}

	c.image2, err = NewIMG(mime2, c.imageFile2)
	if err != nil {
		return errors.Wrap(err, "parse image 2:")
	}

	if c.image1.image.Bounds() != c.image2.image.Bounds() {
		return errors.New("image dimensions are not the same")
	}

	// check if the output filename path exists
	if _, err = os.Stat(filepath.Dir(c.outputFile)); os.IsNotExist(err) {
		return err
	}

	// validate the output file format
	if filepath.Ext(c.outputFile) != ".png" {
		return errors.New("output must be a .png file")
	}
	return nil
}

func getImageFormat(filename string) (string, error) {
	mime, _, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", errors.Wrap(err, "can't detect image's type")
	} else {
		return mime, nil
	}
}

func isMimeSupported(mime string) bool {
	for _, v := range supportedMimes {
		if v == mime {
			return true
		}
	}
	return false
}
