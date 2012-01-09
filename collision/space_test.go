package collision

import (
	"github.com/teomat/mater/vect"
)

func ExampleAddBodies() {
	space := NewSpace(vect.Vect{0, 10})

	{
		//Create a body
		body := NewBody(BodyType_Dynamic)

		//Change body properties
		body.SetMass(2.0)

		//Create a new shape
		circle := NewCircle(vect.Vect{-3, -10}, 1)

		//Add the shape to the body, and the body to the space
		body.AddShape(circle)
		space.AddBody(body)
	}

	{
		body := NewBody(BodyType_Dynamic)

		body.SetMass(4.0)

		box := NewBox(vect.Vect{0, 0}, 1, 1)

		body.Transform.Position = vect.Vect{3, -10}

		//AddShape calls Update on the added shape
		body.AddShape(box)

		body.Transform.SetAngle(0.5)

		//If Settings.AutoUpdateShapes is set, thsi will be called automatically on the next call to space.Step().
		body.UpdateShapes()

		space.AddBody(body)
	}

	{
		body := NewBody(BodyType_Static)

		segment := NewSegment(vect.Vect{-5, 0}, vect.Vect{5, 0}, 0)

		//segment is a *Shape, not a *CircleShape!
		segment.Friction = 0.8

		body.AddShape(segment)

		space.AddBody(body)
	}
}
