package mater

import (
	"bytes"
	"encoding/json"
	"github.com/teomat/mater/collision"
	"github.com/teomat/mater/transform"
	"log"
	"os"
)

var SaveDirectory = "saves/"

func (mater *Mater) SaveScene(path string) error {
	scene := mater.Scene

	path = SaveDirectory + path

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	//encoder := json.NewEncoder(file)
	//err = encoder.Encode(scene)

	dataString, err := json.MarshalIndent(scene, "", "\t")
	if err != nil {
		log.Printf("Error encoding Scene: %v", err)
		return err
	}

	buf := bytes.NewBuffer(dataString)
	n, err := buf.WriteTo(file)
	if err != nil {
		log.Printf("Error after writing %v characters to File: %v", n, err)
		return err
	}

	return nil
}

func (mater *Mater) LoadScene(path string) error {

	var scene *Scene

	path = SaveDirectory + path

	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	scene = new(Scene)
	decoder := json.NewDecoder(file)

	err = decoder.Decode(scene)
	if err != nil {
		log.Printf("Error decoding Scene: %v", err)
		return err
	}

	mater.Scene.Destroy()

	mater.Scene = scene
	scene.Space.Enabled = true

	return nil
}

func (scene *Scene) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	var err error

	buf.WriteByte('{')

	buf.WriteString(`"StaticEntity":`)
	encoder.Encode(&scene.StaticEntity)

	buf.WriteString(`,"Entities":`)
	entities, err := scene.MarshalEntities()
	if err != nil {
		return nil, err
	}
	buf.Write(entities)

	buf.WriteString(`,"Space":`)
	err = encoder.Encode(scene.Space)
	if err != nil {
		log.Printf("Error encoding scene")
		return nil, err
	}

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (scene *Scene) MarshalEntities() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

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

	return buf.Bytes(), nil
}

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

	staticEntity, err := scene.UnmarshalEntity(sd.StaticEntity)
	if err != nil {
		log.Printf("Error decoding static entity")
		return err
	}
	scene.StaticEntity = *staticEntity

	if scene.Entities == nil {
		scene.Entities = make(map[int]*Entity, 32)
	}

	for _, rawEntity := range sd.Entities {
		entity, err := scene.UnmarshalEntity(rawEntity)
		if err != nil {
			log.Printf("Error decoding entity")
			return err
		}
		scene.AddEntity(entity)
	}

	return nil
}

func (scene *Scene) UnmarshalEntity(data []byte) (*Entity, error) {
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
		return nil, err
	}

	entity := new(Entity)
	entity.Scene = scene

	if ed.ID > lastEntityId {
		lastEntityId = ed.ID
	}

	entity.id = ed.ID
	if entity.id <= 0 {
		entity.id = nextId()
	}

	entity.Enabled = ed.Enabled
	entity.Transform = ed.Transform

	entity.Components = make(map[string]Component, len(ed.Components))
	entity.ComponentList = make([]Component, 0, len(ed.Components))

	for _, componentData := range ed.Components {
		name := componentData.Name
		component := NewComponent(name)
		if component == nil {
			continue
		}
		err := component.Unmarshal(entity, componentData.Data)
		if err != nil {
			log.Printf("Error decoding entity")
			return nil, err
		}
		entity.AddComponent(component)
	}

	return entity, nil
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
		data, err := component.Marshal(entity)
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
