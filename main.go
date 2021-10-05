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

type DudeGenerator struct {
	image  *gg.Context
	height int
	width  int
}

func main() {
	imager := DudeGenerator{}
	imager.GenerateBackground(4096, 4096)
	imager.GeneratePointerDudes()
}

func (c *DudeGenerator) GeneratePointerDudes() {
	// setup
	numDudes := 101

	// scaling
	dudeScalingFactor := .9

	// positioning
	horizDisplacementFactor := .2
	vertDisplacementFactor := .6
	heightOffset := 0
	widthOffset := c.width / 70

	// ptr dude
	originalDude := c.LoadImage(pointerGuyPath)
	dudeHeight := originalDude.Bounds().Max.X
	dudeWidth := originalDude.Bounds().Max.Y

	pointerString := "int"

	for dudeNum := numDudes - 1; dudeNum >= 0; dudeNum-- {
		// rescale dudes
		thisScale := math.Pow(dudeScalingFactor, float64(dudeNum+1))
		scaledDude := imaging.Resize(originalDude, int(float64(dudeHeight)*thisScale), int(float64(dudeWidth)*thisScale), imaging.ResampleFilter{})

		// increment offsets
		widthOffset += int(float64(scaledDude.Rect.Dx()) * horizDisplacementFactor)
		heightOffset += int(float64(scaledDude.Rect.Dy()) * vertDisplacementFactor)

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
	c.LoadFont(float64(c.height) / (20.0 * float64(numDudes)))
	c.image.SetRGB255(0x4A, 0xA9, 0xBC)

	strHorizontalBase := float64(c.width) * .45

	c.image.DrawStringAnchored("int", strHorizontalBase, float64(c.height)/20, .5, 0)

	// save image
	c.image.SavePNG(savePath)
}

func (c *DudeGenerator) LoadImage(str string) image.Image {
	img, err := gg.LoadImage(str)
	if err != nil {
		panic(err)
	}
	return img
}

func (c *DudeGenerator) GenerateBackground(width int, height int) {
	dc := gg.NewContext(width, height)

	// background
	bgTile := c.LoadImage(backgroundTilePath)
	bgTile = imaging.Resize(bgTile, width, height, imaging.ResampleFilter{})
	dc.DrawImage(bgTile, 0, 0)

	c.image = dc
	c.height = height
	c.width = width
}

func (c *DudeGenerator) LoadFont(size float64) {
	if err := c.image.LoadFontFace(boldFontPath, size); err != nil {
		panic(err)
	}
}
