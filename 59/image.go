package main

import (
	"code.google.com/p/go-tour/pic"
	"image"
	"image/color"
	"math"
)

type Image struct{}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 256, 256)
}

func avg(x, y int) uint8 {
	return uint8((x + y) / 2)
}

func plus(x, y int) uint8 {
	return uint8(x * y)
}

func pow(x, y int) uint8 {
	return uint8(math.Pow(float64(x), float64(y)))
}

func (i Image) At(x, y int) color.Color {
	//v := avg(x, y)
	//v := plus(x, y)
	v := pow(x, y)
	return color.RGBA{v, v, 255, 255}
}

func main() {
	m := Image{}
	pic.ShowImage(m)
}
