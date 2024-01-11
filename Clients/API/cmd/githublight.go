//
//  githublight.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/29/2023
//  Copyright 2023-2024 Karl Kraft. All rights reserved
//

package main

import (
	"LightServer/arduino"
	"context"
	"fmt"
	"github.com/google/go-github/v57/github"
	"karlkraft.com/GitHubLight"
	"log"
	"math"
	"net"
	"os"
	"time"
)

func main() {
	settings := readSettings()

	addr := fmt.Sprintf("%s:%d", settings.BoxIP, settings.BoxPort)
	conn, err := net.Dial("udp", addr)
	if err != nil {
		log.Fatalf("Could not connect to %s (%v)", addr, err)
		return
	}

	lightTest(conn)

	client := github.NewClient(nil).WithAuthToken(settings.GithubToken)

	// set REVIEW light
	lightPattern := lightCommandForReview(settings, client)
	lightPattern.Send(conn)
	// get merge value
	// check for pulls
	// check for commit strength

}

func lightCommandForReview(settings *GitHubLight.Settings, client *github.Client) arduino.LightCommand {
	options := github.SearchOptions{
		Sort:      "created",
		Order:     "asc",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}
	issues, _, err := client.Search.Issues(context.Background(), "is:open is:pr review-requested:"+settings.Username, &options)
	if err != nil {
		log.Fatalf("Could not fetch issues (%v)", err)
	}
	if len(issues.Issues) > 0 {
		oldestIssue := issues.Issues[0]
		ts := oldestIssue.CreatedAt
		now := time.Now()
		age := math.Abs(ts.Sub(now).Hours())
		log.Printf("Oldest PR is %0.0f hours old.", age)
		log.Printf("%v", *oldestIssue.CommentsURL)
		idx := int(age)
		if idx > 7 {
			idx = 7
		}
		tone := palette[idx]
		return arduino.LightCommand{
			Start:  0,
			Length: 3,
			Red:    tone[0],
			Green:  tone[1],
			Blue:   tone[2],
		}
	} else {
		log.Printf("No PRs need review.")
		return arduino.LightCommand{
			Start:  0,
			Length: 3,
			Red:    0x00,
			Green:  0x00,
			Blue:   0x00,
		}
	}

}

func lightTest(conn net.Conn) {

	for x := 8; x < 9; x++ {
		set := palette[x]
		lightPattern := arduino.LightCommand{
			Start:  0,
			Length: 12,
			Red:    set[0],
			Green:  set[1],
			Blue:   set[2],
		}
		lightPattern.Send(conn)
		time.Sleep(time.Millisecond * 1000)
	}

}

func readSettings() *GitHubLight.Settings {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home directory %v", err)
	}
	return GitHubLight.ReadSettings(dirname + "/.githubLightBox")
}
