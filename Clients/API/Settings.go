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

type Settings struct {
	GithubToken string `toml:"GITHUB_ACCESS_TOKEN"`
	Username    string
}

func ReadSettings(path string) *Settings {
	s := &Settings{}
	meta, err := toml.DecodeFile(path, s)
	if err != nil {
		log.Fatalf("Unable to read file %s (%v)", path, err)
	}
	if !meta.IsDefined("GITHUB_ACCESS_TOKEN") {
		log.Fatalf("Specify GITHUB_ACCESS_TOKEN in configuration (e.g. GITHUB_ACCESS_TOKEN=\"github_pat_xxxx...\"")
	}
	if !meta.IsDefined("Username") {
		log.Fatalf("Specify Username in configuration (e.g. Username=\"KarlKraft\"")
	}
	return s
}
