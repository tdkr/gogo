package helper

import "math"

func GetIntSign(value int32) int32 {
	switch {
	case value > 0:
		return 1
	case value < 0:
		return -1
	case value == 0:
		return 0
	}
	return 0
}

func GetFloatSign(value float32) int32 {
	switch math.Signbit(float64(value)) {
	case true:
		return -1
	case false:
		return 1
	}
	return 0
}
