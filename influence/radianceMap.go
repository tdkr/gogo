package influence

import "github.com/tdkr/gogo/model"

func GetRadianceMap(board *model.Board, sign int32, opts ...option) {
	o := NewOptions(opts...)
	height := board.GetSize()
	width := board.GetSize()
	visited := make(map[int32]interface{})
	result := make(map[int32]float32)

	//getMirroredVertex := func(v int32) int32 {
	//	w := v / height
	//	h := v / width
	//	return w*width + h
	//}

	castRadiance := func(chain map[int32]interface{}) {
		queue := make([][2]int32, len(chain))
		for k, _ := range chain {
			queue = append(queue, [2]int32{k, 0})
		}

		visited := make(map[int32]interface{})

		for len(queue) > 0 {
			item := queue[0]

			mv := item[0]
			if mv == item[0] {
				result[mv] += o.p2 / float32(item[1]/o.p1*6+1)
			} else {
				result[mv] += float32(o.p3)
			}

			for _, nv := range board.GetNeighbors(item[0]) {
				if item[1] >= o.p1 || board.GetSign(nv) == -sign || visited[nv] != nil {
					continue
				}

				visited[nv] = struct {
				}{}
				queue = append(queue, [2]int32{nv, item[1] + 1})
			}
		}
	}

	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {
			v := y*width + x
			if board.GetSign(v) != sign || visited[v] != nil {
				continue
			}

			chain := board.GetChain(v)
			for k, _ := range chain {
				visited[k] = struct {
				}{}
			}

			castRadiance(chain)
		}
	}
}
