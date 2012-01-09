package main

import (
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/vect"
	"encoding/json"
	"os"
	"fmt"
)

var basePath = "saves/examples/"

func saveToFile(i interface{}, fileName string) {
	path := basePath + fileName + ".json"
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
	basePath = "saves/examples/"
	collisionExamples()
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
			{-1,  1},
			{ 1,  1},
			{ 1, -1},
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
