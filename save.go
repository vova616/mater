package mater

import (
	"bytes"
	"json"
	"os"
	"box2d"
	"fmt"
)

var saveDirectory = "saves/"

func (mater *Mater) SaveScene (path string) os.Error{
	scene := mater.Scene

	path = saveDirectory + path

	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	//encoder := json.NewEncoder(file)
	//err = encoder.Encode(scene)

	dataString, err := json.MarshalIndent(scene, "", "\t")
	if err != nil {
		fmt.Printf("Error encoding Scene: %v", err)
		return err
	}

	buf := bytes.NewBuffer(dataString)
	n, err := buf.WriteTo(file)
	if err != nil {
		fmt.Printf("Error after writing %v characters to File: %v", n, err)
		return err
	}

	return nil
}

func (mater *Mater) LoadScene (path string) os.Error {

	var scene *Scene

	path = saveDirectory + path

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	scene = new(Scene)
	decoder := json.NewDecoder(file)

	err = decoder.Decode(scene)
	if err != nil {
		fmt.Printf("Error decoding Scene: %v", err)
		return err
	}

	mater.Scene.Destroy()

	mater.Scene = scene
	scene.World.Enabled = true

	if mater.Scene.Camera == nil {
		cam := mater.DefaultCamera
		mater.Scene.Camera = &cam
	} else {
		mater.Scene.Camera.ScreenSize = mater.ScreenSize
	}

	cl := &MaterContactListener{mater}
	mater.Scene.World.SetContactListener(cl)
	mater.Scene.World.SetContactFilter(cl)
	mater.Dbg.DebugView.Reset(mater.Scene.World)

	return nil
}

func (scene *Scene) MarshalJSON() ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	var err os.Error

	buf.WriteString(`{"Camera":`)
	encoder.Encode(scene.Camera)

	buf.WriteString(`,"LastEntityId":`)
	encoder.Encode(lastEntityId)

	buf.WriteString(`,"Entities":`)
	entities, err := scene.MarshalEntities()
	if err != nil {
		return nil, err
	}
	buf.Write(entities)


	buf.WriteString(`,"World":`)
	world, err := scene.MarshalWorld()
	if err != nil {
		return nil, err
	}
	buf.Write(world)

	buf.WriteByte('}')

	return buf.Bytes(), nil
}

func (scene *Scene) MarshalEntities() ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('[')
	for _, entity := range scene.Entities {

		err := encoder.Encode(entity)
		if err != nil {
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

func (scene *Scene) MarshalWorld() ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	world := scene.World

	buf.WriteByte('{')
	buf.WriteString(`"Gravity":`)
	encoder.Encode(world.Gravity)

	buf.WriteString(`,"Bodies":`)
	buf.WriteByte('[')

	//actual number of serialized bodies may be different than the total number of bodies
	//because of that we keep track how many we actually write
	bodyNum := 0
	for _, body := range world.BodyList() {
		
		//if a bodies UserData is set, we do not serialize it
		if body.UserData != nil {
			continue
		}
		bodyNum++

		err := encoder.Encode(body)
		if err != nil {
			return nil, err
		}

		buf.WriteByte(',')
	}

	if bodyNum > 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}
	buf.WriteByte(']')

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (scene *Scene) UnmarshalJSON(data []byte) os.Error {
	sceneData := struct {
		LastEntityId int
		Camera *Camera
		World *box2d.World
		Entities []json.RawMessage
	}{}

	err := json.Unmarshal(data, &sceneData)
	if err != nil {
		return err
	}

	sd := &sceneData

	lastEntityId = sd.LastEntityId

	scene.Camera = sd.Camera
	scene.World = sd.World

	if scene.Entities == nil {
		scene.Entities = make(map[int]*Entity, 32)
	}

	for _, rawEntity := range sd.Entities {
		err := scene.UnmarshalEntity(rawEntity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (scene *Scene) UnmarshalEntity(data []byte) os.Error {
	entityData := struct {
		ID int
		Enabled bool
		Components map[string]json.RawMessage
	}{}
	ed := &entityData

	err := json.Unmarshal(data, ed)
	if err != nil {
		return err
	}

	entity := new(Entity)
	entity.Scene = scene

	entity.id = ed.ID
	entity.Enabled = ed.Enabled

	entity.Components = make(map[string]Component, len(ed.Components))

	for name, componentData := range ed.Components {
		component := NewComponent(name)
		if component == nil {
			continue
		}
		err := component.Unmarshal(entity, componentData)
		if err != nil {
			return err
		}
		entity.Components[name] = component
		component.Init(entity)
	}

	scene.AddEntity(entity)

	return nil
}

func (entity *Entity) MarshalJSON() ([]byte, os.Error) {
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)

	buf.WriteByte('{')
	buf.WriteString(`"ID":`)
	encoder.Encode(entity.id)

	buf.WriteString(`,"Enabled":`)
	encoder.Encode(entity.Enabled)

	buf.WriteString(`,"Components":`)
	buf.WriteByte('{')
	ccount := 0
	for _, component := range entity.Components {
		name := component.Name()
		data, err := component.Marshal(entity)
		if err != nil {
			return nil, err
		}
		//Don't write anything if the component returns nil
		if data == nil {
			continue
		}

		ccount++
		buf.WriteByte('"')
		buf.WriteString(name)
		buf.WriteString(`":`)		
		buf.Write(data)
		buf.WriteByte(',')
	}

	if ccount > 0 {
		//cut trailing comma
		buf.Truncate(buf.Len() - 1)
	}

	buf.WriteByte('}')
	buf.WriteByte('}')

	return buf.Bytes(), nil
}
