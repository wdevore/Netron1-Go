package simulation

import (
	"fmt"
	"image/color"
)

// KCell is a knowledge cell
type KCell struct {
	knowledgeCenter bool
	col, row        int
	color           color.RGBA

	// A cell can only "increase" its knowledge
	// 1 = orange
	// 2 = green
	// 3 = teal
	// 4 = purple
	knowledge     int
	nextKnowledge int

	// 0 = stupid, 10 = smart
	intelligence int

	// 0 = No knowledge (susceptable)
	// 1 = has knowledge (infected)
	state     int
	nextState int
}

func NewKCellCenter(col, row int, color color.RGBA, knowledge int) KCell {
	k := KCell{
		col:             col,
		row:             row,
		color:           color,
		knowledge:       knowledge,
		knowledgeCenter: true,
		state:           1}
	return k
}

func (k *KCell) toString() string {
	return fmt.Sprintf("[%d,%d] (%d->%d) k: |%d|, KCenter: %t", k.col, k.row, k.state, k.nextState, k.knowledge, k.knowledgeCenter)
}
