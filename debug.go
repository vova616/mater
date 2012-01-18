package mater

type DebugData struct {
	SingleStep bool
	Console    *Console
}

func (dd *DebugData) Init(mater *Mater) {
	dd.Console = new(Console)
	dd.Console.Init(mater)
}
