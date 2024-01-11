//
//  githublight.go
//  GitHubLight
//
//  Created by Karl Kraft on 12/29/2023
//  Copyright 2023-2024 Karl Kraft. All rights reserved
//

package main

import (
	"context"
	"github.com/google/go-github/v57/github"
	"karlkraft.com/GitHubLight"
	"karlkraft.com/GitHubLight/api"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	settings := readSettings()
	client, err := api.NewClient(settings.LightServer)
	if err != nil {
		log.Fatalf("Unable to create LightServer client %v", err)
	}
	clientReport := githubScan(settings)
	_, err = client.ReportPost(context.Background(), &clientReport)
	if err != nil {
		log.Fatalf("%v", err)
	}

}

func githubScan(settings *GitHubLight.Settings) api.ClientReport {
	clientReport := api.ClientReport{
		Clientid: settings.ClientID,
	}
	gh := github.NewClient(nil).WithAuthToken(settings.GithubToken)
	options := github.SearchOptions{
		Sort:      "created",
		Order:     "asc",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 30,
		},
	}

	issues, _, err := gh.Search.Issues(context.Background(), "is:open is:pr review-requested:"+settings.Username, &options)
	if err != nil {
		log.Fatalf("Could not fetch issues (%v)", err)
	}
	reviewIssues := make([]api.ReportsItem, 0)
	for _, issue := range issues.Issues {
		reviewIssues = append(reviewIssues,
			api.ReportsItem{
				Type: api.ReportTupleReportsItem,
				ReportTuple: api.ReportTuple{
					Owner:      owner(issue.RepositoryURL),
					Repository: repository(issue.RepositoryURL),
					Section:    api.ReportTupleSectionReview,
					Age:        age(issue.CreatedAt),
				},
			})
	}
	clientReport.Reports = reviewIssues
	return clientReport
}

func age(ts *github.Timestamp) int {
	elapsed := time.Since(ts.Time)
	return int(elapsed.Seconds())
}

func repository(repoURL *string) string {
	u, err := url.Parse(*repoURL)
	if err != nil {
		return ""
	}
	pathComponents := strings.Split(u.Path, "/")
	return pathComponents[len(pathComponents)-1]
}

func owner(repoURL *string) string {
	u, err := url.Parse(*repoURL)
	if err != nil {
		return ""
	}
	pathComponents := strings.Split(u.Path, "/")
	return pathComponents[len(pathComponents)-2]
}

func readSettings() *GitHubLight.Settings {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to find users home directory %v", err)
	}
	return GitHubLight.ReadSettings(dirname + "/.githubLightBox")
}