package influence

type Options struct {
	discrete      bool
	maxDistance   float32
	minRadiance   float32
	radianceVar1  int32
	radicanceVar2 float32
	radicanceVar3 float32
}

type option func(o *Options)

func Discrete(val bool) option {
	return func(o *Options) {
		o.discrete = val
	}
}

func MaxDistance(val float32) option {
	return func(o *Options) {
		o.maxDistance = val
	}
}

func MinRadiance(val float32) option {
	return func(o *Options) {
		o.minRadiance = val
	}
}

func RadianceVar1(val int32) option {
	return func(o *Options) {
		o.radianceVar1 = val
	}
}

func RadicanceVar2(val float32) option {
	return func(o *Options) {
		o.radicanceVar2 = val
	}
}

func RadicanceVar3(val float32) option {
	return func(o *Options) {
		o.radicanceVar3 = val
	}
}

func NewOptions(opts ...option) *Options {
	o := &Options{
		discrete:      false,
		maxDistance:   6,
		minRadiance:   2,
		radianceVar1:  6,
		radicanceVar2: 1.5,
		radicanceVar3: 2,
	}
	for _, v := range opts {
		v(o)
	}
	return o
}
