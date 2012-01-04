package transform

import (
	"encoding/json"
	"log"

	"mater/vect"
)

func (xf *Transform) MarshalJSON() ([]byte, error) {
	xfData := struct {
		X, Y     float64
		Rotation float64
	}{
		X:        xf.Position.X,
		Y:        xf.Position.Y,
		Rotation: xf.Angle(),
	}

	return json.Marshal(&xfData)
}

func (xf *Transform) UnmarshalJSON(data []byte) error {
	xf.SetIdentity()

	xfData := struct {
		X, Y     float64
		Rotation float64
	}{
		X:        xf.Position.X,
		Y:        xf.Position.Y,
		Rotation: xf.Angle(),
	}

	err := json.Unmarshal(data, &xfData)
	if err != nil {
		log.Printf("Error decoding transform")
		return err
	}

	xf.Position = vect.Vect{xfData.X, xfData.Y}
	xf.SetAngle(xfData.Rotation)

	return nil
}
