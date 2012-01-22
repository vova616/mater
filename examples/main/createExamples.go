package main

import (
	"encoding/json"
	"fmt"
	"github.com/teomat/mater/camera"
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/engine"
	"github.com/teomat/mater/vect"
	"os"
)

var examplesPath = "examples/"

func saveToFile(i interface{}, fileName string) {
	path := Settings.SaveDir + examplesPath + fileName + ".json"
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error opening File: %v", err)
		return
	}
	defer file.Close()

	dataString, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		fmt.Printf("Error encoding: %v", err)
		return
	}

	n, err := file.Write(dataString)
	if err != nil {
		fmt.Printf("Error after writing %v characters to File \"%v\": %v", n, path, err)
		return
	}
}

//Creates example json files
func allExamples() {
	examplesPath = "examples/"
	collisionExamples()
	examplesPath = "tests/"
	collisionTests()
}

//Creates examples for mater/collision.
func collisionExamples() {
	{
		circle := collision.NewCircle(vect.Vect{}, 1.0)

		saveToFile(circle, "collision/shape_circle")
	}
	{
		segment := collision.NewSegment(vect.Vect{1, 1}, vect.Vect{-1, -1}, 0)

		saveToFile(segment, "collision/shape_segment")
	}
	{
		verts := collision.Vertices{
			{-1, -1},
			{-1, 1},
			{1, 1},
			{1, -1},
		}

		poly := collision.NewPolygon(verts, vect.Vect{})

		saveToFile(poly, "collision/shape_polygon")
	}
	{
		box := collision.NewBox(vect.Vect{}, 1, 1)

		saveToFile(box, "collision/shape_box")
	}
	{
		body := collision.NewBody(collision.BodyType_Static)

		saveToFile(body, "collision/body_static")
	}
	{
		body := collision.NewBody(collision.BodyType_Dynamic)

		saveToFile(body, "collision/body_dynamic")
	}
	{
		space := collision.NewSpace()

		saveToFile(space, "collision/space")
	}
}

func collisionTests() {
	scene := new(engine.Scene)
	scene.Init()

	//add a camera to the scene
	{
		cam := new(camera.Camera)
		cam.Scale = vect.Vect{32, 32}
		cam.IsMainCamera = true

		scene.StaticEntity.AddComponent(cam)
	}

	//circle-segment collision
	{
		space := collision.NewSpace()
		scene.Space = space
		space.Gravity = vect.Vect{0, 10}

		{
			body := collision.NewBody(collision.BodyType_Static)
			body.Transform.Position = vect.Vect{-3, 0}
			body.Transform.SetAngle(0.5)

			seg := collision.NewSegment(vect.Vect{5, 0}, vect.Vect{-5, 0}, 0.0)
			body.AddShape(seg)

			space.AddBody(body)
		}

		{
			body := collision.NewBody(collision.BodyType_Static)
			body.Transform.Position = vect.Vect{3, 0}
			body.Transform.SetAngle(-0.5)

			seg := collision.NewSegment(vect.Vect{5, 0}, vect.Vect{-5, 0}, 1.0)
			body.AddShape(seg)

			space.AddBody(body)
		}

		{
			body := collision.NewBody(collision.BodyType_Dynamic)
			body.Transform.Position = vect.Vect{6, -7}
			body.Transform.SetAngle(0)

			circle := collision.NewCircle(vect.Vect{0, 0}, 1.0)
			body.AddShape(circle)

			space.AddBody(body)
		}

		saveToFile(scene, "circle-segment")
	}

	//circle-polygon collision
	{
		space := collision.NewSpace()
		scene.Space = space
		space.Gravity = vect.Vect{0, 10}

		{
			body := collision.NewBody(collision.BodyType_Static)
			body.Transform.Position = vect.Vect{-3, 0}
			body.Transform.SetAngle(0.5)

			verts := collision.Vertices{
				{-5, -1},
				{-5, 1},
				{5, 1},
				{5, -1},
			}

			poly := collision.NewPolygon(verts, vect.Vect{})
			body.AddShape(poly)

			space.AddBody(body)
		}

		{
			body := collision.NewBody(collision.BodyType_Static)
			body.Transform.Position = vect.Vect{3, 0}
			body.Transform.SetAngle(-0.5)

			verts := collision.Vertices{
				{-5, -1},
				{-5, 1},
				{5, 1},
				{5, -1},
			}

			poly := collision.NewPolygon(verts, vect.Vect{})
			body.AddShape(poly)

			space.AddBody(body)
		}

		{
			body := collision.NewBody(collision.BodyType_Dynamic)
			body.Transform.Position = vect.Vect{6, -7}
			body.Transform.SetAngle(0)

			circle := collision.NewCircle(vect.Vect{0, 0}, 1.0)
			body.AddShape(circle)

			space.AddBody(body)
		}

		saveToFile(scene, "circle-polygon")
	}

	//polygon-polygon collision
	{
		space := collision.NewSpace()
		scene.Space = space
		space.Gravity = vect.Vect{0, 10}

		{
			body := collision.NewBody(collision.BodyType_Static)
			body.Transform.Position = vect.Vect{-3, 0}
			body.Transform.SetAngle(0.5)

			verts := collision.Vertices{
				{-5, -1},
				{-5, 1},
				{5, 1},
				{5, -1},
			}

			poly := collision.NewPolygon(verts, vect.Vect{})
			body.AddShape(poly)

			space.AddBody(body)
		}

		{
			body := collision.NewBody(collision.BodyType_Static)
			body.Transform.Position = vect.Vect{3, 0}
			body.Transform.SetAngle(-0.5)

			verts := collision.Vertices{
				{-5, -1},
				{-5, 1},
				{5, 1},
				{5, -1},
			}

			poly := collision.NewPolygon(verts, vect.Vect{})
			body.AddShape(poly)

			space.AddBody(body)
		}

		{
			body := collision.NewBody(collision.BodyType_Dynamic)
			body.Transform.Position = vect.Vect{6, -7}
			body.Transform.SetAngle(0)
			//fixed rotation for now
			body.SetInertia(0)

			box := collision.NewBox(vect.Vect{0, 0}, 1, 1)
			body.AddShape(box)

			space.AddBody(body)
		}

		saveToFile(scene, "polygon-polygon")
	}
}
