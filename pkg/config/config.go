// Copyright © 2025 Prabhjot Singh Sethi, All Rights reserved
// Author: Prabhjot Singh Sethi <prabhjot.sethi@gmail.com>

package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// mongo db config struct
type MongoDB struct {
	Uri string `yaml:"uri,omitempty"`
}

// Base config struct
type BaseConfig struct {
	ConfigDB *MongoDB `yaml:"configDB,omitempty"`
}

// get Config database information, if the struct
// is nil it ensures sending the default mongodb
// config for base development scenarios
func (c *BaseConfig) GetConfigDB() *MongoDB {
	if c.ConfigDB != nil {
		return c.ConfigDB
	}

	return &MongoDB{
		Uri: "mongodb://localhost:27017/?replicaSet=rs0",
	}
}

// Parse YAML Config file from the provided config file path
// returns pointer to config structure and error if failed to
// generate the config struct.
// This also ensures handling scenarios when no config file
// is provided
func ParseConfig(filePath string) (*BaseConfig, error) {
	config := &BaseConfig{}
	// Process config file if file path is provided
	if filePath != "" {
		// open the provided config file
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		// ensure that we close the file before returning from
		// here, following constructs of release the unused
		// resources for garbage collector to kick in
		defer func() {
			_ = file.Close()
		}()

		// Get a new Yaml decoder
		decoder := yaml.NewDecoder(file)
		// decode the provided yaml config from the config file
		if err := decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	return config, nil
}
