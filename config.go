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
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the configuration of this software
type Config struct {
	// Database contains paths to the JSON files containing users and talking lists
	Database struct {
		// The path to the JSON file containing the talking lists
		TalkingListsPath string `yaml:"talking_lists"`

		// The path to the JSON file containing the users
		UsersPath string `yaml:"users"`
	} `yaml:"database"`

	// Authentication contains settings for the authentication system using JSON Web Tokens
	Authentication struct {
		// The Secret to use for signing JSON Web Tokens
		Secret string `yaml:"secret"`

		// Timeout in seconds, after which a JSON Web Token loses its validity
		TimeoutSeconds int `yaml:"timeout_seconds"`
	} `yaml:"authentication"`
}

// Represents the configuration data of this software
var cfg Config

// This function tries to load the current configuration from 'config.yml'.
func cfgLoad() error {
	file, err := os.Open("config.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	return nil
}
