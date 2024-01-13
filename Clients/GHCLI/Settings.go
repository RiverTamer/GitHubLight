//
//  Settings.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/30/2023
//  Copyright 2023-2024 Karl Kraft. All rights reserved
//

package GitHubLight

import (
	"github.com/BurntSushi/toml"
	"log"
)

type CommandLineClientSettings struct {
	Frequency int
	ClientID  string
	Username  string
	Teams     []string
	Repos     []string
}

type Settings struct {
	LightServer       string
	CommandLineClient CommandLineClientSettings
}

func ReadSettings(path string) *Settings {
	s := &Settings{}
	meta, err := toml.DecodeFile(path, s)
	if err != nil {
		log.Fatalf("Unable to read file %s (%v)", path, err)
	}
	if !meta.IsDefined("LightServer") {
		log.Fatalf("Specify LightServer in configuration (e.g. LightServer=\"http://localhost:8080/\"")
	}
	// TODO - take a list of required keys
	// output all the missing
	// then fatal
	return s
}
