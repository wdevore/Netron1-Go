package simulation

import (
	"Netron1-Go/api"
	"Netron1-Go/gui"
	"fmt"
	"image/color"
	"math/rand"
)

type SIRModel struct {
	infectedColor    color.RGBA // cell type = 1
	susceptibleColor color.RGBA // cell type = 2
	removedColor     color.RGBA // cell type = 3

	raster           api.IRasterBuffer
	cellStates       [][]int
	cellNextStates   [][]int
	transmissionRate float32
}

func NewSIRModel() api.IModel {
	o := new(SIRModel)
	// Infected = Blue
	o.infectedColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Susceptible = Skin
	o.susceptibleColor = color.RGBA{R: 255, G: 225, B: 200, A: 255}
	// Removed = Gray
	o.removedColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}

	return o
}

func (s *SIRModel) Name() string {
	return "SIRModel"
}

func (s *SIRModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.transmissionRate = 0.5
	rand.Seed(131)

	s.cellStates = make([][]int, s.raster.Height())
	s.cellNextStates = make([][]int, s.raster.Height())
	for i := range s.cellStates {
		s.cellStates[i] = make([]int, s.raster.Width())
		s.cellNextStates[i] = make([]int, s.raster.Width())
	}
}

// SendEvent receives an event from the host simulation
func (s *SIRModel) SendEvent(event string) {
}

func (s *SIRModel) Properties() api.IProperties {
	return gui.NewProperties(300, 300, 1500, 100, 1)
}

func (s *SIRModel) Reset() {
	fmt.Println(("--- sir reset ---"))
	s.raster.Clear()

	// Start with the center "cell" infected = Blue
	s.raster.SetPixelColor(s.infectedColor)

	w := s.raster.Width()
	h := s.raster.Height()

	cx := w / 2
	cy := h / 2

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cellStates[col][row] = 2     // Susceptible
			s.cellNextStates[col][row] = 0 // Undetermined
		}
	}

	s.cellStates[cx][cy] = 1 // Infected

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

func (s *SIRModel) Step() bool {
	// The current "step" works on the current-state but updates the next-state
	// Once done, the next-state is copied back to the current-state.
	w := s.raster.Width()
	h := s.raster.Height()
	ce := 0
	infected := 0

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			// fmt.Println("row-col: ", row, ",", col, " = ", s.cellStates[col][row])
			if s.cellStates[col][row] == 1 { // Is infected
				// Transition from Infected to Removed
				s.cellNextStates[col][row] = 3

				// Check neighbors.
				ce = col + 1
				if ce < w { // Right
					// If neighbor isn't (infected AND they arent removed) = Suceptible
					// then randomize against transmission rate
					if s.cellStates[ce][row] != 3 {
						chance := rand.Float32()
						if chance < s.transmissionRate {
							// Cell is now infected
							s.cellNextStates[ce][row] = 1
							infected++
						}
					}
				}

				ce = col - 1
				if ce >= 0 { // Left
					if s.cellStates[ce][row] != 3 {
						chance := rand.Float32()
						if chance < s.transmissionRate {
							// Cell is now infected
							s.cellNextStates[ce][row] = 1
							infected++
						}
					}
				}

				ce = row - 1
				if ce >= 0 { // Top
					if s.cellStates[col][ce] != 3 {
						chance := rand.Float32()
						if chance < s.transmissionRate {
							// Cell is now infected
							s.cellNextStates[col][ce] = 1
							infected++
						}
					}
				}

				ce = row + 1
				if ce < h { // Bottom
					if s.cellStates[col][ce] != 3 {
						chance := rand.Float32()
						if chance < s.transmissionRate {
							// Cell is bottom infected
							s.cellNextStates[col][ce] = 1
							infected++
						}
					}
				}
			} else {
				if s.cellStates[col][row] == 0 {
					// This cell is now determined to susceptible
					s.cellNextStates[col][row] = 2
				}
			}
		}
	}

	// Copy next-state to current-state
	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cellStates[col][row] = s.cellNextStates[col][row]
			if s.cellStates[col][row] == 1 {
				s.raster.SetPixelColor(s.infectedColor)
			} else if s.cellStates[col][row] == 2 {
				s.raster.SetPixelColor(s.susceptibleColor)
			} else {
				s.raster.SetPixelColor(s.removedColor)
			}
			s.raster.SetPixel(col, row)
		}
	}

	// fmt.Println("Newly infected: ", infected)
	return infected > 0
}
