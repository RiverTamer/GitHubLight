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
	"karlkraft.com/githublightpull"
	"karlkraft.com/githublightpull/api"
	"log"
	"os"
	"os/exec"
	"regexp"
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
	err := fetchRemotes(path)
	if err != nil {
		log.Printf("!! Could not use git fetch (%v)", err)
		return pullReports
	}

	currentBranchName, err := findCurrentBranch(path)
	if err != nil {
		log.Printf("!! Could not findCurrentBranch() (%v)", err)
		return pullReports
	}
	currentBranchDirty, err := isCurrentBranchDirty(path)
	if err != nil {
		log.Printf("!! Could not isCurrentBranchDirty() (%v)", err)
		return pullReports
	}

	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short) %(push:track)", "refs/heads")
	cmd.Dir = path
	statusData, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not get branch list on %s (%v)", path, err)
		return pullReports
	}
	for _, s := range strings.Split(string(statusData), "\n") {
		if len(s) == 0 {
			continue
		}
		branchName := strings.Fields(s)[0]
		if !strings.Contains(s, "[behind ") && !strings.Contains(s, ", behind ") {
			// Nothing to pull
			log.Printf("✅ " + branchName)
		} else if branchName == currentBranchName {
			if currentBranchDirty {
				// Current branch needs pull but is dirty
				behind, _ := commitsBehind(s)
				log.Printf("🚨 %s (%d)", branchName, behind)
				anItem := api.ReportsItem{
					Type: api.ReportTupleReportsItem,
					ReportTuple: api.ReportTuple{
						Repository: path,
						Section:    api.ReportTupleSectionPull,
						Age:        behind,
						URL:        "http://" + hostName + "/",
						Notes:      branchName,
					},
				}
				pullReports = append(pullReports, anItem)
			} else {
				// current branch can be pulled directly
				behind, _ := commitsBehind(s)
				log.Printf("⤵️ %s (%d)", branchName, behind)
				pullCurrentBranch(path)
			}
		} else if strings.Contains(s, "ahead ") {
			// Branch need to be pulled but also needs a push first
			behind, _ := commitsBehind(s)
			log.Printf("⬆️ %s (%d)", branchName, behind)
			anItem := api.ReportsItem{
				Type: api.ReportTupleReportsItem,
				ReportTuple: api.ReportTuple{
					Repository: path,
					Section:    api.ReportTupleSectionPull,
					Age:        behind,
					URL:        "http://" + hostName + "/",
					Notes:      branchName,
				},
			}
			pullReports = append(pullReports, anItem)
		} else {
			// non-current branch, can be pulled with fetch
			behind, _ := commitsBehind(s)
			log.Printf("⬇️ %s (%d)", branchName, behind)
			updateBranch(path, branchName)
		}
	}
	return pullReports
}

func pullCurrentBranch(path string) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	print(string(output))
}

func updateBranch(path string, name string) {
	cmd := exec.Command("git", "fetch", "-u", "origin", name+":"+name)
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	print(string(output))
}

func commitsBehind(s string) (int, error) {
	const regexMatch = `behind (\d+)` // to match updated_at=123456
	rx, err := regexp.Compile(regexMatch)
	if err != nil {
		return 0, err
	}
	res := rx.FindStringSubmatch(s)
	return strconv.Atoi(res[1])
}

func fetchRemotes(path string) error {
	cmd := exec.Command("git", "fetch")
	cmd.Dir = path
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
func isCurrentBranchDirty(path string) (bool, error) {
	cmd := exec.Command("git", "status", "-s")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return true, err
	}
	s := string(output)
	return len(s) > 0, nil

}

func findCurrentBranch(path string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	s := string(output)
	s = s[:len(s)-1]
	return s, nil
}

func readSettings() *GitHubLight.Settings {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home directory %v", err)
	}
	return GitHubLight.ReadSettings(dirname + "/.githubLightBox")
}
