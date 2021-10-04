package main

import (
	"image"
	_ "image/png"
	"math"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
)

const outPath = "./out"
const savePath = outPath + "/image.png"

const assetsPath = "./assets"
const boldFontPath = assetsPath + "/Montserrat-Bold.ttf"
const pointerGuyPath = assetsPath + "/pointer-dude.jpg"
const backgroundTilePath = assetsPath + "/bg.png"

type Imager struct {
	image  *gg.Context
	height int
	width  int
}

func main() {
	imager := Imager{}
	imager.GenerateBackground(4096, 4096)
	imager.GeneraterPointerDudes()
}

func (c *Imager) GeneraterPointerDudes() {
	// setup
	numDudes := 101
	horizDisplacementFactor := .2
	vertIncrementFactor := .55
	dudeScalingFactor := .9
	fontIncrementFactor := .95

	// ptr dude
	originalDude := c.LoadImage(pointerGuyPath)
	dudeHeight := originalDude.Bounds().Max.Y
	dudeWidth := originalDude.Bounds().Max.Y

	// positioning
	heightOffset := 0
	widthOffset := 0

	pointerString := "int"

	for dudeNum := numDudes; dudeNum >= 0; dudeNum-- {
		// rescale dudes
		thisScale := math.Pow(dudeScalingFactor, float64(dudeNum+1))
		scaledDude := imaging.Resize(originalDude, int(float64(dudeHeight)*thisScale), int(float64(dudeWidth)*thisScale), imaging.ResampleFilter{})

		// increment offsets
		widthOffset += int(float64(scaledDude.Rect.Dx()) * horizDisplacementFactor)
		heightOffset += int(float64(scaledDude.Rect.Dy()) * vertIncrementFactor)

		// calculate vertical and horizontal positions
		hPos := c.width / 2
		vPos := heightOffset

		// if on the right side
		if dudeNum%2 != numDudes%2 {
			hPos += widthOffset
		} else {
			scaledDude = imaging.FlipH(scaledDude)
			hPos -= widthOffset
		}

		// draw dude
		c.image.DrawImageAnchored(scaledDude, hPos, vPos, .5, 0)

		// draw pointer
		pointerString += "*"

		// calculate font size based on height of dude
		fontSize := float64(400) * (float64(scaledDude.Rect.Dx()) / float64(c.height))
		c.LoadFont(fontSize)
		c.image.SetRGB255(0x4A, 0xA9, 0xBC)
		c.image.DrawStringAnchored(pointerString, float64(hPos), float64(vPos)-5, .5, 0)
	}

	// draw the 'int'
	c.LoadFont(17 * math.Pow(fontIncrementFactor, float64(numDudes)))
	c.image.SetRGB255(0x4A, 0xA9, 0xBC)
	c.image.DrawStringAnchored(pointerString, float64(c.width/2), 10, .5, 0)

	// save image
	c.image.SavePNG(savePath)
}

func (c *Imager) LoadImage(str string) image.Image {
	img, err := gg.LoadImage(str)
	if err != nil {
		panic(err)
	}
	return img
}

func (c *Imager) GenerateBackground(width int, height int) {
	dc := gg.NewContext(width, height)

	// background
	bgTile := c.LoadImage(backgroundTilePath)
	bgTile = imaging.Resize(bgTile, width, height, imaging.ResampleFilter{})
	dc.DrawImage(bgTile, 0, 0)

	c.image = dc
	c.height = height
	c.width = width
}

func (c *Imager) LoadFont(size float64) {
	if err := c.image.LoadFontFace(boldFontPath, size); err != nil {
		panic(err)
	}
}
