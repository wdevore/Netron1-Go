package api

import "github.com/veandco/go-sdl2/sdl"

// ISurface is the graph viewer
type ISurface interface {
	// Open(IHost)
	Open()
	Close()

	Run(chToSim, chFromSim chan string)
	Quit()

	Configure()
	SetFont(fontPath string, size int) error

	Raster() IRasterBuffer
	SetDrawColor(color sdl.Color)
	SetPixel(x, y int)
}
