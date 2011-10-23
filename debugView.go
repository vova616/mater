package mater

import (
	. "box2d"
	. "box2d/vector2"
	"box2d/settings"
	"gl"
)

var _tmpVertices Vertices
func init () {
	_tmpVertices = make(Vertices, settings.MaxPolygonVertices)
}

type _contactPoint struct {
	Normal, Position Vector2
	State PointState
}

type DebugView struct {
	contactListener ContactListener
	pointCount, lastPoint int
	points []_contactPoint
	world *World
}

func NewDebugView (world *World) *DebugView{
	dv := new(DebugView)
	dv.world = world
	dv.points = make([]_contactPoint, 32)
	if cl := world.ContactManager().ContactListener; cl != nil {
		dv.contactListener = cl
	}
	world.SetContactListener(dv)
	return dv
}

func (dv *DebugView) Reset (world *World) {
	dv.world = world
	dv.points = make([]_contactPoint, 32)
	dv.contactListener = nil
	if cl := world.ContactManager().ContactListener; cl != nil {
		dv.contactListener = cl
	}
	dv.pointCount = 0
	dv.lastPoint = 0
	world.SetContactListener(dv)
}

func (dv *DebugView) BeginContact (contact *Contact) {
	if dv.contactListener != nil {
		dv.contactListener.BeginContact(contact)
	}
}

func (dv *DebugView) EndContact (contact *Contact) {
	if dv.contactListener != nil {
		dv.contactListener.EndContact(contact)
	}
}

func (dv *DebugView) PostSolve (contact *Contact, impulse *ContactConstraint) {
	if dv.contactListener != nil {
		dv.contactListener.PostSolve(contact, impulse)
	}
}

func (dv *DebugView) PreSolve (contact *Contact, oldManifold *Manifold) {
	if dv.contactListener != nil {
		dv.contactListener.PreSolve(contact, oldManifold)
	}

	manifold := contact.Manifold

	if manifold.PointCount == 0 {
		return
	}

	_, state2 := GetPointStates(oldManifold, &manifold)

	normal, points := contact.GetWorldManifold()

	for i := 0; i < manifold.PointCount; i++ {
		var cp _contactPoint
		cp.Position = points[i]
		cp.Normal = normal
		cp.State = state2[i]

		if dv.lastPoint > len(dv.points) - 1 {
			dv.points = append(dv.points, cp)
		} else {
			dv.points[dv.lastPoint] = cp
		}
		dv.lastPoint++
	}
}

func (dv *DebugView) DrawDebugData () {
	world := dv.world

	gl.Color3f(.4, .4, .4)
	for _, aabb := range(world.GetDynamicTreeNodes()) {
		Render.DrawQuad(aabb.LowerBound, aabb.UpperBound, false)
	}

	//Draw shapes
	for _, b := range(world.BodyList()) {
		xf := &(*b.Transform())
		for _, f := range(b.FixtureList()) {
			if b.Enabled() == false {
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
			DrawShape(f.Shape(), xf)
		}
	}

	const axisScale = .3
	
	if dv.lastPoint != 0 {
		dv.pointCount = dv.lastPoint
	}

	for i := 0; i < dv.pointCount; i++ {
		point := &dv.points[i]

		gl.Color3f(.4, .9, .4)
		p1 := point.Position
		p2 := Add(p1, Scale(point.Normal, axisScale))
		Render.DrawLine(p1, p2)
		
		gl.Begin(gl.POINTS)
		if point.State == PointState_Add {
			gl.Color3f(.3, .95, .3)
			Render.DrawCircle(p1, .1, true)
		} else if point.State == PointState_Persist {
			gl.Color3f(.3, .3, .95)
			Render.DrawCircle(p1, .1, true)
		}
		gl.End()
	}
	dv.lastPoint = 0
}

func DrawShape(shape *Shape, xf *Transform) {
	switch shape.ShapeType() {
		case ShapeType_Circle:
			circle := shape.ShapeClass.(*CircleShape)
			Render.DrawCircle(Add(xf.Position, circle.Position()), circle.Radius(), false)
			break
		case ShapeType_Polygon:
			poly := shape.ShapeClass.(*PolygonShape)
			vertCount := len(poly.Vertices)

			for i := 0; i < vertCount; i++ {
				_tmpVertices[i] = MultiplyTransformVect(xf, &poly.Vertices[i])
			}
			Render.DrawPoly(_tmpVertices, vertCount, false)
			break
	}
}
