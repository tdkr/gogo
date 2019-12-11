package influence

import (
	"github.com/tdkr/gogo/helper"
	"github.com/tdkr/gogo/model"
)

func GetAreaMap(board *model.Board) map[int32]int32 {
	mark := make(map[int32]interface{})
	result := make(map[int32]int32)

	for i := int32(0); i < board.GetLength(); i++ {
		if mark[i] != nil {
			continue
		}

		if sign := board.GetSign(i); sign != model.CellSignNone {
			mark[i] = sign
			result[i] = sign
			continue
		}

		indicator := int32(1)
		sign := int32(0)
		chain := board.GetChain(i)

		for k, _ := range chain {
			if indicator == 0 {
				break
			}

			for _, v := range board.GetNeighbors(k) {
				ns := board.GetSign(v)
				if ns == model.CellSignNone {
					continue
				}

				if sign == 0 {
					sign = helper.GetIntSign(ns)
				} else if sign != helper.GetIntSign(ns) {
					indicator = 0
					break
				}
			}
		}

		for k, _ := range chain {
			result[k] = sign * indicator
			mark[k] = struct {
			}{}
		}
	}

	return result
}
