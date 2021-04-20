package api

// IModel represents a information flow model
type IModel interface {
	Configure(rasterBuffer IRasterBuffer)
	Reset()
	Step() bool

	Name() string
}
