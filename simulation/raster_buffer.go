package simulation

import (
	"Netron1-Go/api"
	"image"
	"image/color"
)

// RasterBuffer provides a memory mapped RGBA and Z buffer
// This buffer must be blitted to another buffer, for example,
// PNG or display buffer (like SDL).
type RasterBuffer struct {
	width  int
	height int

	// Image pixel buffers
	pixels1    *image.RGBA
	pixels2    *image.RGBA
	bounds     image.Rectangle
	bufferFlag bool // false = pixels1, true = pixels2

	alphaBlending bool

	// Pen colors
	ClearColor color.RGBA
	PixelColor color.RGBA
}

// NewRasterBuffer creates a display buffer
func NewRasterBuffer(width, height int) api.IRasterBuffer {
	o := new(RasterBuffer)
	o.width = width
	o.height = height

	o.alphaBlending = false
	o.bufferFlag = false

	o.bounds = image.Rect(0, 0, width, height)
	o.pixels1 = image.NewRGBA(o.bounds)
	o.pixels2 = image.NewRGBA(o.bounds)

	o.ClearColor.R = 127
	o.ClearColor.G = 127
	o.ClearColor.B = 127
	o.ClearColor.A = 255

	return o
}

// Width of buffer
func (rb *RasterBuffer) Width() int {
	return rb.width
}

// Height of buffer
func (rb *RasterBuffer) Height() int {
	return rb.height
}

// EnableAlphaBlending turns on/off per pixel alpha blending
func (rb *RasterBuffer) EnableAlphaBlending(enable bool) {
	rb.alphaBlending = enable
}

// Pixels returns underlying color buffer
func (rb *RasterBuffer) Pixels() *image.RGBA {
	if rb.bufferFlag {
		return rb.pixels2
	} else {
		return rb.pixels1
	}
}

func (rb *RasterBuffer) BackPixels() *image.RGBA {
	if !rb.bufferFlag {
		return rb.pixels2
	} else {
		return rb.pixels1
	}
}

func (rb *RasterBuffer) Swap() {
	rb.bufferFlag = !rb.bufferFlag
}

// Clear clears both color buffers
func (rb *RasterBuffer) Clear() {
	buf := rb.pixels1
	if rb.bufferFlag {
		buf = rb.pixels2
	}

	for y := 0; y < rb.height; y++ {
		for x := 0; x < rb.width; x++ {
			buf.SetRGBA(x, y, rb.ClearColor)
			buf.SetRGBA(x, y, rb.ClearColor)
		}
	}
}

// ClearColorBuffer clears only the color/pixel buffer
func (rb *RasterBuffer) ClearColorBuffer() {
	buf := rb.pixels1
	if rb.bufferFlag {
		buf = rb.pixels2
	}
	/// TODO use image/draw to clear using a SRC
	for y := 0; y < rb.height; y++ {
		for x := 0; x < rb.width; x++ {
			buf.SetRGBA(x, y, rb.ClearColor)
			buf.SetRGBA(x, y, rb.ClearColor)
		}
	}
}

// SetPixel sets a pixel in the buffer
func (rb *RasterBuffer) SetPixel(x, y int) {
	if x < 0 || x > rb.width || y < 0 || y > rb.height {
		return
	}

	buf := rb.pixels1
	if rb.bufferFlag {
		buf = rb.pixels2
	}

	// https://en.wikipedia.org/wiki/Alpha_compositing Alpha blending section
	// Non premultiplied alpha
	if rb.alphaBlending {
		dst := buf.RGBAAt(x, y)
		src := rb.PixelColor
		A := float32(src.A) / 255.0
		dst.R = uint8(float32(src.R)*A + float32(dst.R)*(1.0-A))
		dst.G = uint8(float32(src.G)*A + float32(dst.G)*(1.0-A))
		dst.B = uint8(float32(src.B)*A + float32(dst.B)*(1.0-A))
		dst.A = 255
		buf.SetRGBA(x, y, dst)
	} else {
		buf.SetRGBA(x, y, rb.PixelColor)
	}
}

// GetPixel get pixel
func (rb *RasterBuffer) GetPixel(x, y int) color.RGBA {
	buf := rb.pixels1
	if rb.bufferFlag {
		buf = rb.pixels2
	}
	return buf.RGBAAt(x, y)
}

// SetPixelColor set the current pixel color and sets the pixel
// using SetPixel()
func (rb *RasterBuffer) SetPixelColor(c color.RGBA) {
	rb.PixelColor = c
}

// SetClearColor set the clear buffer
func (rb *RasterBuffer) SetClearColor(c color.RGBA) {
	rb.ClearColor = c
}
