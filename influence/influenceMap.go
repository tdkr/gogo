package influence

import "github.com/tdkr/gogo/model"

func GetInfluenceMap(board *model.Board, opts ...option) {
	o := NewOptions(opts...)

	areaMap := GetAreaMap(board)
}
