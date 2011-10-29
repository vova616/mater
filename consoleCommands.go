package mater

import (
	"fmt"
	"os"
)

type consoleCommand func(mater *Mater, params []string)

var commands = map[string]func(*Mater, []string) {
	"save": command_save,
	"load": command_load,
	"exit": command_quit,
	"quit": command_quit,
}
var commandNames []string

func init () {
	commandNames := make([]string, 0, len(commands))
	for key, _ := range commands {
		commandNames = append(commandNames, key)
	}
}

func command_save (mater *Mater, params []string) {
	if len(params) < 1 {
		fmt.Printf("Usage: save <filename>\n")
		return
	}

	path := params[0]
	fmt.Printf("Saving  to %v\n", path)
	mater.SaveScene(path)
}

func command_load (mater *Mater, params []string) {
	if len(params) < 1 {
		fmt.Printf("Usage: load <filename>\n")
		return
	}

	path := params[0]
	fmt.Printf("Loading  from %v\n", path)
	mater.Paused = true
	mater.LoadScene(path)
}

func command_quit (mater *Mater, params []string) {
	os.Exit(0)
}