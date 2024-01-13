//
//  main.go
//  GitHubLight
//
//  Created by Karl Kraft on 1/12/2024
//  Copyright 2024 Karl Kraft. All rights reserved
//

package main

import (
	"context"
	"github.com/ogen-go/ogen/json"
	GitHubLight "karlkraft.com/ghcli"
	"karlkraft.com/ghcli/api"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"
)

type PullRequestList []PullRequest

type PullRequest struct {
	IsDraft        bool             `json:"isDraft"`
	Author         AuthorInfo       `json:"author"`
	Created        time.Time        `json:"createdAt"`
	Updated        time.Time        `json:"updatedAt"`
	MergeStatus    string           `json:"mergeStateStatus"` // e.g. BLOCKED
	ReviewDecision string           `json:"reviewDecision"`   // e.g. REVIEW_REQUIRED
	Requests       []ReviewRequests `json:"reviewRequests"`
	Reviews        []Review         `json:"reviews"`
	URL            string           `json:"url"`
}

type Review struct {
	Author     AuthorInfo `json:"author"`
	ReviewType string     `json:"state"` // e.g. COMMENTED, APPROVED

}

type ReviewRequests struct {
	Name string `json:"name"`
}

type AuthorInfo struct {
	Bot   bool   `json:"is_bot"`
	Login string `json:"login"`
}

func main() {
	settings := readSettings()
	client, err := api.NewClient(settings.LightServer)
	if err != nil {
		log.Fatalf("Unable to create LightServer client %v", err)
	}
	for true {
		log.Println("Scanning")
		clientReport := githubScan(settings)
		_, err = client.ReportPost(context.Background(), &clientReport)
		if err != nil {
			log.Fatalf("%v", err)
		}
		time.Sleep(time.Minute * 5)
	}

}

func githubScan(settings *GitHubLight.Settings) api.ClientReport {
	clientReport := api.ClientReport{
		Clientid: settings.CommandLineClient.ClientID,
	}
	reportIssues := make([]api.ReportsItem, 0)

	for _, repoName := range settings.CommandLineClient.Repos {
		repoComponents := strings.Split(repoName, "/")
		owner := repoComponents[0]
		repoShortName := repoComponents[1]

		prListCommand := exec.Command("gh", "pr", "list", "-R", repoName, "-S", "is:open is:pr", "--json", "author,comments,createdAt,isDraft,mergeStateStatus,reviewDecision,reviewRequests,reviews,state,updatedAt,url")
		output, err := prListCommand.CombinedOutput()
		if err != nil {
			log.Printf("Unable to fetch repo %s  (%v)\n", repoName, err)
			continue
		}
		res := make(PullRequestList, 0)
		err = json.Unmarshal([]byte(output), &res)
		if err != nil {
			log.Printf("Unable to parse JSON (%v)\n", repoName, err)
			continue
		}
		for _, pr := range res {
			if pr.IsDraft {
				continue
			}
			if pr.ReviewDecision == "APPROVED" {
				if pr.Author.Login == settings.CommandLineClient.Username {
					anItem := api.ReportsItem{
						Type: api.ReportTupleReportsItem,
						ReportTuple: api.ReportTuple{
							Owner:      owner,
							Repository: repoShortName,
							Section:    api.ReportTupleSectionMerge,
							Age:        age(pr.Created, pr.Updated),
							Reference:  pr.URL,
						},
					}
					reportIssues = append(reportIssues, anItem)
				} else {
					continue
				}
			} else {
				commentMade := false
				for _, review := range pr.Reviews {
					if review.Author.Login == settings.CommandLineClient.Username {
						commentMade = true
					}
				}
				commentDesired := false
				for _, request := range pr.Requests {
					if request.Name == settings.CommandLineClient.Username {
						commentDesired = true
					}
					if slices.Contains(settings.CommandLineClient.Teams, request.Name) {
						commentDesired = true
					}
				}

				if commentDesired && !commentMade {
					anItem := api.ReportsItem{
						Type: api.ReportTupleReportsItem,
						ReportTuple: api.ReportTuple{
							Owner:      owner,
							Repository: repoShortName,
							Section:    api.ReportTupleSectionReview,
							Age:        age(pr.Created, pr.Updated),
							Reference:  pr.URL,
						},
					}
					reportIssues = append(reportIssues, anItem)
				}
			}
		}
	}
	clientReport.Reports = reportIssues
	return clientReport
}

func age(created time.Time, updated time.Time) int {
	createdAge := time.Since(created)
	updatedAge := time.Since(updated)
	return int(min(createdAge.Seconds(), updatedAge.Seconds()))
}

func readSettings() *GitHubLight.Settings {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home directory %v", err)
	}
	return GitHubLight.ReadSettings(dirname + "/.githubLightBox")
}
