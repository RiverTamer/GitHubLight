//
//  main.go
//  PullCheck
//
//  Created by Karl Kraft on 1/12/2024
//  Copyright 2024 Karl Kraft. All rights reserved
//

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"karlkraft.com/githublightpull"
	"karlkraft.com/githublightpull/api"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	var path string
	var clientid string

	flag.StringVar(&path, "p", "path", "repository or path to report")
	flag.StringVar(&clientid, "c", "clientid", "client id to report")
	flag.Parse()

	settings := readSettings()
	client, err := api.NewClient(settings.LightServer)
	if err != nil {
		log.Fatalf("Unable to create LightServer client %v", err)
	}
	clientReport := api.ClientReport{
		Clientid: clientid,
		Reports:  make([]api.ReportsItem, 0),
	}
	if _, err := os.Stat(path + "/.git"); !os.IsNotExist(err) {
		clientReport.Reports = processFolder(path)
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		for _, file := range entries {
			if (file.Type() & os.ModeDir) > 0 {
				clientReport.Reports = append(clientReport.Reports, processFolder(path+"/"+file.Name())...)
			}
		}
	}
	_, err = client.ReportPost(context.Background(), &clientReport)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func processFolder(path string) []api.ReportsItem {
	log.Printf("Scanning %s", path)
	hostName, _ := os.Hostname()
	pullReports := make([]api.ReportsItem, 0)
	cmd := exec.Command("git", "fetch")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not use git fetch %v", err)
		return nil
	}

	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatalf("Could not open %s\n", path)
	}

	branches, err := r.Branches()
	if err != nil {
		log.Fatalf("Could not get branch list")
	}
	_ = branches.ForEach(func(branch *plumbing.Reference) error {
		components := strings.Split(branch.Name().String(), "/")
		branchName := components[len(components)-1]
		log.Printf("  Branch %s", branchName)
		revCommand := exec.Command("git", "rev-list", "--left-right", "--count", fmt.Sprintf("%s...origin/%s", branchName, branchName))
		revCommand.Dir = path
		output, err = revCommand.CombinedOutput()
		if err != nil {
			log.Printf("  ^^ Error (may not exist on remote")
			return nil
		}
		mergeCountFields := strings.Fields(string(output))
		commitsToMerge, _ := strconv.Atoi(mergeCountFields[1])
		if commitsToMerge > 0 {
			//log.Printf("Branch %s in path %s needs to be updated\n", branchName, path)
			anItem := api.ReportsItem{
				Type: api.ReportTupleReportsItem,
				ReportTuple: api.ReportTuple{
					Owner:      hostName,
					Repository: path,
					Section:    api.ReportTupleSectionPull,
					Age:        commitsToMerge,
					Reference:  branchName,
				},
			}
			pullReports = append(pullReports, anItem)
		}
		return nil
	})
	return pullReports
}

func readSettings() *GitHubLight.Settings {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home directory %v", err)
	}
	return GitHubLight.ReadSettings(dirname + "/.githubLightBox")
}
