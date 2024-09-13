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
	"fmt"
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: githublight /path/to/config.toml")
		os.Exit(-2)
	}
	settings := GitHubLight.ReadSettings(os.Args[1])

	client, err := api.NewClient(settings.LightServer)
	if err != nil {
		log.Fatalf("Unable to create LightServer client %v", err)
	}
	for true {
		log.Println("Scanning")
		clientReport := githubScan(settings)
		_, err = client.ReportPost(context.Background(), &clientReport)
		if err != nil {
			log.Println("Unable to scan %v", err)
		}
		time.Sleep(time.Minute * 5)
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
	} else {
		fmt.Println("Review requested: ", len(issues.Issues))
	}

	reportIssues := make([]api.ReportsItem, 0)
	for _, issue := range issues.Issues {
		reportIssues = append(reportIssues,
			api.ReportsItem{
				Type: api.ReportTupleReportsItem,
				ReportTuple: api.ReportTuple{
					Repository: owner(issue.RepositoryURL) + "/" + repository(issue.RepositoryURL),
					Section:    api.ReportTupleSectionReview,
					Age:        age(issue.CreatedAt),
					URL:        *issue.HTMLURL,
					Notes:      *issue.Title,
				},
			})
	}

	issues, _, err = gh.Search.Issues(context.Background(), "is:open is:pr review:approved author:"+settings.Username, &options)
	if err != nil {
		log.Fatalf("Could not fetch issues (%v)", err)
	} else {
		fmt.Println("Merge needed: ", len(issues.Issues))
	}

	for _, issue := range issues.Issues {
		reportIssues = append(reportIssues,
			api.ReportsItem{
				Type: api.ReportTupleReportsItem,
				ReportTuple: api.ReportTuple{
					Repository: owner(issue.RepositoryURL) + "/" + repository(issue.RepositoryURL),
					Section:    api.ReportTupleSectionMerge,
					Age:        age(issue.UpdatedAt),
					URL:        *issue.HTMLURL,
					Notes:      *issue.Title,
				},
			})
	}

	//
	clientReport.Reports = reportIssues
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
