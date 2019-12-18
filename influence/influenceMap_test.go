package influence

import (
	"testing"
)

func TestGetInfluenceMap(t *testing.T) {
	t.Run("should return same dimensions of input data", func(t *testing.T) {
		result := GetInfluenceMap(UnfinishedBoard)
		if len(result) != len(UnfinishedBoard) || len(result[0]) != len(UnfinishedBoard[0]) {
			t.Errorf("dimensions not match")
		}
	})

	t.Run("should have same sign as stones on stone vertices", func(t *testing.T) {
		result := GetInfluenceMap(UnfinishedBoard)

		for y := 0; y < len(UnfinishedBoard); y++ {
			succeed := true
			for x := 0; x < len(UnfinishedBoard[y]); x++ {
				if UnfinishedBoard[y][x] != 0 && result[y][x] != UnfinishedBoard[y][x] {
					t.Errorf("stone sign not match, result:%+v, data:%+v", result, UnfinishedBoard)
					succeed = false
					break
				}
			}
			if !succeed {
				break
			}
		}
	})

	t.Run("should return a number between -1 and 1", func(t *testing.T) {
		result := GetInfluenceMap(UnfinishedBoard)

		for y := 0; y < len(UnfinishedBoard); y++ {
			succeed := true
			for x := 0; x < len(UnfinishedBoard[y]); x++ {
				if result[y][x] < -1 || result[y][x] > 1 {
					t.Errorf("should return a number between -1 and 1 : %+v", result)
					succeed = false
					break
				}
			}
			if !succeed {
				break
			}
		}
	})

	t.Run("should return -1, 0, 1 if discrete is set to true", func(t *testing.T) {
		result := GetInfluenceMap(UnfinishedBoard, Discrete(true))

		for y := 0; y < len(UnfinishedBoard); y++ {
			succeed := true
			for x := 0; x < len(UnfinishedBoard[y]); x++ {
				if result[y][x] != -1 || result[y][x] != 0 || result[y][x] != 1 {
					t.Errorf("should return -1, 0, 1 : %+v", result[y][x])
					succeed = false
					break
				}
			}
			if !succeed {
				break
			}
		}
	})
}
