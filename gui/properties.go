package gui

import "Netron1-Go/api"

const (
	// SurfaceScale scales the view
	SurfaceScale = 300
	width        = SurfaceScale
	height       = SurfaceScale
	windowPosX   = 1500
	windowPosY   = 100
	fps          = 30.0
	framePeriod  = 1.0 / fps * 1000.0
)

type Properties struct {
	width, height, windowPosX, windowPosY int
	scale                                 int
}

func NewProperties(w, h, wpx, wpy, scale int) api.IProperties {
	return &Properties{width: w, height: h, windowPosX: wpx, windowPosY: wpy, scale: scale}
}

func (p *Properties) Width() int {
	return p.width
}

func (p *Properties) Height() int {
	return p.height
}

func (p *Properties) WindowPosX() int {
	return p.windowPosX
}

func (p *Properties) WindowPosY() int {
	return p.windowPosY
}

func (p *Properties) Scale() int {
	return p.scale
}
