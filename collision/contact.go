package collision

import (
	"mater/vect"
	"unsafe"
)

type FeaturePair struct {
	InEdge1 uint8
	OutEdge1 uint8
	InEdge2 uint8
	OutEdge2 uint8
}

//unsafe pointer magic because go doesn't have unions
func (fp *FeaturePair) Value() int32 {
	return *(*int32)(unsafe.Pointer(fp))
}

type Contact struct {
	Position vect.Vect
	Normal vect.Vect
	R1, R2 vect.Vect

	Separation float64
	Pn float64	// accumulated normal impulse
	Pt float64	// accumulated tangent impulse
	Pnb float64	// accumulated normal impulse for position bias
	MassNormal, MassTangent float64
	Bias float64
	Feature FeaturePair
}
