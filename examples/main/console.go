package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

type Console struct {
	Reader  *bufio.Reader
	Command chan string
	Buffer  *bytes.Buffer
}

//Initializes the debug console and starts a goroutine to read from stdin
func (console *Console) Init() {
	console.Command = make(chan string)
	console.Reader = bufio.NewReader(os.Stdin)
	console.Buffer = bytes.NewBuffer(nil)

	go func() {
		for {
			console.Command <- console.Read()
		}
	}()
}

//Reads up to a newline from stdin
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

//Tries to run the given string as a command.
//Everything after the first space is passed to the command as parameters.
func (console *Console) ExecuteCommand(param string) {
	params := strings.Split(param, " ")
	commandName := params[0]
	params = params[1:]

	if commandFunc, ok := commands[commandName]; ok {
		commandFunc(params)
	} else {
		fmt.Printf("Command %v not found\n", commandName)
	}
}
