package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
	"log"
	"math"
)

type PolygonAxis struct {
	//The axis normal.
	N vect.Vect
	//Dunno what this is.
	D float64
}

//Don't modify directly or you'll fuck shit up.
//Seriously.
type PolygonShape struct {
	//The raw vertices of the polygon.
	Verts  Vertices
	//The transformed vertices.
	TVerts Vertices
	//The axes of the polygon.
	Axes   []PolygonAxis
	//The transformed axes of the polygon
	TAxes  []PolygonAxis
	//The number of vertices.
	NumVerts int
}

//Creates a new PolygonShape with the given vertices offset by offset.
//Returns nil if the given vertices are not valid
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

//Sets the vertices offset by the offset and calculates the PolygonAxes.
func (poly *PolygonShape) SetVerts(verts Vertices, offset vect.Vect) {

	if verts == nil {
		log.Printf("Error: no vertices passed!")
		return
	}

	if verts.ValidatePolygon() == false {
		log.Printf("Warning: vertices not valid")
	}

	numVerts := len(verts)
	oldnumVerts := len(poly.Verts)
	poly.NumVerts = numVerts

	if oldnumVerts < numVerts {
		//create new slices
		poly.Verts = make(Vertices, numVerts)
		poly.TVerts = make(Vertices, numVerts)
		poly.Axes = make([]PolygonAxis, numVerts)
		poly.TAxes = make([]PolygonAxis, numVerts)

	} else {
		//reuse old slices
		poly.Verts = poly.Verts[:numVerts]
		poly.TVerts = poly.TVerts[:numVerts]
		poly.Axes = poly.Axes[:numVerts]
		poly.TAxes = poly.TAxes[:numVerts]
	}

	for i := 0; i < numVerts; i++ {
		a := vect.Add(offset, verts[i])
		b := vect.Add(offset, verts[(i + 1) % numVerts])
		n := vect.Normalize(vect.Perp(vect.Sub(b, a)))

		poly.Verts[i] = a
		poly.Axes[i].N = n
		poly.Axes[i].D = vect.Dot(n, a)
	}
}

func (poly *PolygonShape) ShapeType() ShapeType {
	return ShapeType_Polygon
}

//Calculates the transformed vertices and axes and the bounding box.
func (poly *PolygonShape) Update(xf transform.Transform) aabb.AABB {
	//transform axes
	{
		src := poly.Axes
		dst := poly.TAxes

		for i := 0; i < poly.NumVerts; i++ {
			n := xf.RotateVect(src[i].N)
			dst[i].N = n
			dst[i].D = vect.Dot(xf.Position, n) + src[i].D
		}
	}
	//transform verts
	{
		inf := math.Inf(1)
		aabb := aabb.AABB{
			Lower: vect.Vect{ inf,  inf},
			Upper: vect.Vect{-inf, -inf},
		}

		src := poly.Verts
		dst := poly.TVerts

		for i := 0; i < poly.NumVerts; i++ {
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
