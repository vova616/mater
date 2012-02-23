package collision

import (
	"github.com/teomat/mater/vect"
)

/*type FeaturePair struct {
	InEdge1  uint8
	OutEdge1 uint8
	InEdge2  uint8
	OutEdge2 uint8
}

// unsafe pointer magic because go doesn't have unions
func (fp *FeaturePair) Value() int32 {
	return *(*int32)(unsafe.Pointer(fp))
}*/

// Contact point between 2 shapes.
type Contact struct {
	Position vect.Vect
	Normal   vect.Vect
	Dist     float64

	R1, R2   vect.Vect
	nMass float64
	tMass float64
	bounce float64

	jnAcc float64
	jtAcc float64
	jBias float64

	bias float64

	hash hashValue
}

func (con *Contact) reset(pos, norm vect.Vect, dist float64, hash hashValue) {
	con.Position = pos
	con.Normal = norm
	con.Dist = dist
	con.hash = hash

	con.jnAcc = 0.0
	con.jtAcc = 0.0
	con.jBias = 0.0
}
