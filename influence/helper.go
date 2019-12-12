package influence

import "github.com/tdkr/gogo/model"

func getNeighbors(vec model.Vector2) []*model.Vector2 {
	return []*model.Vector2{
		{vec.X - 1, vec.Y},
		{vec.X + 1, vec.Y},
		{vec.X, vec.Y - 1},
		{vec.X, vec.Y + 1},
	}
}
