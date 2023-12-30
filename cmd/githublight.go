//
//  githublight.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/29/2023
//  Copyright 2023 Karl Kraft. All rights reserved.
//

package main

import (
	"fmt"
	"karlkraft.com/GitHubLight"
	"log"
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

	//client := github.NewClient(nil).WithAuthToken(githubAccessToken)
	//options := github.SearchOptions{
	//	Sort:      "committer-date",
	//	Order:     "desc",
	//	TextMatch: false,
	//	ListOptions: github.ListOptions{
	//		Page:    1,
	//		PerPage: 30,
	//	},
	//}
	//commits, _, err := client.Search.Commits(context.Background(), "author:KarlKraft", &options)
	//if err != nil {
	//	log.Fatalf("Could not fetch commits. %v", err)
	//	return
	//}
	//log.Printf("Total commit count is %d", *commits.Total)

}

func lightTest(conn net.Conn) {

	colors := [2][3]uint8{
		{0x40, 0x40, 0x40},
		{0x0, 0x0, 0x0},
	}

	for x := 0; x < 2; x++ {
		set := colors[x]
		lightPattern := GitHubLight.LightCommand{
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
