package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/teomat/mater/collision"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var SaveDirectory = "saves/"

func saveSpace(path string) error {
	path = Settings.SaveDir + path

	file, err := os.Create(path)
	if err != nil {
		log.Printf("Error opening File: %v", err)
		return err
	}
	defer file.Close()

	dataString, err := json.MarshalIndent(space, "", "\t")
	if err != nil {
		log.Printf("Error encoding Space: %v", err)
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

func loadSpace(path string) error {

	var newSpace *collision.Space

	path = Settings.SaveDir + path

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading File: %v", err)
		return err
	}

	newSpace = new(collision.Space)

	err = json.Unmarshal(data, newSpace)
	if err != nil {
		log.Printf("Error decoding Space")
		printSyntaxError(string(data), err)
		return err
	}

	space = newSpace
	space.Enabled = true

	return nil
}

func printSyntaxError(js string, err error) {
	syntax, ok := err.(*json.SyntaxError)
	if !ok {
		fmt.Println(err)
		return
	}

	start, end := strings.LastIndex(js[:syntax.Offset], "\n")+1, len(js)
	if idx := strings.Index(js[start:], "\n"); idx >= 0 {
		end = start + idx
	}

	line, pos := strings.Count(js[:start], "\n"), int(syntax.Offset)-start-1
	line = line + 1

	fmt.Printf("Error in line %d: %s \n", line, err)

	if start > 0 && start < end {
		fmt.Printf("%s\n%s^", strings.Replace(js[start:end], "\t", " ", -1), strings.Repeat(" ", pos))
	}
}
