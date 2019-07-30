package service

import (
	"image"
	"image/png"
	"io"
	"os"
)

type IMG struct {
	mime  string
	image image.Image
}

func NewIMG(mime string, filename string) (*IMG, error) {
	img := new(IMG)
	img.mime = mime
	var reader io.Reader
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	img.image, err = png.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}
