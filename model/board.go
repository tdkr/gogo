package model

import "math"

const (
	StoneSignWhite = -1
	StoneSignNone  = 0
	StoneSignBlack = 1
)

type Board struct {
	size        int32
	width       int32
	height      int32
	captures    []int32
	arrangement [][]int32
}

func NewBoard(opts ...option) *Board {
	o := NewOptions(opts...)
	board := &Board{
		width:  19,
		height: 19,
	}
	if o.width != 0 {
		board.width = o.width
	}
	if o.height != 0 {
		board.height = o.height
	}
	if o.arrangement != nil {
		board.arrangement = o.arrangement
	}
	if o.captures != nil {
		board.captures = o.captures
	}
	return board
}

func (board *Board) Get(vec *Vector2) int32 {
	return board.arrangement[vec.Y][vec.X]
}

func (board *Board) Get2(x, y int32) int32 {
	return board.arrangement[x][y]
}

func (board *Board) Set(vec *Vector2, sign int32) {
	board.arrangement[vec.Y][vec.X] = sign
}

func (board *Board) Set2(x, y, sign int32) {
	board.arrangement[x][y] = sign
}

func (board *Board) Clone() *Board {
	return nil
}

func (board *Board) Diff(b *Board) []*Vector2 {
	if board.width != b.width || board.height != b.height {
		return nil
	}

	result := make([]*Vector2, 0)

	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			if board.Get2(x, y) != b.Get2(x, y) {
				result = append(result, NewVec2(x, y))
			}
		}
	}
	return result
}

func (board *Board) HasVertex(v *Vector2) bool {
	return 0 <= v.X && v.X < board.width && 0 <= v.Y && v.Y < board.height
}

func (board *Board) HasVertex2(x, y int32) bool {
	return 0 <= x && y < board.width && 0 <= y && y < board.height
}

func (board *Board) Clear() {
	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			board.Set2(x, y, StoneSignNone)
		}
	}
}

func (board *Board) IsSquare() bool {
	return board.width == board.height
}

func (board *Board) IsEmpty() bool {
	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			if board.Get2(x, y) != StoneSignNone {
				return false
			}
		}
	}
	return true
}

func (board *Board) GetDistance(v1 *Vector2, v2 *Vector2) int32 {
	return int32(math.Abs(float64(v1.X-v2.X)) + math.Abs(float64(v1.Y-v2.Y)))
}

func (board *Board) GetNeighbors(vec *Vector2, ignoreBoard bool) []*Vector2 {
	if !ignoreBoard && !board.HasVertex(vec) {
		return []*Vector2{}
	}

	if ignoreBoard {
		return []*Vector2{
			NewVec2(vec.X-1, vec.Y),
			NewVec2(vec.X+1, vec.Y),
			NewVec2(vec.X, vec.Y-1),
			NewVec2(vec.X, vec.Y+1),
		}
	}

	result := make([]*Vector2, 0, 4)
	if vec.X > 0 {
		result = append(result, NewVec2(vec.X-1, vec.Y))
		if vec.X < board.width-1 {
			result = append(result, NewVec2(vec.X+1, vec.Y))
		}
	}
	if vec.Y > 0 {
		result = append(result, NewVec2(vec.X, vec.Y-1))
		if vec.Y < board.height-1 {
			result = append(result, NewVec2(vec.X, vec.Y+1))
		}
	}
	return result
}

func (board *Board) getConnectedComponentInner(vec *Vector2, signs map[int32]interface{}, result *Stack) {
	if !board.HasVertex(vec) {
		return
	}

	result.Push(vec)

	for _, v := range board.GetNeighbors(vec, false) {
		if result.Find(v) > 0 {
			continue
		}
		sign := board.Get(v)
		if signs[sign] == nil {
			continue
		}
		result.Push(v)
		board.getConnectedComponentInner(v, signs, result)
	}
}

func (board *Board) GetConnectedComponent(vec *Vector2, signs map[int32]interface{}) *Stack {
	result := NewStack()
	board.getConnectedComponentInner(vec, signs, result)
	return result
}

func (board *Board) GetChain(vec *Vector2) *Stack {
	sign := board.Get(vec)
	result := NewStack()
	board.getConnectedComponentInner(vec, map[int32]interface{}{
		sign: struct {
		}{},
	}, result)
	return result
}

func (board *Board) hasLibertiesInner(vec *Vector2, visited map[int32]interface{}, sign int32) bool {
	visited[vec.HashCode()] = struct {
	}{}

	neighbors := board.GetNeighbors(vec, false)
	for _, v := range neighbors {
		s := board.Get(v)
		if s == StoneSignNone {
			return true
		}
		if s == sign && visited[v.HashCode()] == nil {
			if board.hasLibertiesInner(v, visited, sign) {
				return true
			}
		}
	}

	return false
}

func (board *Board) HasLiberties(vec *Vector2) bool {
	if !board.HasVertex(vec) {
		return false
	}
	visited := map[int32]interface{}{

	}
	return board.hasLibertiesInner(vec, visited, board.Get(vec))
}

func (board *Board) GetRelatedChain(vec *Vector2) []*Vector2 {
	if !board.HasVertex(vec) {
		return []*Vector2{}
	}

	sign := board.Get(vec)
	if sign == StoneSignNone {
		return []*Vector2{}
	}

	signs := make(map[int32]interface{})
	signs[sign] = struct {
	}{}
	area := board.GetConnectedComponent(vec, signs).Nodes()

	cnt := 0
	i := 0
	for i < len(area)-cnt {
		v := area[i]
		if board.Get(v) == sign {
			area[i] = area[len(area)-1-cnt]
			cnt++
		} else {
			i++
		}
	}

	return area[:len(area)-cnt]
}

//komi        float32 贴目
//handicap    float32 让子
func (board *Board) GetScore(areaMap [][]int32, komi float32, handicap int32) *Score {
	score := NewScore(komi, handicap, board.captures)

	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			sign := areaMap[x][y]
			if sign == StoneSignNone {
				continue
			}
			scoreIdx := 0
			switch sign {
			case StoneSignBlack:
				scoreIdx = ScoreIndexBlack
			case StoneSignWhite:
				scoreIdx = ScoreIndexWhite
			}
			score.Area[scoreIdx]++
			if board.Get2(x, y) == StoneSignNone {
				score.Territory[scoreIdx]++
			}
		}
	}

	return score
}

func (board *Board) Vec2Coord(vec *Vector2) int32 {
	if !board.HasVertex(vec) {
		return -1
	}

	return vec.X + vec.Y*board.height
}

func (board *Board) Coord2Vec(coord int32) *Vector2 {
	x := coord % board.height
	y := coord / board.height
	return NewVec2(x, y)
}

//TODO
func (board *Board) IsValid() bool {
	return true
}

func (board *Board) MakeMove(sign int32, vec *Vector2) *Board {
	move := board.Clone()

	if sign == StoneSignNone || !board.HasVertex(vec) {
		return move
	}

	move.Set(vec, sign)

	deadNeighbors := make([]*Vector2, 0, 0)
	for _, v := range move.GetNeighbors(vec, false) {
		if move.Get(v) == -sign && !move.HasLiberties(v) {
			deadNeighbors = append(deadNeighbors, v)
		}
	}

	for _, v := range deadNeighbors {
		if move.Get(v) == StoneSignNone {
			continue
		}

		for _, cv := range move.GetChain(v).Nodes() {
			move.Set(cv, 0)
			move.captures[(-sign)+1/2]++
		}
	}

	move.Set(vec, sign)

	// Detect suicide
	if len(deadNeighbors) == 0 && !move.HasLiberties(vec) {
		for _, cv := range move.GetChain(vec).Nodes() {
			move.Set(cv, 0)
			move.captures[(sign+1)/2]++
		}
	}

	return move
}

func (board *Board) GetHandicapPlacement(count int32, tygemFlag bool) []*Vector2 {
	result := make([]*Vector2, 0)

	if board.width <= 6 || board.height <= 6 || count < 2 {
		return result
	}

	near := NewVec2(2, 2)
	if board.width >= 13 {
		near.X = 3
	}
	if board.height >= 13 {
		near.Y = 3
	}
	far := NewVec2(board.width-near.X-1, board.height-near.Y-1)
	middle := NewVec2((board.width-1)/2, (board.height-1)/2)

	if tygemFlag {
		result = []*Vector2{
			near, far, far, near,
		}
	} else {
		result = []*Vector2{
			near, far, near, far,
		}
	}

	if board.width%2 != 0 && board.height%2 != 0 && board.width != 7 && board.height != 7 {
		if count == 5 {
			result = append(result, middle)
		}

		result = append(result, NewVec2(near.X, middle.Y))
		result = append(result, NewVec2(far.X, middle.Y))

		if count == 7 {
			result = append(result, middle)
		}

		result = append(result, NewVec2(middle.X, near.Y))
		result = append(result, NewVec2(middle.X, far.Y))
		result = append(result, NewVec2(middle.X, middle.Y))
	} else if board.width%2 != 0 && board.width != 7 {
		result = append(result, NewVec2(middle.X, near.Y))
		result = append(result, NewVec2(middle.X, far.Y))
	} else if board.height%2 != 0 && board.height != 7 {
		result = append(result, NewVec2(near.X, middle.Y))
		result = append(result, NewVec2(far.X, middle.Y))
	}

	return result[:count]
}

//TODO
func (board *Board) GenerateAscii() string {
	return ""
}

//TODO
func (board *Board) GetPositionHash() string {
	return ""
}

//TODO
func (board *Board) GetHash() string {
	return ""
}
