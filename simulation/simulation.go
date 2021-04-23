package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
	"time"
)

const gifOutputPath = "/media/RAMDisk/netron/"

type Simulation struct {
	debug     int
	loop      bool
	running   bool
	paused    bool
	completed bool

	raster  api.IRasterBuffer
	surface api.ISurface

	model api.IModel

	discName   string
	discNameId int
	gAni       *gif.GIF
	// Set to true if you want an animated gif generated.
	// Note: rendering will considerably slower.
	// To save gif use the "t" key to stop simulation. Once
	// stopped the gif is saved.
	enableGif bool
}

func NewSimulation() api.ISimulation {
	o := new(Simulation)
	o.running = false
	o.loop = true
	o.enableGif = false
	return o
}

func (s *Simulation) Initialize(rasterBuffer api.IRasterBuffer, surface api.ISurface) {
	s.raster = rasterBuffer
	s.surface = surface
}

// Boot is the simulation bootstrap. The simulation isn't
// running until told to do so.
func (s *Simulation) Start(inChan chan string, outChan chan string) {
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
				s.reset()
				s.surface.Update(true)
				outChan <- "Started"
			case "step":
				s.model.Step()
				if s.enableGif {
					s.addGif()
				}
				s.surface.Update(true)
				outChan <- "Stepped"
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
				s.surface.Update(true)
			case "stop":
				outChan <- "Stopped"
				s.running = false
				s.completed = false
				if s.enableGif {
					s.saveGif()
				}
			case "status":
				outChan <- fmt.Sprintf("Status: %d", s.debug)
			default:
				s.model.SendEvent(cmd)
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
				// Save image to disc
				if s.enableGif {
					s.addGif()
				}
				s.surface.Update(true)
				if !s.running {
					outChan <- "Complete"
					s.completed = true
				}
			}

			time.Sleep(1 * time.Millisecond)
		}
	}
}

func (s *Simulation) Configure(model api.IModel) {
	s.model = model //NewSISCityModel()
	s.model.Configure(s.raster)
	s.model.Reset()
	s.surface.Update(true)
}

func (s *Simulation) reset() {
	s.discNameId = 0
	s.debug = 0
	s.paused = false
	s.completed = false
	s.gAni = &gif.GIF{}
	s.gAni.Image = []*image.Paletted{}
	s.gAni.Delay = []int{}

	s.model.Reset()
}

func (s *Simulation) save() {
	s.saveGif()
}

func (s *Simulation) saveGif() {
	fn := fmt.Sprintf(gifOutputPath+"%s%d.gif", s.model.Name(), s.discNameId)
	f, err := os.Create(fn)

	if err != nil {
		log.Fatalln(err)
	}

	err = gif.EncodeAll(f, s.gAni)

	if err != nil {
		log.Fatalln(err)
	}
}

func (s *Simulation) addGif() {
	img := s.raster.Pixels()

	gifOps := &gif.Options{NumColors: 256, Drawer: draw.FloydSteinberg}
	pimg := image.NewPaletted(img.Bounds(), palette.Plan9[:gifOps.NumColors])
	if gifOps.Quantizer != nil {
		pimg.Palette = gifOps.Quantizer.Quantize(make(color.Palette, 0, gifOps.NumColors), img)
	}

	gifOps.Drawer.Draw(pimg, img.Bounds(), img, image.Point{})

	s.gAni.Image = append(s.gAni.Image, pimg)
	spf := 1 / float64(15)
	s.gAni.Delay = append(s.gAni.Delay, int(spf*100))

	// err = gif.Encode(f, s.raster.Pixels(), gifOps)

	s.discNameId++
}

func (s *Simulation) savePng() {
	fn := fmt.Sprintf("/media/RAMDisk/netron/%s%d.png", s.model.Name(), s.discNameId)
	f, err := os.Create(fn)

	if err != nil {
		log.Fatalln(err)
	}

	err = png.Encode(f, s.raster.Pixels())

	if err != nil {
		log.Fatalln(err)
	}
	s.discNameId++
}
