package api

// ISimulation is simulation host
type ISimulation interface {
	Initialize(rasterBuffer IRasterBuffer, surface ISurface)
	Configure(model IModel)
	Start(inChan chan string, outChan chan string)
}
