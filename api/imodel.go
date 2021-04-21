package api

// IModel represents a information flow model
type IModel interface {
	Configure(rasterBuffer IRasterBuffer)
	Reset()
	Step() bool
	SendEvent(string)
	Name() string
}
