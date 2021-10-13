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
