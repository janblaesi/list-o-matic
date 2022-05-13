//     __    _      __        ____        __  ___      __  _
//    / /   (_)____/ /_      / __ \      /  |/  /___ _/ /_(_)____
//   / /   / / ___/ __/_____/ / / /_____/ /|_/ / __ `/ __/ / ___/
//  / /___/ (__  ) /_/_____/ /_/ /_____/ /  / / /_/ / /_/ / /__
// /_____/_/____/\__/      \____/     /_/  /_/\__,_/\__/_/\___/
//
// Copyright 2021-2022 Jan Blaesi
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
// THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
// CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

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
