package api

import (
	"image"
	"image/color"
)

// IRasterBuffer api for color and depth buffer
type IRasterBuffer interface {
	EnableAlphaBlending(enable bool)

	Width() int
	Height() int

	Pixels() *image.RGBA
	BackPixels() *image.RGBA

	SetClearColor(c color.RGBA)
	Clear()
	Swap()

	SetPixelColor(c color.RGBA)
	SetPixel(x, y int)
	GetPixel(x, y int) color.RGBA
}
