package influence

import (
	"reflect"
	"testing"
)

func TestGetRadianceMap(t *testing.T) {
	type args struct {
		board [][]float32
		sign  float32
		opts  []option
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
			if got := GetRadianceMap(tt.args.board, tt.args.sign, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRadianceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
