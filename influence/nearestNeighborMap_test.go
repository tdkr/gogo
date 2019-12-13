package influence

import (
	"testing"
)

func TestGetNearestNeighborMap(t *testing.T) {
	t.Run("should return same dimensions of input data", func(t *testing.T) {
		result := GetNearestNeighborMap(UnfinishedBoard, 1)
		if len(UnfinishedBoard) != len(result) || len(UnfinishedBoard[0]) != len(result[0]) {
			t.Errorf("dimensions not match, %d, %d, %d, %d", len(UnfinishedBoard), len(result), len(UnfinishedBoard[0]), len(result[0]))
		}
	})
	
	t.Run("only stone positions of the same color should have value 0", func(t *testing.T) {
		sign := float32(-1)
		result := GetNearestNeighborMap(UnfinishedBoard, sign)

		for y := 0; y < len(UnfinishedBoard); y++ {
			for x := 0; x < len(UnfinishedBoard[y]); x++ {
				if UnfinishedBoard[y][x] == sign {
					if result[y][x] != 0 {
						t.Errorf("should be 0!")
					}
				} else {
					if result[y][x] == 0 {
						t.Errorf("should not be 0!")
					}
				}
			}
		}
	})
}
