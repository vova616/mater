package collision

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
	"strings"
	"strconv"
	"math"
)

// float64 wrapper that can be used to marshal +/-Inf and NaN to json
type InfFloat float64

func (f InfFloat) MarshalJSON() ([]byte, error) {
	str := strconv.FormatFloat(float64(f), 'g', -1, 64)
	if math.IsInf(float64(f), 0) || math.IsNaN(float64(f)) {
		str = strconv.Quote(str)
	}
	return []byte(str), nil
}

func (f *InfFloat) UnmarshalJSON(data []byte) error {
	str := string(data)
	if data[0] == '"' {
		//Inf or NaN
		var err error
		str, err = strconv.Unquote(str)
		if err != nil {
			log.Printf("Error decoding quoted InfFloat")
			return err
		}
	}

	f64, err := strconv.ParseFloat(str, 64)
	*f = InfFloat(f64)
	if err != nil {
		log.Printf("Error decoding InfFloat")
	}
	return err
}

//START VERTICES REGION
/*func (verts Vertices) MarshalJSON ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('[')

	for _, v := range verts {
		encoder.Encode(v.X)
		buf.WriteByte(',')
		encoder.Encode(v.Y)
		buf.WriteByte(',')
	}

	if len(verts) > 0 {
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte(']')
	return
}

func (verts *Vertices) UnmarshalJSON(data []byte) error {
	vertData := []float64{}
	err := json.Unmarshal(data, &vertData)
	if err != nil {
		log.Printf("Error decoding vertices!")
		return err
	}

	if len(vertData) % 2 != 0 {
		log.Printf("Error: Need at least 2 values for each Vertex")
		return errors.New("Need at least 2 values for each Vertex")
	}

	v := make(Vertices, len(vertData) / 2)
	*verts = v

	for i := 0; i < len(vertData) / 2; i ++ {
		v[i].X = vertData[i]
		v[i].Y = vertData[i + 1]
	}

	return nil
}*/
//END VERTICES REGION

//START SPACE REGION

// Serializes gravity and bodies to json.
// Bodies with UserData != nil are not serialized.
func (space *Space) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')

	buf.WriteString(`"Gravity":`)
	err := encoder.Encode(space.Gravity)
	if err != nil {
		log.Printf("Error decoding gravity")
		return nil, err
	}

	buf.WriteString(`,"Bodies":`)
	buf.WriteByte('[')

	bodyNum := 0
	for _, body := range space.Bodies {
		if body.UserData != nil {
			continue
		}

		bodyNum++

		err := encoder.Encode(body)
		if err != nil {
			log.Printf("Error encoding body: %v", body)
			return nil, err
		}
		buf.WriteByte(',')
	}

	//check if we serialized any bodies
	if bodyNum != 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteByte(']')

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (space *Space) UnmarshalJSON(data []byte) error {
	spaceData := struct {
		Gravity                     vect.Vect
		Bodies []*Body
	}{
		Gravity: space.Gravity,
	}

	err := json.Unmarshal(data, &spaceData)
	if err != nil {
		log.Printf("Error decoding space")
		return err
	}
	space.init()
	space.Gravity = spaceData.Gravity

	for _, body := range spaceData.Bodies {
		space.AddBody(body)
	}

	return nil
}

//END SPACE REGION

//START BODY REGION

func (body *Body) MarshalJSON() ([]byte, error) {
	if body.IsStatic() {
		bodyData := struct {
			Type      string
			Transform *transform.Transform
			Shapes    []*Shape
			Enabled   bool
		}{
			Type:      body.BodyType().ToString(),
			Transform: &body.Transform,
			Shapes:    body.Shapes,
			Enabled:   body.Enabled,
		}
		return json.Marshal(&bodyData)
	} else {
		bodyData := struct {
			Type            string
			Transform       *transform.Transform
			Shapes          []*Shape
			Enabled         bool
			Mass            InfFloat
			Inertia         InfFloat
			Velocity        vect.Vect
			AngularVelocity float64
			Force           vect.Vect
			Torque          float64
			IgnoreGravity   bool
		}{
			Type:            body.BodyType().ToString(),
			Transform:       &body.Transform,
			Shapes:          body.Shapes,
			Enabled:         body.Enabled,
			Mass:            InfFloat(body.mass),
			Inertia:         InfFloat(body.i),
			Velocity:        body.Velocity,
			AngularVelocity: body.AngularVelocity,
			Force:           body.Force,
			Torque:          body.Torque,
			IgnoreGravity:   body.IgnoreGravity,
		}
		return json.Marshal(&bodyData)
	}
	return nil, nil
}

func (body *Body) UnmarshalJSON(data []byte) error {
	if body.Shapes == nil {
		//body probably not initialized
		body.init()
	}
	bodyData := struct {
		Type            string
		Transform       *transform.Transform
		Shapes          []*Shape
		Enabled         bool
		Mass            InfFloat
		Inertia         InfFloat
		Velocity        vect.Vect
		AngularVelocity float64
		Force           vect.Vect
		Torque          float64
		IgnoreGravity   bool
	}{ //initializing everything to the bodies default values
		Type:            body.BodyType().ToString(),
		Transform:       &body.Transform,
		Shapes:          body.Shapes,
		Enabled:         body.Enabled,
		Mass:            InfFloat(body.mass),
		Inertia:         InfFloat(body.i),
		Velocity:        body.Velocity,
		AngularVelocity: body.AngularVelocity,
		Force:           body.Force,
		Torque:          body.Torque,
		IgnoreGravity:   body.IgnoreGravity,
	}

	err := json.Unmarshal(data, &bodyData)
	if err != nil {
		log.Printf("Error decoding body")
		return err
	}


	m := float64(bodyData.Mass)
	i := float64(bodyData.Inertia)

	if m == 0.0 {
		body.mass = 0.0
		body.invMass = 0.0
	} else {
		body.mass = m
		body.invMass = 1.0 / m
	}

	if i == 0.0 {
		body.i = 0.0
		body.invI = 0.0
	} else {
		body.i = i
		body.invI = 1.0 / i
	}

	var bodyType BodyType
	bodyType.FromString(bodyData.Type)
	body.SetBodyType(bodyType)

	body.Transform = *bodyData.Transform
	body.Velocity = bodyData.Velocity
	body.AngularVelocity = bodyData.AngularVelocity

	body.Force = bodyData.Force
	body.Torque = bodyData.Torque
	body.IgnoreGravity = bodyData.IgnoreGravity	

	body.Enabled = bodyData.Enabled

	for _, shape := range bodyData.Shapes {
		body.AddShape(shape)
	}

	return nil
}

//END BODY REGION

//START SHAPE REGION

func (shape *Shape) MarshalJSON() ([]byte, error) {
	if shape.ShapeClass == nil {
		log.Printf("Error: shape.ShapeClass not set")
		return nil, errors.New("shape.ShapeClass not set")
	}

	return shape.ShapeClass.marshalShape(shape)
}

func (shape *Shape) UnmarshalJSON(data []byte) error {
	shapeType := struct {
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
		return circle.unmarshalShape(shape, data)
	case "segment":
		segment := new(SegmentShape)
		shape.ShapeClass = segment
		return segment.unmarshalShape(shape, data)
	case "polygon":
		poly := new(PolygonShape)
		shape.ShapeClass = poly
		return poly.unmarshalShape(shape, data)
	case "box":
		box := new(BoxShape)
		shape.ShapeClass = box
		return box.unmarshalShape(shape, data)

	}

	log.Printf("Error: unknown shapetype: %v", shapeType.ShapeType)
	return errors.New("Unknown shapetype")
}

//END SHAPE REGION

//START CIRCLESHAPE REGION

func (circle *CircleShape) marshalShape(shape *Shape) ([]byte, error) {

	if shape.ShapeClass != circle {
		log.Printf("Error: circleshape and shape.ShapeClass don't match")
		return nil, errors.New("Wrong parent shape")
	}

	circleData := struct {
		ShapeType             string
		Friction, Restitution float64
		Position              vect.Vect
		Radius                float64
	}{
		ShapeType:   "Circle",
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Position:    circle.Position,
		Radius:      circle.Radius,
	}

	return json.Marshal(&circleData)
}

func (circle *CircleShape) unmarshalShape(shape *Shape, data []byte) error {
	if shape.ShapeClass != circle {
		log.Printf("Error: circleshape and shape.ShapeClass don't match")
		return errors.New("Wrong parent shape")
	}

	circleData := struct {
		Friction, Restitution float64
		Position              vect.Vect
		Radius                float64
	}{
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Position:    circle.Position,
		Radius:      circle.Radius,
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

//START SEGMENTSHAPE REGION

func (segment *SegmentShape) marshalShape(shape *Shape) ([]byte, error) {
	if shape.ShapeClass != segment {
		log.Printf("Error: segmentshape and shape.ShapeClass don't match")
		return nil, errors.New("Wrong parent shape")
	}

	segmentData := struct {
		ShapeType             string
		Friction, Restitution float64
		A, B                  vect.Vect
		Radius                float64
	}{
		ShapeType:   "Segment",
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		A:           segment.A,
		B:           segment.B,
		Radius:      segment.Radius,
	}

	return json.Marshal(&segmentData)
}

func (segment *SegmentShape) unmarshalShape(shape *Shape, data []byte) error {
	if shape.ShapeClass != segment {
		log.Printf("Error: segmentshape and shape.ShapeClass don't match")
		return errors.New("Wrong parent shape")
	}
	segmentData := struct {
		Friction, Restitution float64
		A, B                  vect.Vect
		Radius                float64
	}{
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		A:           segment.A,
		B:           segment.B,
		Radius:      segment.Radius,
	}

	err := json.Unmarshal(data, &segmentData)
	if err != nil {
		log.Printf("Error decoding SegmentShape")
		return err
	}

	shape.Friction = segmentData.Friction
	shape.Restitution = segmentData.Restitution
	segment.A = segmentData.A
	segment.B = segmentData.B
	segment.Radius = segmentData.Radius

	return nil
}

//END SEGMENTSHAPE REGION

//BEGIN POLYSHAPE REGION

func (poly *PolygonShape) marshalShape(shape *Shape) ([]byte, error) {
	if shape.ShapeClass != poly {
		log.Printf("Error: polyshape and shape.ShapeClass don't match")
		return nil, errors.New("Wrong parent shape")
	}

	polyData := struct {
		ShapeType string
		Friction, Restitution float64
		Vertices Vertices
	}{
		ShapeType: "Polygon",
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Vertices: poly.Verts,
	}

	return json.Marshal(&polyData)
}

func (poly *PolygonShape) unmarshalShape(shape *Shape, data []byte) error {
	if shape.ShapeClass != poly {
		log.Printf("Error: polyshape and shape.ShapeClass don't match")
		return errors.New("Wrong parent shape")
	}

	polyData := struct {
		Vertices Vertices
		Friction, Restitution float64
	}{
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Vertices: poly.Verts,
	}

	err := json.Unmarshal(data, &polyData)
	if err != nil {
		log.Printf("Error decoding PolygonShape")
		return err
	}

	poly.SetVerts(polyData.Vertices, vect.Vect{})
	return nil
}

//END POLYSHAPE REGION

//START BOXSHAPE REGION

func (box *BoxShape) marshalShape(shape *Shape) ([]byte, error) {
	if shape.ShapeClass != box {
		log.Printf("Error: boxshape and shape.ShapeClass don't match")
		return nil, errors.New("Wrong parent shape")
	}

	boxData := struct {
		ShapeType string
		Friction, Restitution float64
		Width float64
		Height float64
		Position vect.Vect
	}{
		ShapeType: "Box",
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Width: box.Width,
		Height: box.Height,
		Position: box.Position,
	}

	return json.Marshal(&boxData)
}

func (box *BoxShape) unmarshalShape(shape *Shape, data []byte) error {
	if shape.ShapeClass != box {
		log.Printf("Error: boxshape and shape.ShapeClass don't match")
		return errors.New("Wrong parent shape")
	}

	boxData := struct {
		Friction, Restitution float64
		Width float64
		Height float64
		Position vect.Vect
	}{
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Width: box.Width,
		Height: box.Height,
		Position: box.Position,
	}

	err := json.Unmarshal(data, &boxData)
	if err != nil {
		log.Printf("Error decoding BoxShape")
		return err
	}

	box.Width = boxData.Width
	box.Height = boxData.Height
	box.Position = boxData.Position
	if box.Polygon == nil {
		box.Polygon = new(PolygonShape)
	}
	box.UpdatePoly()
	return nil
}

//END BOXSHAPE REGION
