package influence

import (
	"github.com/tdkr/gogo/model"
)

type castItem struct {
	vec   *model.Vector2
	level int32
}

func GetRadianceMap(board [][]float32, sign float32, opts ...option) [][]float32 {
	o := NewOptions(opts...)
	width, height := GetMatrixSize(board)
	visited := make(map[int32]interface{})
	result := NewFloatMatrix(len(board), len(board[0]), 0)

	getMirroredVertex := func(v *model.Vector2) *model.Vector2 {
		if isValidVertex(board, int(v.X), int(v.Y)) {
			return model.NewVec2(v.X, v.Y)
		}
		rv := model.NewVec2(v.X, v.Y)
		if rv.X < 0 {
			rv.X = -rv.X - 1
		}
		if rv.X >= width {
			rv.X = width*2 - rv.X - 1
		}
		if rv.Y < 0 {
			rv.Y = -rv.Y - 1
		}
		if rv.Y >= height {
			rv.Y = height*2 - rv.Y - 1
		}
		return rv
	}

	castRadiance := func(chain *model.VecStack) {
		queue := make([]*castItem, chain.Size())
		for i, v := range chain.Nodes() {
			queue[i] = &castItem{
				vec:   v,
				level: 0,
			}
		}

		visited := make(map[int32]interface{})

		for len(queue) > 0 {
			item := queue[0]
			queue[0] = queue[len(queue)-1]
			queue = queue[:len(queue)-1]

			mv := getMirroredVertex(item.vec)
			if mv.Equals(item.vec) {
				result[mv.Y][mv.X] += o.radicanceVar2 / float32(item.level/o.radianceVar1*6+1)
			} else {
				result[mv.Y][mv.X] += o.radicanceVar3
			}

			for _, nv := range getNeighbors(item.vec) {
				if item.level >= o.radianceVar1 ||
					isValidVertex(board, int(item.vec.X), int(item.vec.Y)) && board[item.vec.Y][item.vec.X] == -sign ||
					visited[nv.HashCode()] != nil {
					continue
				}

				visited[nv.HashCode()] = struct {
				}{}
				queue = append(queue, &castItem{
					vec:   nv,
					level: item.level + 1,
				})
			}
		}
	}

	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {
			v := model.NewVec2(x, y)

			if !isValidVertex(board, int(v.X), int(v.Y)) || board[v.Y][v.X] != sign || visited[v.HashCode()] != nil {
				continue
			}

			chain := GetChain(board, v)
			for _, cv := range chain.Nodes() {
				visited[cv.HashCode()] = struct {
				}{}
			}

			castRadiance(chain)
		}
	}

	return result
}
