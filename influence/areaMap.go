package influence

import (
	"github.com/tdkr/gogo/model"
)

func GetAreaMap(board [][]float32) [][]float32 {
	result := NewFloatMatrix(board, model.StoneSignNone)

	width, height := GetMatrixSize(board)

	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {
			v := model.NewVec2(x, y)

			if result[y][x] != model.StoneSignNone {
				continue
			}

			if sign := board[y][x]; sign != model.StoneSignEmpty {
				result[y][x] = sign
				continue
			}

			chain := GetChain(board, v)
			sign := float32(0)
			indicator := float32(1)

			for _, cv := range chain.Nodes() {
				if indicator == 0 {
					break
				}

				for _, nv := range getNeighbors(cv) {
					if !isValidVertex(board, int(nv.X), int(nv.Y)) {
						continue
					}

					val := board[nv.Y][nv.X]
					if val == model.StoneSignEmpty {
						continue
					}

					if sign == 0 {
						sign = GetFloatSign(val)
					} else if sign != GetFloatSign(val) {
						indicator = 0
						break
					}
				}
			}

			for _, cv := range chain.Nodes() {
				result[cv.Y][cv.X] = sign * indicator
			}
		}
	}

	return result
}
