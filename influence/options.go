package influence

type Options struct {
	discrete    bool
	maxDistance int32
	minRadiance int32
	p1          int32
	p2          float32
	p3          int32
}

type option func(o *Options)

func Discrete(val bool) option {
	return func(o *Options) {
		o.discrete = val
	}
}

func MaxDistance(val int32) option {
	return func(o *Options) {
		o.maxDistance = val
	}
}

func MinRadiance(val int32) option {
	return func(o *Options) {
		o.minRadiance = val
	}
}

func P1(val int32) option {
	return func(o *Options) {
		o.p1 = val
	}
}

func P2(val float32) option {
	return func(o *Options) {
		o.p2 = val
	}
}

func P3(val int32) option {
	return func(o *Options) {
		o.p3 = val
	}
}

func NewOptions(opts ...option) *Options {
	o := &Options{}
	for _, v := range opts {
		v(o)
	}
	return o
}
