package influence

import (
	"fmt"
	"math"

	"github.com/tdkr/gogo/model"
)

func GetInfluenceMap(board [][]float32, opts ...option) [][]float32 {
	o := NewOptions(opts...)
	width, height := GetMatrixSize(board)

	areaMap := GetAreaMap(board)

	pnnMap := GetNearestNeighborMap(board, 1)
	nnnMap := GetNearestNeighborMap(board, -1)

	prMap := GetRadianceMap(board, 1)
	nrMap := GetRadianceMap(board, -1)

	inf := float32(math.MaxFloat32)
	max := -inf
	min := inf

	result := DuplicateFloatMatrix(areaMap)

	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {
			if result[y][x] != 0 {
				continue
			}

			s := GetFloatSign(nnnMap[y][x] - pnnMap[y][x])
			faraway := false
			dim := false
			if s == 0 {
				faraway = true
				dim = true
			} else if s > 0 {
				faraway = pnnMap[y][x] > o.maxDistance
				dim = float32(math.Round(float64(prMap[y][x]))) < o.minRadiance
			} else {
				faraway = nnnMap[y][x] > o.maxDistance
				dim = float32(math.Round(float64(nrMap[y][x]))) < o.minRadiance
			}

			if faraway || dim {
				result[y][x] = 0
			} else {
				if s > 0 {
					result[y][x] = float32(s) * prMap[y][x]
				} else {
					result[y][x] = float32(s) * nrMap[y][x]
				}
			}

			if result[y][x] > max {
				max = result[y][x]
			}
			if result[y][x] < min {
				min = result[y][x]
			}

			if o.discrete {
				result[y][x] = GetFloatSign(result[y][x])
			}

			fmt.Println("iterate1", x, y, result[y][x], max, min)
		}
	}

	// Postprocessing

	for x := int32(0); x < width; x++ {
		for y := int32(0); y < height; y++ {
			if areaMap[y][x] != 0 {
				continue
			}

			v := model.NewVec2(x, y)

			// Prevent single point areas

			mSign := GetFloatSign(result[y][x])

			if mSign != 0 {
				cnt1, cnt2 := 0, 0
				for _, v := range getNeighbors(v) {
					if isValidVertex(board, int(v.X), int(v.Y)) {
						cnt1++
						if GetFloatSign(board[v.Y][v.X]) == mSign {
							cnt2++
						}
					}
				}

				if cnt1 >= 2 && cnt2 == cnt1 {
					result[y][x] = 0
					continue
				}
			}

			// Fix ragged areas

			if mSign != 0 {
				posNeighbors := make([]*model.Vector2, 0)
				for _, nv := range getNeighbors(v) {
					if isValidVertex(board, int(nv.X), int(nv.Y)) && GetFloatSign(result[nv.Y][nv.X]) == mSign {
						posNeighbors = append(posNeighbors, nv)
					}
				}

				if len(posNeighbors) == 1 {
					pv := posNeighbors[0]

					if board[pv.Y][pv.X] == mSign {
						result[y][x] = 0
						continue
					}
				}
			}

			// Fix empty pillars

			distance := MinInt(x, y, width-x-1, height-y-1)

			if distance <= 2 && mSign == 0 {
				signedNeighbors := make([]*model.Vector2, 0)
				for _, nv := range getNeighbors(v) {
					if isValidVertex(result, int(nv.X), int(nv.Y)) && result[nv.Y][nv.X] != 0 {
						signedNeighbors = append(signedNeighbors, nv)
					}
				}

				if len(signedNeighbors) >= 2 {
					v1, v2 := signedNeighbors[0], signedNeighbors[1]
					s := GetFloatSign(result[v1.Y][v1.X])

					if len(signedNeighbors) >= 3 || v1.X == v2.X || v1.Y == v2.Y {
						flag := true
						for _, sv := range signedNeighbors {
							if result[sv.Y][sv.X] != s {
								flag = false
								break
							}
						}
						if flag {
							if o.discrete {
								result[y][x] = s
							} else {
								result[y][x] = result[signedNeighbors[0].Y][signedNeighbors[0].X]
								for i := 1; i < len(signedNeighbors); i++ {
									sv := signedNeighbors[i]
									if result[sv.Y][sv.X] > result[y][x] {
										result[y][x] = result[sv.Y][sv.X]
									}
								}
							}
							mSign = s
						}
					}
				}
			}

			// Normalize

			if !o.discrete {
				if mSign > 0 {
					result[y][x] = MinFloat(result[y][x]/max, 1)
				} else {
					result[y][x] = MaxFloat(-result[y][x]/min, -1)
				}
			}

			// fmt.Println("iterate2", x, y, result[y][x], max, min)
		}
	}

	return GetAreaMap(result)
}
