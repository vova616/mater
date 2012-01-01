package collision

import (
	"mater/vect"
	"mater/transform"
	"bytes"
	"json"
	"os"
	"log"
)

func (space *Space) MarshalJSON() ([]byte, os.Error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')

	buf.WriteString(`"Gravity":`)
	encoder.Encode(space.Gravity)

	buf.WriteString(`,"StaticBodies":`)
	buf.WriteByte('[')

	staticBodyNum := 0
	for _, body := range space.StaticBodies {
		if body.UserData != nil {
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
	for _, body := range space.DynamicBodies {
		if body.UserData != nil {
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
		log.Printf("Error unmarshaling space")
		return err
	}

	space.Gravity = spaceData.Gravity

	for _, body := range spaceData.DynamicBodies {
		body.isStatic = false
		space.AddBody(body)
	}

	for _, body := range spaceData.StaticBodies {
		body.isStatic = true
		space.AddBody(body)
	}

	return nil
}

func (body *Body) MarshalJSON() ([]byte, os.Error) {
	if body.IsStatic() {
		bodyData := struct{
			Transform transform.Transform
			Friction, Mass, Inertia float64
			Shapes []*Shape
			Enabled bool
		}{
			Transform: body.Transform,
			Friction: body.Friction,
			Mass: body.mass,
			Inertia: body.i,
			Shapes: body.Shapes,
			Enabled: body.Enabled,
		}

		return json.Marshal(&bodyData)
	} else {
		bodyData := struct{
			Transform transform.Transform
			Friction, Mass, Inertia float64
			Shapes []*Shape
			Enabled bool
			Velocity vect.Vect
			AngularVelocity float64
			Force, Torque vect.Vect
		}{
			Transform: body.Transform,
			Friction: body.Friction,
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
