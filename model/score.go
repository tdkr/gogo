package model

const (
	ScoreIndexBlack = 0
	ScoreIndexWhite = 1
)

type Score struct {
	Area      []int32
	Territory []int32
	Captures  []int32
	Komi      float32
	Handicap  int32
}

func NewScore(komi float32, handicap int32, captures []int32) *Score {
	s := &Score{
		Komi:     komi,
		Handicap: handicap,
		Captures: captures,
	}
	return s
}

func (s *Score) AreaScore() float32 {
	return float32(s.Area[ScoreIndexBlack]) - float32(s.Area[ScoreIndexWhite]) - s.Komi - float32(s.Handicap)
}

func (s *Score) TerritoryScore() float32 {
	return float32(s.Territory[ScoreIndexBlack]) -
		float32(s.Territory[ScoreIndexWhite]) +
		float32(s.Captures[ScoreIndexBlack]) -
		float32(s.Captures[ScoreIndexWhite]) - s.Komi
}
