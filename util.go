// List-O-Matic Talking List Management System
// Copyright (C) 2021 Jan Blaesi
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
