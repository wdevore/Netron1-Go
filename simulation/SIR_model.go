package simulation

import (
	"Netron1-Go/api"
	"image/color"
	"math/rand"
)

type SIRModel struct {
	infectedColor    color.RGBA // cell type = 1
	susceptibleColor color.RGBA // cell type = 2
	removedColor     color.RGBA // cell type = 3

	raster           api.IRasterBuffer
	cells            [][]int
	transmissionRate float32
}

func NewSIRModel() api.IModel {
	o := new(SIRModel)
	// Infected = Blue
	o.infectedColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Susceptible = Skin
	o.susceptibleColor = color.RGBA{R: 255, G: 219, B: 172, A: 255}
	// Removed = Gray
	o.removedColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}

	return o
}

func (s *SIRModel) GetInfectedColor() color.RGBA {
	return s.infectedColor
}

func (s *SIRModel) GetSusceptibleColor() color.RGBA {
	return s.susceptibleColor
}

func (s *SIRModel) GetRemovedColor() color.RGBA {
	return s.removedColor
}

func (s *SIRModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.transmissionRate = 0.45

	s.cells = make([][]int, s.raster.Height())
	for i := range s.cells {
		s.cells[i] = make([]int, s.raster.Width())
	}
}

func (s *SIRModel) Reset() {
	s.raster.Clear()

	// Start with the center "cell" infected = Blue
	s.raster.SetPixelColor(s.infectedColor)

	w := s.raster.Width()
	h := s.raster.Height()

	cx := w / 2
	cy := h / 2

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cells[col][row] = 2 // Susceptible
		}
	}

	s.raster.SetPixel(cx, cy)
	s.cells[cx][cy] = 1 // Infected
}

func (s *SIRModel) Step() bool {
	w := s.raster.Width()
	h := s.raster.Height()
	ce := 0

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cells[col][row] == 1 { // Is infected
				// Transition Infected to Removed
				s.cells[col][row] = 3

				// Check neighbors.
				ce = col + 1
				if ce < w { // Right
					// If neighbor isn't infected AND they arent removed
					// then randomize against transmission rate
					if s.cells[ce][row] != 3 {
						chance := rand.Float32()
						if chance > s.transmissionRate {
							// Cell is now infected
							s.cells[ce][row] = 1
						}
					}
				}

				ce = col - 1
				if ce >= 0 { // Left
					// If neighbor isn't (infected AND they arent removed) = Suceptible
					// then randomize against transmission rate
					if s.cells[ce][row] != 3 {
						chance := rand.Float32()
						if chance > s.transmissionRate {
							// Cell is now infected
							s.cells[ce][row] = 1
						}
					}
				}

				ce = row - 1
				if ce >= 0 { // Top
					// If neighbor isn't (infected AND they arent removed) = Suceptible
					// then randomize against transmission rate
					if s.cells[col][ce] != 3 {
						chance := rand.Float32()
						if chance > s.transmissionRate {
							// Cell is now infected
							s.cells[col][ce] = 1
						}
					}
				}

				ce = row + 1
				if ce < h { // Bottom
					// If neighbor isn't (infected AND they arent removed) = Suceptible
					// then randomize against transmission rate
					if s.cells[col][ce] != 3 {
						chance := rand.Float32()
						if chance > s.transmissionRate {
							// Cell is now infected
							s.cells[col][ce] = 1
						}
					}
				}
			}
		}
	}

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cells[col][row] == 1 {
				s.raster.SetPixelColor(s.infectedColor)
			} else if s.cells[col][row] == 2 {
				s.raster.SetPixelColor(s.susceptibleColor)
			} else {
				s.raster.SetPixelColor(s.removedColor)
			}
			s.raster.SetPixel(col, row)
		}
	}

	return true
}
