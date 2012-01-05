package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
	"log"
	"math"
)

type PolygonAxis struct {
	n vect.Vect
	d float64
}

type PolygonShape struct {
	verts  Vertices
	tVerts Vertices
	axes   []PolygonAxis
	tAxes  []PolygonAxis
	numVerts int
}

func NewPolygon(verts Vertices, offset vect.Vect) *Shape {
	if verts == nil {
		log.Printf("Error: no vertices passed!")
		return nil
	}

	shape := new(Shape)
	poly := &PolygonShape{}

	poly.SetVerts(verts, offset)

	shape.ShapeClass = poly
	return shape
}

func (poly *PolygonShape) SetVerts(verts Vertices, offset vect.Vect) {

	if verts == nil {
		log.Printf("Error: no vertices passed!")
		return
	}

	if verts.ValidatePolygon() == false {
		log.Printf("Warning: vertices not valid")
	}

	numVerts := len(verts)
	oldnumVerts := len(poly.verts)
	poly.numVerts = numVerts

	if oldnumVerts < numVerts {
		//create new slices
		poly.verts = make(Vertices, numVerts)
		poly.tVerts = make(Vertices, numVerts)
		poly.axes = make([]PolygonAxis, numVerts)
		poly.tAxes = make([]PolygonAxis, numVerts)

	} else {
		//reuse old slices
		poly.verts = poly.verts[:numVerts]
		poly.tVerts = poly.tVerts[:numVerts]
		poly.axes = poly.axes[:numVerts]
		poly.tAxes = poly.tAxes[:numVerts]
	}

	for i := 0; i < numVerts; i++ {
		a := vect.Add(offset, verts[i])
		b := vect.Add(offset, verts[(i + 1) % numVerts])
		n := vect.Normalize(vect.Perp(vect.Sub(b, a)))

		poly.verts[i] = a
		poly.axes[i].n = n
		poly.axes[i].d = vect.Dot(n, a)
	}
}

func (poly *PolygonShape) ShapeType() ShapeType {
	return ShapeType_Polygon
}

func (poly *PolygonShape) Update(xf transform.Transform) aabb.AABB {
	//transform axes
	{
		src := poly.axes
		dst := poly.tAxes

		for i := 0; i < poly.numVerts; i++ {
			n := xf.RotateVect(src[i].n)
			dst[i].n = n
			dst[i].d = vect.Dot(xf.Position, n) + src[i].d
		}
	}
	//transform verts
	{
		inf := math.Inf(1)
		aabb := aabb.AABB{
			Lower: vect.Vect{-inf, -inf},
			Upper: vect.Vect{ inf,  inf},
		}

		src := poly.verts
		dst := poly.tVerts

		for i := 0; i < poly.numVerts; i++ {
			v := xf.TransformVect(src[i])

			dst[i] = v
			aabb.Lower.X = math.Min(aabb.Lower.X, v.X)
			aabb.Upper.X = math.Max(aabb.Upper.X, v.X)
			aabb.Lower.Y = math.Min(aabb.Lower.Y, v.Y)
			aabb.Upper.Y = math.Max(aabb.Upper.Y, v.Y)
		}

		return aabb
	}
}

func (poly *PolygonShape) TestPoint(xf transform.Transform, point vect.Vect) bool {
	panic("Not yet implemented!")
}
