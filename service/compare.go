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

type MatchMode int

const (
	// keep matching pixels in original color
	MATCH_MODE_KEEP_ORIGINAL MatchMode = iota
	// replace matching pixels with black
	MATCH_MODE_BLACK
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

func (c *Comparer) Compare(mode MatchMode) error {
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
				// define the resulting color depending on the match mode
				var matchedColor uint8
				if mode == MATCH_MODE_KEEP_ORIGINAL {
					matchedColor = uint8(r1 >> 8)
				} else {
					matchedColor = 0
				}
				newColor = color.RGBA{
					R: matchedColor,
					G: matchedColor,
					B: matchedColor,
					A: math.MaxUint8,
				}
			} else if r1 < r2 {
				// normalize r2 - r1 difference to math.MaxUint8
				normalized := uint8(math.MaxUint8 * float64((r2-r1)>>8) / float64(math.MaxUint8))

				//highlight pixels with green if second image pixel is brighter
				newColor = color.RGBA{
					R: 0,
					G: normalized,
					B: 0,
					A: math.MaxUint8,
				}
			} else {
				// normalize r1 - r2 difference to math.MaxUint8
				normalized := uint8(math.MaxUint8 * float64((r1-r2)>>8) / float64(math.MaxUint8))

				//highlight pixels with red if second image pixel is darker
				newColor = color.RGBA{
					R: normalized,
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

// saves the generated result image into provided output filepath
func (c *Comparer) saveResult(resultImage image.Image) error {
	// save the result image file
	newFile, err := os.Create(c.outputFile)
	if err != nil {
		log.Printf("failed creating %s: %s", c.outputFile, err)
		panic(err.Error())
	}
	defer newFile.Close()

	// create the png encoder with best compression
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
		BufferPool:       nil,
	}
	// write the output file
	err = encoder.Encode(newFile, resultImage)
	if err != nil {
		return errors.Wrap(err, "can't save the output file")
	}
	return nil
}

// Performs input files validation
// Returns error in case of
// - can't detect input images format
// - provided images' mime type is not supported
// - error when loading input images
// - output path not valid
// - output file is not a .png etension
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

// detects mime type for the image at provided filename
func getImageFormat(filename string) (string, error) {
	mime, _, err := mimetype.DetectFile(filename)
	if err != nil {
		return "", errors.Wrap(err, "can't detect image's type")
	} else {
		return mime, nil
	}
}

// checks whether mime type is supported
func isMimeSupported(mime string) bool {
	for _, v := range supportedMimes {
		if v == mime {
			return true
		}
	}
	return false
}
