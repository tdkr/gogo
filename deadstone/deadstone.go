package deadstone

import (
	"github.com/tdkr/gogo/influence"
	"github.com/tdkr/gogo/model"
	"math/rand"
)

func playTillEnd(board *model.Board, sign int32, rnd rand.Rand) {
	illegalVerts := make([]*model.Vector2, 0, 0)
	finished := [2]bool{false, false}
	freeVerts := board.GetVertexBySign(model.StoneSignEmpty)

	for len(freeVerts) > 0 && (!finished[0] || !finished[1]) {
		makeMove := false

		for len(freeVerts) > 0 {
			rndIndex := rand.Int31n(int32(len(freeVerts)))
			rndVec := freeVerts[rndIndex]

			freeVerts[rndIndex] = freeVerts[len(freeVerts)-1]
			freeVerts = freeVerts[:len(freeVerts)-1]

			if deadVerts := board.MakePseudoMove(sign, rndVec); deadVerts != nil {
				freeVerts = append(freeVerts, deadVerts...)

				if sign < 0 {
					finished[0] = false
				} else {
					finished[1] = false
				}

				makeMove = true
				break
			} else {
				illegalVerts = append(illegalVerts, rndVec)
			}
		}

		if sign > 0 {
			finished[0] = !makeMove
		} else {
			finished[1] = !makeMove
		}

		freeVerts = append(freeVerts, illegalVerts...)
		sign = -sign
	}

	// Patch holes

	for i := int32(0); i < board.Width(); i++ {
		for j := int32(0); j < board.Height(); j++ {
			v := model.NewVec2(i, j)

			if board.Get(v) != model.StoneSignEmpty {
				continue
			}

			sign := int32(0)
			for _, nv := range board.GetNeighbors(v, false) {
				if s := board.Get(nv); s != model.StoneSignEmpty {
					sign = s
					break
				}
			}
			if sign != 0 {
				board.Set(v, sign)
			}
		}
	}
}

func getProbalityMap(board *model.Board, iterations int32, rand rand.Rand) [][]float32 {
	result := influence.NewFloatMatrix(int(board.Width()), int(board.Height()), 0)

	whiteSigns := influence.NewFloatMatrix(int(board.Width()), int(board.Height()), 0)
	blackSigns := influence.NewFloatMatrix(int(board.Width()), int(board.Height()), 0)

	for i := int32(0); i < iterations; i++ {
		sign := int32(0)
		if i > iterations/2 {
			sign = model.StoneSignWhite
		} else {
			sign = model.StoneSignBlack
		}

		dupBoard := board.Clone()
		playTillEnd(dupBoard, sign, rand)

		for i := int32(0); i < dupBoard.Width(); i++ {
			for j := int32(0); j < dupBoard.Height(); j++ {
				v := model.NewVec2(i, j)
				s := dupBoard.Get(v)
				if s == model.StoneSignWhite {
					whiteSigns[j][i] += 1
				} else if s == model.StoneSignBlack {
					blackSigns[j][i] += 1
				}
			}
		}
	}

	for i := int32(0); i < board.Height(); i++ {
		for j := int32(0); j < board.Width(); j++ {
			w := whiteSigns[i][j]
			b := blackSigns[i][j]
			if w+b != 0 {
				result[i][j] = b*2.0/(w+b) - 1.0
			}
		}
	}

	return result
}

func Guess(board *model.Board, finished bool, iteration int32, rnd rand.Rand) *model.VecStack {
	var floating *model.VecStack = nil
	if finished {
		floating = board.GetFloatingStones()
		for _, v := range floating.Nodes() {
			board.Set(v, 0)
		}
	} else {
		floating = model.NewVecStack()
	}

	probMap := getProbalityMap(board, iteration, rnd)
	result := model.NewVecStack()
	visited := make(map[int32]interface{})

	for i := int32(0); i < board.Width(); i++ {
		for j := int32(0); j < board.Height(); j++ {
			v := model.NewVec2(i, j)
			sign := board.Get(v)
			if sign == model.StoneSignEmpty || visited[v.HashCode()] != nil {
				continue
			}

			chain := board.GetChain(v)
			prob := float32(0)
			for _, cv := range chain.Nodes() {
				prob += probMap[cv.Y][cv.X]
			}
			prob /= float32(chain.Size())

			dead := int32(influence.GetFloatSign(prob)) == -sign

			for _, cv := range chain.Nodes() {
				if dead {
					result.Push(cv)
				}
				visited[cv.HashCode()] = struct {
				}{}
			}
		}
	}

	if !finished {
		return result
	}

	// Preserve life & death status of related chains

	visited = make(map[int32]interface{})
	newResult := floating

	for _, v := range result.Nodes() {
		if visited[v.HashCode()] != nil {
			continue
		}
		related := board.GetRelatedChain(v)
		deadProb := float32(0)
		deadCnt := float32(0)
		for _, rv := range related {
			if result.Find(rv) > 0 {
				deadCnt++
			}
		}
		deadProb = deadCnt / float32(len(related))
		dead := deadProb > 0.5
		for _, cv := range related {
			if dead {
				newResult.Push(cv)
			}
			visited[cv.HashCode()] = struct {
			}{}
		}
	}

	return newResult
}
