//
//  main.go
//  PullCheck
//
//  Created by Karl Kraft on 1/12/2024
//  Copyright 2024-2025 Karl Kraft. All rights reserved
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
	var verbose bool
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	flag.StringVar(&path, "p", "path", "repository or path to report")
	flag.StringVar(&clientid, "c", "clientid", "client id to report")
	flag.BoolVar(&verbose, "v", false, "Report every step")
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
	var dumpKey bool
	clientReport.Reports, dumpKey = recursiveProcessFolder(path, verbose)
	if dumpKey {
		displayKey()
	}

	_, err = client.ReportPost(context.Background(), &clientReport)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func recursiveProcessFolder(path string, verbose bool) ([]api.ReportsItem, bool) {
	dumpKey := false
	reports := make([]api.ReportsItem, 0)
	if _, err := os.Stat(path + "/.git"); !os.IsNotExist(err) {
		return processFolder(path, verbose)
	} else {
		entries, err := os.ReadDir(path)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
		for _, file := range entries {
			if (file.Type() & os.ModeDir) > 0 {
				var newReports []api.ReportsItem
				var b bool
				newReports, b = recursiveProcessFolder(path+"/"+file.Name(), verbose)
				dumpKey = dumpKey || b
				reports = append(reports, newReports...)
			}
		}
		return reports, dumpKey
	}
}

//goland:noinspection ALL
func processFolder(path string, verbose bool) ([]api.ReportsItem, bool) {
	reportPath := path
	pullReports := make([]api.ReportsItem, 0)
	hostName, _ := os.Hostname()
	err := fetchRemotes(path)
	if err != nil {
		log.Printf("!! Could not use git fetch (%v)", err)
		return pullReports, true
	}

	currentBranchName, err := findCurrentBranch(path)
	if err != nil {
		log.Printf("!! Could not findCurrentBranch() (%v)", err)
		return pullReports, true
	}
	currentBranchDirty, err := isCurrentBranchDirty(path)
	if err != nil {
		log.Printf("!! Could not isCurrentBranchDirty() (%v)", err)
		return pullReports, true
	}

	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short) %(push:track)", "refs/heads")
	cmd.Dir = path
	statusData, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not get branch list on %s (%v)", path, err)
		return pullReports, true
	}
	for _, s := range strings.Split(string(statusData), "\n") {
		if len(s) == 0 {
			continue
		}
		branchName := strings.Fields(s)[0]
		if !strings.Contains(s, "[behind ") && !strings.Contains(s, ", behind ") {
			// Nothing to pull
			if verbose {
				if len(reportPath) > 0 {
					log.Printf("Scanning %s", path)
					reportPath = ""
				}
				log.Printf("✅ " + branchName)
			}
		} else if branchName == currentBranchName {
			if currentBranchDirty {
				if len(reportPath) > 0 {
					log.Printf("Scanning %s", path)
					reportPath = ""
				}
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
				if len(reportPath) > 0 {
					log.Printf("Scanning %s", path)
					reportPath = ""
				}
				// the current branch can be pulled directly
				behind, _ := commitsBehind(s)
				log.Printf("⤵️ %s (%d)", branchName, behind)
				pullCurrentBranch(path)
			}
		} else if strings.Contains(s, "ahead ") {
			if len(reportPath) > 0 {
				log.Printf("Scanning %s", path)
				reportPath = ""
			}
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
			if len(reportPath) > 0 {
				log.Printf("Scanning %s", path)
				reportPath = ""
			}
			// non-current branch, can be pulled with fetch
			behind, _ := commitsBehind(s)
			log.Printf("⬇️ %s (%d)", branchName, behind)
			updateBranch(path, branchName)
		}
	}
	return pullReports, len(reportPath) == 0
}

func displayKey() {
	log.Printf("✅ branch is synced")
	log.Printf("⬆️ branch needs to be pushed")
	log.Printf("⬇️ branch is behind and was pulled")
	log.Printf("⬇️ branch is behind and was fetched")
	log.Printf("🚨 branch cannot be synced (uncommitted changes)")

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
