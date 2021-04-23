package api

// IProperties defines the properties for the model
type IProperties interface {
	Width() int
	Height() int
	WindowPosX() int
	WindowPosY() int
	Scale() int
}
