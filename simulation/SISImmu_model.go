package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"image/color"
	"math/rand"
)

type SISimmuModel struct {
	undetermenedColor color.RGBA // cell type = 0
	infectedColor     color.RGBA // cell type = 1
	susceptibleColor  color.RGBA // cell type = 2
	removedColor      color.RGBA // cell type = 3

	raster api.IRasterBuffer
	// cellStates     [][]int
	// cellNextStates [][]int
	cells [][]Cell

	acceptibleRate float64

	// The chance that the person picks up meditation again
	repeatRate float64
	// The chance they will drop meditation.
	dropRate float64
	// The chance they will try meditation
	pickupRate float64

	// The chance that someone spontaneously starts meditating.
	spontaneousRate float64

	// The chance the person has no interest at all
	immunityRate float64
}

func NewSISimmuModel() api.IModel {
	o := new(SISimmuModel)
	o.undetermenedColor = color.RGBA{R: 200, G: 255, B: 200, A: 255}
	// Infected = Blue
	o.infectedColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Susceptible = Skin
	o.susceptibleColor = color.RGBA{R: 255, G: 225, B: 200, A: 255}
	// Removed = Gray
	o.removedColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}

	return o
}

func (s *SISimmuModel) Name() string {
	return "SISimmuModel"
}

func (s *SISimmuModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.acceptibleRate = 0.26 // 26 = below threshold
	s.repeatRate = 0.2
	s.dropRate = 0.9
	s.pickupRate = 0.5
	s.spontaneousRate = 0.25
	s.immunityRate = 0.01

	rand.Seed(13163)

	s.cells = make([][]Cell, s.raster.Height())
	for i := range s.cells {
		s.cells[i] = make([]Cell, s.raster.Width())
	}
}

// SendEvent receives an event from the host simulation
func (s *SISimmuModel) SendEvent(event string) {
}

func (s *SISimmuModel) Reset() {
	fmt.Println(("--- sir reset ---"))
	s.raster.Clear()

	w := s.raster.Width()
	h := s.raster.Height()

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			// Is this person have no interest ever
			if rand.Float64() < s.immunityRate {
				s.cells[col][row].state = 3 // Immune
				s.cells[col][row].nextState = 3
			} else {
				s.cells[col][row].state = 2 // Susceptible
				s.cells[col][row].nextState = 2
			}
		}
	}

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cells[col][row].state == 1 {
				s.raster.SetPixelColor(s.infectedColor)
			} else if s.cells[col][row].state == 3 {
				s.raster.SetPixelColor(s.removedColor)
			} else {
				s.raster.SetPixelColor(s.susceptibleColor)
			}
			s.raster.SetPixel(col, row)
		}
	}
}

func (s *SISimmuModel) Step() bool {
	// The current "step" works on the current-state but updates the next-state
	// Once done, the next-state is copied back to the current-state.
	w := s.raster.Width()
	h := s.raster.Height()
	ce := 0
	infected := 0

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cells[col][row].state == 1 { // Is infected
				// How likely will they drop it.
				if rand.Float64() < s.dropRate {
					s.cells[col][row].nextState = 2 // Suceptible
				}

				// Check neighbors.
				ce = col + 1
				if ce < w { // Right
					// If neighbor isn't (infected AND they are acceptible) = Suceptible
					if s.cells[ce][row].state == 2 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is now infected
							s.cells[ce][row].nextState = 1
							infected++
						}
					}
				}

				ce = col - 1
				if ce >= 0 { // Left
					if s.cells[ce][row].state == 2 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is now infected
							s.cells[ce][row].nextState = 1
							infected++
						}
					}
				}

				ce = row - 1
				if ce >= 0 { // Top
					if s.cells[col][ce].state == 2 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is now infected
							s.cells[col][ce].nextState = 1
							infected++
						}
					}
				}

				ce = row + 1
				if ce < h { // Bottom
					if s.cells[col][ce].state == 2 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is bottom infected
							s.cells[col][ce].nextState = 1
							infected++
						}
					}
				}
			}
		}
	}

	if rand.Float64() < s.spontaneousRate {
		// pick a col and row to infect
		col := int(rand.Float64()*float64(s.raster.Width())) - 1
		if col < 0 {
			col = 0
		}
		row := int(rand.Float64()*float64(s.raster.Width())) - 1
		if row < 0 {
			row = 0
		}
		if s.cells[col][row].state != 3 {
			s.cells[col][row].nextState = 1
			infected++
		}
	}

	// Copy next-state to current-state
	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cells[col][row].state = s.cells[col][row].nextState
			if s.cells[col][row].state == 1 {
				s.raster.SetPixelColor(s.infectedColor)
			} else if s.cells[col][row].state == 2 {
				s.raster.SetPixelColor(s.susceptibleColor)
			} else if s.cells[col][row].state == 3 {
				s.raster.SetPixelColor(s.removedColor)
			}
			s.raster.SetPixel(col, row)
		}
	}

	// fmt.Println("Newly infected: ", infected)
	return true //infected > 0
}
