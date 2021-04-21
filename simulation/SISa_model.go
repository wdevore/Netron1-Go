package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"image/color"
	"math/rand"
)

type SISaModel struct {
	undetermenedColor color.RGBA // cell type = 0
	infectedColor     color.RGBA // cell type = 1
	susceptibleColor  color.RGBA // cell type = 2
	removedColor      color.RGBA // cell type = 3

	raster         api.IRasterBuffer
	cellStates     [][]int
	cellNextStates [][]int

	acceptibleRate float64

	// The chance that the person picks up meditation again
	repeatRate float64
	// The chance they will drop meditation.
	dropRate float64
	// The chance they will try meditation
	pickupRate float64

	// The chance that someone spontaneously starts meditating.
	spontaneousRate float64
}

func NewSISaModel() api.IModel {
	o := new(SISaModel)
	o.undetermenedColor = color.RGBA{R: 200, G: 255, B: 200, A: 255}
	// Infected = Blue
	o.infectedColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Susceptible = Skin
	o.susceptibleColor = color.RGBA{R: 255, G: 225, B: 200, A: 255}
	// Removed = Gray
	o.removedColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}

	return o
}

func (s *SISaModel) Name() string {
	return "SISaModel"
}

func (s *SISaModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.acceptibleRate = 0.26
	s.repeatRate = 0.2
	s.dropRate = 0.9
	s.pickupRate = 0.5
	s.spontaneousRate = 0.5

	rand.Seed(13163)

	s.cellStates = make([][]int, s.raster.Height())
	s.cellNextStates = make([][]int, s.raster.Height())
	for i := range s.cellStates {
		s.cellStates[i] = make([]int, s.raster.Width())
		s.cellNextStates[i] = make([]int, s.raster.Width())
	}
}

// SendEvent receives an event from the host simulation
func (s *SISaModel) SendEvent(event string) {
}

func (s *SISaModel) Reset() {
	fmt.Println(("--- sir reset ---"))
	s.raster.Clear()

	// Start with the center "cell" infected = Blue
	s.raster.SetPixelColor(s.infectedColor)

	w := s.raster.Width()
	h := s.raster.Height()

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cellStates[col][row] = 2     // Susceptible
			s.cellNextStates[col][row] = 0 // Undetermined
		}
	}

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cellStates[col][row] == 1 {
				s.raster.SetPixelColor(s.infectedColor)
			} else if s.cellStates[col][row] == 3 {
				s.raster.SetPixelColor(s.removedColor)
			} else {
				s.raster.SetPixelColor(s.susceptibleColor)
			}
			s.raster.SetPixel(col, row)
		}
	}
}

func (s *SISaModel) Step() bool {
	// The current "step" works on the current-state but updates the next-state
	// Once done, the next-state is copied back to the current-state.
	w := s.raster.Width()
	h := s.raster.Height()
	ce := 0
	infected := 0

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cellStates[col][row] == 1 { // Is infected
				// How likely will they drop it.
				if rand.Float64() < s.dropRate {
					s.cellNextStates[col][row] = 2 // Suceptible
				}

				// Check neighbors.
				ce = col + 1
				if ce < w { // Right
					// If neighbor isn't (infected AND they are acceptible) = Suceptible
					if s.cellStates[ce][row] != 3 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is now infected
							s.cellNextStates[ce][row] = 1
							infected++
						}
					}
				}

				ce = col - 1
				if ce >= 0 { // Left
					if s.cellStates[ce][row] != 3 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is now infected
							s.cellNextStates[ce][row] = 1
							infected++
						}
					}
				}

				ce = row - 1
				if ce >= 0 { // Top
					if s.cellStates[col][ce] != 3 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is now infected
							s.cellNextStates[col][ce] = 1
							infected++
						}
					}
				}

				ce = row + 1
				if ce < h { // Bottom
					if s.cellStates[col][ce] != 3 {
						if rand.Float64() < s.acceptibleRate {
							// Cell is bottom infected
							s.cellNextStates[col][ce] = 1
							infected++
						}
					}
				}
			} else {
				if s.cellStates[col][row] == 0 {
					// This person hasn't experienced meditation yet.
					if rand.Float64() < s.pickupRate {
						s.cellNextStates[col][row] = 2 // Suceptible
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
		s.cellNextStates[col][row] = 1
		infected++
	}

	// Copy next-state to current-state
	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cellStates[col][row] = s.cellNextStates[col][row]
			if s.cellStates[col][row] == 1 {
				s.raster.SetPixelColor(s.infectedColor)
			} else if s.cellStates[col][row] == 2 {
				s.raster.SetPixelColor(s.susceptibleColor)
			} else if s.cellStates[col][row] == 3 {
				s.raster.SetPixelColor(s.removedColor)
			} else {
				s.raster.SetPixelColor(s.undetermenedColor)
			}
			s.raster.SetPixel(col, row)
		}
	}

	// fmt.Println("Newly infected: ", infected)
	return true //infected > 0
}
