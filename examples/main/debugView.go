package main

import (
	"github.com/banthar/Go-OpenGL/gl"
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
)

func DrawDebugData(space *collision.Space) {
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

	gl.Color3f(0, 1, 0.5)
	for _, b := range space.Bodies {
		DrawTransform(&b.Transform, 0.2)
	}

	if Settings.DrawAABBs {
		for _, b := range space.Bodies {
			gl.Color3f(.3, .7, .7)
			for _, s := range b.Shapes {
				DrawQuad(s.AABB.Lower, s.AABB.Upper, false)
			}
		}
	}

	const contactRadius = 0.2
	const contactNormalScale = 0.5

	for arb := space.ContactManager.ArbiterList.Arbiter; arb != nil; arb = arb.Next {
		for i := 0; i < arb.NumContacts; i++ {
			con := arb.Contacts[i]
			gl.Color3f(0, 0, 1)
			p1 := con.Position
			p2 := vect.Add(p1, vect.Mult(con.Normal, contactNormalScale))
			//p2 := vect.Add(p1, vect.Mult(con.Normal, con.Separation))
			DrawLine(p1, p2)
			gl.Color3f(0, 1, 0)
			DrawCircle(con.Position, contactRadius, false)
		}
	}

	if Settings.DrawTreeNodes {
		for _, node := range space.GetDynamicTreeNodes() {
			gl.Color3f(0.0, .7, .7)
			DrawQuad(node.AABB().Lower, node.AABB().Upper, false)
		}
	}
}

func DrawShape(shape *collision.Shape) {
	switch shape.ShapeType() {
	case collision.ShapeType_Circle:
		circle := shape.ShapeClass.(*collision.CircleShape)
		DrawCircle(circle.Tc, circle.Radius, false)
		const circleMarkerSize = 0.08
		{
			p1 := vect.Add(circle.Tc, vect.Vect{0, circleMarkerSize})
			p2 := vect.Sub(circle.Tc, vect.Vect{0, circleMarkerSize})
			DrawLine(p1, p2)
		}
		{
			p1 := vect.Add(circle.Tc, vect.Vect{circleMarkerSize, 0})
			p2 := vect.Sub(circle.Tc, vect.Vect{circleMarkerSize, 0})
			DrawLine(p1, p2)
		}
		break
	case collision.ShapeType_Segment:
		segment := shape.ShapeClass.(*collision.SegmentShape)
		a := segment.Ta
		b := segment.Tb
		r := segment.Radius
		DrawLine(a, b)
		if segment.Radius > 0.0 {
			DrawCircle(a, r, false)
			DrawCircle(b, r, false)

			verts := [4]vect.Vect{
				vect.Add(a, vect.Vect{0, r}),
				vect.Add(a, vect.Vect{0, -r}),
				vect.Add(b, vect.Vect{0, -r}),
				vect.Add(b, vect.Vect{0, r}),
			}
			DrawPoly(verts[:], 4, false)

		}
		if Settings.DrawNormals {
			n := segment.Tn
			DrawLine(a, vect.Add(a, n))
			DrawLine(b, vect.Add(b, n))
		}
	case collision.ShapeType_Polygon:
		poly := shape.ShapeClass.(*collision.PolygonShape)
		verts := poly.TVerts
		DrawPoly(verts, poly.NumVerts, false)
		if Settings.DrawNormals {
			axes := poly.TAxes
			for i, v := range verts {
				a := axes[i]
				v1 := v
				v2 := verts[(i+1)%len(verts)]
				DrawLine(v1, vect.Add(v1, a.N))
				DrawLine(v2, vect.Add(v2, a.N))
			}
		}

	case collision.ShapeType_Box:
		poly := shape.ShapeClass.(*collision.BoxShape).Polygon
		verts := poly.TVerts
		DrawPoly(verts, poly.NumVerts, false)
	}
}

func DrawTransform(xf *transform.Transform, radius float64) {
	DrawCircle(xf.Position, radius, false)
	p := xf.RotateVect(vect.Vect{0, -radius})
	p.Add(xf.Position)
	DrawLine(xf.Position, p)
}
