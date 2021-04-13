package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"time"
)

type Simulation struct {
	debug     int
	loop      bool
	running   bool
	paused    bool
	completed bool

	raster api.IRasterBuffer

	model api.IModel
}

func NewSimulation() *Simulation {
	o := new(Simulation)
	o.running = false
	o.loop = true
	return o
}

func (s *Simulation) Initialize(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
}

// Boot is the simulation bootstrap. The simulation isn't
// running until told to do so.
func (s *Simulation) Start(inChan chan string, outChan chan string) {
	s.configure()

	for s.loop {
		select {
		case cmd := <-inChan:
			switch cmd {
			case "exit":
				if s.running {
					outChan <- "Terminated"
				} else {
					if !s.completed {
						outChan <- "Exited"
					}
				}
				s.loop = false
			case "run":
				if s.running {
					// Starting a run while the sim is already running
					// isn't allowed.
					outChan <- "Running"
					continue
				}

				s.running = true
				outChan <- "Started"
			case "pause":
				if !s.running {
					outChan <- "Not Running"
					continue
				}

				if s.paused {
					outChan <- "Already Paused"
					continue
				}

				outChan <- "Paused"
				s.paused = true
			case "resume":
				if !s.paused {
					outChan <- "Not Paused"
					continue
				}

				outChan <- "Resumed"
				s.paused = false
			case "reset":
				outChan <- "Reset"
				s.running = false
				s.reset()
			case "stop":
				outChan <- "Stopped"
				s.running = false
				s.completed = false
			case "status":
				outChan <- fmt.Sprintf("Status: %d", s.debug)
			}
		default:
			if !s.running {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			if s.paused {
			} else {
				// The sim is running, make a step
				s.running = s.model.Step()
				if !s.running {
					outChan <- "Complete"
					s.completed = true
				}
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *Simulation) configure() {
	s.model = NewSIRModel()
	s.model.Configure(s.raster)
}

func (s *Simulation) reset() {
	s.debug = 0
	s.paused = false
	s.completed = false

	s.model.Reset()
}

// func (s *Simulation) step() bool {
// 	return s.model.Step()
// 	// s.debug++
// 	// if s.debug > 100 {
// 	// 	return false
// 	// }

// 	// s.raster.SetPixel(s.debug, s.debug)
// 	// return true
// }
