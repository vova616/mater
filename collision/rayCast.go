package collision

import (
	"github.com/teomat/mater/vect"
)

// Input for raycast queries.
type RayCastInput struct {
	MaxFraction    float64
	Point1, Point2 vect.Vect
}

// Output for raycast queries.
type RayCastOutput struct {
	Fraction float64
	Normal   vect.Vect
}
