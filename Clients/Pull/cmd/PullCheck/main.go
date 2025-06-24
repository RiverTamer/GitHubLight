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
	keysToDump := ""
	clientReport.Reports, keysToDump = recursiveProcessFolder(path, verbose)
	displayKeys(keysToDump)

	_, err = client.ReportPost(context.Background(), &clientReport)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}

func recursiveProcessFolder(path string, verbose bool) ([]api.ReportsItem, string) {
	keysToDump := ""
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
				var newKeys string
				newReports, newKeys = recursiveProcessFolder(path+"/"+file.Name(), verbose)
				keysToDump += newKeys
				reports = append(reports, newReports...)
			}
		}
		return reports, keysToDump
	}
}

//goland:noinspection ALL
func processFolder(path string, verbose bool) ([]api.ReportsItem, string) {
	keysToDump := ""
	if verbose {
		log.Printf("Scanning %s", path)
		keysToDump = " "
	}
	pullReports := make([]api.ReportsItem, 0)
	hostName, _ := os.Hostname()
	err := fetchRemotes(path)
	if err != nil {
		log.Printf("!! Could not use git fetch (%v)", err)
		return pullReports, ""
	}

	currentBranchName, err := findCurrentBranch(path)
	if err != nil {
		log.Printf("!! Could not findCurrentBranch() (%v)", err)
		return pullReports, ""
	}
	currentBranchDirty, err := isCurrentBranchDirty(path)
	if err != nil {
		log.Printf("!! Could not isCurrentBranchDirty() (%v)", err)
		return pullReports, ""
	}

	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short) %(push:track)", "refs/heads")
	cmd.Dir = path
	statusData, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Could not get branch list on %s (%v)", path, err)
		return pullReports, ""
	}
	for _, s := range strings.Split(string(statusData), "\n") {
		if len(s) == 0 {
			continue
		}
		branchName := strings.Fields(s)[0]
		if !strings.Contains(s, "[behind ") && !strings.Contains(s, ", behind ") {
			// Nothing to pull
			if verbose {
				if len(keysToDump) == 0 {
					log.Printf("Scanning %s", path)
				}
				log.Printf("✅ " + branchName)
				keysToDump += "✅"
			}
		} else if branchName == currentBranchName {
			if currentBranchDirty {
				if len(keysToDump) == 0 {
					log.Printf("Scanning %s", path)
				}
				// Current branch needs pull but is dirty
				behind, _ := commitsBehind(s)
				log.Printf("🚨 %s (%d)", branchName, behind)
				keysToDump += "🚨"
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
				if len(keysToDump) == 0 {
					log.Printf("Scanning %s", path)
				}
				// the current branch can be pulled directly
				behind, _ := commitsBehind(s)
				log.Printf("⤵️ %s (%d)", branchName, behind)
				keysToDump += "⤵️"
				pullCurrentBranch(path)
			}
		} else if strings.Contains(s, "ahead ") {
			if len(keysToDump) == 0 {
				log.Printf("Scanning %s", path)
			}
			// Branch need to be pulled but also needs a push first
			behind, _ := commitsBehind(s)
			log.Printf("⬆️ %s (%d)", branchName, behind)
			keysToDump += "⬆️"
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
			if len(keysToDump) == 0 {
				log.Printf("Scanning %s", path)
			}
			// non-current branch, can be pulled with fetch
			behind, _ := commitsBehind(s)
			log.Printf("⬇️ %s (%d)", branchName, behind)
			keysToDump += "⬇️"
			updateBranch(path, branchName)
		}
	}
	return pullReports, keysToDump
}

func displayKeys(keysToDump string) {
	if strings.Contains(keysToDump, "✅") {
		log.Printf("✅ branch is synced")
	}
	if strings.Contains(keysToDump, "⬆️") {
		log.Printf("⬆️ branch needs to be pushed")
	}
	if strings.Contains(keysToDump, "⬇️") {
		log.Printf("⬇️ branch is behind and was pulled")
	}
	if strings.Contains(keysToDump, "⬇️") {
		log.Printf("⬇️ branch is behind and was fetched")
	}
	if strings.Contains(keysToDump, "🚨") {
		log.Printf("🚨 branch cannot be synced (uncommitted changes)")
	}

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
