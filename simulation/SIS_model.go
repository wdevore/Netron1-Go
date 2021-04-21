package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"image/color"
	"math/rand"
)

type SISModel struct {
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
}

func NewSISModel() api.IModel {
	o := new(SISModel)
	o.undetermenedColor = color.RGBA{R: 200, G: 255, B: 200, A: 255}
	// Infected = Blue
	o.infectedColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Susceptible = Skin
	o.susceptibleColor = color.RGBA{R: 255, G: 225, B: 200, A: 255}
	// Removed = Gray
	o.removedColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}

	return o
}

func (s *SISModel) Name() string {
	return "SISModel"
}

func (s *SISModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.acceptibleRate = 0.28
	s.repeatRate = 0.2
	s.dropRate = 0.9
	s.pickupRate = 0.5

	rand.Seed(131)

	s.cellStates = make([][]int, s.raster.Height())
	s.cellNextStates = make([][]int, s.raster.Height())
	for i := range s.cellStates {
		s.cellStates[i] = make([]int, s.raster.Width())
		s.cellNextStates[i] = make([]int, s.raster.Width())
	}
}

func (s *SISModel) Reset() {
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
	s.cellStates[cx][cy+1] = 1
	s.cellStates[cx][cy+2] = 1
	s.cellStates[cx][cy-1] = 1
	s.cellStates[cx][cy-2] = 1

	s.cellStates[cx-1][cy] = 1
	s.cellStates[cx-1][cy+1] = 1
	s.cellStates[cx-1][cy+2] = 1
	s.cellStates[cx-1][cy-1] = 1
	s.cellStates[cx-1][cy-2] = 1

	s.cellStates[cx-2][cy] = 1
	s.cellStates[cx-2][cy+1] = 1
	s.cellStates[cx-2][cy+2] = 1
	s.cellStates[cx-2][cy-1] = 1
	s.cellStates[cx-2][cy-2] = 1

	s.cellStates[cx+1][cy] = 1
	s.cellStates[cx+1][cy+1] = 1
	s.cellStates[cx+1][cy+2] = 1
	s.cellStates[cx+1][cy-1] = 1
	s.cellStates[cx+1][cy-2] = 1

	s.cellStates[cx+2][cy] = 1
	s.cellStates[cx+2][cy+1] = 1
	s.cellStates[cx+2][cy+2] = 1
	s.cellStates[cx+2][cy-1] = 1
	s.cellStates[cx+2][cy-2] = 1

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

// SendEvent receives an event from the host simulation
func (s *SISModel) SendEvent(event string) {
}

func (s *SISModel) Step() bool {
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
	return infected > 0
}
