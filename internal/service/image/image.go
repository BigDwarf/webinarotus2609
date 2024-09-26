package image

import (
	"github.com/golang/freetype"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

type Service struct {
}

func NewImageService() *Service {
	return &Service{}
}

func (i *Service) DrawText(inputFileName, topText, bottomText string) (string, error) {
	imgFile, err := os.Open(inputFileName)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	fontBytes, err := os.ReadFile("/home/arslan/GolandProjects/webinar2609/roboto.ttf")
	if err != nil {
		return "", err
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return "", err
	}

	originalImage, _, err := image.Decode(imgFile)
	if err != nil {
		return "", err
	}
	// Create a new RGBA image
	rgba := image.NewRGBA(originalImage.Bounds())
	draw.Draw(rgba, rgba.Bounds(), originalImage, image.Point{}, draw.Src)

	// Create a freetype context for drawing text
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(30) // Adjust font size as needed
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.White))

	pt := freetype.Pt(10, 40)
	_, err = c.DrawString(topText, pt)
	if err != nil {
		return "", err
	}

	pt = freetype.Pt(10, rgba.Bounds().Dy()-10)
	_, err = c.DrawString(bottomText, pt)
	if err != nil {
		return "", err
	}

	outFile, err := os.CreateTemp("", "meme")
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, rgba, nil)
	if err != nil {
		return "", err
	}

	return outFile.Name(), nil
}
