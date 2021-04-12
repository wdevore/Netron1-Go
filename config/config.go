package config

import (
	"Netron1-Go/api"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// API is the runtime configuration
var API api.IConfig

// The JSON data structure
type configJSON struct {
	ErrLog    string `json:"ErrLog"`
	InfoLog   string `json:"InfoLog"`
	ExitState string `json:"ExitState"`
	LogRoot   string `json:"LogRoot"`
	DataRoot  string `json:"DataRoot"`
}

type configuration struct {
	dirty      bool
	conf       configJSON
	path       string
	configFile string
}

// NewConfig construct an IConfig object
func NewConfig(configFile string) (api.IConfig, error) {
	o := new(configuration)
	o.dirty = false
	o.configFile = configFile

	// dir, err := filepath.Abs(filepath.Dir(""))
	confPath, err := filepath.Abs("")
	if err != nil {
		return nil, err
	}

	o.path = confPath

	jsonFile, err := os.Open(confPath + "/" + configFile)
	if err != nil {
		log.Fatalln("ERROR:", err)
		return nil, err
	}

	defer jsonFile.Close()

	err = json.NewDecoder(jsonFile).Decode(&o.conf)

	if err != nil {
		log.Fatalln("JSON ERROR:", err)
		return nil, err
	}

	return o, nil
}

// Save persists the current config to json file.
func (c *configuration) Save() error {
	if !c.dirty {
		fmt.Print("nothing to save")
		return nil
	}

	indentedJSON, _ := json.MarshalIndent(c.conf, "", "  ")

	err := ioutil.WriteFile(c.path+"/"+c.configFile, indentedJSON, 0644)
	if err != nil {
		// log.Fatalln("ERROR:", err)
		return err
	}

	return nil
}

// ErrLogFileName is the name of the error log file.
func (c *configuration) ErrLogFileName() string {
	return c.conf.ErrLog
}

// InfoLogFileName is the name of the info log file.
func (c *configuration) InfoLogFileName() string {
	return c.conf.InfoLog
}

// LogRoot is the base path to where log files are located.
func (c *configuration) LogRoot() string {
	return c.conf.LogRoot
}

func (c *configuration) DataRoot() string {
	return c.conf.DataRoot
}

// ExitState indicates what the last state the
// simulation was in when deuron exited.
// Values:
//   Terminated = user quit simulation while it was inprogress
//   Completed = sim terminated on its own
//   Crashed = sim died
//   Paused = user paused simulation and exited
//   Exited = user exited when no simulation was running
func (c *configuration) ExitState() string {
	return c.conf.ExitState
}

// SetExitState sets a value upon deuron exit.
func (c *configuration) SetExitState(state string) {
	c.conf.ExitState = state
	c.dirty = true
}
