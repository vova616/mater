package collision

import (
	"math"
)

type settings struct {
	AccumulateImpulses bool
	Iterations         int
	AutoClearForces    bool
	AABBExtension      float64
	AABBMultiplier     float64
	CollisionSlop      float64
	CollisionBias      float64
	PositionCorrection bool
}

var Settings settings = settings{
	AccumulateImpulses: true,
	Iterations:         1,
	AutoClearForces:    true,
	AABBExtension:      0.1,
	AABBMultiplier:     2.0,
	CollisionSlop:      0.1,
	CollisionBias:      math.Pow(1.0-0.1, 60.0),
	PositionCorrection: true,
}
