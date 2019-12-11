package influence

import (
	"github.com/tdkr/gogo/model"
	"math"
)

func GetNearestNeighborMap(board *model.Board, sign int32) map[int32]int32 {
	var maxInt32 int32 = math.MaxInt32
	var min int32 = maxInt32
	result := make(map[int32]int32)
	width := board.GetSize()
	height := board.GetSize()

	f := func(x int32, y int32) {
		vertex := y*width + x
		if board.GetSign(vertex) == sign {
			min = 0
		} else {
			min++
		}

		if min < result[vertex] {
			result[vertex] = min
		} else {
			min = result[vertex]
		}
	}

	for y := int32(0); y < height; y++ {
		min = maxInt32

		for x := int32(0); x < width; x++ {
			old := maxInt32

			f(x, y)
			old = min

			for ny := y + 1; ny < height; ny++ {
				f(x, ny)
			}
			min = old

			for ny := y - 1; ny >= 0; ny-- {
				f(x, ny)
			}
			min = old
		}
	}

	for y := height - 1; y >= 0; y-- {
		min = maxInt32

		for x := width - 1; x >= 0; x-- {
			old := maxInt32

			f(x, y)
			old = min

			for ny := y + 1; ny < height; ny++ {
				f(x, ny)
			}
			min = old

			for ny := y - 1; ny >= 0; ny-- {
				f(x, ny)
			}
			min = old
		}
	}

	return result
}
