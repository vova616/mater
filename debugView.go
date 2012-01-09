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

	//Draw shapes
	for _, b := range space.Bodies {
		if b.Enabled == false {
			//Inactive
			gl.Color3f(.5, .8, .5)
		} else if b.IsStatic() {
			//Static
			gl.Color3f(1, 1, 1)
		} else {
			//Normal
			gl.Color3f(1, 0, 0)
		}
		for _, s := range b.Shapes {
			DrawShape(s)

		}
	}

	//Draw aabbs
	const drawAABBs = true
	if drawAABBs {
		for _, b := range space.Bodies {
			gl.Color3f(.3, .7, .7)
			for _, s := range b.Shapes {
				render.DrawQuad(s.AABB.Lower, s.AABB.Upper, false)
			}
		}
	}

	const contactRadius = 0.2
	const contactNormalScale = 0.5

	for _, arb := range space.Arbiters {
		for i := 0; i < arb.NumContacts; i++ {
			con := arb.Contacts[i]
			gl.Color3f(0, 0, 1)
			p1 := con.Position
			p2 := vect.Add(p1, vect.Mult(con.Normal, contactNormalScale))
			//p2 := vect.Add(p1, vect.Mult(con.Normal, con.Separation))
			render.DrawLine(p1, p2)
			gl.Color3f(0, 1, 0)
			render.DrawCircle(con.Position, contactRadius, false)
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
	xf := shape.Body.Transform
	switch shape.ShapeType() {
		case collision.ShapeType_Circle:
			circle := shape.ShapeClass.(*collision.CircleShape)
			render.DrawCircle(vect.Add(xf.Position, xf.RotateVect(circle.Position)), circle.Radius, false)
			break
		case collision.ShapeType_Segment:
			segment := shape.ShapeClass.(*collision.SegmentShape)
			a := segment.Ta
			b := segment.Tb
			r := segment.Radius
			render.DrawLine(a, b)
			if segment.Radius > 0.0 {
				render.DrawCircle(a, r, false)
				render.DrawCircle(b, r, false)

				verts := [4]vect.Vect{
					vect.Add(a, vect.Vect{0, r}),
					vect.Add(a, vect.Vect{0, -r}),
					vect.Add(b, vect.Vect{0, -r}),
					vect.Add(b, vect.Vect{0, r}),
				}
				render.DrawPoly(verts[:], 4, false)

			}
			//Normal:
			/*n := segment.Normal()
			render.DrawLine(a, vect.Add(a, n))
			render.DrawLine(b, vect.Add(b, n))*/
		case collision.ShapeType_Polygon:
			poly := shape.ShapeClass.(*collision.PolygonShape)
			verts := poly.TVerts
			render.DrawPoly(verts, poly.NumVerts, false)
			//Normals
			/*axes := poly.TAxes
			for i, v := range verts {
				a := axes[i]
				v1 := v
				v2 := verts[(i + 1) % len(verts)]
				render.DrawLine(v1, vect.Add(v1, a.N))
				render.DrawLine(v2, vect.Add(v2, a.N))
			}*/

		case collision.ShapeType_Box:
			poly := shape.ShapeClass.(*collision.BoxShape).Polygon
			verts := poly.TVerts
			render.DrawPoly(verts, poly.NumVerts, false)
	}
}
