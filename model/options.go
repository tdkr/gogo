package model

type Options struct {
	width       int32
	height      int32
	arrangement [][]int32
	captures    []int32
}

type option func(o *Options)

func Width(val int32) option {
	return func(o *Options) {
		o.width = val
	}
}

func Height(val int32) option {
	return func(o *Options) {
		o.height = val
	}
}

func arrangement(val [][]int32) option {
	return func(o *Options) {
		o.arrangement = val
	}
}

func captures(val []int32) option {
	return func(o *Options) {
		o.captures = val
	}
}

func NewOptions(opts ...option) *Options {
	o := &Options{}
	for _, v := range opts {
		v(o)
	}
	return o
}
