package api

import "image/color"

// IModel represents a information flow model
type IModel interface {
	GetInfectedColor() color.RGBA
	GetSusceptibleColor() color.RGBA
	GetRemovedColor() color.RGBA

	Configure(rasterBuffer IRasterBuffer)
	Reset()
	Step() bool
}
