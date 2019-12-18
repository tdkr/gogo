package influence

import (
	"math"

	"github.com/tdkr/gogo/model"
)

func getNeighbors(vec model.Vector2) []model.Vector2 {
	return []model.Vector2{
		{vec.X - 1, vec.Y},
		{vec.X + 1, vec.Y},
		{vec.X, vec.Y - 1},
		{vec.X, vec.Y + 1},
	}
}

func getChainInner(board [][]float32, vec model.Vector2, result *model.VecStack, visited map[int32]interface{}, sign float32) {
	result.Push(vec)
	visited[vec.HashCode()] = struct {
	}{}

	for _, v := range getNeighbors(vec) {
		if !isValidVertex(board, int(v.X), int(v.Y)) || board[v.Y][v.X] != sign || visited[v.HashCode()] != nil {
			continue
		}

		getChainInner(board, v, result, visited, sign)
	}
}

func GetChain(board [][]float32, v model.Vector2) *model.VecStack {
	sign := board[v.Y][v.X]
	result := model.NewVecStack()
	visited := make(map[int32]interface{})
	getChainInner(board, v, result, visited, sign)
	return result
}

func GetIntSign(value int32) int32 {
	switch {
	case value > 0:
		return 1
	case value < 0:
		return -1
	case value == 0:
		return 0
	}
	return 0
}

func GetFloatSign(value float32) float32 {
	switch {
	case value > 0:
		return 1
	case value < 0:
		return -1
	case value == 0:
		return 0
	}
	return 0
}

func MinInt(vals ...int32) int32 {
	min := int32(math.MaxInt32)
	for _, v := range vals {
		if min > v {
			min = v
		}
	}
	return min
}

func MinFloat(vals ...float32) float32 {
	min := float32(math.MaxFloat32)
	for _, v := range vals {
		if min > v {
			min = v
		}
	}
	return min
}

func MaxInt(vals ...int32) int32 {
	max := int32(math.MinInt32)
	for _, v := range vals {
		if max < v {
			max = v
		}
	}
	return max
}

func MaxFloat(vals ...float32) float32 {
	max := float32(-math.MaxFloat32)
	for _, v := range vals {
		if max < v {
			max = v
		}
	}
	return max
}

func GetMatrixSize(data [][]float32) (int32, int32) {
	h := len(data)
	if h == 0 {
		return 0, 0
	} else {
		return int32(len(data[0])), int32(h)
	}
}

func NewFloatMatrix(w, h int, val float32) [][]float32 {
	ret := make([][]float32, h)
	for i := 0; i < h; i++ {
		ret[i] = make([]float32, w)
		for j := 0; j < w; j++ {
			ret[i][j] = val
		}
	}
	return ret
}

func isValidVertex(data [][]float32, x, y int) bool {
	if x < 0 || y < 0 {
		return false
	}
	if y >= len(data) {
		return false
	}
	return x < len(data[y])
}

func CloneFloatMatrix(m [][]float32) [][]float32 {
	ret := make([][]float32, len(m))
	for i, v := range m {
		ret[i] = make([]float32, len(v))
		for i2, v2 := range v {
			ret[i][i2] = v2
		}
	}
	return ret
}
