package mater
/*
import (
	"json"
	"bytes"
	"os"
	"github.com/abneptis/GoUUID"
	"box2d"
)

func (scene *Scene) MarshalJSON() ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')

	buf.WriteString(`"World":`)
	if err := encoder.Encode(scene.World); err != nil {
		return nil, err
	}

	buf.WriteString(`, "Entities":`)
	if err := encoder.Encode(scene.Entities); err != nil {
		return nil, err
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (scene *Scene) UnmarshalJSON(data []byte) os.Error {
	sd := &struct{
		World *box2d.World
		Entities []* struct{
			BodyId *uuid.UUID
			Enabled bool
			EntityClass EntityClass
		}
	}{}

	if err := json.Unmarshal(data, sd); err != nil {
		return err
	}

	scene.World = sd.World
	bodyMap := make(map[uuid.UUID]*box2d.Body)

	for _, body := range scene.World.BodyList() {
		uuid := body.BodyId()
		bodyMap[uuid] = body
	}

	for _, ed := range sd.Entities {
		entity := new(Entity)
		entity.Enabled = ed.Enabled
		entity.EntityClass = ed.EntityClass
		entity.Body = bodyMap[*ed.BodyId]
		scene.Entities = append(scene.Entities, entity)
	}

	sd = nil
	bodyMap = nil

	return nil
}

func (entity *Entity) MarshalJSON() ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')

	buf.WriteString(`"BodyId":`)
	if entity.Body == nil {
		buf.WriteString(`null`)
	} else {
		encoder.Encode(entity.Body.BodyId())
	}
	
	buf.WriteString(`,"Enabled":`)
	encoder.Encode(entity.Enabled)

	buf.WriteString(`,"EntityClass":`)
	if err := encoder.Encode(entity.EntityClass); err != nil {
		return nil, err
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}*/