package simulation

import (
	"Netron1-Go/api"
	"fmt"
	"image/color"
	"math/rand"
)

// The Knowledge model is based on "information" and "sequences".
// The idea is that a cell can only gain knowledge if a neighbor
// cell has a higher level knowledge

// In this model "knowledge" spreads via a sequence:
// Orange -> Green -> Teal -> Purple

// For example, Orange can't gain Purple directly but Teal can.

// type KnowledgeCenter struct {
// 	col, row  int
// 	color     color.RGBA
// 	knowledge int
// }

type SISKnowledgeModel struct {
	blueColor   color.RGBA
	orangeColor color.RGBA
	greenColor  color.RGBA
	tealColor   color.RGBA
	purpleColor color.RGBA
	susColor    color.RGBA // cell type = 2

	raster api.IRasterBuffer
	cells  [][]KCell

	knowledgeCenters []KCell

	acceptableRate float64
	// The chance they will drop meditation.
	dropRate float64
}

func NewSISKnowledgeModel() api.IModel {
	o := new(SISKnowledgeModel)

	o.blueColor = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	o.orangeColor = color.RGBA{R: 255, G: 127, B: 0, A: 255}
	o.greenColor = color.RGBA{R: 0, G: 255, B: 100, A: 255}
	o.tealColor = color.RGBA{R: 0, G: 200, B: 200, A: 255}
	o.purpleColor = color.RGBA{R: 255, G: 0, B: 255, A: 255}
	o.susColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	return o
}

func (s *SISKnowledgeModel) Name() string {
	return "SISKnowledgeModel"
}

func (s *SISKnowledgeModel) Configure(rasterBuffer api.IRasterBuffer) {
	s.raster = rasterBuffer
	s.acceptableRate = 0.23 // 0.22
	s.dropRate = 0.4        // 0.6

	rand.Seed(131)

	s.knowledgeCenters = []KCell{} //make([]KnowledgeCenter, 4)

	s.cells = make([][]KCell, s.raster.Height())
	for i := range s.cells {
		s.cells[i] = make([]KCell, s.raster.Width())
	}
}

// SendEvent receives an event from the host simulation
func (s *SISKnowledgeModel) SendEvent(event string) {
}

func (s *SISKnowledgeModel) Reset() {
	fmt.Println(("--- sir reset ---"))
	s.raster.Clear()

	w := s.raster.Width()
	h := s.raster.Height()

	// Initialize population
	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			s.cells[col][row].state = 0 // No knowledge
			s.cells[col][row].nextState = 0
			s.cells[col][row].row = row
			s.cells[col][row].col = col
		}
	}

	// Initial knowledge crowd
	px := 150
	py := 150
	radius := 4
	for col := px; col < px+radius; col += 1 {
		for row := py; row < py+radius; row += 1 {
			s.cells[col][row].state = 1     // has knowledge
			s.cells[col][row].knowledge = 1 // orange knowledge
		}
	}

	// Create 4 knowledge centers with different skills levels
	distance := 15

	s.knowledgeCenters = append(s.knowledgeCenters, NewKCellCenter(px, px, color.RGBA{R: 255, G: 127, B: 0, A: 255}, 1))
	s.knowledgeCenters = append(s.knowledgeCenters, NewKCellCenter(px+distance, px, color.RGBA{R: 0, G: 255, B: 100, A: 255}, 2))
	s.knowledgeCenters = append(s.knowledgeCenters, NewKCellCenter(px+distance, px+distance, color.RGBA{R: 0, G: 200, B: 200, A: 255}, 3))
	s.knowledgeCenters = append(s.knowledgeCenters, NewKCellCenter(px, px+distance, color.RGBA{R: 255, G: 0, B: 255, A: 255}, 4))
	// Mark them on the grid
	for _, k := range s.knowledgeCenters {
		s.cells[k.col][k.row].state = 1
		s.cells[k.col][k.row].knowledge = k.knowledge
		s.cells[k.col][k.row].knowledgeCenter = k.knowledgeCenter
	}

	s.drawKnowledgeCenters()
}

func (s *SISKnowledgeModel) Step() bool {
	// The current "step" works on the current-state but updates the next-state
	// Once done, the next-state is copied back to the current-state.
	w := s.raster.Width()
	h := s.raster.Height()
	ce := 0
	knowledged := 0

	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			cenC := &s.cells[col][row] // Center
			// if s.regionKnowledge(col, row) == 2 {
			// 	fmt.Println("C: " + cenC.toString())
			// }

			// Knowledge centers can only pass knowledge on to outsiders
			if cenC.state == 1 && !cenC.knowledgeCenter {
				// s.drawMap(col, row)
				// Check neighbors to see if the center cell either transfers
				// knowledge or receives it.
				ce = col + 1
				if ce < w { // Right
					nei := &s.cells[ce][row]
					if nei.state == 0 {
						// The neighbor has NO knowledge. If they are receptive
						// then the neighbor gains the center's knowledge.
						if rand.Float64() < s.acceptableRate {
							nei.nextKnowledge = cenC.knowledge
							nei.nextState = 1
							knowledged++
						}
					} else {
						// The neighbor has knowledge. If it is higher then
						// take on that knowledge.
						cenC.nextKnowledge = s.gainKnowledge(cenC.knowledge, nei.knowledge)
					}
				}

				ce = col - 1
				if ce >= 0 { // Left
					nei := &s.cells[ce][row]
					if nei.state == 0 {
						if rand.Float64() < s.acceptableRate {
							nei.nextKnowledge = cenC.knowledge
							nei.nextState = 1
							knowledged++
						}
					} else {
						cenC.nextKnowledge = s.gainKnowledge(cenC.knowledge, nei.knowledge)
					}
				}

				ce = row - 1
				if ce >= 0 { // Top
					nei := &s.cells[col][ce]
					if nei.state == 0 {
						if rand.Float64() < s.acceptableRate {
							nei.nextKnowledge = cenC.knowledge
							nei.nextState = 1
							knowledged++
						}
					} else {
						cenC.nextKnowledge = s.gainKnowledge(cenC.knowledge, nei.knowledge)
					}
				}

				ce = row + 1
				if ce < h { // Bottom
					nei := &s.cells[col][ce]
					if nei.state == 0 {
						if rand.Float64() < s.acceptableRate {
							nei.nextKnowledge = cenC.knowledge
							nei.nextState = 1
							knowledged++
						}
					} else {
						cenC.nextKnowledge = s.gainKnowledge(cenC.knowledge, nei.knowledge)
					}
				}

				// Knowledge centers retain their knowledge, everyone
				// else may lose their knowledge.
				if !cenC.knowledgeCenter {
					if rand.Float64() < s.dropRate {
						cenC.nextState = 0 // Loses knowledge, but retains skill
					}
				}
			}
		}
	}

	// Copy next-state to current-state
	s.draw(w, h)
	s.drawKnowledgeCenters()

	// fmt.Println("Newly infected: ", infected)
	return knowledged > 0
}

func (s *SISKnowledgeModel) gainKnowledge(cenK, neiK int) int {
	// if nei.knowledge > cenC.knowledge {
	// 	cenC.nextKnowledge = nei.knowledge
	// }

	if cenK < neiK {
		if cenK == 1 && neiK == 2 {
			// fmt.Println(k, ",", rgk)
			cenK = neiK // Increase from orange to green
		} else if cenK == 2 && neiK == 3 {
			// fmt.Println(k, ",", rgk)
			cenK = neiK // Increase from green to teal
		} else if cenK == 3 && neiK == 4 {
			// fmt.Println(k, ",", rgk)
			cenK = neiK // Increase from teal to purple
		}
	}

	return cenK
}

func (s *SISKnowledgeModel) drawMap(c, r int) {
	fmt.Println("========================================")
	for col := c - 5; col < c+5; col += 1 {
		fmt.Print("|")
		for row := r - 5; row < r+5; row += 1 {
			if s.cells[col][row].knowledgeCenter {
				fmt.Print(s.cells[col][row].knowledge, "+")
			} else {
				if col == c && row == r {
					fmt.Print(s.cells[col][row].knowledge, ".")
				} else {
					fmt.Print(s.cells[col][row].knowledge, " ")
				}
			}
		}
		fmt.Println("|")
	}
}

func (s *SISKnowledgeModel) regionKnowledge(c, r int) int {
	for _, k := range s.knowledgeCenters {
		if k.col == c && k.row == r {
			return k.knowledge
		}
	}

	return 0
}

func (s *SISKnowledgeModel) drawKnowledgeCenters() {
	radius := 2

	for _, k := range s.knowledgeCenters {
		s.raster.SetPixelColor(k.color)
		for col := k.col; col < k.col+radius; col += 1 {
			for row := k.row; row < k.row+radius; row += 1 {
				s.raster.SetPixel(col, row)
			}
		}
	}
}

func (s *SISKnowledgeModel) draw(w, h int) {
	for col := 0; col < w; col += 1 {
		for row := 0; row < h; row += 1 {
			if s.cells[col][row].state == 1 {
				switch s.cells[col][row].knowledge {
				case 1:
					s.raster.SetPixelColor(s.orangeColor)
				case 2:
					s.raster.SetPixelColor(s.greenColor)
				case 3:
					s.raster.SetPixelColor(s.tealColor)
				case 4:
					s.raster.SetPixelColor(s.purpleColor)
				}
			} else {
				s.raster.SetPixelColor(s.susColor)
			}

			s.raster.SetPixel(col, row)

			if !s.cells[col][row].knowledgeCenter {
				s.cells[col][row].state = s.cells[col][row].nextState
				s.cells[col][row].knowledge = s.cells[col][row].nextKnowledge
			}
		}
	}
}

func (s *SISKnowledgeModel) drawCell(col, row int) {
	if s.cells[col][row].state == 1 {
		switch s.cells[col][row].knowledge {
		case 1:
			s.raster.SetPixelColor(s.orangeColor)
		case 2:
			s.raster.SetPixelColor(s.greenColor)
		case 3:
			s.raster.SetPixelColor(s.tealColor)
		case 4:
			s.raster.SetPixelColor(s.purpleColor)
		}
	} else {
		s.raster.SetPixelColor(s.susColor)
	}

	s.raster.SetPixel(col, row)
}
