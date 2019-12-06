package logic

import (
	"github.com/tdkr/gogo/model"
	"math/rand"
)

func playTillEnd(board *model.Board, sign int32, rand rand.Rand) {
	illegalVts := make([]int32, 0)
	finished := [2]bool{false, false}
	freeVts := board.GetVertexBySign(model.CellSignNone)

	for len(freeVts) > 0 && (!finished[0] || !finished[1]) {
		madeMove := false

		for len(freeVts) > 0 {
			rndIndex := rand.Int31n(int32(len(freeVts)))
			freeVertex := freeVts[rndIndex]

			freeVts = append(freeVts[:rndIndex], freeVts[rndIndex+1:]...)

			if ret := board.MakePseudoMove(sign, freeVertex); ret != nil {
				freeVts = append(freeVts, ret...)

				if sign < 0 {
					finished[0] = false
				} else {
					finished[1] = false
				}

				madeMove = true
				break
			} else {
				illegalVts = append(illegalVts, freeVertex)
			}
		}

		if sign > 0 {
			finished[0] = !madeMove
		} else {
			finished[1] = !madeMove
		}

		freeVts = append(freeVts, illegalVts...)
		sign = -sign
	}

	// Patch Holes

	for _, v := range board.GetVertexBySign(model.CellSignNone) {
		sign := int32(0)

		for _, v := range board.GetNeighbors(v) {
			s := board.GetSign(v)

			if s == model.CellSignBlack || s == model.CellSignWhite {
				sign = s
				break
			}
		}

		if sign != model.CellSignNone {
			board.SetSign(v, sign)
		}
	}
}

func getProbabilityMap(board *model.Board, iteration int32, rand rand.Rand) map[int32]float32 {
	mark := make(map[int32][2]int32)

	for i := int32(0); i < iteration; i++ {
		sign := int32(model.CellSignBlack)
		if i < iteration/2 {
			sign = model.CellSignWhite
		}

		dupBoard := board.Clone()
		playTillEnd(dupBoard, sign, rand)

		for i := int32(0); i < board.GetSize()*board.GetSize(); i++ {
			sign := board.GetSign(i)

			if sign == model.CellSignNone {
				continue
			} else if sign == model.CellSignWhite {
				mark[i][0] += 1
			} else if sign == model.CellSignBlack {
				mark[i][1] += 1
			}
		}
	}

	result := make(map[int32]float32, 0)
	for i, v := range mark {
		t := v[0] + v[1]
		if t == 0 {
			continue
		}
		result[i] = float32(v[0])*2.0/float32(t) - 1.0
	}

	return result
}

func Estimate(board *model.Board, finished bool, iterations int32, rand rand.Rand) []int32 {
	result := make([]int32, 0)

	var floating []int32 = nil
	if finished {
		floating = board.GetFloatingStones()

		for _, v := range floating {
			board.SetSign(v, 0)
		}
	} else {
		floating = make([]int32, 0)
	}

	probMap := getProbabilityMap(board, iterations, rand)
	visited := make(map[int32]interface{})

	for i, v := range probMap {
		sign := board.GetSign(i)
		if sign == model.CellSignNone || visited[i] != nil {
			continue
		}

		chain := board.GetChain(i)
		probability := 0.0
		for i,v := range chain {

		}
	}

	return result
}
