package collision

import (
	"mater/vect"
	"mater/transform"
	"bytes"
	"json"
	"os"
	"log"
	"strings"
)

//START SPACE REGION
func (space *Space) MarshalJSON() ([]byte, os.Error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')

	buf.WriteString(`"Gravity":`)
	err := encoder.Encode(space.Gravity)
	if err != nil {
		log.Printf("Error decoding gravity")
		return nil, err
	}

	buf.WriteString(`,"StaticBodies":`)
	buf.WriteByte('[')

	staticBodyNum := 0
	for _, body := range space.Bodies {
		if body.UserData != nil {
			continue
		}

		if body.IsStatic() == false {
			continue
		}

		staticBodyNum++

		err := encoder.Encode(body)
		if err != nil {
			log.Printf("Error encoding body: %v", body)
			return nil, err
		}
		buf.WriteByte(',')
	}

	//check if we serialized any bodies
	if staticBodyNum != 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteByte(']')

	buf.WriteString(`,"DynamicBodies":`)
	buf.WriteByte('[')

	dynamicBodyNum := 0
	for _, body := range space.Bodies {
		if body.UserData != nil {
			continue
		}

		if body.IsStatic() == true {
			continue
		}

		dynamicBodyNum++

		err := encoder.Encode(body)
		if err != nil {
			log.Printf("Error encoding body: %v", body)
			return nil, err
		}
		buf.WriteByte(',')
	}

	//check if we serialized any bodies
	if dynamicBodyNum != 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte(']')

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (space *Space) UnmarshalJSON(data []byte) os.Error {
	spaceData := struct{
		Gravity vect.Vect
		DynamicBodies, StaticBodies []*Body
	}{
		Gravity: space.Gravity,
	}
	
	err := json.Unmarshal(data, &spaceData)
	if err != nil {
		log.Printf("Error decoding space")
		return err
	}

	space.Gravity = spaceData.Gravity

	for _, body := range spaceData.DynamicBodies {
		body.BodyType = BodyType_Dynamic
		space.AddBody(body)
	}

	for _, body := range spaceData.StaticBodies {
		body.BodyType = BodyType_Static
		space.AddBody(body)
	}

	return nil
}
//END SPACE REGION

//START BODY REGION
func (body *Body) MarshalJSON() ([]byte, os.Error) {
	if body.IsStatic() {
		bodyData := struct{
			Transform *transform.Transform
			Mass, Inertia float64
			Shapes []*Shape
			Enabled bool
		}{
			Transform: &body.Transform,
			Mass: body.mass,
			Inertia: body.i,
			Shapes: body.Shapes,
			Enabled: body.Enabled,
		}

		return json.Marshal(&bodyData)
	} else {
		bodyData := struct{
			Transform *transform.Transform
			Mass, Inertia float64
			Shapes []*Shape
			Enabled bool
			Velocity vect.Vect
			AngularVelocity float64
			Force vect.Vect
			Torque float64
		}{
			Transform: &body.Transform,
			Mass: body.mass,
			Inertia: body.i,
			Shapes: body.Shapes,
			Enabled: body.Enabled,
			Velocity: body.Velocity,
			AngularVelocity: body.AngularVelocity,
			Force: body.Force, 
			Torque: body.Torque,
		}

		return json.Marshal(&bodyData)
	}
	return nil, nil
}

func (body *Body) UnmarshalJSON(data []byte) os.Error {
	if body.Shapes == nil {
		//body probably not initialized
		body.init()
	}
	bodyData := struct{
		Transform *transform.Transform
		Mass, Inertia float64
		Shapes []*Shape
		Enabled bool
		Velocity vect.Vect
		AngularVelocity float64
		Force vect.Vect
		Torque float64
	}{//initializing everything to the bodies default values
		Transform: &body.Transform,
		Mass: body.mass,
		Inertia: body.i,
		Shapes: body.Shapes,
		Enabled: body.Enabled,
		Velocity: body.Velocity,
		AngularVelocity: body.AngularVelocity,
		Force: body.Force, 
		Torque: body.Torque,
	}

	err := json.Unmarshal(data, &bodyData)
	if err != nil {
		log.Printf("Error decoding body")
		return err
	}

	body.Transform = *bodyData.Transform
	body.Velocity = bodyData.Velocity
	body.AngularVelocity = bodyData.AngularVelocity

	body.Force = bodyData.Force
	body.Torque = bodyData.Torque

	body.SetMass(bodyData.Mass)
	body.SetInertia(bodyData.Inertia)

	body.Enabled = bodyData.Enabled

	for _, shape := range bodyData.Shapes {
		body.AddShape(shape)
	}

	return nil
}
//END BODY REGION

//START SHAPE REGION
func (shape *Shape) MarshalJSON() ([]byte, os.Error) {
	if shape.ShapeClass == nil {
		log.Printf("Error: shape.ShapeClass not set")
		return nil, os.NewError("shape.ShapeClass not set")
	}

	return shape.ShapeClass.MarshalShape(shape)
}

func (shape *Shape) UnmarshalJSON(data []byte) (os.Error) {
	shapeType := struct{
		ShapeType string
	}{}

	err := json.Unmarshal(data, &shapeType)
	if err != nil {
		log.Printf("Error: could not find shapetype")
		return err
	}

	switch strings.ToLower(shapeType.ShapeType) {
		case "circle":
			circle := new(CircleShape)
			shape.ShapeClass = circle
			return circle.UnmarshalShape(shape, data)
	}

	log.Printf("Error: unknown shapetype: %v", shapeType.ShapeType)
	return os.NewError("Unknown shapetype")
}
//END SHAPE REGION

//START CIRCLESHAPE REGION
func (circle *CircleShape) MarshalShape(shape *Shape) ([]byte, os.Error) {

	if shape.ShapeClass != circle {
		log.Printf("Error: circleshape and shape.ShapeClass don't match")
		return nil, os.NewError("Wrong parent shape")
	}

	circleData := struct {
		ShapeType string
		Friction, Restitution float64
		Position vect.Vect
		Radius float64
	}{
		ShapeType: "Circle",
		Friction: shape.Friction,
		Restitution: shape.Restitution,
		Position: circle.Position,
		Radius: circle.Radius,
	}

	return json.Marshal(&circleData)
}

func (circle *CircleShape) UnmarshalShape(shape *Shape, data []byte) os.Error {
	if shape.ShapeClass != circle {
		log.Printf("Error: circleshape and shape.ShapeClass don't match")
		return os.NewError("Wrong parent shape")
	}

	circleData := struct {
		Friction, Restitution float64
		Position vect.Vect
		Radius float64
	}{
		Friction: shape.Friction,
		Restitution: shape.Restitution,
		Position: circle.Position,
		Radius: circle.Radius,
	}

	err := json.Unmarshal(data, &circleData)
	if err != nil {
		log.Printf("Error decoding CircleShape")
		return err
	}

	shape.Friction = circleData.Friction
	shape.Restitution = circleData.Restitution
	circle.Position = circleData.Position
	circle.Radius = circleData.Radius

	return nil
}
//END CIRCLESHAPE REGION
