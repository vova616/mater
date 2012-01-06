package collision

import (
	"mater/vect"
	"mater/aabb"
	"mater/transform"
)

//Convenience wrapper around PolygonShape.
type BoxShape struct {
	Polygon *PolygonShape
	verts [4]vect.Vect
	Width float64
	Height float64
	Position vect.Vect
}

func NewBoxShape(pos vect.Vect, w, h float64) *Shape {
	shape := new(Shape)
	box := &BoxShape{
		Polygon: &PolygonShape{},
		Width: w,
		Height: h,
		Position: pos,
	}

	hw := w / 2.0
	hh := h / 2.0
	box.verts = [4]vect.Vect{
		{-hw, -hh},
		{-hw,  hh},
		{ hw,  hh},
		{ hw, -hh},
	}

	poly := box.Polygon
	poly.SetVerts(box.verts[:], box.Position)

	shape.ShapeClass = box
	return shape
}

//Recalculates the internal Polygon with new Width, Height and Position.
func (box *BoxShape) UpdatePoly() {
	hw := box.Width / 2.0
	hh := box.Height / 2.0
	box.verts = [4]vect.Vect{
		{-hw, -hh},
		{-hw,  hh},
		{ hw,  hh},
		{ hw, -hh},
	}

	poly := box.Polygon
	poly.SetVerts(box.verts[:], box.Position)
}

func (box *BoxShape) ShapeType() ShapeType {
	return ShapeType_Box
}

//Recalculates the transformed vertices and axes and the bounding box.
func (box *BoxShape) Update(xf transform.Transform) aabb.AABB {
	return box.Polygon.Update(xf)
}

func (box *BoxShape) TestPoint(xf transform.Transform, point vect.Vect) bool {
	return box.Polygon.TestPoint(xf, point)
}
