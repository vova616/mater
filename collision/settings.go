package collision

type settings struct {
	AccumulateImpulses bool
	PositionCorrection bool
	Iterations         int
	AutoUpdateShapes   bool
	AutoClearForces    bool
	AABBExtension      float64
	AABBMultiplier     float64
}

var Settings settings = settings{
	AccumulateImpulses: true,
	PositionCorrection: true,
	Iterations:         1,
	AutoUpdateShapes:   true,
	AutoClearForces:    true,
	AABBExtension:      0.1,
	AABBMultiplier:     2.0,
}
