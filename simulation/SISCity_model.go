package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"image/color"
	"math/rand"
)

// The city model is based on "zones". Each zone has an epic center
// that is the most connected and active.
// In the city more people tend to meditate verse the suburbs.

type SISCityModel struct {
	undetermenedColor color.RGBA // cell type = 0
	infectedColor     color.RGBA // cell type = 1
	susceptibleColor  color.RGBA // cell type = 2
	removedColor      color.RGBA // cell type = 3

	degree5Color color.RGBA
	degree6Color color.RGBA
	degree7Color color.RGBA
	degree8Color color.RGBA

	raster api.IRasterBuffer
	cells  [][]Cell

	acceptibleRate float64
	// The chance they will drop meditation.
	dropRate float64

	// degree goes from 4 to 8
	degree int
}

func NewSISCityModel() api.IModel {
	o := new(SISCityModel)
	o.undetermenedColor = color.RGBA{R: 200, G: 255, B: 200, A: 255}
	// Infected = Blue
	o.infectedColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	// Susceptible = Skin
	o.susceptibleColor = color.RGBA{R: 255, G: 225, B: 200, A: 255}
	// Removed = Gray
	o.removedColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}

	o.degree5Color = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	o.degree6Color = color.RGBA{R: 175, G: 175, B: 175, A: 255}
	o.degree7Color = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	o.degree8Color = color.RGBA{R: 125, G: 125, B: 125, A: 255}

	return o
}

func (s *SISCityModel) GetInfectedColor() color.RGBA {
	return s.infectedColor
}

func (s *SISCityModel) GetSusceptibleColor() color.RGBA {
	return s.susceptibleColor
}

func (s *SISCityModel) GetRemovedColor() color.RGBA {
	return s.removedColor
}

func (s *SISCityModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.acceptibleRate = 0.27
	s.dropRate = 0.9

	rand.Seed(131)

	s.cells = make([][]Cell, s.raster.Height())
	for i := range s.cells {
		s.cells[i] = make([]Cell, s.raster.Width())
	}
}

func (s *SISCityModel) Reset() {
	fmt.Println(("--- sir reset ---"))
	s.raster.Clear()

	w := s.raster.Width()
	h := s.raster.Height()

	// Create two zones. Each zone is a square.

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cells[col][row].state = 2     // Susceptible
			s.cells[col][row].nextState = 2 // Undetermined
		}
	}

	px := 250
	py := 250
	radius := 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].state = 1
		}
	}

	px = 210
	py = 210
	radius = 40
	// Create largest area first
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 5
		}
	}

	px += 5
	py += 5
	radius -= 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 6
		}
	}

	px += 5
	py += 5
	radius -= 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 7
		}
	}

	px += 5
	py += 5
	radius -= 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 8
		}
	}

	// ------------------------------------------------
	px = 260
	py = 260
	radius = 40
	// Create largest area first
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 5
		}
	}

	px += 5
	py += 5
	radius -= 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 6
		}
	}

	px += 5
	py += 5
	radius -= 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 7
		}
	}

	px += 5
	py += 5
	radius -= 10
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].degree = 8
		}
	}

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.drawCell(col, row)
		}
	}
}

func (s *SISCityModel) Step() bool {
	// The current "step" works on the current-state but updates the next-state
	// Once done, the next-state is copied back to the current-state.
	w := s.raster.Width()
	h := s.raster.Height()
	ce := 0
	infected := 0

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			switch s.cells[col][row].state {
			case 1:
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

				if s.cells[col][row].degree > 4 {
					// top/right
					ce = row - 1
					re := col + 1
					if ce >= 0 { // Top
						if s.cells[re][ce].state == 2 {
							if rand.Float64() < s.acceptibleRate {
								// Cell is now infected
								s.cells[col][ce].nextState = 1
								infected++
							}
						}
					}
				}
				if s.cells[col][row].degree > 5 {
					// bottom/right
					ce = row + 1
					re := col + 1
					if ce >= 0 { // Top
						if s.cells[re][ce].state == 2 {
							if rand.Float64() < s.acceptibleRate {
								// Cell is now infected
								s.cells[col][ce].nextState = 1
								infected++
							}
						}
					}
				}
				if s.cells[col][row].degree > 6 {
					// bottom/left
					ce = row + 1
					re := col - 1
					if ce >= 0 { // Top
						if s.cells[re][ce].state == 2 {
							if rand.Float64() < s.acceptibleRate {
								// Cell is now infected
								s.cells[col][ce].nextState = 1
								infected++
							}
						}
					}
				}
				if s.cells[col][row].degree > 7 {
					// top/left
					ce = row - 1
					re := col - 1
					if ce >= 0 { // Top
						if s.cells[re][ce].state == 2 {
							if rand.Float64() < s.acceptibleRate {
								// Cell is now infected
								s.cells[col][ce].nextState = 1
								infected++
							}
						}
					}
				}
			}
		}
	}

	// Copy next-state to current-state
	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.drawCell(col, row)
			s.cells[col][row].state = s.cells[col][row].nextState
		}
	}

	// fmt.Println("Newly infected: ", infected)
	return infected > 0
}

func (s *SISCityModel) drawCell(col, row int) {
	switch s.cells[col][row].degree {
	case 5:
		s.raster.SetPixelColor(s.degree5Color)
	case 6:
		s.raster.SetPixelColor(s.degree6Color)
	case 7:
		s.raster.SetPixelColor(s.degree7Color)
	case 8:
		s.raster.SetPixelColor(s.degree8Color)
	}

	switch s.cells[col][row].state {
	case 1:
		s.raster.SetPixelColor(s.infectedColor)
	case 2:
		if s.cells[col][row].degree < 5 {
			s.raster.SetPixelColor(s.susceptibleColor)
		}
		// default:
		// 	s.raster.SetPixelColor(s.undetermenedColor)
	}

	s.raster.SetPixel(col, row)
}
