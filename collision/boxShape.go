package collision

import (
	"github.com/teomat/mater/aabb"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
)

// Convenience wrapper around PolygonShape.
type BoxShape struct {
	// The polygon that represents this box. Do not touch!
	Polygon *PolygonShape
	verts   [4]vect.Vect
	// The width of the box. Call UpdatePoly() if changed.
	Width float64
	// The height of the box. Call UpdatePoly() if changed.
	Height float64
	// The center of the box. Call UpdatePoly() if changed.
	Position vect.Vect
}

// Creates a new BoxShape with given position, width and height.
func NewBox(pos vect.Vect, w, h float64) *Shape {
	shape := new(Shape)
	box := &BoxShape{
		Polygon:  &PolygonShape{},
		Width:    w,
		Height:   h,
		Position: pos,
	}

	hw := w / 2.0
	hh := h / 2.0
	box.verts = [4]vect.Vect{
		{-hw, -hh},
		{-hw, hh},
		{hw, hh},
		{hw, -hh},
	}

	poly := box.Polygon
	poly.SetVerts(box.verts[:], box.Position)

	shape.ShapeClass = box
	return shape
}

// Recalculates the internal Polygon with the Width, Height and Position.
func (box *BoxShape) UpdatePoly() {
	hw := box.Width / 2.0
	hh := box.Height / 2.0
	box.verts = [4]vect.Vect{
		{-hw, -hh},
		{-hw, hh},
		{hw, hh},
		{hw, -hh},
	}

	poly := box.Polygon
	poly.SetVerts(box.verts[:], box.Position)
}

// Returns ShapeType_Box. Needed to implemet the ShapeClass interface.
func (box *BoxShape) ShapeType() ShapeType {
	return ShapeType_Box
}

// Recalculates the transformed vertices, axes and the bounding box.
func (box *BoxShape) update(xf transform.Transform) aabb.AABB {
	return box.Polygon.update(xf)
}

// Returns true if the given point is located inside the box.
func (box *BoxShape) TestPoint(point vect.Vect) bool {
	return box.Polygon.TestPoint(point)
}
