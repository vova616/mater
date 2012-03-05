package collision

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/teomat/mater/transform"
	"github.com/teomat/mater/vect"
	"log"
	"math"
	"strconv"
	"strings"
	"unsafe"
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
		Gravity *vect.Vect
		Bodies  []*Body
	}{
		Gravity: &space.Gravity,
	}

	err := json.Unmarshal(data, &spaceData)
	if err != nil {
		log.Printf("Error decoding space")
		return err
	}
	space.init()
	//space.Gravity = spaceData.Gravity

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
		Velocity        *vect.Vect
		AngularVelocity float64
		Force           *vect.Vect
		Torque          float64
		IgnoreGravity   bool
	}{ //initializing everything to the bodies default values
		Type:            body.BodyType().ToString(),
		Transform:       &body.Transform,
		Shapes:          body.Shapes,
		Enabled:         body.Enabled,
		Mass:            InfFloat(body.mass),
		Inertia:         InfFloat(body.i),
		Velocity:        &body.Velocity,
		AngularVelocity: body.AngularVelocity,
		Force:           &body.Force,
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
	//body.Velocity = bodyData.Velocity
	body.AngularVelocity = bodyData.AngularVelocity

	//body.Force = bodyData.Force
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

	shapeData := struct {
		ShapeType   string
		Friction    float64
		Restitution float64
		Sensor      bool
		Surface_v   vect.Vect

		CollisionCat string
		CollidesWith string
	}{
		ShapeType:   shape.ShapeType().ToString(),
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Sensor:      shape.IsSensor,
		Surface_v:   shape.Surface_v,

		//encoded as bitstrings
		CollisionCat: strconv.FormatUint(uint64(shape.CollisionCat), 2),
		CollidesWith: strconv.FormatUint(uint64(shape.CollidesWith), 2),
	}

	data, err := json.Marshal(&shapeData)
	if err != nil {
		log.Printf("Error encoding shape")
		return nil, err
	}

	scData, err := shape.ShapeClass.marshalShape(shape)
	if err != nil {
		log.Printf("Error encoding shape")
		return nil, err
	}

	//TODO: find an nicer solution for this
	data[len(data)-1] = ','

	return append(data, scData[1:]...), nil
}

func (shape *Shape) UnmarshalJSON(data []byte) error {
	shape.init()

	shapeData := struct {
		ShapeType   string
		Friction    float64
		Restitution float64
		Sensor      bool
		Surface_v   vect.Vect

		CollisionCat string
		CollidesWith string
	}{
		Friction:    shape.Friction,
		Restitution: shape.Restitution,
		Sensor:      shape.IsSensor,
		Surface_v:   shape.Surface_v,

		CollisionCat: strconv.FormatUint(uint64(shape.CollisionCat), 2),
		CollidesWith: strconv.FormatUint(uint64(shape.CollidesWith), 2),
	}

	err := json.Unmarshal(data, &shapeData)
	if err != nil {
		log.Printf("Error: could not find shapetype")
		return err
	}

	shape.Friction = shapeData.Friction
	shape.Restitution = shapeData.Restitution
	shape.IsSensor = shapeData.Sensor
	shape.Surface_v = shapeData.Surface_v

	colCat, err := strconv.ParseUint(shapeData.CollisionCat, 2, int(unsafe.Sizeof(CollisionCategory(0))))
	if err != nil {
		log.Printf("Error decoding CollisionCat")
		return err
	}
	colWith, err := strconv.ParseUint(shapeData.CollidesWith, 2, int(unsafe.Sizeof(CollisionCategory(0))))
	if err != nil {
		log.Printf("Error decoding CollidesWith")
		return err
	}

	shape.CollisionCat = CollisionCategory(colCat)
	shape.CollidesWith = CollisionCategory(colWith)

	switch strings.ToLower(shapeData.ShapeType) {
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

	log.Printf("Error: unknown shapetype: %v", shapeData.ShapeType)
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
		Position vect.Vect
		Radius   float64
	}{
		Position: circle.Position,
		Radius:   circle.Radius,
	}

	return json.Marshal(&circleData)
}

func (circle *CircleShape) unmarshalShape(shape *Shape, data []byte) error {
	if shape.ShapeClass != circle {
		log.Printf("Error: circleshape and shape.ShapeClass don't match")
		return errors.New("Wrong parent shape")
	}

	circleData := struct {
		Position *vect.Vect
		Radius   float64
	}{
		Position: &circle.Position,
		Radius:   circle.Radius,
	}

	err := json.Unmarshal(data, &circleData)
	if err != nil {
		log.Printf("Error decoding CircleShape")
		return err
	}

	//circle.Position = circleData.Position
	circle.Radius = circleData.Radius
	circle.Shape = shape
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
		A, B   vect.Vect
		Radius float64
	}{
		A:      segment.A,
		B:      segment.B,
		Radius: segment.Radius,
	}

	return json.Marshal(&segmentData)
}

func (segment *SegmentShape) unmarshalShape(shape *Shape, data []byte) error {
	if shape.ShapeClass != segment {
		log.Printf("Error: segmentshape and shape.ShapeClass don't match")
		return errors.New("Wrong parent shape")
	}
	segmentData := struct {
		A, B   *vect.Vect
		Radius float64
	}{
		A:      &segment.A,
		B:      &segment.B,
		Radius: segment.Radius,
	}

	err := json.Unmarshal(data, &segmentData)
	if err != nil {
		log.Printf("Error decoding SegmentShape")
		return err
	}

	//segment.A = segmentData.A
	//segment.B = segmentData.B
	segment.Radius = segmentData.Radius
	segment.Shape = shape

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
		Vertices Vertices
	}{
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
		Vertices *Vertices
	}{}

	err := json.Unmarshal(data, &polyData)
	if err != nil {
		log.Printf("Error decoding PolygonShape")
		return err
	}

	poly.SetVerts(*polyData.Vertices, vect.Vect{})
	poly.Shape = shape
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
		Width    float64
		Height   float64
		Position vect.Vect
	}{
		Width:    box.Width,
		Height:   box.Height,
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
		Width    float64
		Height   float64
		Position *vect.Vect
	}{
		Width:    box.Width,
		Height:   box.Height,
		Position: &box.Position,
	}

	err := json.Unmarshal(data, &boxData)
	if err != nil {
		log.Printf("Error decoding BoxShape")
		return err
	}

	box.Width = boxData.Width
	box.Height = boxData.Height
	//box.Position = boxData.Position
	if box.Polygon == nil {
		box.Polygon = new(PolygonShape)
	}
	box.Shape = shape
	box.Polygon.Shape = shape
	box.UpdatePoly()
	return nil
}

//END BOXSHAPE REGION
