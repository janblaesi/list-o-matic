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
