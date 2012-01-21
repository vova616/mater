package engine

import (
	"bytes"
	"encoding/json"
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/transform"
	"log"
)

func (scene *Scene) UnmarshalJSON(data []byte) error {
	sceneData := struct {
		Space        *collision.Space
		StaticEntity json.RawMessage
		Entities     []json.RawMessage
	}{}

	err := json.Unmarshal(data, &sceneData)
	if err != nil {
		log.Printf("Error decoding scene")
		return err
	}

	sd := &sceneData

	scene.Space = sd.Space

	//BEGIN REGION ENTITY DATA
	type entityData struct {
		ID         int
		Enabled    bool
		Transform  *transform.Transform
		Components []Component
	}

	unmarshalEntityData := func(entity *entityData, data []byte) error {
		type compData struct {
			Name string
			Data json.RawMessage
		}

		entityData := struct {
			ID         int
			Enabled    bool
			Transform  *transform.Transform
			Components []compData
		}{}
		ed := &entityData

		err := json.Unmarshal(data, ed)
		if err != nil {
			log.Printf("Error decoding entity")
			return err
		}

		if ed.ID > lastEntityId {
			lastEntityId = ed.ID
		}

		entity.ID = ed.ID
		if entity.ID <= 0 {
			entity.ID = nextId()
		}

		entity.Enabled = ed.Enabled
		entity.Transform = ed.Transform

		entity.Components = make([]Component, 0, len(ed.Components))

		for _, componentData := range ed.Components {
			name := componentData.Name
			component := NewComponent(name)
			if component == nil {
				continue
			}
			err := component.Unmarshal(componentData.Data)
			if err != nil {
				log.Printf("Error decoding entity")
				return err
			}
			entity.Components = append(entity.Components, component)
		}

		return nil
	}
	//END REGION ENTITY DATA

	//actual decoding starts here
	staticEntity := new(entityData)
	err = unmarshalEntityData(staticEntity, sd.StaticEntity)
	if err != nil {
		log.Printf("Error decoding static entity")
		return err
	}

	scene.Entities = make(map[int]*Entity, len(sd.Entities))

	entities := make([]entityData, len(sd.Entities))
	for i, rawEntity := range sd.Entities {
		entity := &entities[i]

		err := unmarshalEntityData(entity, rawEntity)
		if err != nil {
			log.Printf("Error decoding entity")
			return err
		}
	}
	//after this point everything is decoded and we can use the callbacks in scene.Callbacks

	//creating the entities starts here
	createEntityFromData := func(entity *Entity, ed *entityData) {
		entity.Scene = scene
		entity.id = ed.ID
		entity.Enabled = ed.Enabled
		entity.Transform = ed.Transform
		entity.Components = make(map[string]Component, len(ed.Components))
		entity.ComponentList = make([]Component, 0, len(ed.Components))
		for _, component := range ed.Components {
			entity.AddComponent(component)
		}
	}

	createEntityFromData(&scene.StaticEntity, staticEntity)
	
	for _, entityData := range entities {
		entity := new(Entity)
		createEntityFromData(entity, &entityData)
		scene.AddEntity(entity)
	}

	return nil
}

func (scene *Scene) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	var err error

	buf.WriteByte('{')

	buf.WriteString(`"StaticEntity":`)
	encoder.Encode(scene.StaticEntity)

	buf.WriteString(`,"Entities":`)
	buf.WriteByte('[')
	for _, entity := range scene.Entities {

		err := encoder.Encode(entity)
		if err != nil {
			log.Printf("Error encoding entity")
			return nil, err
		}

		buf.WriteByte(',')
	}
	if len(scene.Entities) > 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte(']')
	if err != nil {
		return nil, err
	}

	buf.WriteString(`,"Space":`)
	err = encoder.Encode(scene.Space)
	if err != nil {
		log.Printf("Error encoding scene")
		return nil, err
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (entity *Entity) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')
	buf.WriteString(`"ID":`)
	encoder.Encode(entity.id)

	buf.WriteString(`,"Enabled":`)
	encoder.Encode(entity.Enabled)

	buf.WriteString(`,"Transform":`)
	encoder.Encode(entity.Transform)

	buf.WriteString(`,"Components":`)
	buf.WriteByte('[')
	ccount := 0
	for _, component := range entity.ComponentList {
		name := component.Name()
		data, err := component.Marshal()
		if err != nil {
			log.Printf("Error encoding entity")
			return nil, err
		}
		//Don't write anything if the component returns nil
		if data == nil {
			continue
		}
		ccount++

		buf.WriteByte('{')

		buf.WriteString(`"Name":`)
		encoder.Encode(name)

		buf.WriteString(`,"Data":`)
		buf.Write(data)

		buf.WriteByte('}')

		buf.WriteByte(',')
	}

	if ccount > 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte(']')
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
