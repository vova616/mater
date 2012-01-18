package collision

import (
	"github.com/teomat/mater/vect"
)

type RayCastInput struct {
	MaxFraction    float64
	Point1, Point2 vect.Vect
}

type RayCastOutput struct {
	Fraction float64
	Normal   vect.Vect
}
