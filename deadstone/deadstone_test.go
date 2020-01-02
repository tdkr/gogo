package deadstone

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/tdkr/gogo/model"
)

var finishedBoard = [][]int32{
	{0, 0, 0, -1, -1, -1, 1, 0, 1, 1, -1, -1, 0, -1, 0, -1, -1, 1, 0},
	{0, 0, -1, 0, -1, 1, 1, 1, 0, 1, -1, 0, -1, -1, -1, -1, 1, 1, 0},
	{0, 0, -1, -1, -1, 1, 1, 0, 0, 1, 1, -1, -1, 1, -1, 1, 0, 1, 0},
	{0, 0, 0, 0, -1, -1, 1, 0, 1, -1, 1, 1, 1, 1, 1, 0, 1, 0, 0},
	{0, 0, 0, 0, -1, 0, -1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0},
	{0, 0, -1, 0, 0, -1, -1, 1, 0, -1, -1, 1, -1, -1, 0, 1, 0, 0, 1},
	{0, 0, 0, -1, -1, 1, 1, 1, 1, 1, 1, 1, 1, -1, -1, -1, 1, 1, 1},
	{0, 0, -1, 1, 1, 0, 1, -1, -1, 1, 0, 1, -1, 0, 1, -1, -1, -1, 1},
	{0, 0, -1, -1, 1, 1, 1, 0, -1, 1, -1, -1, 0, -1, -1, 1, 1, 1, 1},
	{0, 0, -1, 1, 1, -1, -1, -1, -1, 1, 1, 1, -1, -1, -1, -1, 1, -1, -1},
	{-1, -1, -1, -1, 1, 1, 1, -1, 0, -1, 1, -1, -1, 0, -1, 1, 1, -1, 0},
	{-1, 1, -1, 0, -1, -1, -1, -1, -1, -1, 1, -1, 0, -1, -1, 1, -1, 0, -1},
	{1, 1, 1, 1, -1, 1, 1, 1, -1, 1, 0, 1, -1, 0, -1, 1, -1, -1, 0},
	{0, 1, -1, 1, 1, -1, -1, 1, -1, 1, 1, 1, -1, 1, -1, 1, 1, -1, 1},
	{0, 0, -1, 1, 0, 0, 1, 1, -1, -1, 0, 1, -1, 1, -1, 1, -1, 0, -1},
	{0, 0, 1, 0, 1, 0, 1, 1, 1, -1, -1, 1, -1, -1, 1, -1, -1, -1, 0},
	{0, 0, 0, 0, 1, 1, 0, 1, -1, 0, -1, -1, 1, 1, 1, 1, -1, -1, -1},
	{0, 0, 1, 1, -1, 1, 1, -1, 0, -1, -1, 1, 1, 1, 1, 0, 1, -1, 1},
	{0, 0, 0, 1, -1, -1, -1, -1, -1, 0, -1, -1, 1, 1, 0, 1, 1, 1, 0},
}

var unfinishedBoard = [][]int32{
	{0, -1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	{1, -1, -1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0},
	{-1, 0, -1, -1, 1, 1, 0, 1, 0, 1, 1, 0, -1, 0, -1, 0, -1, 1, 0},
	{-1, 0, 0, 0, -1, -1, -1, 1, 0, -1, 0, 0, 0, 0, 0, -1, 1, 0, 0},
	{0, -1, 0, -1, 1, 1, 1, -1, -1, 0, 0, -1, 0, 0, -1, 1, 1, 0, 0},
	{-1, 0, -1, -1, 0, 0, 0, -1, 0, 0, -1, 0, 0, 0, -1, 1, 0, 1, 0},
	{-1, -1, 1, -1, 1, 1, 1, -1, 0, 1, 1, -1, 0, 0, 1, 1, -1, 0, 0},
	{0, 1, 0, -1, 0, -1, 1, -1, 0, 1, -1, 1, 0, 0, 0, 0, -1, 0, 1},
	{0, 0, 0, -1, 0, -1, 1, 1, 0, 1, -1, -1, 0, 0, 0, 1, 1, -1, 0},
	{0, 0, 1, 1, 1, -1, -1, 1, 0, 0, 1, 0, -1, -1, 1, 1, -1, -1, -1},
	{0, -1, -1, -1, 1, 0, 0, 1, 0, 1, 0, 0, -1, 1, 0, 1, 1, -1, 0},
	{0, -1, 1, 1, 1, 1, 0, -1, 1, 0, 0, 0, -1, 1, 0, 0, 1, -1, 0},
	{0, 0, 0, 0, 0, -1, 0, -1, -1, 0, 0, 0, -1, 1, 0, 0, 1, -1, 0},
	{0, -1, -1, 0, -1, 0, 0, 0, 0, 0, 0, 0, -1, -1, 1, 1, -1, -1, 0},
	{0, 0, 1, -1, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, -1, -1, 1, 1, -1},
	{0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, -1, 0, 1, 0, 1, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 1, 0, -1, 0, 0, -1, 1, 0, 0, 1, 0, 1, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

func Test_playTillEnd(t *testing.T) {
	type args struct {
		board *model.Board
		sign  int32
		rnd   *rand.Rand
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				board: model.NewBoard(model.Height(19), model.Width(19), model.Arrangement(unfinishedBoard)),
				sign:  -1,
				rnd:   rand.New(rand.NewSource(time.Now().UnixNano())),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dup := tt.args.board
			playTillEnd(tt.args.board, tt.args.sign, tt.args.rnd)
			fmt.Println("playTillEnd", tt.args.sign, tt.args.board, dup)
		})
	}
}

func Test_getProbabilityMap(t *testing.T) {
	type args struct {
		board      *model.Board
		iterations int32
		rand       *rand.Rand
	}
	tests := []struct {
		name string
		args args
		want [][]float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProbabilityMap(tt.args.board, tt.args.iterations, tt.args.rand); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getProbabilityMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGuess(t *testing.T) {
	type args struct {
		board     *model.Board
		finished  bool
		iteration int32
		rnd       *rand.Rand
	}
	tests := []struct {
		name string
		args args
		want *model.VecStack
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Guess(tt.args.board, tt.args.finished, tt.args.iteration, tt.args.rnd); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Guess() = %v, want %v", got, tt.want)
			}
		})
	}
}
