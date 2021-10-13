package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Dump a JSON object to a file
func dumpJsonToFile(obj interface{}, filename string) error {
	// MarshalIndent will return pretty-printed JSON, so a user may edit
	// the output file when the application is shut down, but should be really cautious
	objJson, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}

	fileHandle, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	if _, err := fileHandle.Write(objJson); err != nil {
		return err
	}

	return nil
}

// Read a JSON object from a file
func parseJsonFromFile(obj interface{}, filename string) error {
	fileHandle, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	objJson, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return err
	}

	err = json.Unmarshal(objJson, obj)
	if err != nil {
		return err
	}

	return nil
}
