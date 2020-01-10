package model

import (
	"bytes"
	"crypto/md5"
	"math"
	"strconv"
)

const (
	StoneSignNone  = -2
	StoneSignWhite = -1
	StoneSignEmpty = 0
	StoneSignBlack = 1
)

type Board struct {
	width       int32
	height      int32
	captures    []int32
	arrangement [][]int32
}

func NewBoard(opts ...option) *Board {
	o := NewOptions(opts...)
	board := &Board{
		width:    19,
		height:   19,
		captures: []int32{0, 0},
	}
	if o.width != 0 {
		board.width = o.width
	}
	if o.height != 0 {
		board.height = o.height
	}
	board.arrangement = make([][]int32, board.height)
	for i := int32(0); i < board.height; i++ {
		board.arrangement[i] = make([]int32, board.width)
	}
	if o.arrangement != nil {
		for i, v := range o.arrangement {
			for i2, v2 := range v {
				board.arrangement[i][i2] = v2
			}
		}
	}
	if o.captures != nil {
		copy(board.captures, o.captures)
	}
	return board
}

func (board *Board) Width() int32 {
	return board.width
}

func (board *Board) Height() int32 {
	return board.height
}

func (board *Board) Get(vec Vector2) int32 {
	return board.arrangement[vec.Y][vec.X]
}

func (board *Board) Get2(x, y int32) int32 {
	return board.arrangement[y][x]
}

func (board *Board) Set(vec Vector2, sign int32) {
	board.arrangement[vec.Y][vec.X] = sign
}

func (board *Board) Set2(x, y, sign int32) {
	board.arrangement[x][y] = sign
}

func (board *Board) Clone() *Board {
	dup := NewBoard(Width(board.width), Height(board.height), Arrangement(board.arrangement), Captures(board.captures))
	return dup
}

func (board *Board) Diff(b *Board) []Vector2 {
	if b == nil || board == nil || board.width != b.width || board.height != b.height {
		return nil
	}

	result := make([]Vector2, 0)

	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			if board.Get2(x, y) != b.Get2(x, y) {
				result = append(result, Vect2(x, y))
			}
		}
	}
	return result
}

func (board *Board) HasVertex(v Vector2) bool {
	return 0 <= v.X && v.X < board.width && 0 <= v.Y && v.Y < board.height
}

func (board *Board) HasVertex2(x, y int32) bool {
	return 0 <= x && y < board.width && 0 <= y && y < board.height
}

func (board *Board) Clear() {
	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			board.Set2(x, y, StoneSignEmpty)
		}
	}
}

func (board *Board) IsSquare() bool {
	return board.width == board.height
}

func (board *Board) IsEmpty() bool {
	for x := int32(0); x < board.width; x++ {
		for y := int32(0); y < board.height; y++ {
			if board.Get2(x, y) != StoneSignEmpty {
				return false
			}
		}
	}
	return true
}

func (board *Board) GetDistance(v1 Vector2, v2 Vector2) int32 {
	return int32(math.Abs(float64(v1.X-v2.X)) + math.Abs(float64(v1.Y-v2.Y)))
}

func (board *Board) GetNeighbors(vec Vector2, ignoreBoard bool) []Vector2 {
	if !ignoreBoard && !board.HasVertex(vec) {
		return []Vector2{}
	}

	if ignoreBoard {
		return []Vector2{
			Vect2(vec.X-1, vec.Y),
			Vect2(vec.X+1, vec.Y),
			Vect2(vec.X, vec.Y-1),
			Vect2(vec.X, vec.Y+1),
		}
	}

	result := make([]Vector2, 0, 4)
	if vec.X >= 0 {
		if vec.X > 0 {
			result = append(result, Vect2(vec.X-1, vec.Y))
		}
		if vec.X < board.width-1 {
			result = append(result, Vect2(vec.X+1, vec.Y))
		}
	}
	if vec.Y >= 0 {
		if vec.Y > 0 {
			result = append(result, Vect2(vec.X, vec.Y-1))
		}
		if vec.Y < board.height-1 {
			result = append(result, Vect2(vec.X, vec.Y+1))
		}
	}
	return result
}

func (board *Board) getConnectedComponentInner(vec Vector2, signs map[int32]interface{}, result *VecStack) {
	if !board.HasVertex(vec) {
		return
	}

	result.Push(vec)

	for _, v := range board.GetNeighbors(vec, false) {
		if result.Find(v) >= 0 || signs[board.Get(v)] == nil {
			continue
		}

		board.getConnectedComponentInner(v, signs, result)
	}
}

func (board *Board) GetConnectedComponent(vec Vector2, signs map[int32]interface{}) *VecStack {
	result := NewVecStack()
	if signs == nil {
		signs = make(map[int32]interface{})
		signs[board.Get(vec)] = struct {
		}{}
	}
	board.getConnectedComponentInner(vec, signs, result)
	return result
}

func (board *Board) GetChain(vec Vector2) *VecStack {
	return board.GetConnectedComponent(vec, nil)
}

func (board *Board) hasLibertiesInner(vec Vector2, visited map[int32]interface{}, sign int32) bool {
	visited[vec.HashCode()] = struct {
	}{}

	neighbors := board.GetNeighbors(vec, false)
	for _, v := range neighbors {
		s := board.Get(v)
		if s == StoneSignEmpty {
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

func (board *Board) HasLiberties(vec Vector2) bool {
	if !board.HasVertex(vec) {
		return false
	}
	visited := map[int32]interface{}{

	}
	return board.hasLibertiesInner(vec, visited, board.Get(vec))
}

func (board *Board) GetRelatedChain(vec Vector2) []Vector2 {
	if !board.HasVertex(vec) {
		return []Vector2{}
	}

	sign := board.Get(vec)
	if sign == StoneSignEmpty {
		return []Vector2{}
	}

	area := board.GetConnectedComponent(vec, nil).Nodes()

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
			if sign == StoneSignEmpty {
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
			if board.Get2(x, y) == StoneSignEmpty {
				score.Territory[scoreIdx]++
			}
		}
	}

	return score
}

func (board *Board) Vec2Coord(vec Vector2) int32 {
	if !board.HasVertex(vec) {
		return -1
	}

	return vec.X + vec.Y*board.height
}

func (board *Board) Coord2Vec(coord int32) Vector2 {
	x := coord % board.height
	y := coord / board.height
	return Vect2(x, y)
}

//TODO
func (board *Board) IsValid() bool {
	return true
}

func (board *Board) MakeMove(sign int32, vec Vector2) *Board {
	move := board.Clone()

	if sign == StoneSignEmpty || !board.HasVertex(vec) {
		return move
	}

	move.Set(vec, sign)

	deadNeighbors := make([]Vector2, 0, 0)
	for _, v := range move.GetNeighbors(vec, false) {
		if move.Get(v) == -sign && !move.HasLiberties(v) {
			deadNeighbors = append(deadNeighbors, v)
		}
	}

	for _, v := range deadNeighbors {
		if move.Get(v) == StoneSignEmpty {
			continue
		}

		for _, cv := range move.GetChain(v).Nodes() {
			move.Set(cv, 0)
			move.captures[(-sign+1)/2]++
		}
	}

	// Detect suicide
	if len(deadNeighbors) == 0 && !move.HasLiberties(vec) {
		for _, cv := range move.GetChain(vec).Nodes() {
			move.Set(cv, 0)
			move.captures[(sign+1)/2]++
		}
	}

	return move
}

func (board *Board) GetHandicapPlacement(count int32, tygemFlag bool) []Vector2 {
	result := make([]Vector2, 0)

	if board.width <= 6 || board.height <= 6 || count < 2 {
		return result
	}

	near := Vect2(2, 2)
	if board.width >= 13 {
		near.X = 3
	}
	if board.height >= 13 {
		near.Y = 3
	}
	far := Vect2(board.width-near.X-1, board.height-near.Y-1)
	middle := Vect2((board.width-1)/2, (board.height-1)/2)

	if tygemFlag {
		result = []Vector2{
			near, far, far, near,
		}
	} else {
		result = []Vector2{
			near, far, near, far,
		}
	}

	if board.width%2 != 0 && board.height%2 != 0 && board.width != 7 && board.height != 7 {
		if count == 5 {
			result = append(result, middle)
		}

		result = append(result, Vect2(near.X, middle.Y))
		result = append(result, Vect2(far.X, middle.Y))

		if count == 7 {
			result = append(result, middle)
		}

		result = append(result, Vect2(middle.X, near.Y))
		result = append(result, Vect2(middle.X, far.Y))
		result = append(result, Vect2(middle.X, middle.Y))
	} else if board.width%2 != 0 && board.width != 7 {
		result = append(result, Vect2(middle.X, near.Y))
		result = append(result, Vect2(middle.X, far.Y))
	} else if board.height%2 != 0 && board.height != 7 {
		result = append(result, Vect2(near.X, middle.Y))
		result = append(result, Vect2(far.X, middle.Y))
	}

	return result[:count]
}

//TODO
func (board *Board) GenerateAscii() string {
	return ""
}

func (board *Board) GetPositionHash() [16]byte {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteByte('{')
	for _, v := range board.arrangement {
		buffer.WriteByte('{')
		for _, v2 := range v {
			buffer.WriteString(strconv.Itoa(int(v2)))
		}
		buffer.WriteByte('}')
	}
	buffer.WriteByte('}')
	return md5.Sum(buffer.Bytes())
}

//TODO
func (board *Board) GetHash() string {
	return ""
}

func (board *Board) Equals(b *Board) bool {
	if board == nil || b == nil || board.height != b.height || board.width != b.width {
		return false
	}
	for i := int32(0); i < board.width; i++ {
		for j := int32(0); j < board.height; j++ {
			if board.arrangement[i][j] != b.arrangement[i][j] {
				return false
			}
		}
	}
	return true
}

func (board *Board) GetVertexBySign(sign int32) []Vector2 {
	ret := make([]Vector2, 0)
	for i := int32(0); i < board.width; i++ {
		for j := int32(0); j < board.height; j++ {
			if board.arrangement[i][j] == sign {
				ret = append(ret, Vect2(i, j))
			}
		}
	}
	return ret
}

func (board *Board) MakePseudoMove(sign int32, vec Vector2) []Vector2 {
	neighbors := board.GetNeighbors(vec, false)

	isEye := true
	for _, v := range neighbors {
		if board.Get(v) != sign {
			isEye = false
			break
		}
	}
	if isEye {
		return nil
	}

	checkCapture := false
	checkMultiDeadChains := false
	board.Set(vec, sign)

	if !board.HasLiberties(vec) {
		isPointChain := true
		for _, v := range neighbors {
			if board.Get(v) == sign {
				isPointChain = false
				break
			}
		}

		checkMultiDeadChains = isPointChain
		checkCapture = !isPointChain
	}

	dead := make([]Vector2, 0)
	deadChains := 0

	for _, v := range neighbors {
		if board.Get(v) != -sign || board.HasLiberties(v) {
			continue
		}

		chain := board.GetChain(v)
		deadChains += 1

		for _, cv := range chain.Nodes() {
			board.Set(cv, 0)
			dead = append(dead, cv)
		}
	}

	if (checkMultiDeadChains && deadChains <= 1) ||
		(checkCapture && len(dead) == 0) {
		for _, v := range dead {
			board.Set(v, -sign)
		}
		board.Set(vec, 0)
		return nil
	}

	return dead
}

func (board *Board) GetFloatingStones() *VecStack {
	visited := make(map[int32]interface{})
	result := NewVecStack()

	for i := int32(0); i < board.width; i++ {
		for j := int32(0); j < board.height; j++ {
			v := Vect2(i, j)

			if board.Get(v) != StoneSignEmpty || visited[v.HashCode()] != nil {
				continue
			}

			posArea := board.GetConnectedComponent(v, map[int32]interface{}{
				StoneSignWhite: struct {
				}{},
				StoneSignEmpty: struct {
				}{},
			})
			negArea := board.GetConnectedComponent(v, map[int32]interface{}{
				StoneSignEmpty: struct {
				}{},
				StoneSignBlack: struct {
				}{},
			})
			posDead := make([]Vector2, 0)
			negDead := make([]Vector2, 0)
			posDiff, negDiff := 0, 0
			for _, v := range posArea.Nodes() {
				if board.Get(v) == StoneSignWhite {
					posDead = append(posDead, v)
				} else if negArea.Find(v) < 0 {
					posDiff++
				}
			}
			for _, v := range negArea.Nodes() {
				if board.Get(v) == StoneSignBlack {
					negDead = append(negDead, v)
				} else if posArea.Find(v) < 0 {
					negDiff++
				}
			}

			favorNeg := negDiff <= 1 && len(negDead) <= len(posDead)
			favorPos := posDiff <= 1 && len(posDead) <= len(negDead)

			var actualArea *VecStack = nil
			var actualDead []Vector2 = nil
			if !favorNeg && favorPos {
				actualArea = posArea
				actualDead = posDead
			} else if favorNeg && !favorPos {
				actualArea = negArea
				actualDead = negDead
			} else {
				actualArea = board.GetChain(v)
			}
			for _, v := range actualArea.Nodes() {
				visited[v.HashCode()] = struct {
				}{}
			}
			if actualDead != nil {
				for _, v := range actualDead {
					result.Push(v)
				}
			}
		}
	}
	return result
}

func (board *Board) SetCapture(vec Vector2) {
	sign := board.Get(vec)
	if board.Get(vec) == StoneSignEmpty {
		return
	}
	board.captures[(sign+1)/2]++
	board.Set(vec, StoneSignEmpty)
}

func (board *Board) Captures() []int32 {
	return board.captures
}
