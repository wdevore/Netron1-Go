package api

// IConfig holds configuration and runtime properties.
type IConfig interface {
	ErrLogFileName() string
	InfoLogFileName() string
	LogRoot() string

	ExitState() string
	SetExitState(string)

	DataRoot() string
	Save() error
}
