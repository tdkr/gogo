package model

const (
	CellSignWhite = -1
	CellSignNone  = 0
	CellSignBlack = 1
)

type BoardCell struct {
	X int32
	Y int32
}

type Board struct {
	size        int32
	width       int32
	height      int32
	captures    [2]int32
	arrangement []int32
}

func NewBoard(size int32) *Board {
	board := &Board{}
	return board
}

func (board *Board) GetNeighbors(vertex int32) []int32 {
	result := []int32{}
	if vertex >= board.size {
		result = append(result, vertex-board.size)
	}
	if vertex < (board.size-1)*board.size {
		result = append(result, vertex+board.size)
	}

	mod := vertex % board.width
	if mod > 0 {
		result = append(result, vertex-1)
	}
	if mod < board.width-1 {
		result = append(result, vertex+1)
	}

	return result
}

func (board *Board) SetSign(vertex int32, sign int32) {
	board.arrangement[vertex] = sign
}

func (board *Board) GetSign(vertex int32) (sign int32) {
	sign = board.arrangement[vertex]
	return
}

func (board *Board) hasVertex(vertex int32) bool {
	return vertex >= 0 && vertex < board.size*board.size
}

func (board *Board) getConnectedComponentInner(vertex int32, signs []int32, result []int32) {
	result = append(result, vertex)

	for _, v := range board.GetNeighbors(vertex) {
		sign := board.GetSign(v)
		idx := sign + 1 // -1,0,1 -> 0,1,2
		if signs[idx] > 0 && result[v] == 0 {
			board.getConnectedComponentInner(v, signs, result)
		}
	}
}

func (board *Board) getConnectedComponent(vertex int32, signs []int32) []int32 {
	result := make([]int32, 0)
	board.getConnectedComponentInner(vertex, signs, result)
	return result
}

func filterSlice(s []int32, f func(int32) bool) []int32 {
	results := make([]int32, 0, len(s))
	for _, v := range s {
		if f(v) {
			results = append(results, v)
		}
	}
	return results
}

func findSlice(s []int32, val int32) bool {
	for _, v := range s {
		if v == val {
			return true
		}
	}
	return false
}

func (board *Board) getChain(vertex int32) []int32 {
	sign := board.GetSign(vertex)
	signs := []int32{0, 0, 0}
	signs[sign+1] = 1
	return board.getConnectedComponent(vertex, signs)
}

func (board *Board) GetFloatingStones() []int32 {
	result := make([]int32, 0)
	visited := make([]interface{}, 0)
	for i, v := range board.arrangement {
		if v != 0 || visited[v] != nil {
			continue
		}

		posArea := board.getConnectedComponent(int32(i), []int32{1, 1, 0}) // -1, 0 (sign + 1 as index 0, 1, 2)
		negArea := board.getConnectedComponent(int32(i), []int32{0, 1, 1}) // 0, 1

		posDead := filterSlice(posArea, func(i int32) bool {
			return board.arrangement[i] == CellSignWhite
		})
		negDead := filterSlice(negArea, func(i int32) bool {
			return board.arrangement[i] == CellSignBlack
		})

		posDiff := len(filterSlice(posArea, func(i int32) bool {
			return !findSlice(posDead, i) && !findSlice(negArea, i)
		}))
		negDiff := len(filterSlice(negArea, func(i int32) bool {
			return !findSlice(negDead, i) && !findSlice(posArea, i)
		}))

		favorNeg := negDiff <= 1 && len(negDead) <= len(posDead)
		favorPos := posDiff <= 1 && len(posDead) <= len(negDead)

		var actualArea []int32 = nil
		var actualDead []int32 = nil
		if !favorNeg && favorPos {
			actualArea = posArea
			actualDead = posDead
		} else if favorNeg && !favorPos {
			actualArea = negArea
			actualDead = negDead
		} else {
			actualArea = board.getChain(v)
			actualDead = []int32{}
		}

		for _, v := range actualArea {
			visited[v] = struct {
			}{}
		}
		result = append(result, actualDead...)
	}

	return result
}
