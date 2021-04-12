package api

import (
	"image"
	"image/color"
)

// IRasterBuffer api for color and depth buffer
type IRasterBuffer interface {
	EnableAlphaBlending(enable bool)
	Pixels() *image.RGBA
	SetClearColor(c color.RGBA)
	Clear()
	SetPixel(x, y int)
	SetPixelColor(c color.RGBA)
}
