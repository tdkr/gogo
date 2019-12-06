package model

func filterMap(m map[int32]interface{}, f func(k int32, v interface{}) bool) map[int32]interface{} {
	result := make(map[int32]interface{})
	for k, v := range m {
		if f(k, v) {
			result[k] = v
		}
	}
	return result
}
