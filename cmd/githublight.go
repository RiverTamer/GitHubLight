//
//  githublight.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/29/2022
//  Copyright 2022-2023 Karl Kraft. All rights reserved.
//

package main

import (
	"context"
	"github.com/google/go-github/v57/github"
	"karlkraft.com/GitHubLight/lifx"
	"log"
	"os"
)

func main() {
	githubAccessToken := os.Getenv("GITHUBLIGHT_ACCESS_TOKEN")
	if len(githubAccessToken) == 0 {
		log.Fatalf("You need to define GITHUBLIGHT_ACCESS_TOKEN.")
	}
	lifxIP := os.Getenv("LIFX_IP")
	if len(lifxIP) == 0 {
		log.Fatalf("You need to define LIFX_IP.")
	}
	packet := lifx.NewSetColorLIFXPacket(0.333333, 1.0, 1.0, 3500, 0)
	packet.Dump()

	client := github.NewClient(nil).WithAuthToken(githubAccessToken)
	options := github.SearchOptions{
		Sort:      "committer-date",
		Order:     "desc",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}
	commits, _, err := client.Search.Commits(context.Background(), "author:KarlKraft", &options)
	if err != nil {
		log.Fatalf("Could not fetch commits. %v", err)
		return
	}
	log.Printf("Total commit count is %d", *commits.Total)

}
