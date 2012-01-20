package main

import (
	"github.com/teomat/mater/engine"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Console struct {
	Scene   *engine.Scene
	Reader  *bufio.Reader
	Command chan string
	Buffer  *bytes.Buffer
}

func (console *Console) Init(scene *engine.Scene) {
	if console.Scene != nil {
		return
	}

	console.Command = make(chan string)
	console.Scene = scene
	console.Reader = bufio.NewReader(os.Stdin)
	console.Buffer = bytes.NewBuffer(nil)

	go console.ProcessInput()
}

func (console *Console) ProcessInput() {
	for {
		console.Command <- console.Read()
	}
}

func (console *Console) Read() string {
	command := ""
	var char byte
	console.Buffer.Truncate(0)
	for char != '\n' {
		char, _ = console.Reader.ReadByte()
		console.Buffer.WriteByte(char)
	}
	console.Buffer.Truncate(console.Buffer.Len() - 1)
	command = console.Buffer.String()
	return command
}

func (console *Console) ExecuteCommand(param string) {
	params := strings.Split(param, " ")
	command := params[0]
	params = params[1:]

	if commandFunc, ok := commands[command]; ok {
		commandFunc(console.Scene, params)
	} else {
		fmt.Printf("Command not found\n")
	}
}
