package model

const (
	CellSignWhite = -1
	CellSignNone  = 0
	CellSignBlack = 1
)

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

func (board *Board) getConnectedComponentInner(vertex int32, signs []int32, result map[int32]interface{}) {
	result[vertex] = struct {
	}{}

	for _, v := range board.GetNeighbors(vertex) {
		sign := board.GetSign(v)
		idx := sign + 1 // -1,0,1 -> 0,1,2
		if signs[idx] > 0 && result[v] == nil {
			board.getConnectedComponentInner(v, signs, result)
		}
	}
}

func (board *Board) getConnectedComponent(vertex int32, signs []int32) map[int32]interface{} {
	result := make(map[int32]interface{})
	board.getConnectedComponentInner(vertex, signs, result)
	return result
}

func (board *Board) GetChain(vertex int32) map[int32]interface{} {
	sign := board.GetSign(vertex)
	signs := []int32{0, 0, 0}
	signs[sign+1] = 1 // -1, 0, 1 -> 0, 1, 2
	return board.getConnectedComponent(vertex, signs)
}

func (board *Board) GetFloatingStones() []int32 {
	result := make([]int32, 0)
	visited := make([]interface{}, 0)
	for i, v := range board.arrangement {
		if v != 0 || visited[v] != nil {
			continue
		}

		// 白子与空子片区
		posArea := board.getConnectedComponent(int32(i), []int32{1, 1, 0}) // -1, 0 (sign + 1 as index 0, 1, 2)

		// 黑子与空子片区
		negArea := board.getConnectedComponent(int32(i), []int32{0, 1, 1}) // 0, 1

		// 白子列表
		posDead := filterMap(posArea, func(k int32, v interface{}) bool {
			return board.arrangement[k] == CellSignWhite
		})

		// 黑子列表
		negDead := filterMap(negArea, func(k int32, v interface{}) bool {
			return board.arrangement[k] == CellSignBlack
		})

		// 白空片区中的空子数量
		posDiff := len(filterMap(posArea, func(k int32, v interface{}) bool {
			return posDead[k] == nil && negArea[k] == nil
		}))

		// 黑空片区中的空子数量
		negDiff := len(filterMap(negArea, func(k int32, v interface{}) bool {
			return negDead[k] == nil && posArea[k] == nil
		}))

		favorNeg := negDiff <= 1 && len(negDead) <= len(posDead)
		favorPos := posDiff <= 1 && len(posDead) <= len(negDead)

		var actualArea map[int32]interface{} = nil
		var actualDead map[int32]interface{} = nil
		if !favorNeg && favorPos {
			actualArea = posArea
			actualDead = posDead
		} else if favorNeg && !favorPos {
			actualArea = negArea
			actualDead = negDead
		} else {
			actualArea = board.GetChain(v)
			actualDead = make(map[int32]interface{})
		}
		for k, _ := range actualArea {
			visited[k] = struct {
			}{}
		}
		for k, _ := range actualDead {
			result = append(result, k)
		}
	}

	return result
}

func (board *Board) GetArrangement() []int32 {
	return board.arrangement
}

func (board *Board) GetSize() int32 {
	return board.size
}

func (board *Board) Clone() *Board {
	return nil
}

func (board *Board) hasLibertiesInner(vertex int32, visited map[int32]interface{}, sign int32) bool {
	visited[vertex] = struct {
	}{}

	for _, v := range board.GetNeighbors(vertex) {
		x := board.GetSign(v)
		if x == CellSignNone {
			return true
		}
		if x == sign && visited[v] == nil {
			if board.hasLibertiesInner(v, visited, sign) {
				return true
			}
		}
	}

	return false
}

// 是否有气
func (board *Board) hasLiberties(vertex int32) bool {
	visited := make(map[int32]interface{})
	sign := board.GetSign(vertex)
	return board.hasLibertiesInner(vertex, visited, sign)
}

func (board *Board) GetVertexBySign(sign int32) []int32 {
	result := make([]int32, 0, board.size*board.size)
	for i, v := range board.arrangement {
		if v == sign {
			result = append(result, int32(i))
		}
	}
	return result
}

func (board *Board) MakePseudoMove(sign int32, vertex int32) []int32 {
	neighbors := board.GetNeighbors(vertex)
	checkCapture := false
	checkMultiDeadChains := false

	isClose := true
	for _, v := range neighbors {
		s := board.GetSign(v)
		if s != sign {
			isClose = false
			break
		}
	}
	if isClose {
		return nil
	}

	board.SetSign(vertex, sign)

	if !board.hasLiberties(vertex) {
		isPointChain := true

		for _, v := range neighbors {
			if board.GetSign(v) == sign {
				isPointChain = false
				break
			}
		}

		if isPointChain {
			checkMultiDeadChains = true
		} else {
			checkCapture = true
		}
	}

	dead := make([]int32, 0)
	deadChains := 0

	for _, v := range neighbors {
		if board.GetSign(v) != -sign || board.hasLiberties(v) {
			continue
		}

		chain := board.GetChain(v)
		deadChains += 1

		for k, _ := range chain {
			board.SetSign(k, 0)
			dead = append(dead, k)
		}
	}

	if checkMultiDeadChains && deadChains <= 1 || checkCapture && len(dead) == 0 {
		for _, v := range dead {
			board.SetSign(v, -sign)
		}

		board.SetSign(vertex, 0)
		return nil
	}

	return dead
}

func (board *Board) GetRelatedChains(vertex int32) map[int32]interface{} {
	sign := board.GetSign(vertex)
	signs := make([]int32, 3)
	signs[1] = 1
	signs[sign+1] = 1
	return board.getConnectedComponent(vertex, signs)
}

func (board *Board) GetLength() int32 {
	return board.size * board.size
}
