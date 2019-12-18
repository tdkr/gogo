package influence

import (
	"math"
)

func GetNearestNeighborMap(board [][]float32, sign float32) [][]float32 {
	var inf float32 = math.MaxFloat32
	var min float32 = inf
	result := NewFloatMatrix(len(board), len(board[0]), inf)
	width, height := GetMatrixSize(board)

	f := func(x int32, y int32) {
		if board[y][x] == sign {
			min = 0
		} else {
			min++
		}

		if min < result[y][x] {
			result[y][x] = min
		} else {
			min = result[y][x]
		}
	}

	for y := int32(0); y < height; y++ {
		min = inf

		for x := int32(0); x < width; x++ {
			old := inf

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
		min = inf

		for x := width - 1; x >= 0; x-- {
			old := inf

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
