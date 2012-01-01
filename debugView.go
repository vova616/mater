package mater

import (
	"mater/render"
	"mater/collision"
	"mater/vect"
	"gl"
)

var _tmpVertices []vect.Vect
func init () {
	_tmpVertices = make([]vect.Vect, 8)
}

type DebugView struct {
	space *collision.Space
}

func NewDebugView (space *collision.Space) *DebugView{
	dv := new(DebugView)
	dv.space = space
	return dv
}

func (dv *DebugView) Reset (space *collision.Space) {
	dv.space = space
}

func (dv *DebugView) DrawDebugData () {
	space := dv.space

	//Draw static shapes
	for _, b := range space.StaticBodies {
		for _, s := range b.Shapes {
			if b.Enabled == false {
				//Inactive
				gl.Color3f(.5, .8, .5)
			} else {
				//Static
				gl.Color3f(1, 1, 1)
			}
			DrawShape(s)
		}
	}

	/*

	if b.Enabled == false {
				//Inactive
				gl.Color3f(.5, .8, .5)
			} else if b.IsStatic() {
				//Static
				gl.Color3f(1, 1, 1)
			} else if b.Awake() == false {
				//Sleeping
				gl.Color3f(.5, .5, .5)
			} else {
				//Default
				gl.Color3f(1, 0, 0)
			}
			*/
			/*
	const axisScale = .3
	
	if dv.lastPoint != 0 {
		dv.pointCount = dv.lastPoint
	}

	for i := 0; i < dv.pointCount; i++ {
		point := &dv.points[i]

		gl.Color3f(.4, .9, .4)
		p1 := point.Position
		p2 := Add(p1, Scale(point.Normal, axisScale))
		render.DrawLine(p1, p2)
		
		gl.Begin(gl.POINTS)
		if point.State == PointState_Add {
			gl.Color3f(.3, .95, .3)
			render.DrawCircle(p1, .1, true)
		} else if point.State == PointState_Persist {
			gl.Color3f(.3, .3, .95)
			render.DrawCircle(p1, .1, true)
		}
		gl.End()
	}
	dv.lastPoint = 0*/
}

func DrawShape(shape *collision.Shape) {
	switch shape.ShapeType() {
		case collision.ShapeType_Circle:
			circle := shape.ShapeClass.(*collision.CircleShape)
			render.DrawCircle(vect.Add(shape.Body.Transform.Position, circle.Position), circle.Radius, false)
			break
		/*case ShapeType_Polygon:
			poly := shape.ShapeClass.(*PolygonShape)
			vertCount := len(poly.Vertices)

			for i := 0; i < vertCount; i++ {
				_tmpVertices[i] = MultiplyTransformVect(xf, &poly.Vertices[i])
			}
			render.DrawPoly(_tmpVertices, vertCount, false)
			break*/
	}
}
