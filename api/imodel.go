package api

// ISimulation is simulation
type IModel interface {
	Properties() IProperties
	Configure(rasterBuffer IRasterBuffer)
	Reset()
	Step() bool
	SendEvent(string)
	Name() string
}
