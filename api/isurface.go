package api

// ISurface is the graph viewer
type ISurface interface {
	// Open(IHost)
	Open(IModel)
	Close()

	Run(chToSim, chFromSim chan string)
	Quit()

	SetFont(fontPath string, size int) error

	Update(bool)
	Raster() IRasterBuffer
}
